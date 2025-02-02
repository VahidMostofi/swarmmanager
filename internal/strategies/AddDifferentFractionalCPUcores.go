package strategies

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
)

// AddDifferentFractionalCPUcores ...
type AddDifferentFractionalCPUcores struct {
	ServiceToAmount   map[string]float64
	MaxServiceIncease map[string]float64
	Agreements        []Agreement
	MultiContainer    bool
}

// GetInitialConfig ...
func (c *AddDifferentFractionalCPUcores) GetInitialConfig(workload loadgenerator.Workload) (map[string]swarm.SimpleSpecs, error) {
	return make(map[string]swarm.SimpleSpecs), nil
}

// Configure ....
func (c *AddDifferentFractionalCPUcores) Configure(info history.Information, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	// isChanged := false

	// newSpecs := make(map[string]swarm.ServiceSpecs)
	// for key, value := range currentState {
	// 	newSpecs[key] = value
	// }

	// initialCPUCount := make(map[string]float64)
	// newCPUCount := make(map[string]float64)
	// for key := range currentState {
	// 	initialCPUCount[key] = currentState[key].CPULimits * float64(currentState[key].ReplicaCount)
	// 	newCPUCount[key] = currentState[key].CPULimits * float64(currentState[key].ReplicaCount)
	// }

	// for service := range currentState {
	// 	doMonitor := false
	// 	for _, srv := range servicesToMonitor {
	// 		if srv == currentState[service].Name {
	// 			doMonitor = true
	// 			break
	// 		}
	// 	}
	// 	if !doMonitor || currentState[service].Name == "gateway" { //TODO second part of the condition
	// 		continue
	// 	}
	// 	isServiceChanged := false
	// 	for _, ag := range c.Agreements {
	// 		if isServiceChanged {
	// 			break
	// 		}

	// 		var whatToCompareTo float64
	// 		if ag.PropertyToConsider == "ResponseTimesMean" {
	// 			whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].ResponseTimesMean)
	// 		} else if ag.PropertyToConsider == "ResponseTimes90Percentile" {
	// 			whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].ResponseTimes90Percentile)
	// 		} else if ag.PropertyToConsider == "ResponseTimes95Percentile" {
	// 			whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].ResponseTimes95Percentile)
	// 		} else if ag.PropertyToConsider == "ResponseTimes99Percentile" {
	// 			whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].ResponseTimes99Percentile)
	// 		} else if ag.PropertyToConsider == "RTToleranceIntervalUBoundc90p95" {
	// 			whatToCompareTo = *(values[currentState[service].Name].ResponseTimes["total"].RTToleranceIntervalUBoundConfidence90p95)
	// 		} else {
	// 			return nil, false, fmt.Errorf("Configurer Agent: the PropertyToConsider is unknown: %s", ag.PropertyToConsider)
	// 		}
	// 		log.Println("Configurer Agent:", currentState[service].Name, ag.PropertyToConsider, "is", whatToCompareTo, "and should be less than or equal to", ag.Value)
	// 		if ag.Value < whatToCompareTo {

	// 			if !c.MultiContainer {
	// 				log.Println("Configurer Agent:", currentState[service].Name, "change CPU count from", currentState[service].CPULimits, "to", currentState[service].CPULimits+c.ServiceToAmount[service+".service"])

	// 				temp := currentState[service]
	// 				temp.CPULimits += c.ServiceToAmount[service+".service"]
	// 				temp.CPUReservation += c.ServiceToAmount[service+".service"]
	// 				newSpecs[service] = temp

	// 				log.Println("Configurer Agent:", newSpecs["gateway"].Name, "change CPU count from", newSpecs["gateway"].CPULimits, "to", newSpecs["gateway"].CPULimits+c.ServiceToAmount[service+".gateway"])
	// 				temp = newSpecs["gateway"]
	// 				temp.CPULimits += c.ServiceToAmount[service+".gateway"]
	// 				temp.CPUReservation += c.ServiceToAmount[service+".gateway"]
	// 				newSpecs["gateway"] = temp
	// 			} else {
	// 				newCPUCount[service] += c.ServiceToAmount[service+".service"]
	// 				log.Println("Configurer Agent:", currentState[service].Name, "change CPU count from", initialCPUCount[service], "to", newCPUCount[service])

	// 				newCPUCount["gateway"] += c.ServiceToAmount[service+".gateway"]
	// 				log.Println(c.ServiceToAmount[service+".gateway"])
	// 				log.Println("Configurer Agent:", currentState["gateway"].Name, "change CPU count from", initialCPUCount["gateway"], "to", newCPUCount["gateway"])
	// 			}

	// 			isServiceChanged = true
	// 			isChanged = true
	// 		}
	// 	}
	// }

	// for key := range newSpecs {
	// 	if !c.MultiContainer {
	// 		if newSpecs[key].CPULimits-initialCPUCount[key] > c.MaxServiceIncease[key] {
	// 			log.Println("Configurer Agent:", newSpecs[key].Name, "cpu count has increased", newSpecs[key].CPULimits-initialCPUCount[key], "changing the increase to", c.MaxServiceIncease[key])
	// 			temp := newSpecs[key]
	// 			temp.CPULimits = initialCPUCount[key] + c.MaxServiceIncease[key]
	// 			temp.CPUReservation = initialCPUCount[key] + c.MaxServiceIncease[key]
	// 			newSpecs[key] = temp
	// 		}
	// 	} else {
	// 		if newCPUCount[key]-initialCPUCount[key] > c.MaxServiceIncease[key] {
	// 			log.Println("Configurer Agent:", newSpecs[key].Name, "cpu count has increased", newCPUCount[key]-initialCPUCount[key], "changing the increase to", c.MaxServiceIncease[key])
	// 			newCPUCount[key] = initialCPUCount[key] + c.MaxServiceIncease[key]
	// 		}
	// 	}
	// }

	// if !c.MultiContainer {
	// 	for key := range newSpecs {
	// 		temp := newSpecs[key]
	// 		temp.EnvironmentVariables = utils.UpdateENVWorkerCounts(newSpecs[key].EnvironmentVariables, int(math.Ceil(newSpecs[key].CPULimits)))
	// 		newSpecs[key] = temp
	// 		log.Println("Configurer Agent:", newSpecs[key].Name, "has cpu value", newSpecs[key].CPULimits, "change worker count to", int(math.Ceil(newSpecs[key].CPULimits)))
	// 	}
	// } else {
	// 	for key := range newSpecs {
	// 		replicaCount := int(math.Ceil(newCPUCount[key]))
	// 		temp := newSpecs[key]
	// 		temp.CPULimits = float64(newCPUCount[key] / float64(replicaCount))
	// 		temp.CPUReservation = float64(newCPUCount[key] / float64(replicaCount))
	// 		temp.ReplicaCount = replicaCount
	// 		newSpecs[key] = temp
	// 		log.Println("Configurer Agent:", newSpecs[key].Name, "is going to use", newCPUCount[key], "CPU cores, which is", replicaCount, "containers each with", newSpecs[key].CPULimits, "cores")
	// 	}
	// }

	// return newSpecs, isChanged, nil
	return nil, false, nil
}

