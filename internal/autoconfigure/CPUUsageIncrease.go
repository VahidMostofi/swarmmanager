package autoconfigure

import (
	"fmt"
	"log"

	"github.com/VahidMostofi/swarmmanager/internal/swarm"
)

//TODO I NEED TO REstart every single container
// CPUUsageIncrease ...
// TODO: add options to work based on mean, or percentiles ...
// TODO: increase if the utilization is less than something. Make sure you don't increase the count for db or jaeger
type CPUUsageIncrease struct {
	Threshold       float64
	ValueToConsider string // CPUUsageMean,CPUUsage90Percentile 70-95, 99
}

// Configure ....
func (c *CPUUsageIncrease) Configure(values map[string]ServiceInfo, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	isChanged := false
	if c.Threshold < 1 {
		return nil, isChanged, fmt.Errorf("the Threshold value is not set for CPUUsageIncrease")
	}
	log.Println("Configurer Agent:", "increase replica count based on", c.ValueToConsider, "and the threshold is", c.Threshold)
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
		var whatToCompareTo float64
		if c.ValueToConsider == "CPUUsageMean" {
			whatToCompareTo = values[currentState[service].Name].CPUUsageMean
		} else if c.ValueToConsider == "CPUUsage90Percentile" {
			whatToCompareTo = values[currentState[service].Name].CPUUsage90Percentile
		} else if c.ValueToConsider == "CPUUsage95Percentile" {
			whatToCompareTo = values[currentState[service].Name].CPUUsage95Percentile
		} else if c.ValueToConsider == "CPUUsage99Percentile" {
			whatToCompareTo = values[currentState[service].Name].CPUUsage99Percentile
		} else {
			return nil, false, fmt.Errorf("the PropertyToConsider is unknown: %s", c.ValueToConsider)
		}
		log.Println("Configurer Agent:", currentState[service].Name, c.ValueToConsider, "is", whatToCompareTo, "it should be less than or equal to", c.Threshold)
		if values[currentState[service].Name].CPUUsage90Percentile > c.Threshold {
			log.Println("Configurer Agent:", currentState[service].Name, "change replica count from", currentState[service].ReplicaCount, "to", currentState[service].ReplicaCount+1)
			temp := currentState[service]
			temp.ReplicaCount++
			currentState[service] = temp
			isChanged = true
		}
	}
	return currentState, isChanged, nil
}
