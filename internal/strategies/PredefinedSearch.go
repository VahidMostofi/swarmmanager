package strategies

import (
	"fmt"
	"log"

	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/VahidMostofi/swarmmanager/internal/utils"
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

// OnFeedbackCallback ...
func (c *PredefinedSearch) OnFeedbackCallback(map[string]history.ServiceInfo) error { return nil }

func (c *PredefinedSearch) GetInitialConfig() (map[string]swarm.SimpleSpecs, error) {
	return make(map[string]swarm.SimpleSpecs), nil
}

// Configure ...
// this is not stable! //TODO
func (c *PredefinedSearch) Configure(values map[string]history.ServiceInfo, currentSpecs map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	isChanged := false
	log.Println("Configurer Agent: configure at step", c.Step)
	if c.Step == 0 { // the input configuration
		c.PreviousSpecs = append(c.PreviousSpecs, currentSpecs)
		c.Step++
		return c.Configure(values, currentSpecs, servicesToMonitor)
	}
	if c.Step == 1 { // 1 container with multiple (= replica count) cores
		currentSpecs = clone(c.PreviousSpecs[0])
		for key := range currentSpecs {
			if !contains(servicesToMonitor, currentSpecs[key].Name) {
				continue
			}
			isChanged = true
			temp := currentSpecs[key]
			temp.CPUReservation = float64(temp.ReplicaCount)
			temp.CPULimits = float64(temp.ReplicaCount)
			temp.EnvironmentVariables = utils.UpdateENVWorkerCounts(temp.EnvironmentVariables, temp.ReplicaCount)
			temp.ReplicaCount = 1
			currentSpecs[key] = temp
		}
		c.PreviousSpecs = append(c.PreviousSpecs, currentSpecs)
		c.Step++
		return currentSpecs, isChanged, nil
	}
	if c.Step == 2 {
		currentSpecs = clone(c.PreviousSpecs[0])
		for key := range currentSpecs {
			if !contains(servicesToMonitor, currentSpecs[key].Name) {
				continue
			}
			isChanged = true
			temp := currentSpecs[key]
			temp.CPUReservation = float64(temp.ReplicaCount)
			temp.CPULimits = float64(temp.ReplicaCount)
			temp.EnvironmentVariables = utils.UpdateENVWorkerCounts(temp.EnvironmentVariables, temp.ReplicaCount*2)
			temp.ReplicaCount = 1
			currentSpecs[key] = temp
		}
		c.PreviousSpecs = append(c.PreviousSpecs, currentSpecs)
		c.Step++
		return currentSpecs, isChanged, nil
	}
	if c.Step == 3 {
		currentSpecs = clone(c.PreviousSpecs[0])
		for key := range currentSpecs {
			if !contains(servicesToMonitor, currentSpecs[key].Name) {
				continue
			}
			isChanged = true
			temp := currentSpecs[key]
			temp.CPUReservation = float64(temp.ReplicaCount)
			temp.CPULimits = float64(temp.ReplicaCount)
			temp.EnvironmentVariables = utils.UpdateENVWorkerCounts(temp.EnvironmentVariables, temp.ReplicaCount+1)
			temp.ReplicaCount = 1
			currentSpecs[key] = temp
		}
		c.PreviousSpecs = append(c.PreviousSpecs, currentSpecs)
		c.Step++
		return currentSpecs, isChanged, nil
	}
	if c.Step == 4 {
		currentSpecs = clone(c.PreviousSpecs[0])
		for key := range currentSpecs {
			if !contains(servicesToMonitor, currentSpecs[key].Name) {
				continue
			}
			if currentSpecs[key].ReplicaCount%2 != 0 {
				continue
			}
			isChanged = true
			temp := currentSpecs[key]
			temp.CPUReservation = float64(temp.ReplicaCount / 2)
			temp.CPULimits = float64(temp.ReplicaCount / 2)
			temp.EnvironmentVariables = utils.UpdateENVWorkerCounts(temp.EnvironmentVariables, temp.ReplicaCount/2)
			temp.ReplicaCount = 2
			currentSpecs[key] = temp
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
		for key := range currentSpecs {
			if !contains(servicesToMonitor, currentSpecs[key].Name) {
				continue
			}
			if currentSpecs[key].ReplicaCount%2 != 0 {
				continue
			}
			isChanged = true
			temp := currentSpecs[key]
			temp.CPUReservation = float64(temp.ReplicaCount / 2)
			temp.CPULimits = float64(temp.ReplicaCount / 2)
			temp.EnvironmentVariables = utils.UpdateENVWorkerCounts(temp.EnvironmentVariables, 1+temp.ReplicaCount/2)
			temp.ReplicaCount = 2
			currentSpecs[key] = temp
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
		for key := range currentSpecs {
			if !contains(servicesToMonitor, currentSpecs[key].Name) {
				continue
			}
			isChanged = true
			temp := currentSpecs[key]
			temp.CPUReservation = float64(temp.ReplicaCount) - 0.5
			temp.CPULimits = float64(temp.ReplicaCount) - 0.5
			temp.EnvironmentVariables = utils.UpdateENVWorkerCounts(temp.EnvironmentVariables, temp.ReplicaCount)
			temp.ReplicaCount = 1
			currentSpecs[key] = temp
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

func clone(m map[string]swarm.ServiceSpecs) map[string]swarm.ServiceSpecs {
	clone := make(map[string]swarm.ServiceSpecs)
	for key, specs := range m {
		clone[key] = specs
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
