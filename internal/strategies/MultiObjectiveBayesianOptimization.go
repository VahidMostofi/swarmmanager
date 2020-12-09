package strategies

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/VahidMostofi/swarmmanager/internal/utils"
)

var wg sync.WaitGroup

// MultiObjectiveBayesianOptimization ...
type MultiObjectiveBayesianOptimization struct {
	InitialConfig    map[string]int
	index            int
	cmd              *exec.Cmd
	stdin            io.WriteCloser
	configCh         chan map[string]serviceConfig
	PythonPath       string
	PythonScriptPath string
}

type serviceConfig struct {
	CPUAmount      float64 `json:"cpu_count"`
	ContainerCount int     `json:"container_count"`
	WorkerCount    int     `json:"worker_count"`
}

type dataToSend struct {
	Feedbacks []float64 `json:"feedbacks"`
}

func (c *MultiObjectiveBayesianOptimization) Write(p []byte) (int, error) {
	config := make(map[string]serviceConfig)
	err := json.Unmarshal(p, &config)
	if err != nil {
		log.Println("MOBO: non json response:", string(p))
		if strings.Trim(string(p), "\n") == "done" {
			log.Println("MOBO: Python is done")
			c.configCh <- config
			return len(p), nil
		} else {
			log.Println("MOBO: from python:", string(p))
			return len(p), nil
		}
	}
	c.configCh <- config
	return len(p), nil
}

// GetInitialConfig ...
func (c *MultiObjectiveBayesianOptimization) GetInitialConfig(workload loadgenerator.Workload) (map[string]swarm.SimpleSpecs, error) {
	if c.InitialConfig == nil {
		c.InitialConfig = make(map[string]int)
		for _, serviceName := range configs.GetConfig().TestBed.ServicesToConfigure {
			c.InitialConfig[serviceName] = 1
		}
	}
	config := make(map[string]swarm.SimpleSpecs)
	for key := range c.InitialConfig {
		temp := config[key]
		temp.CPU = 1
		temp.Replica = 1
		temp.Worker = 1
		config[key] = temp
	}
	return config, nil
}

// Configure ...
func (c *MultiObjectiveBayesianOptimization) Configure(info history.Information, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	isChanged := false
	if c.index == 0 {
		c.configCh = make(chan map[string]serviceConfig)
		log.Println("MOBO: first iteration of configurer")
		ctx, _ := context.WithCancel(context.Background())
		c.cmd = exec.CommandContext(ctx, c.PythonPath, "-W", "ignore", c.PythonScriptPath)
		stdin, err := c.cmd.StdinPipe()
		if err != nil {
			panic(err)
		}
		// defer stdin.Close()
		c.stdin = stdin
		c.cmd.Stdout = c
		c.cmd.Stderr = os.Stderr
		err = c.cmd.Start()
		log.Println("MOBO: started python program")
		if err != nil {
			panic(err)
		}
	} else {
		values := make([]float64, 0)
		keys := make([]string, 0)
		for key := range info.RequestResponseTimes {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			rst := info.RequestResponseTimes[key]
			// log.Println("response time for", serviceName, "is", *rst.ResponseTimes95Percentile)
			if os.Args[1] == "theory" {
				values = append(values, *rst.ResponseTimesMean)
			} else {
				values = append(values, *rst.ResponseTimes95Percentile)
			}
		}
		// fmt.Println("values-response-times", keys, values)
		feedbacks := &dataToSend{Feedbacks: values}
		b, err := json.Marshal(feedbacks)
		if err != nil {
			panic(fmt.Errorf("error while converting feedbacks to json: %w", err))
		}
		log.Println("MOBO: sending feedback:", string(b)+"\n")
		io.WriteString(c.stdin, string(b)+"\n")
		log.Println("MOBO: sent feedback:", string(b))
	}
	var config map[string]serviceConfig
	log.Println("MOBO: waiting for config")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case config = <-c.configCh:
				if len(config) > 0 {
					log.Println("MOBO: got the config")
				}
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
	if len(config) == 0 {
		return nil, false, nil
	}
	newSpecs := make(map[string]swarm.ServiceSpecs)
	for key, specs := range currentState {
		specs.Name = key
		newSpecs[key] = currentState[key]
		doMonitor := false
		for _, serviceName := range servicesToMonitor {
			if specs.Name == serviceName {
				doMonitor = true
				break
			}
		}
		if !doMonitor {
			continue
		}
		temp := newSpecs[key]
		temp.CPULimits = config[specs.Name].CPUAmount
		temp.CPUReservation = config[specs.Name].CPUAmount
		temp.ReplicaCount = config[specs.Name].ContainerCount
		temp.EnvironmentVariables = utils.UpdateENVWorkerCounts(temp.EnvironmentVariables, config[specs.Name].WorkerCount)
		newSpecs[key] = temp
		isChanged = true
	}
	c.index++
	return newSpecs, isChanged, nil
}

// OnFeedbackCallback ...
func (c *MultiObjectiveBayesianOptimization) OnFeedbackCallback(map[string]history.ServiceInfo) error {
	return nil
}
