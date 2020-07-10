package autoconfigure

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/VahidMostofi/swarmmanager/internal/swarm"
)

// PredefinedSearch ...
type PredefinedSearch struct {
	PreviousSpecs []map[string]swarm.ServiceSpecs
	Step          int
}

// GetNewPredefinedSearcher ...
func GetNewPredefinedSearcher() *PredefinedSearch {
	return &PredefinedSearch{}
}

// Configure ...
func (c *PredefinedSearch) Configure(values map[string]ServiceInfo, currentSpecs map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	isChanged := false
	log.Println("Configurer Agent: configure at step", c.Step)
	if c.Step == 0 {
		c.PreviousSpecs = append(c.PreviousSpecs, currentSpecs)
		c.Step++
		return c.Configure(values, currentSpecs, servicesToMonitor)
	}
	if c.Step == 1 {
		currentSpecs = clone(c.PreviousSpecs[0])
		for serviceID := range currentSpecs {
			if !contains(servicesToMonitor, currentSpecs[serviceID].Name) {
				continue
			}
			isChanged = true
			temp := currentSpecs[serviceID]
			temp.CPUReservation = float64(temp.ReplicaCount)
			temp.CPULimits = float64(temp.ReplicaCount)
			temp.EnvironmentVariables = updateENVWorkerCounts(temp.EnvironmentVariables, temp.ReplicaCount)
			temp.ReplicaCount = 1
			currentSpecs[serviceID] = temp
		}
		c.PreviousSpecs = append(c.PreviousSpecs, currentSpecs)
		c.Step++
		return currentSpecs, isChanged, nil
	}
	if c.Step == 2 {
		currentSpecs = clone(c.PreviousSpecs[0])
		for serviceID := range currentSpecs {
			if !contains(servicesToMonitor, currentSpecs[serviceID].Name) {
				continue
			}
			isChanged = true
			temp := currentSpecs[serviceID]
			temp.CPUReservation = float64(temp.ReplicaCount)
			temp.CPULimits = float64(temp.ReplicaCount)
			temp.EnvironmentVariables = updateENVWorkerCounts(temp.EnvironmentVariables, temp.ReplicaCount*2)
			temp.ReplicaCount = 1
			currentSpecs[serviceID] = temp
		}
		c.PreviousSpecs = append(c.PreviousSpecs, currentSpecs)
		c.Step++
		return currentSpecs, isChanged, nil
	}
	if c.Step == 3 {
		currentSpecs = clone(c.PreviousSpecs[0])
		for serviceID := range currentSpecs {
			if !contains(servicesToMonitor, currentSpecs[serviceID].Name) {
				continue
			}
			isChanged = true
			temp := currentSpecs[serviceID]
			temp.CPUReservation = float64(temp.ReplicaCount)
			temp.CPULimits = float64(temp.ReplicaCount)
			temp.EnvironmentVariables = updateENVWorkerCounts(temp.EnvironmentVariables, temp.ReplicaCount+1)
			temp.ReplicaCount = 1
			currentSpecs[serviceID] = temp
		}
		c.PreviousSpecs = append(c.PreviousSpecs, currentSpecs)
		c.Step++
		return currentSpecs, isChanged, nil
	}
	if c.Step == 4 {
		currentSpecs = clone(c.PreviousSpecs[0])
		for serviceID := range currentSpecs {
			if !contains(servicesToMonitor, currentSpecs[serviceID].Name) {
				continue
			}
			if currentSpecs[serviceID].ReplicaCount%2 != 0 {
				continue
			}
			isChanged = true
			temp := currentSpecs[serviceID]
			temp.CPUReservation = float64(temp.ReplicaCount / 2)
			temp.CPULimits = float64(temp.ReplicaCount / 2)
			temp.EnvironmentVariables = updateENVWorkerCounts(temp.EnvironmentVariables, temp.ReplicaCount/2)
			temp.ReplicaCount = 2
			currentSpecs[serviceID] = temp
		}
		c.PreviousSpecs = append(c.PreviousSpecs, currentSpecs)
		c.Step++
		if !isChanged {
			return c.Configure(values, currentSpecs, servicesToMonitor)
		}
		return currentSpecs, isChanged, nil
	}
	if c.Step == 5 {
		currentSpecs = clone(c.PreviousSpecs[0])
		for serviceID := range currentSpecs {
			if !contains(servicesToMonitor, currentSpecs[serviceID].Name) {
				continue
			}
			if currentSpecs[serviceID].ReplicaCount%2 != 0 {
				continue
			}
			isChanged = true
			temp := currentSpecs[serviceID]
			temp.CPUReservation = float64(temp.ReplicaCount / 2)
			temp.CPULimits = float64(temp.ReplicaCount / 2)
			temp.EnvironmentVariables = updateENVWorkerCounts(temp.EnvironmentVariables, 1+temp.ReplicaCount/2)
			temp.ReplicaCount = 2
			currentSpecs[serviceID] = temp
		}
		c.PreviousSpecs = append(c.PreviousSpecs, currentSpecs)
		c.Step++
		if !isChanged {
			return c.Configure(values, currentSpecs, servicesToMonitor)
		}
		return currentSpecs, isChanged, nil
	}
	if c.Step == 6 {
		currentSpecs = clone(c.PreviousSpecs[0])
		for serviceID := range currentSpecs {
			if !contains(servicesToMonitor, currentSpecs[serviceID].Name) {
				continue
			}
			isChanged = true
			temp := currentSpecs[serviceID]
			temp.CPUReservation = float64(temp.ReplicaCount) - 0.5
			temp.CPULimits = float64(temp.ReplicaCount) - 0.5
			temp.EnvironmentVariables = updateENVWorkerCounts(temp.EnvironmentVariables, temp.ReplicaCount)
			temp.ReplicaCount = 1
			currentSpecs[serviceID] = temp
		}
		c.PreviousSpecs = append(c.PreviousSpecs, currentSpecs)
		c.Step++
		if !isChanged {
			return c.Configure(values, currentSpecs, servicesToMonitor)
		}
		return currentSpecs, isChanged, nil
	}
	if c.Step == 7 {
		return currentSpecs, isChanged, nil
	}
	return nil, false, fmt.Errorf("PredefinedSearch: it shouldn't be here")
}

func updateENVWorkerCounts(envs []string, count int) []string {
	newEnvs := make([]string, len(envs))
	for i, env := range envs {
		if strings.Contains(env, "WorkerCount") {
			newEnv := "WorkerCount=" + strconv.Itoa(count)
			newEnvs[i] = newEnv
		} else {
			newEnvs[i] = envs[i]
		}
	}
	return newEnvs
}

func clone(m map[string]swarm.ServiceSpecs) map[string]swarm.ServiceSpecs {
	clone := make(map[string]swarm.ServiceSpecs)
	for serviceID, specs := range m {
		clone[serviceID] = specs
	}
	return clone
}

func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
