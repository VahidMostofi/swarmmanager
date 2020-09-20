package strategies

// import (
// 	"fmt"
// 	"log"
// 	"math"

// 	"github.com/VahidMostofi/swarmmanager/internal/history"
// 	"github.com/VahidMostofi/swarmmanager/internal/swarm"
// 	"github.com/VahidMostofi/swarmmanager/internal/utils"
// )

// // PerPathEUBasedScaling ...
// type PerPathEUBasedScaling struct {
// 	Path2Service2EUtilization           map[string]map[string]float64 //EUtilization stands for Estimated Utilization
// 	NormalizedPath2Service2EUtilization map[string]map[string]float64
// 	StepSize                            float64
// 	Agreements                          []Agreement
// 	MultiContainer                      bool
// 	path2StepSize                       map[string]float64
// 	initialized                         bool
// 	// path2BestBadConfig                  map[string]*pathConfigPerformance
// 	// path2CheapestGoodConfig             map[string]*pathConfigPerformance
// }

// // Init ...
// func (c *PerPathEUBasedScaling) Init() error {
// 	c.initialized = true
// 	c.path2StepSize = make(map[string]float64)
// 	for path := range c.Path2Service2EUtilization {
// 		c.path2StepSize[path] = c.StepSize
// 	}

// 	for path := range c.NormalizedPath2Service2EUtilization {
// 		var sum float64
// 		for _, value := range c.NormalizedPath2Service2EUtilization[path] {
// 			sum += value
// 		}
// 		if math.Abs(sum-1) > 1e-4 {
// 			return fmt.Errorf("path %s is not normalized, the sum is %f", path, sum)
// 		}
// 	}
// 	return nil
// }

// func (c *PerPathEUBasedScaling) getReconfiguredConfiguration(service2totalResource map[string]float64) map[string]swarm.SimpleSpecs {
// 	reconfiguredSpecs := make(map[string]swarm.SimpleSpecs)
// 	if c.MultiContainer {
// 		for service, totalCPU := range service2totalResource {
// 			replicaCount := int(math.Ceil(totalCPU))
// 			reconfiguredSpecs[service] = swarm.SimpleSpecs{
// 				CPU:     float64(totalCPU / float64(replicaCount)),
// 				Replica: replicaCount,
// 				Worker:  1,
// 			}
// 		}
// 	} else {
// 		for service, totalCPU := range service2totalResource {
// 			reconfiguredSpecs[service] = swarm.SimpleSpecs{
// 				CPU:     totalCPU,
// 				Replica: 1,
// 				Worker:  int(math.Ceil(totalCPU)),
// 			}
// 		}
// 	}
// 	return reconfiguredSpecs
// }

// func (c *PerPathEUBasedScaling) updateSpecsFromSimpleSpecs(current map[string]swarm.ServiceSpecs, delta map[string]swarm.SimpleSpecs) map[string]swarm.ServiceSpecs {
// 	for service, dChange := range delta {
// 		temp := current[service]
// 		temp.EnvironmentVariables = utils.UpdateENVWorkerCounts(temp.EnvironmentVariables, dChange.Worker)
// 		temp.CPULimits = dChange.CPU
// 		temp.CPUReservation = dChange.CPU
// 		temp.ReplicaCount = dChange.Replica
// 		current[service] = temp
// 	}
// 	return current
// }

// // GetInitialConfig ...
// func (c *PerPathEUBasedScaling) GetInitialConfig() (map[string]swarm.SimpleSpecs, error) {
// 	if !c.initialized {
// 		return nil, fmt.Errorf("PerPathEUBasedScaling: the configurer need to be initialized, call Init()")
// 	}
// 	totalAllocatedResources := make(map[string]float64) // total allocated CPU to each service initially

// 	for _, service2EUtilization := range c.Path2Service2EUtilization {
// 		for service, estimatedUtilization := range service2EUtilization {
// 			if current, ok := totalAllocatedResources[service]; ok {
// 				totalAllocatedResources[service] = current + estimatedUtilization
// 			} else {
// 				totalAllocatedResources[service] = estimatedUtilization
// 			}
// 		}
// 	}
// 	initialConfig := c.getReconfiguredConfiguration(totalAllocatedResources)
// 	log.Println("Configurer Agent: providing initial config:", initialConfig)
// 	return initialConfig, nil
// }

// // Configure ....
// func (c *PerPathEUBasedScaling) Configure(values map[string]history.ServiceInfo, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
// 	isChanged := false

// 	newSpecs := make(map[string]swarm.ServiceSpecs)
// 	for key, value := range currentState {
// 		newSpecs[key] = value
// 	}

