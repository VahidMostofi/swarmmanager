package strategies

// import (
// 	"fmt"
// 	"log"

// 	"github.com/VahidMostofi/swarmmanager/internal/history"
// 	"github.com/VahidMostofi/swarmmanager/internal/swarm"
// )

// // HybridCPUUtilResponseTimeSimpleIncrease ...
// type HybridCPUUtilResponseTimeSimpleIncrease struct {
// 	Agreements []Agreement
// }

// // OnFeedbackCallback ...
// func (c *HybridCPUUtilResponseTimeSimpleIncrease) OnFeedbackCallback(map[string]history.ServiceInfo) error {
// 	return nil
// }

// // GetInitialConfig ...
// func (c *HybridCPUUtilResponseTimeSimpleIncrease) GetInitialConfig() (map[string]swarm.SimpleSpecs, error) {
// 	return make(map[string]swarm.SimpleSpecs), nil
// }

// // Configure ...
// func (c *HybridCPUUtilResponseTimeSimpleIncrease) Configure(values map[string]history.ServiceInfo, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
// 	isChanged := false

// 	newSpecs := make(map[string]swarm.ServiceSpecs)
// 	for key, value := range currentState {
// 		newSpecs[key] = value
// 	}

// 	initialReplicaCounts := make(map[string]int) // I'm keeping this, but we can use currentSpecs as it is not chaning
// 	for key := range currentState {
// 		initialReplicaCounts[key] = currentState[key].ReplicaCount
// 	}

// 	for service := range currentState {
// 		doMonitor := false
// 		for _, srv := range servicesToMonitor {
// 			if srv == currentState[service].Name {
// 				doMonitor = true
// 				break
// 			}
// 		}
// 		if !doMonitor || currentState[service].Name == "gateway" { //TODO second part of the condition
// 			continue
// 		}
// 		isServiceChanged := false
// 		for _, ag := range c.Agreements {
// 			if isServiceChanged {
// 				break
// 			}

// 			var whatToCompareTo float64
// 			if ag.PropertyToConsider == "ResponseTimesMean" {
// 				whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].ResponseTimesMean)
// 			} else if ag.PropertyToConsider == "ResponseTimes90Percentile" {
// 				whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].ResponseTimes90Percentile)
// 			} else if ag.PropertyToConsider == "ResponseTimes95Percentile" {
// 				whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].ResponseTimes95Percentile)
// 			} else if ag.PropertyToConsider == "ResponseTimes99Percentile" {
// 				whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].ResponseTimes99Percentile)
// 			} else if ag.PropertyToConsider == "RTToleranceIntervalUBoundc90p95" {
// 				whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].RTToleranceIntervalUBoundConfidence90p95)
// 			} else {
// 				return nil, false, fmt.Errorf("HybridCPUUtilResponseTimeSimpleIncrease: the PropertyToConsider is unknown: %s", ag.PropertyToConsider)
// 			}
// 			log.Println("Configurer Agent:", currentState[service].Name, ag.PropertyToConsider, "is", whatToCompareTo, "and should be less than or equal to", ag.Value)
// 			if ag.Value < whatToCompareTo {
// 				log.Println("Configurer Agent:", currentState[service].Name, "change replica count from", currentState[service].ReplicaCount, "to", currentState[service].ReplicaCount+1)

// 				temp := currentState[service]
// 				temp.ReplicaCount++
// 				newSpecs[service] = temp

// 				log.Println("Configurer Agent:", newSpecs["gateway"].Name, "change replica count from", newSpecs["gateway"].ReplicaCount, "to", newSpecs["gateway"].ReplicaCount+1)
// 				temp = newSpecs["gateway"]
// 				temp.ReplicaCount++
// 				newSpecs["gateway"] = temp

// 				isServiceChanged = true
// 				isChanged = true
// 			}
// 		}

// 	}

// 	CPUThreshold := 80.0

// 	for key := range newSpecs {
// 		if newSpecs[key].ReplicaCount-initialReplicaCounts[key] > 1 {
// 			log.Println("Configurer Agent:", newSpecs[key].Name, "replica count has increased", newSpecs[key].ReplicaCount-initialReplicaCounts[key], "changing the increase to 1")
// 			temp := newSpecs[key]
// 			temp.ReplicaCount = initialReplicaCounts[key] + 1
// 			newSpecs[key] = temp
// 		}
// 	}

// 	for key := range newSpecs {
// 		if newSpecs[key].ReplicaCount-initialReplicaCounts[key] > 0 {
// 			if values[key].CPUUsage90Percentile < CPUThreshold {
// 				temp := newSpecs[key]
// 				log.Println("Configurer Agent:", newSpecs[key].Name, "replica count has increased by 1 but the CPU utilization is", values[key].CPUUsage90Percentile, "which is less than", CPUThreshold, "so we change it back to its initial value")
// 				temp.ReplicaCount = initialReplicaCounts[key]
// 				newSpecs[key] = temp
// 			}
// 		}
// 	}

// 	isChanged = false
// 	for key := range newSpecs {
// 		if newSpecs[key].ReplicaCount != currentState[key].ReplicaCount {
// 			isChanged = true
// 			break
// 		}
// 	}

// 	return newSpecs, isChanged, nil
// }
