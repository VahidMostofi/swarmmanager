/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"os"

	"github.com/VahidMostofi/swarmmanager/cmd"

	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
)

func main() {
	if os.Args[1] == "theory" {
		fmt.Println("VAHID")
		startTheory()
		return
	}
	cmd.Execute()
}

type theoryWorkload struct {
}

func (t theoryWorkload) GetThroughput() float64 {
	return 10
}

func (t theoryWorkload) GetRequestProportion() map[string]float64 {
	return map[string]float64{
		"get_book": 0.7,
		"login":    0.3,
	}
}

func getMeanResponseTime(requestName string, alfas map[string]float64) float64 {
	if requestName == "login" {
		return 1000 * (((0.025 / alfas["gateway"]) / (1 - (((0.3 * 10 * 0.025) / alfas["gateway"]) + ((0.7 * 10 * 0.07) / alfas["gateway"])))) + ((0.18 / alfas["auth"]) / (1 - ((0.3 * 10 * 0.18) / alfas["auth"]))))
	} else if requestName == "get_book" {
		return 1000 * (((0.07 / alfas["gateway"]) / (1 - (((0.3 * 10 * 0.025) / alfas["gateway"]) + ((0.7 * 10 * 0.07) / alfas["gateway"])))) + ((0.1 / alfas["books"]) / (1 - ((0.7 * 10 * 0.1) / alfas["books"]))))
	}
	return -1
}

func getUtilization(serviceName string, alfas map[string]float64) float64 {
	if serviceName == "gateway" {
		return (((0.3 * 10 * 0.025) / alfas["gateway"]) + ((0.7 * 10 * 0.07) / alfas["gateway"]))
	} else if serviceName == "auth" {
		return ((0.3 * 10 * 0.18) / alfas["auth"])
	} else if serviceName == "books" {
		return ((0.7 * 10 * 0.1) / alfas["books"])
	}
	return -1
}

func startTheory() {
	strategy := strategies.PerPathActualUtilization{
		StepSize:          0.25,
		Agreements:        []strategies.Agreement{{"ResponseTimesMean", 250}},
		MultiContainer:    true,
		DemandsFilePath:   "/home/vahid/Dropbox/data/swarm-manager-data/demands/simple_three_theory.yaml",
		ConstantInit:      true,
		ConstantInitValue: 1.0,
	}
	services := []string{"auth", "books", "gateway"}
	strategy.Init()
	t := theoryWorkload{}
	currentConfig, err := strategy.GetInitialConfig(t)
	if err != nil {
		panic(err)
	}
	fmt.Println("initial configs:", currentConfig)

	currentState := make(map[string]swarm.ServiceSpecs)
	stepCount := 1
	for {
		alphas := make(map[string]float64)
		for serviceName, simpleConfig := range currentConfig {
			currentState[serviceName] = swarm.ServiceSpecs{
				ReplicaCount: simpleConfig.Replica,
				CPULimits:    simpleConfig.CPU,
			}
			alphas[serviceName] = simpleConfig.CPU * float64(simpleConfig.Replica)
		}

		info := history.Information{}
		info.RequestResponseTimes = make(map[string]history.ResponseTimeStats)
		for _, requestName := range []string{"login", "get_book"} {
			mrt := getMeanResponseTime(requestName, alphas)
			if mrt < 0 {
				panic("mrt is less than 0!!!!!!!!!!")
			}
			info.RequestResponseTimes[requestName] = history.ResponseTimeStats{ResponseTimesMean: &mrt}
		}

		info.ServicesInfo = make(map[string]history.ServiceInfo)
		for _, serviceName := range services {
			info.ServicesInfo[serviceName] = history.ServiceInfo{
				CPUUsageMean: getUtilization(serviceName, alphas),
			}
		}
		newState, isChanged, err := strategy.Configure(info, currentState, services)
		if err != nil {
			panic(err)
		}
		for serviceName, serviceState := range newState {
			ss := swarm.SimpleSpecs{}
			ss.CPU = serviceState.CPULimits
			ss.Replica = serviceState.ReplicaCount
			currentConfig[serviceName] = ss
		}
		if !isChanged {
			break
			// fmt.Println("DONE")
		}
		stepCount++
	}
	var totalCPU float64 = 0
	for serviceName, ss := range currentConfig {
		fmt.Println(serviceName, ss.CPU*float64(ss.Replica))
		totalCPU += ss.CPU * float64(ss.Replica)
	}
	fmt.Println(stepCount, totalCPU)
}
