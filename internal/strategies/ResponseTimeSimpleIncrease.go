package strategies

import (
	"fmt"
	"log"

	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
)

// Agreement ...
type Agreement struct {
	PropertyToConsider string // ResponseTimesMean,ResponseTimes90Percentile,ResponseTimes95Percentile,ResponseTimes99Percentile,RTToleranceIntervalUBoundc90p95
	Value              float64
}

// ResponseTimeSimpleIncrease ...
type ResponseTimeSimpleIncrease struct {
	Agreements []Agreement
}

// OnFeedbackCallback ...
func (c *ResponseTimeSimpleIncrease) OnFeedbackCallback(map[string]history.ServiceInfo) error {
	return nil
}

// GetInitialConfig ...
func (c *ResponseTimeSimpleIncrease) GetInitialConfig() (map[string]swarm.SimpleSpecs, error) {
	return make(map[string]swarm.SimpleSpecs), nil
}

// Configure ...
func (c *ResponseTimeSimpleIncrease) Configure(values map[string]history.ServiceInfo, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	isChanged := false

	newSpecs := make(map[string]swarm.ServiceSpecs)
	for key, value := range currentState {
		newSpecs[key] = value
	}

	initialReplicaCounts := make(map[string]int) // I'm keeping this, but we can use currentSpecs as it is not chaning
	for key := range currentState {
		initialReplicaCounts[key] = currentState[key].ReplicaCount
	}

	for service := range currentState {
		doMonitor := false
		for _, srv := range servicesToMonitor {
			if srv == currentState[service].Name {
				doMonitor = true
				break
			}
		}
		if !doMonitor || currentState[service].Name == "gateway" { //TODO second part of the condition
			continue
		}
		isServiceChanged := false
		for _, ag := range c.Agreements {
			if isServiceChanged {
				break
			}

			var whatToCompareTo float64
			if ag.PropertyToConsider == "ResponseTimesMean" { //TODO all conditions need to be checked
				whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].ResponseTimesMean)
			} else if ag.PropertyToConsider == "ResponseTimes90Percentile" {
				whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].ResponseTimes90Percentile)
			} else if ag.PropertyToConsider == "ResponseTimes95Percentile" {
				whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].ResponseTimes95Percentile)
			} else if ag.PropertyToConsider == "ResponseTimes99Percentile" {
				whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].ResponseTimes99Percentile)
			} else if ag.PropertyToConsider == "RTToleranceIntervalUBoundc90p95" {
				whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].RTToleranceIntervalUBoundConfidence90p95)
			} else {
				return nil, false, fmt.Errorf("ResponseTimeSimpleIncrease: the PropertyToConsider is unknown: %s", ag.PropertyToConsider)
			}
			log.Println("Configurer Agent:", currentState[service].Name, ag.PropertyToConsider, "is", whatToCompareTo, "and should be less than or equal to", ag.Value)
			if ag.Value < whatToCompareTo {
				log.Println("Configurer Agent:", currentState[service].Name, "change replica count from", currentState[service].ReplicaCount, "to", currentState[service].ReplicaCount+1)

				temp := currentState[service]
				temp.ReplicaCount++
				newSpecs[service] = temp

				log.Println("Configurer Agent:", newSpecs["gateway"].Name, "change replica count from", newSpecs[service].ReplicaCount, "to", newSpecs[service].ReplicaCount+1)
				temp = newSpecs["gateway"]
				temp.ReplicaCount++
				newSpecs["gateway"] = temp

				isServiceChanged = true
				isChanged = true
			}
		}

	}

	for key := range newSpecs {
		if newSpecs[key].ReplicaCount-initialReplicaCounts[key] > 1 {
			log.Println("Configurer Agent:", newSpecs[key].Name, "replica count has increased", newSpecs[key].ReplicaCount-initialReplicaCounts[key], "changing the increase to 1")
			temp := newSpecs[key]
			temp.ReplicaCount = initialReplicaCounts[key] + 1
			newSpecs[key] = temp
		}
	}

	return newSpecs, isChanged, nil
}
