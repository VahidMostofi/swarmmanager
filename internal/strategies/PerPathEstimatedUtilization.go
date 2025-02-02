package strategies

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"

	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/VahidMostofi/swarmmanager/internal/utils"
	"gopkg.in/yaml.v3"
)

// PerPathEstimatedUtilization ...
type PerPathEstimatedUtilization struct {
	RequestToServiceToEU           map[string]map[string]float64
	NormalizedRequestToServiceToEU map[string]map[string]float64
	StepSize                       float64
	Agreements                     []Agreement
	MultiContainer                 bool
	path2StepSize                  map[string]float64 //TODO feature for future
	initialized                    bool
	DemandsFilePath                string
	demands                        map[string]map[string]float64
}

// Init ...
func (c *PerPathEstimatedUtilization) Init() error {
	c.initialized = true
	c.path2StepSize = make(map[string]float64)

	// read file file
	b, err := ioutil.ReadFile(c.DemandsFilePath)
	if err != nil {
		log.Panic("Configurer Agent: cant find the demand files at: %s", c.DemandsFilePath)
	}
	c.demands = make(map[string]map[string]float64)
	yaml.Unmarshal(b, &c.demands)

	return nil
}

func (c *PerPathEstimatedUtilization) getReconfiguredConfiguration(service2totalResource map[string]float64) map[string]swarm.SimpleSpecs {
	reconfiguredSpecs := make(map[string]swarm.SimpleSpecs)
	if c.MultiContainer {
		for service, totalCPU := range service2totalResource {
			replicaCount := int(math.Ceil(totalCPU))
			reconfiguredSpecs[service] = swarm.SimpleSpecs{
				CPU:     round2(float64(totalCPU / float64(replicaCount))),
				Replica: replicaCount,
				Worker:  1,
			}
		}
	} else {
		for service, totalCPU := range service2totalResource {
			reconfiguredSpecs[service] = swarm.SimpleSpecs{
				CPU:     totalCPU,
				Replica: 1,
				Worker:  int(math.Ceil(totalCPU)),
			}
		}
	}
	return reconfiguredSpecs
}

func (c *PerPathEstimatedUtilization) updateSpecsFromSimpleSpecs(current map[string]swarm.ServiceSpecs, delta map[string]swarm.SimpleSpecs) map[string]swarm.ServiceSpecs {
	for service, dChange := range delta {
		temp := current[service]
		temp.EnvironmentVariables = utils.UpdateENVWorkerCounts(temp.EnvironmentVariables, dChange.Worker)
		temp.CPULimits = dChange.CPU
		temp.CPUReservation = dChange.CPU
		temp.ReplicaCount = dChange.Replica
		current[service] = temp
	}
	return current
}

// GetInitialConfig ...
func (c *PerPathEstimatedUtilization) GetInitialConfig(workload loadgenerator.Workload) (map[string]swarm.SimpleSpecs, error) {
	if !c.initialized {
		return nil, fmt.Errorf("Configurer Agent: the configurer need to be initialized, call Init()")
	}

	c.RequestToServiceToEU = make(map[string]map[string]float64)
	c.NormalizedRequestToServiceToEU = make(map[string]map[string]float64)

	for requestName, requestProb := range workload.GetRequestProportion() {
		c.RequestToServiceToEU[requestName] = make(map[string]float64)
		for serviceName := range c.demands {
			c.RequestToServiceToEU[requestName][serviceName] = workload.GetThroughput() * requestProb * c.demands[serviceName][requestName] / float64(1000)
		}
	}

	for requestName := range c.RequestToServiceToEU {
		for serviceName, eu := range c.RequestToServiceToEU[requestName] {
			if eu > 0 {
				fmt.Println(requestName, serviceName, "eu", eu, "demand", c.demands[serviceName][requestName])
			}
		}
	}
	// fmt.Println(workload.GetThroughput())

	for requestName := range c.RequestToServiceToEU {
		var sum float64 = 0
		for _, eu := range c.RequestToServiceToEU[requestName] {
			sum += eu
		}

		c.NormalizedRequestToServiceToEU[requestName] = make(map[string]float64)

		for serviceName, eu := range c.RequestToServiceToEU[requestName] {
			c.NormalizedRequestToServiceToEU[requestName][serviceName] = eu / sum
		}
	}

	for requestName := range c.RequestToServiceToEU {
		c.path2StepSize[requestName] = c.StepSize
		log.Println(requestName, "step size is", c.path2StepSize[requestName])
	}

	totalAllocatedResources := make(map[string]float64) // total allocated CPU to each service initially

	for _, service2EUtilization := range c.RequestToServiceToEU {
		for serviceName, eu := range service2EUtilization { //eu is estimated utilization
			if current, ok := totalAllocatedResources[serviceName]; ok {
				totalAllocatedResources[serviceName] = current + eu
			} else {
				totalAllocatedResources[serviceName] = eu
			}
		}
	}
	initialConfig := c.getReconfiguredConfiguration(totalAllocatedResources)
	log.Println("Configurer Agent: providing initial config:", initialConfig)
	return initialConfig, nil
}

