package strategies

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"os"
// 	"os/exec"
// 	"strings"
// 	"sync"

// 	"github.com/VahidMostofi/swarmmanager/internal/history"
// 	"github.com/VahidMostofi/swarmmanager/internal/swarm"
// 	"github.com/VahidMostofi/swarmmanager/internal/utils"
// )

// // PythonPath is the path to python interpretor
// const PythonPath = "/home/vahid/envs/data/bin/python"

// // ScriptPath is the path to python script
// const ScriptPath = "/home/vahid/Desktop/projects/swarmmanager/scripts/mobo_CPU_split_mc.py"

// var wg sync.WaitGroup

// // MultiObjectiveBayesianOptimization ...
// type MultiObjectiveBayesianOptimization struct {
// 	InitialConfig    map[string]int
// 	ServicesToReport []string //TODO
// 	PropertyToReport []string //TODO
// 	index            int
// 	cmd              *exec.Cmd
// 	stdin            io.WriteCloser
// 	configCh         chan map[string]serviceConfig
// }

// type serviceConfig struct {
// 	CPUAmount      float64 `json:"cpu_count"`
// 	ContainerCount int     `json:"container_count"`
// 	WorkerCount    int     `json:"worker_count"`
// }

// type dataToSend struct {
// 	Feedbacks []float64 `json:"feedbacks"`
// }

// // GetnewMOBOConfigurer ...
// func GetnewMOBOConfigurer(initialConfig map[string]int) Configurer {
// 	log.Println("MOBO: creating a MOBO configurer with: ", initialConfig)
// 	return &MultiObjectiveBayesianOptimization{
// 		InitialConfig: initialConfig,
// 	}
// }

// func (c *MultiObjectiveBayesianOptimization) Write(p []byte) (int, error) {
// 	config := make(map[string]serviceConfig)
// 	err := json.Unmarshal(p, &config)
// 	if err != nil {
// 		log.Println("MOBO: non json response:", string(p))
// 		if strings.Trim(string(p), "\n") == "done" {
// 			log.Println("MOBO: Python is done")
// 			c.configCh <- config
// 			return len(p), nil
// 		} else {
// 			log.Println("MOBO: from python:", string(p))
// 			return len(p), nil
// 		}
// 	}
// 	c.configCh <- config
// 	return len(p), nil
// }

// // GetInitialConfig ...
// func (c *MultiObjectiveBayesianOptimization) GetInitialConfig() (map[string]swarm.SimpleSpecs, error) {
// 	config := make(map[string]swarm.SimpleSpecs)
// 	for key := range c.InitialConfig {
// 		temp := config[key]
// 		temp.CPU = 1
// 		temp.Replica = 1
// 		temp.Worker = 1
// 		config[key] = temp
// 	}
// 	return config, nil
// }

// // Configure ...
// func (c *MultiObjectiveBayesianOptimization) Configure(info history.Information, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
// 	isChanged := false
// 	if c.index == 0 {
// 		c.configCh = make(chan map[string]serviceConfig)
// 		log.Println("MOBO: first iteration of configurer")
// 		ctx, _ := context.WithCancel(context.Background())
// 		c.cmd = exec.CommandContext(ctx, PythonPath, "-W", "ignore", ScriptPath)
// 		stdin, err := c.cmd.StdinPipe()
// 		if err != nil {
// 			panic(err)
// 		}
// 		// defer stdin.Close()
// 		c.stdin = stdin
// 		c.cmd.Stdout = c
// 		c.cmd.Stderr = os.Stderr
// 		err = c.cmd.Start()
// 		log.Println("MOBO: started python program")
// 		if err != nil {
// 			panic(err)
// 		}
// 	} else {
// 		feedbacks := &dataToSend{Feedbacks: []float64{
// 			*values["auth"].ResponseTimes["total"].RTToleranceIntervalUBoundConfidence90p95,
// 			*values["books"].ResponseTimes["total"].RTToleranceIntervalUBoundConfidence90p95}}
// 		b, err := json.Marshal(feedbacks)
// 		if err != nil {
// 			panic(fmt.Errorf("error while converting feedbacks to json: %w", err))
// 		}
// 		log.Println("MOBO: sending feedback:", string(b)+"\n")
// 		io.WriteString(c.stdin, string(b)+"\n")
// 		log.Println("MOBO: sent feedback:", string(b))
// 	}
// 	var config map[string]serviceConfig
// 	log.Println("MOBO: waiting for config")
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() {
// 		for {
// 			select {
// 			case config = <-c.configCh:
// 				if len(config) > 0 {
// 					log.Println("MOBO: got the config")
// 				}
// 				wg.Done()
// 				return
// 			}
// 		}
// 	}()
// 	wg.Wait()
// 	if len(config) == 0 {
// 		return nil, false, nil
// 	}
// 	newSpecs := make(map[string]swarm.ServiceSpecs)
// 	for key, specs := range currentState {
// 		newSpecs[key] = currentState[key]
// 		doMonitor := false
// 		for _, serviceName := range servicesToMonitor {
// 			if specs.Name == serviceName {
// 				doMonitor = true
// 				break
// 			}
// 		}
// 		if !doMonitor {
// 			continue
// 		}
// 		temp := newSpecs[key]
// 		temp.CPULimits = config[specs.Name].CPUAmount
// 		temp.CPUReservation = config[specs.Name].CPUAmount
// 		temp.ReplicaCount = config[specs.Name].ContainerCount
// 		temp.EnvironmentVariables = utils.UpdateENVWorkerCounts(temp.EnvironmentVariables, config[specs.Name].WorkerCount)
// 		newSpecs[key] = temp
// 		isChanged = true
// 	}
// 	c.index++
// 	return newSpecs, isChanged, nil
// }

// // OnFeedbackCallback ...
// func (c *MultiObjectiveBayesianOptimization) OnFeedbackCallback(map[string]history.ServiceInfo) error {
// 	return nil
// }
