package main

import (
	"fmt"
	"time"

	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"gopkg.in/yaml.v3"
)

func main() {
	configurer := &strategies.MultiObjectiveBayesianOptimization{}
	time.Sleep(1 * time.Second)
	currentSpecs := make(map[string]swarm.ServiceSpecs)
	currentSpecs["gateway"] = swarm.ServiceSpecs{
		Name: "gateway",
		EnvironmentVariables: []string{"JAEGER_COLLECTOR_ENDPOINT=http://jaeger:14268/api/traces",
			"JAEGER_AGENT_HOST=jaeger",
			"WorkerCount=1",
			"JWT_KEY=someKeyIsGoodAndSomeOfThemBNoGEo1ioD!",
		},
		CPULimits:      1,
		CPUReservation: 1,
		ReplicaCount:   1,
	}
	currentSpecs["auth"] = swarm.ServiceSpecs{
		Name: "auth",
		EnvironmentVariables: []string{"JAEGER_COLLECTOR_ENDPOINT=http://jaeger:14268/api/traces",
			"JAEGER_AGENT_HOST=jaeger",
			"WorkerCount=1",
			"JWT_KEY=someKeyIsGoodAndSomeOfThemBNoGEo1ioD!",
		},
		CPULimits:      1,
		CPUReservation: 1,
		ReplicaCount:   1,
	}
	currentSpecs["books"] = swarm.ServiceSpecs{
		Name: "books",
		EnvironmentVariables: []string{"JAEGER_COLLECTOR_ENDPOINT=http://jaeger:14268/api/traces",
			"JAEGER_AGENT_HOST=jaeger",
			"WorkerCount=1",
			"JWT_KEY=someKeyIsGoodAndSomeOfThemBNoGEo1ioD!",
		},
		CPULimits:      1,
		CPUReservation: 1,
		ReplicaCount:   1,
	}
	for {
		newSpecs, isChanged, err := configurer.Configure(nil, currentSpecs, []string{"gateway", "books", "auth"})
		fmt.Println(isChanged)
		if err != nil {
			panic(err)
		}
		if !isChanged {
			break
		}
		b, _ := yaml.Marshal(newSpecs)
		fmt.Println(string(b))
		time.Sleep(1 * time.Second)
	}

}