// Configure ....
func (c *PerPathEstimatedUtilization) Configure(info history.Information, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	isChanged := false

	newSpecs := make(map[string]swarm.ServiceSpecs)
	for key, value := range currentState {
		newSpecs[key] = value
	}

	initialCPUCount := make(map[string]float64)
	newCPUCount := make(map[string]float64)
	for key := range currentState {
		initialCPUCount[key] = currentState[key].CPULimits * float64(currentState[key].ReplicaCount)
		newCPUCount[key] = currentState[key].CPULimits * float64(currentState[key].ReplicaCount)
	}

	for requestName, requestResponseTimes := range info.RequestResponseTimes {
		ag := c.Agreements[0]
		if len(c.Agreements) > 1 {
			return nil, false, fmt.Errorf("Configurer Agent: only works with 1 agreement.")
		}

		whatToCompareTo, err := findWhatToCompareToForAgreement(ag, &requestResponseTimes) // this is for example the 95 percentile of the response times
		if err != nil {
			return nil, false, fmt.Errorf("Configurer Agent: cant figure out what to compare in SLA: %w", err)
		}
		log.Println("Configurer Agent:", "request (path) is", requestName, ",", ag.PropertyToConsider, "is", whatToCompareTo, "and should be less than or equal to", ag.Value)

		if ag.Value < whatToCompareTo { // the path is not meeting the SLA, so we need to add more resources
			for serviceName, neu := range c.NormalizedRequestToServiceToEU[requestName] {
				increaseValue := neu * c.StepSize
				if increaseValue > 0 {
					log.Println("Configurer Agent:", serviceName, "is part of", requestName, "stepSize for path(request)", requestName, "is", c.path2StepSize[requestName], "normalized CPU share for this service is", neu)
					prev := newCPUCount[serviceName]
					newCPUCount[serviceName] += increaseValue
					log.Println("Configurer Agent:", "updating total CPU count from", prev, "to", newCPUCount[serviceName])
					isChanged = true
				}
			}
		}
	}

	service2simpleConfig := c.getReconfiguredConfiguration(newCPUCount)
	newSpecs = c.updateSpecsFromSimpleSpecs(newSpecs, service2simpleConfig)
	return newSpecs, isChanged, nil
}

// OnFeedbackCallback ...
func (c *PerPathEstimatedUtilization) OnFeedbackCallback(map[string]history.ServiceInfo) error {
	return nil
}

func findWhatToCompareToForAgreement(ag Agreement, responseTimes *history.ResponseTimeStats) (float64, error) {
	var whatToCompareTo float64
	if ag.PropertyToConsider == "ResponseTimesMean" {
		whatToCompareTo = *(responseTimes.ResponseTimesMean)
	} else if ag.PropertyToConsider == "ResponseTimes90Percentile" {
		whatToCompareTo = *(responseTimes.ResponseTimes90Percentile)
	} else if ag.PropertyToConsider == "ResponseTimes95Percentile" {
		whatToCompareTo = *(responseTimes.ResponseTimes95Percentile)
	} else if ag.PropertyToConsider == "ResponseTimes99Percentile" {
		whatToCompareTo = *(responseTimes.ResponseTimes99Percentile)
	} else if ag.PropertyToConsider == "RTToleranceIntervalUBoundc90p95" {
		whatToCompareTo = *(responseTimes.RTToleranceIntervalUBoundConfidence90p95)
	} else {
		return 0, fmt.Errorf("Agreement's property is %s", ag.PropertyToConsider)
	}
	return whatToCompareTo, nil
}
