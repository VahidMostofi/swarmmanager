package strategies

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"sort"
	"strconv"

	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/VahidMostofi/swarmmanager/internal/utils"
	"gopkg.in/yaml.v3"
)

// BottleNeckOnlyVersion2 ...
type BottleNeckOnlyVersion2 struct {
	RequestToServiceToEU  map[string]map[string]float64
	StepSize              float64
	Agreements            []Agreement
	MultiContainer        bool
	resource2StepSize     map[string]float64 //TODO feature for future
	path2LastTimeDecrease map[string]bool
	path2TriedBackward    map[string]float64
	initialized           bool
	DemandsFilePath       string
	demands               map[string]map[string]float64
	ConstantInit          bool
	ConstantInitValue     float64
	MinimumStepSize       float64
	MinimumCPUValue       float64 // A control value for running the theory experiments. We cant have alpha less than 1 in theory experiments, or the formulas become unsable
	AfterFound            int
	ChangedPreviously     []string // To make sure we don't repeat ourselves
	cache                 map[string]float64
	bestWhichMeets        float64
	iterationCount        int
	stage1Iterations      int
	stage2Iterations      int
}

// Init ...
func (c *BottleNeckOnlyVersion2) Init() error {
	c.initialized = true
	c.resource2StepSize = make(map[string]float64)
	c.path2LastTimeDecrease = make(map[string]bool)
	c.path2TriedBackward = make(map[string]float64)
	c.MinimumCPUValue = 0.5
	c.ChangedPreviously = make([]string, 0)
	c.cache = make(map[string]float64)
	c.bestWhichMeets = 10000000000

	// read file file
	b, err := ioutil.ReadFile(c.DemandsFilePath)
	if err != nil {
		log.Panic("Configurer Agent: cant find the demand files at: %s", c.DemandsFilePath)
	}
	c.demands = make(map[string]map[string]float64)
	yaml.Unmarshal(b, &c.demands)
	log.Println("Min step size is", c.MinimumStepSize)
	return nil
}