// 	initialCPUCount := make(map[string]float64)
// 	newCPUCount := make(map[string]float64)
// 	for key := range currentState {
// 		initialCPUCount[key] = currentState[key].CPULimits * float64(currentState[key].ReplicaCount)
// 		newCPUCount[key] = currentState[key].CPULimits * float64(currentState[key].ReplicaCount)
// 	}

// 	for path := range c.Path2Service2EUtilization {
// 		ag := c.Agreements[0]
// 		if len(c.Agreements) > 1 {
// 			log.Panic("only works with one agreement!")
// 		}

// 		var whatToCompareTo float64
// 		if ag.PropertyToConsider == "ResponseTimesMean" {
// 			whatToCompareTo = *(values[currentState[path].Name].ResponseTimes["total"].ResponseTimesMean)
// 		} else if ag.PropertyToConsider == "ResponseTimes90Percentile" {
// 			whatToCompareTo = *(values[currentState[path].Name].ResponseTimes["total"].ResponseTimes90Percentile)
// 		} else if ag.PropertyToConsider == "ResponseTimes95Percentile" {
// 			whatToCompareTo = *(values[currentState[path].Name].ResponseTimes["total"].ResponseTimes95Percentile)
// 		} else if ag.PropertyToConsider == "ResponseTimes99Percentile" {
// 			whatToCompareTo = *(values[currentState[path].Name].ResponseTimes["total"].ResponseTimes99Percentile)
// 		} else if ag.PropertyToConsider == "RTToleranceIntervalUBoundc90p95" {
// 			whatToCompareTo = *(values[currentState[path].Name].ResponseTimes["total"].RTToleranceIntervalUBoundConfidence90p95)
// 		} else {
// 			return nil, false, fmt.Errorf("Configurer Agent: the PropertyToConsider is unknown: %s", ag.PropertyToConsider)
// 		}
// 		log.Println("Configurer Agent:", currentState[path].Name, ag.PropertyToConsider, "is", whatToCompareTo, "and should be less than or equal to", ag.Value)

// 		if ag.Value < whatToCompareTo {
// 			// Updating bestBadConfig
// 			// if (c.path2BestBadConfig[path] == nil) || (c.path2BestBadConfig[path] != nil && c.path2BestBadConfig[path].responseTime > whatToCompareTo) {
// 			// 	var totalCPUCount float64
// 			// 	for service := range c.Path2Service2EUtilization {
// 			// 		totalCPUCount += initialCPUCount[service]
// 			// 	}
// 			// 	c.path2BestBadConfig[path] = &pathConfigPerformance{
// 			// 		name:          path,
// 			// 		responseTime:  whatToCompareTo,
// 			// 		totalCPUCount: totalCPUCount,
// 			// 	}
// 			// }

// 			// Reconfiguring
// 			for service, normalizedEU := range c.NormalizedPath2Service2EUtilization[path] {
// 				log.Println("Configurer Agent:", service, "is part of", path, "stepSize for path", path, "is", c.path2StepSize[path], "normalized CPU share for this service is", normalizedEU)
// 				prev := newCPUCount[service]
// 				newCPUCount[service] += normalizedEU * c.path2StepSize[path]
// 				log.Println("Configurer Agent:", "updating total CPU count from", prev, "to", newCPUCount[service])
// 				isChanged = true
// 			}
// 		} //else {
// 		// 	// Updating cheapsetGoodConfig
// 		// 	var totalCPUCount float64
// 		// 	for service := range c.Path2Service2EUtilization {
// 		// 		totalCPUCount += initialCPUCount[service]
// 		// 	}
// 		// 	if (c.path2CheapestGoodConfig[path] == nil) || (c.path2CheapestGoodConfig[path] != nil && c.path2CheapestGoodConfig[path].totalCPUCount > totalCPUCount) {
// 		// 		c.path2CheapestGoodConfig[path] = &pathConfigPerformance{
// 		// 			name:          path,
// 		// 			responseTime:  whatToCompareTo,
// 		// 			totalCPUCount: totalCPUCount,
// 		// 		}
// 		// 	}
// 		// }

// 	}
// 	service2simpleConfig := c.getReconfiguredConfiguration(newCPUCount)
// 	newSpecs = c.updateSpecsFromSimpleSpecs(newSpecs, service2simpleConfig)
// 	return newSpecs, isChanged, nil
// }

// // OnFeedbackCallback ...
// func (c *PerPathEUBasedScaling) OnFeedbackCallback(map[string]history.ServiceInfo) error {
// 	return nil
// }

// // type pathConfigPerformance struct {
// // 	name          string
// // 	responseTime  float64
// // 	totalCPUCount float64
// // }
