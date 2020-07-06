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
	Threshold float64
}

// Configure ....
func (c *CPUUsageIncrease) Configure(values map[string]ServiceInfo, currentState map[string]swarm.ServiceSpecs) (map[string]swarm.ServiceSpecs, bool, error) {
	isChanged := false
	if c.Threshold < 1 {
		return nil, isChanged, fmt.Errorf("the Threshold value is not set for CPUUsageIncrease")
	}
	log.Println("Configurer Agent: the threshold is", c.Threshold)
	for service := range currentState {
		log.Println("Configurer Agent: 90 percentile of CPU usage", currentState[service].Name, values[currentState[service].Name].CPUUsage90Percentile)
		if values[currentState[service].Name].CPUUsage90Percentile > c.Threshold {
			log.Println("Configurer Agent: 90 percentile CPU usage for ", currentState[service].Name, " is more than ", c.Threshold, "% change replica count from", currentState[service].ReplicaCount, "to", currentState[service].ReplicaCount+1)
			temp := currentState[service]
			temp.ReplicaCount++
			currentState[service] = temp
			isChanged = true
		}
	}
	return currentState, isChanged, nil
}
