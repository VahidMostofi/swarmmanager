package autoconfigure

import (
	"fmt"

	"github.com/VahidMostofi/swarmmanager/internal/swarm"
)

//TODO I NEED TO REstart every single container
// CPUUsageIncrease ...
// TODO: add options to work based on mean, or percentiles ...
type CPUUsageIncrease struct {
	Threshold float64
}

// Configure ....
func (c *CPUUsageIncrease) Configure(values map[string]ServiceInfo, currentState map[string]swarm.ServiceSpecs) (map[string]swarm.ServiceSpecs, error) {
	if c.Threshold < 1 {
		return nil, fmt.Errorf("the Threshold value is not set for CPUUsageIncrease")
	}
	for service := range currentState {
		fmt.Println(currentState[service].Name, values[service].CPUUsage90Percentile)
		if values[currentState[service].Name].CPUUsage90Percentile > c.Threshold {
			fmt.Println(service, "mean CPU usage for ", currentState[service].Name, " is more than ", c.Threshold, "% change replica count from", currentState[service].ReplicaCount, "to", currentState[service].ReplicaCount+1)
			temp := currentState[service]
			temp.ReplicaCount++
			currentState[service] = temp
		}
	}
	return currentState, nil
}
