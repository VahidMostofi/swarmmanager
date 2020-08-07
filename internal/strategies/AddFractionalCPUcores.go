package strategies

import (
	"fmt"
	"log"
	"math"

	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/VahidMostofi/swarmmanager/internal/utils"
)

// AddFractionalCPUcores ...
// TODO: add options to work based on mean, or percentiles ...
// TODO: increase if the utilization is less than something. Make sure you don't increase the count for db or jaeger
type AddFractionalCPUcores struct {
	EachStepIncrease float64
	Agreements       []Agreement
}

// GetInitialConfig ...
func (c *AddFractionalCPUcores) GetInitialConfig() (map[string]swarm.SimpleSpecs, error) {
	return make(map[string]swarm.SimpleSpecs), nil
}

// Configure ....
func (c *AddFractionalCPUcores) Configure(values map[string]history.ServiceInfo, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	isChanged := false

	newSpecs := make(map[string]swarm.ServiceSpecs)
	for key, value := range currentState {
		newSpecs[key] = value
	}

	initialCPUCount := make(map[string]float64)
	for key := range currentState {
		initialCPUCount[key] = currentState[key].CPULimits
		if currentState[key].CPULimits != currentState[key].CPUReservation {
			panic(fmt.Errorf("these two musth be equal"))
		}
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
			if ag.PropertyToConsider == "ResponseTimesMean" {
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
				return nil, false, fmt.Errorf("Configurer Agent: the PropertyToConsider is unknown: %s", ag.PropertyToConsider)
			}
			log.Println("Configurer Agent:", currentState[service].Name, ag.PropertyToConsider, "is", whatToCompareTo, "and should be less than or equal to", ag.Value)
			if ag.Value < whatToCompareTo {
				// Actually you need to use these codes, up to isChanged = true
				log.Println("Configurer Agent:", currentState[service].Name, "change CPU count from", currentState[service].CPULimits, "to", currentState[service].CPULimits+c.EachStepIncrease)

				temp := currentState[service]
				temp.CPULimits += c.EachStepIncrease
				temp.CPUReservation += c.EachStepIncrease
				newSpecs[service] = temp

				log.Println("Configurer Agent:", newSpecs["gateway"].Name, "change CPU count from", newSpecs["gateway"].CPULimits, "to", newSpecs["gateway"].CPULimits+c.EachStepIncrease)
				temp = newSpecs["gateway"]
				temp.CPULimits += c.EachStepIncrease
				temp.CPUReservation += c.EachStepIncrease
				newSpecs["gateway"] = temp

				isServiceChanged = true
				isChanged = true

				// if currentState[service].Name == "books" {
				// 	log.Println("Configurer Agent:", currentState[service].Name, "change CPU count from", currentState[service].CPULimits, "to", currentState[service].CPULimits+0.47)
				// 	temp := currentState[service]
				// 	temp.CPULimits += 0.47
				// 	temp.CPUReservation += 0.47
				// 	newSpecs[service] = temp

				// 	log.Println("Configurer Agent:", newSpecs["gateway"].Name, "change CPU count from", newSpecs["gateway"].CPULimits, "to", newSpecs["gateway"].CPULimits+0.29)
				// 	temp = newSpecs["gateway"]
				// 	temp.CPULimits += 0.29
				// 	temp.CPUReservation += 0.29
				// 	newSpecs["gateway"] = temp
				// }

				// if currentState[service].Name == "auth" {
				// 	log.Println("Configurer Agent:", currentState[service].Name, "change CPU count from", currentState[service].CPULimits, "to", currentState[service].CPULimits+0.24)
				// 	temp := currentState[service]
				// 	temp.CPULimits += 0.24
				// 	temp.CPUReservation += 0.24
				// 	newSpecs[service] = temp

				// 	log.Println("Configurer Agent:", newSpecs["gateway"].Name, "change CPU count from", newSpecs["gateway"].CPULimits, "to", newSpecs["gateway"].CPULimits+0.47)
				// 	temp = newSpecs["gateway"]
				// 	temp.CPULimits += 0.47
				// 	temp.CPUReservation += 0.47
				// 	newSpecs["gateway"] = temp
				// }
				// isServiceChanged = true
				// isChanged = true

			}
		}

	}

	for key := range newSpecs {
		if newSpecs[key].CPULimits-initialCPUCount[key] > c.EachStepIncrease {
			log.Println("Configurer Agent:", newSpecs[key].Name, "cpu count has increased", newSpecs[key].CPULimits-initialCPUCount[key], "changing the increase to", c.EachStepIncrease)
			temp := newSpecs[key]
			temp.CPULimits = initialCPUCount[key] + c.EachStepIncrease
			temp.CPUReservation = initialCPUCount[key] + c.EachStepIncrease
			newSpecs[key] = temp
		}
	}

	for key := range newSpecs {
		temp := newSpecs[key]
		temp.EnvironmentVariables = utils.UpdateENVWorkerCounts(newSpecs[key].EnvironmentVariables, int(math.Ceil(newSpecs[key].CPULimits)))
		newSpecs[key] = temp
		log.Println("Configurer Agent:", newSpecs[key].Name, "has cpu value", newSpecs[key].CPULimits, "change worker count to", int(math.Ceil(newSpecs[key].CPULimits)))
	}

	return newSpecs, isChanged, nil
}

// OnFeedbackCallback ...
func (c *AddFractionalCPUcores) OnFeedbackCallback(map[string]history.ServiceInfo) error { return nil }
