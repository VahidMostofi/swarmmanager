package autoconfigure

import (
	"fmt"
	"log"

	"github.com/VahidMostofi/swarmmanager/internal/swarm"
)

// Agreement ...
type Agreement struct {
	PropertyToConsider string // ResponseTimesMean,ResponseTimes90Percentile,ResponseTimes95Percentile,ResponseTimes99Percentile
	Value              float64
}

// ResponseTimeSimpleIncrease ...
type ResponseTimeSimpleIncrease struct {
	Agreements []Agreement
}

// Configure ...
func (rti *ResponseTimeSimpleIncrease) Configure(values map[string]ServiceInfo, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	isChanged := false

	initialReplicaCounts := make(map[string]int)
	for serviceID := range currentState {
		initialReplicaCounts[serviceID] = currentState[serviceID].ReplicaCount
	}

	for service := range currentState {
		doMonitor := false
		for _, srv := range servicesToMonitor {
			if srv == currentState[service].Name {
				doMonitor = true
				break
			}
		}
		if !doMonitor {
			continue
		}
		isServiceChanged := false
		for _, ag := range rti.Agreements {
			if isServiceChanged {
				break
			}

			var whatToCompareTo float64
			if ag.PropertyToConsider == "ResponseTimesMean" {
				whatToCompareTo = values[currentState[service].Name].ResponseTimesMean
			} else if ag.PropertyToConsider == "ResponseTimes90Percentile" {
				whatToCompareTo = values[currentState[service].Name].ResponseTimes90Percentile
			} else if ag.PropertyToConsider == "ResponseTimes95Percentile" {
				whatToCompareTo = values[currentState[service].Name].ResponseTimes95Percentile
			} else if ag.PropertyToConsider == "ResponseTimes99Percentile" {
				whatToCompareTo = values[currentState[service].Name].ResponseTimes99Percentile
			} else {
				return nil, false, fmt.Errorf("ResponseTimeSimpleIncrease: the PropertyToConsider is unknown: %s", ag.PropertyToConsider)
			}
			log.Println("Configurer Agent:", currentState[service].Name, ag.PropertyToConsider, "is", whatToCompareTo, "and should be less than or equal to", ag.Value)
			if ag.Value < whatToCompareTo {
				log.Println("Configurer Agent:", currentState[service].Name, "change replica count from", currentState[service].ReplicaCount, "to", currentState[service].ReplicaCount+1)

				temp := currentState[service]
				temp.ReplicaCount++
				currentState[service] = temp

				var gatewayID string
				for _, value := range currentState {
					if value.Name == "gateway" {
						gatewayID = value.ID
					}
				}
				log.Println("Configurer Agent:", currentState[gatewayID].Name, "change replica count from", currentState[service].ReplicaCount, "to", currentState[service].ReplicaCount+1)
				temp = currentState[gatewayID]
				temp.ReplicaCount++
				currentState[gatewayID] = temp

				isServiceChanged = true
				isChanged = true
			}
		}

	}

	for serviceID := range currentState {
		if currentState[serviceID].ReplicaCount-initialReplicaCounts[serviceID] > 1 {
			log.Println("Configurer Agent:", currentState[serviceID].Name, "replica count has increased", currentState[serviceID].ReplicaCount-initialReplicaCounts[serviceID], "changing the increase to 1")
			temp := currentState[serviceID]
			temp.ReplicaCount = initialReplicaCounts[serviceID] + 1
			currentState[serviceID] = temp
		}
	}

	return currentState, isChanged, nil
}