// OnFeedbackCallback ...
func (c *AddDifferentFractionalCPUcores) OnFeedbackCallback(map[string]history.ServiceInfo) error {
	return nil
}

func roundMap(values map[string]float64) map[string]float64 {
	newValues := make(map[string]float64)
	for key, value := range values {
		newValues[key] = math.Round(value*100) / 100
	}
	return newValues
}

// GetFractionalCPUIncreaseValues ...
func GetFractionalCPUIncreaseValues(workload, indicator string, stepSize float64) (map[string]float64, map[string]float64, error) {
	vus, err := strconv.ParseFloat(strings.Split(workload, "_")[0], 64)
	if err != nil {
		return nil, nil, fmt.Errorf("Cant parse number of VUS in workload: %s", strings.Split(workload, "_")[0])
	}
	sleepTime, err := strconv.ParseFloat(strings.Split(workload, "_")[3], 64)
	if err != nil {
		return nil, nil, fmt.Errorf("Cant parse number of sleepTime in workload: %s", strings.Split(workload, "_")[3])
	}
	authProb, err := strconv.ParseFloat(strings.Split(workload, "_")[2], 64)
	if err != nil {
		return nil, nil, fmt.Errorf("Cant parse number of authProb in workload: %s", strings.Split(workload, "_")[2])
	}
	if authProb >= 1 {
		return nil, nil, fmt.Errorf("authProb can't be more than 1, its: %f", authProb)
	}
	booksProb := 1 - authProb

	X := vus / sleepTime
	demands := map[string]float64{
		"auth.service":  74,
		"auth.gateway":  8,
		"books.service": 62,
		"books.gateway": 42,
	}
	values := make(map[string]float64)
	maxIncrease := make(map[string]float64)

	if strings.ToLower(indicator) == "demand" {
		var sumAll float64
		for _, value := range demands {
			sumAll += value
		}
		for key, demand := range demands {
			values[key] = (demand / sumAll) * stepSize
		}

		maxIncrease["auth"] = values["auth.service"]
		maxIncrease["books"] = values["books.service"]
		maxIncrease["gateway"] = values["auth.gateway"] + values["books.gateway"]
		return roundMap(values), roundMap(maxIncrease), nil
	} else if strings.ToLower(indicator) == "utilization" {
		maxIncrease["auth"] = X * (demands["auth.service"] * authProb)
		maxIncrease["books"] = X * (demands["books.service"] * booksProb)
		maxIncrease["gateway"] = X * (demands["auth.gateway"]*authProb + demands["books.gateway"]*booksProb)
		sumAll := maxIncrease["auth"] + maxIncrease["books"] + maxIncrease["gateway"]

		maxIncrease["auth"] = stepSize * (maxIncrease["auth"] / sumAll)
		maxIncrease["books"] = stepSize * (maxIncrease["books"] / sumAll)
		maxIncrease["gateway"] = stepSize * (maxIncrease["gateway"] / sumAll)

		values["auth.service"] = maxIncrease["auth"]
		values["auth.gateway"] = maxIncrease["gateway"]
		values["books.service"] = maxIncrease["books"]
		values["books.gateway"] = maxIncrease["gateway"]
		return roundMap(values), roundMap(maxIncrease), nil
	} else {
		return nil, nil, fmt.Errorf("Unknown type of indicator: %s possible values are: demand,utilization", indicator)
	}
}