func round1(value float64) float64 {
	return math.Round(value*10) / 10
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

func (c *BottleNeckOnlyVersion2) getReconfiguredConfiguration(service2totalResource map[string]float64) map[string]swarm.SimpleSpecs {
	reconfiguredSpecs := make(map[string]swarm.SimpleSpecs)
	temp := make(map[string]float64)
	for service, cpu := range service2totalResource {
		temp[service] = round1(cpu)
	}
	service2totalResource = temp
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

func (c *BottleNeckOnlyVersion2) updateSpecsFromSimpleSpecs(current map[string]swarm.ServiceSpecs, delta map[string]swarm.SimpleSpecs) map[string]swarm.ServiceSpecs {
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
func (c *BottleNeckOnlyVersion2) GetInitialConfig(workload loadgenerator.Workload) (map[string]swarm.SimpleSpecs, error) {
	if !c.initialized {
		return nil, fmt.Errorf("Configurer Agent: the configurer need to be initialized, call Init()")
	}

	c.RequestToServiceToEU = make(map[string]map[string]float64)

	for requestName, requestProb := range workload.GetRequestProportion() {
		c.RequestToServiceToEU[requestName] = make(map[string]float64)
		for serviceName := range c.demands {
			c.RequestToServiceToEU[requestName][serviceName] = workload.GetThroughput() * requestProb * c.demands[serviceName][requestName] / float64(1000)
		}
	}

	for requestName := range c.RequestToServiceToEU {
		c.path2LastTimeDecrease[requestName] = false
		c.path2TriedBackward[requestName] = c.StepSize
	}

	totalAllocatedResources := make(map[string]float64) // total allocated CPU to each service initially
	if c.ConstantInit {
		for _, service2EUtilization := range c.RequestToServiceToEU {
			for serviceName := range service2EUtilization { //eu is estimated utilization
				totalAllocatedResources[serviceName] = c.ConstantInitValue
				c.resource2StepSize[serviceName] = c.StepSize
			}
		}
	} else {
		for _, service2EUtilization := range c.RequestToServiceToEU {
			for serviceName, eu := range service2EUtilization { //eu is estimated utilization
				if current, ok := totalAllocatedResources[serviceName]; ok {
					totalAllocatedResources[serviceName] = current + eu
				} else {
					totalAllocatedResources[serviceName] = eu
				}
			}
		}
	}
	initialConfig := c.getReconfiguredConfiguration(totalAllocatedResources)
	log.Println("Configurer Agent: providing initial config:", initialConfig)
	return initialConfig, nil
}

// Configure ....
func (c *BottleNeckOnlyVersion2) Configure(info history.Information, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	isChanged := false
	c.iterationCount++
	newSpecs := make(map[string]swarm.ServiceSpecs)
	for key, value := range currentState {
		newSpecs[key] = value
	}

	initialCPUCount := make(map[string]float64)
	newCPUCount := make(map[string]float64)
	for key := range currentState {
		initialCPUCount[key] = round1(currentState[key].CPULimits * float64(currentState[key].ReplicaCount))
		newCPUCount[key] = round1(currentState[key].CPULimits * float64(currentState[key].ReplicaCount))
	}
	allMeet := true
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
			allMeet = false
			var serviceWithMaxCPUUtil string
			var maxCPUUtil float64 = 0
			for serviceName, eu := range c.RequestToServiceToEU[requestName] {
				if eu > 0 {
					// log.Println("Configurer Agent:", serviceName, "is involved in", requestName, "the mean CPU Util is", info.ServicesInfo[serviceName].CPUUsageMean)
					if info.ServicesInfo[serviceName].CPUUsageMean > maxCPUUtil {
						maxCPUUtil = info.ServicesInfo[serviceName].CPUUsageMean
						serviceWithMaxCPUUtil = serviceName
					}
				}
			}
			log.Println("Configurer Agent:", serviceWithMaxCPUUtil, "has the max mean CPU Utilization")
			increaseValue := c.resource2StepSize[serviceWithMaxCPUUtil]

			if increaseValue > 0 {
				log.Println("Configurer Agent:", serviceWithMaxCPUUtil, "is part of", requestName, "stepSize for path(request)", requestName, "is", c.resource2StepSize[serviceWithMaxCPUUtil])
				prev := newCPUCount[serviceWithMaxCPUUtil]
				newCPUCount[serviceWithMaxCPUUtil] += increaseValue
				log.Println("Configurer Agent:", "updating total CPU count from", prev, "to", newCPUCount[serviceWithMaxCPUUtil])
				isChanged = true
			} else {
				return nil, false, fmt.Errorf("Configurer Agent: the increaseValue must be positive")
			}
		}
	}

	if allMeet || c.AfterFound > 0 {
		// fmt.Println("all meet is ", allMeet)
		marginalRequests := make(map[string]float64)
		resourcesWhichAreMax := make(map[string]float64)
		resource2requests := make(map[string][]string)
		service2MaxResponseTime := make(map[string]float64)
		for requestName, requestResponseTimes := range info.RequestResponseTimes {
			if *requestResponseTimes.ResponseTimes95Percentile > c.Agreements[0].Value*0.9 { //TODO not always mean
				// if *requestResponseTimes.ResponseTimesMean > c.Agreements[0].Value*0.9 { //TODO not always mean
				marginalRequests[requestName] = *requestResponseTimes.ResponseTimesMean
			}
			_, max := c.getMinMaxUtilizationsInRequestPath(requestName, info)
			resourcesWhichAreMax[max] = 0

			for serviceName, eu := range c.RequestToServiceToEU[requestName] {
				if _, ok := resource2requests[serviceName]; !ok {
					resource2requests[serviceName] = make([]string, 0)
				}

				if eu > 0 {
					if _, ok := service2MaxResponseTime[serviceName]; !ok {
						service2MaxResponseTime[serviceName] = 0
					}
					resource2requests[serviceName] = append(resource2requests[serviceName], requestName)
					service2MaxResponseTime[serviceName] = math.Max(service2MaxResponseTime[serviceName], *requestResponseTimes.ResponseTimesMean)
				}
			}

		}
		backwardCandidates := make([]string, 0)

		for _, serviceName := range servicesToMonitor {
			flag := true
			for criticalRequest := range marginalRequests {
				fmt.Println(criticalRequest)
				_, ok := resourcesWhichAreMax[serviceName]
				if utils.ContainsString(resource2requests[serviceName], criticalRequest) && ok {
					flag = false
				}
			}

			if _, exists := resourcesWhichAreMax[serviceName]; exists {
				flag = false
			}

			if math.Abs(newCPUCount[serviceName]-c.MinimumCPUValue) < 1e-5 {
				flag = false
			}

			if flag {
				backwardCandidates = append(backwardCandidates, serviceName)
			}
		}

		fmt.Println("backward candidates: ", backwardCandidates)
		if len(backwardCandidates) == 0 {
			// fmt.Println("no candidate for moving backward")
		} else {
			sort.Slice(backwardCandidates, func(i int, j int) bool {
				return c.resource2StepSize[backwardCandidates[i]] > c.resource2StepSize[backwardCandidates[j]]
				// return service2MaxResponseTime[backwardCandidates[i]] < service2MaxResponseTime[backwardCandidates[j]]
			})
			pruneCount := int(len(newCPUCount) / 3)
			if pruneCount <= 0 {
				pruneCount = 1
			}
			// fmt.Println(pruneCount)
			for cIdx := 0; cIdx < pruneCount; cIdx++ {
				if cIdx == len(backwardCandidates) {
					break
				}
				serviceToDecrease := backwardCandidates[cIdx]
				// fmt.Println(minCPUUtil)
				log.Println("Configurer Agent:", serviceToDecrease, "has the min mean CPU Utilization")
				c.resource2StepSize[serviceToDecrease] /= 2
				c.resource2StepSize[serviceToDecrease] = math.Max(c.resource2StepSize[serviceToDecrease], c.MinimumStepSize)
				decreaseValue := c.resource2StepSize[serviceToDecrease]
				c.path2TriedBackward[serviceToDecrease] = c.resource2StepSize[serviceToDecrease]
				if decreaseValue > 0 {
					// log.Println("Configurer Agent:", serviceToDecrease, "is part of", requestName, "stepSize for path(request)", requestName, "is", c.resource2StepSize[serviceToDecrease])
					prev := newCPUCount[serviceToDecrease]
					newCPUCount[serviceToDecrease] -= decreaseValue
					newCPUCount[serviceToDecrease] = math.Max(c.MinimumCPUValue, newCPUCount[serviceToDecrease])
					log.Println("Configurer Agent:", serviceToDecrease, "updating total CPU count from", prev, "to", newCPUCount[serviceToDecrease])
					isChanged = true
				} else {
					return nil, false, fmt.Errorf("Configurer Agent: the increaseValue must be positive")
				}
			}
		}
	}
	// isChanged = true

	var totalPrev float64 = 0
	var totalNew float64 = 0
	for service := range initialCPUCount {
		newCPUCount[service] = math.Min(newCPUCount[service], initialCPUCount[service]+c.StepSize)
		newCPUCount[service] = math.Max(newCPUCount[service], initialCPUCount[service]-c.StepSize/2)
		totalNew += newCPUCount[service]
		totalPrev += initialCPUCount[service]
	}

	if allMeet {
		if c.stage1Iterations == 0 { //need to know how many steps in stage one
			c.stage1Iterations = c.iterationCount
		}
	}

	if c.stage1Iterations > 0 { // we are in stage 2
		c.stage2Iterations++ // tracking iterations in stage 2

		fmt.Println("iterations", c.stage2Iterations, c.stage1Iterations)
		if c.stage2Iterations >= c.stage1Iterations {
			// fmt.Println("stage1", c.stage1Iterations, "stage2", c.stage2Iterations)
			return nil, false, nil
		}
	}

	if totalPrev < c.bestWhichMeets {
		c.bestWhichMeets = totalPrev
	}

	if allMeet || c.AfterFound > 0 {
		c.AfterFound++
		// if allMeet {
		// 	fmt.Println("meet all with", totalPrev)
		// }
	}
	// if allMeet && math.Abs(totalNew-totalPrev) <= c.MinimumStepSize*float64(len(servicesToMonitor))*1.01 { //minvalue and half of the min value
	// if c.AfterFound > 100000 { //minvalue and half of the min value
	// 	return nil, false, nil
	// }

	service2simpleConfig := c.getReconfiguredConfiguration(newCPUCount)
	newSpecs = c.updateSpecsFromSimpleSpecs(newSpecs, service2simpleConfig)
	h := c.hash(newSpecs, servicesToMonitor)
	if _, exists := c.cache[h]; exists {
		// return nil, false, nil
	} else {
		c.cache[h] = totalNew
	}
	// if totalNew > c.bestWhichMeets {
	// 	return nil, false, nil
	// }
	return newSpecs, isChanged, nil
}

func (c *BottleNeckOnlyVersion2) getMinMaxUtilizationsInRequestPath(requestName string, info history.Information) (string, string) {

	var minU float64 = 1000000
	var maxU float64 = 1000000
	var minResource = ""
	var maxResource = ""
	for serviceName, eu := range c.RequestToServiceToEU[requestName] {
		if eu > 0 {
			// log.Println("Configurer Agent:", serviceName, "is involved in", requestName, "the mean CPU Util is", info.ServicesInfo[serviceName].CPUUsageMean)
			if info.ServicesInfo[serviceName].CPUUsageMean < minU {
				minU = info.ServicesInfo[serviceName].CPUUsageMean
				minResource = serviceName
			}
			if info.ServicesInfo[serviceName].CPUUsageMean > maxU {
				maxU = info.ServicesInfo[serviceName].CPUUsageMean
				maxResource = serviceName
			}
		}
	}
	return minResource, maxResource
}

// OnFeedbackCallback ...
func (c *BottleNeckOnlyVersion2) OnFeedbackCallback(map[string]history.ServiceInfo) error {
	return nil
}

func (c *BottleNeckOnlyVersion2) hash(specs map[string]swarm.ServiceSpecs, servicesToMonitor []string) string {
	code := ""
	cpus := ""
	var keys []string
	for _, srv := range servicesToMonitor {
		keys = append(keys, srv)
	}
	sort.Strings(keys)
	for _, srv := range keys {
		count := specs[srv].CPULimits * float64(specs[srv].ReplicaCount)
		countStr := strconv.FormatFloat(count, 'f', 1, 64)
		cpus += countStr + "_"
	}
	code += cpus
	return code
}
