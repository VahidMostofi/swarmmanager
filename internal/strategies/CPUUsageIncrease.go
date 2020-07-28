package strategies

import (
	"fmt"
	"log"

	"github.com/VahidMostofi/swarmmanager/internal/history"
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

func (c *CPUUsageIncrease) GetInitialConfig() (map[string]swarm.SimpleSpecs, error) {
	return make(map[string]swarm.SimpleSpecs), nil
}

// Configure ....
func (c *CPUUsageIncrease) Configure(values map[string]history.ServiceInfo, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	isChanged := false
	if c.Threshold < 1 {
		return nil, isChanged, fmt.Errorf("the Threshold value is not set for CPUUsageIncrease")
	}
	log.Println("Configurer Agent:", "increase replica count based on", c.ValueToConsider, "and the threshold is", c.Threshold)
	newSpecs := make(map[string]swarm.ServiceSpecs)

	for key := range currentState {
		newSpecs[key] = currentState[key]
		doMonitor := false
		for _, srv := range servicesToMonitor {
			if srv == currentState[key].Name {
				doMonitor = true
				break
			}
		}
		if !doMonitor {
			continue
		}
		var whatToCompareTo float64
		if c.ValueToConsider == "CPUUsageMean" {
			whatToCompareTo = values[currentState[key].Name].CPUUsageMean
		} else if c.ValueToConsider == "CPUUsage90Percentile" {
			whatToCompareTo = values[currentState[key].Name].CPUUsage90Percentile
		} else if c.ValueToConsider == "CPUUsage95Percentile" {
			whatToCompareTo = values[currentState[key].Name].CPUUsage95Percentile
		} else if c.ValueToConsider == "CPUUsage99Percentile" {
			whatToCompareTo = values[currentState[key].Name].CPUUsage99Percentile
		} else {
			return nil, false, fmt.Errorf("the PropertyToConsider is unknown: %s", c.ValueToConsider)
		}
		log.Println("Configurer Agent:", currentState[key].Name, c.ValueToConsider, "is", whatToCompareTo, "it should be less than or equal to", c.Threshold)
		if whatToCompareTo > c.Threshold {
			log.Println("Configurer Agent:", currentState[key].Name, "change replica count from", currentState[key].ReplicaCount, "to", currentState[key].ReplicaCount+1)
			temp := currentState[key]
			temp.ReplicaCount++
			newSpecs[key] = temp
			isChanged = true
		}
	}
	return newSpecs, isChanged, nil
}

// OnFeedbackCallback ...
func (c *CPUUsageIncrease) OnFeedbackCallback(map[string]history.ServiceInfo) error { return nil }
