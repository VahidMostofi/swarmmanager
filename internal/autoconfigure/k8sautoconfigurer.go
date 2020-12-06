package autoconfigure

import (
	"fmt"
	"log"
	"time"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/caching"
	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/k8s"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	"github.com/VahidMostofi/swarmmanager/internal/resource/collector"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/montanaflynn/stats"
)

// K8sAutoConfigurer ...
type K8sAutoConfigurer struct {
	Workload        string
	LoadGenerator   loadgenerator.LoadGenerator
	Database        caching.Database
	ConfigurerAgent strategies.Configurer
	K8sConnector    k8s.Connector
	K8sCPUCollector *collector.K8sResourceCollector
	TimingConfigs
}

// NewK8sAutoConfigurer ...
func NewK8sAutoConfigurer(lg loadgenerator.LoadGenerator, c strategies.Configurer, workload string, database caching.Database) *K8sAutoConfigurer {
	k := &K8sAutoConfigurer{
		LoadGenerator:   lg,
		ConfigurerAgent: c,
		Workload:        workload,
		Database:        database,
		TimingConfigs: TimingConfigs{
			IterationDuration:                 time.Duration(configs.GetConfig().Test.Duration),
			WaitAfterServicesAreReadyDuration: 10,
			WaitAfterLoadGeneratorStopped:     time.Duration(configs.GetConfig().Test.WaitAfterLoadGeneratorDone),
			WaitAfterLoadGeneratorStarted:     10,
		},
		K8sCPUCollector: &collector.K8sResourceCollector{},
	}

	k.K8sConnector = k8s.GetNewConnector("ssh", configs.GetConfig().Host.Host)

	return k
}

func simpleToFull(v map[string]swarm.SimpleSpecs) map[string]swarm.ServiceSpecs {
	res := make(map[string]swarm.ServiceSpecs)
	for key, value := range v {
		res[key] = swarm.ServiceSpecs{
			ReplicaCount:   value.Replica,
			CPULimits:      value.CPU,
			CPUReservation: value.CPU,
			Name:           key,
		}
	}
	return res
}

func toHumanReadable(m map[string]swarm.ServiceSpecs) map[string]swarm.ServiceSpecs {
	m2 := make(map[string]swarm.ServiceSpecs)
	for _, value := range m {
		flag := false
		for _, str := range configs.GetConfig().TestBed.ServicesToConfigure {
			if str == value.Name {
				flag = true
				break
			}
		}
		if flag {
			m2[value.Name] = value
		}
	}
	return m2
}

// Start ...
func (a *K8sAutoConfigurer) Start(name string, command string) {
	fmt.Println("lets start!")

	stackHistory := &history.ExecutionDetails{
		Name:     name,
		Workload: a.Workload,
		History:  make([]history.Information, 0),
		Command:  command,
		Config:   configs.GetConfig(),
	}

	currentCConfigSimple, err := a.ConfigurerAgent.GetInitialConfig(a.LoadGenerator.GetWorkload())
	currentConfig := simpleToFull(currentCConfigSimple)
	if err != nil {
		log.Panic(err)
	}

	var start int64
	var end int64
	var iteration int
	newSpecs := make(map[string]swarm.ServiceSpecs)
	a.LoadGenerator.Stop()
	for {
		iteration++
		log.Printf("ITERATION %d\n", iteration)

		info, err := a.Database.Retrieve(string(a.Workload), currentConfig)

		if err == nil {
			log.Println("Autoconfigurer: information is found for this configuration/workload")
		} else {
			a.K8sConnector.ApplyConfig(currentConfig)
			log.Println("wait to all pods be ready")
			for !a.K8sConnector.AreAllPodsRunning() {
				time.Sleep(1 * time.Second)
			}
			log.Println("all pods are ready")
			time.Sleep(a.WaitAfterServicesAreReadyDuration * time.Second)
			a.K8sCPUCollector = &collector.K8sResourceCollector{}
			a.K8sCPUCollector.Start()
			log.Println("start recording CPU stats")
			go a.LoadGenerator.Start()
			time.Sleep(20 * time.Second)
			log.Println("load generator started")
			start = time.Now().UnixNano() / 1e3
			time.Sleep(a.IterationDuration * time.Second)
			end = time.Now().UnixNano() / 1e3
			a.LoadGenerator.Stop()
			a.K8sCPUCollector.Stop()
			log.Println("stopping load generator")
			time.Sleep(a.WaitAfterLoadGeneratorStopped * time.Second)
			var lgFeedback map[string]interface{} = nil
			lgFeedback, err = a.LoadGenerator.GetFeedback()

			serviceInfo := make(map[string]history.ServiceInfo)
			for _, name := range configs.GetConfig().TestBed.ServicesToConfigure {
				cpuMean, err := stats.Mean(a.K8sCPUCollector.GetCPUValues()[name])
				if err != nil {
					panic(err)
				}
				serviceInfo[name] = history.ServiceInfo{
					Start:         start,
					End:           end,
					CPUUsageMean:  cpuMean,
					NumberOfCores: currentCConfigSimple[name].CPU,
					ReplicaCount:  currentCConfigSimple[name].Replica,
				}
			}

			info = history.Information{
				ServicesInfo:          serviceInfo,
				Specs:                 toHumanReadable(currentConfig),
				Workload:              a.Workload,
				RequestResponseTimes:  make(map[string]history.ResponseTimeStats),
				LoadGeneratorFeedback: lgFeedback,
			}

			for _, request := range []string{"get_book", "login", "edit_book"} {
				responseTimesProperties := lgFeedback["metrics"].(map[string]interface{})[request+"_duration_trend"].(map[string]interface{})
				p95 := responseTimesProperties["p(95)"].(float64)
				p90 := responseTimesProperties["p(90)"].(float64)
				mean := responseTimesProperties["avg"].(float64)
				if mean < 1 {
					log.Printf("ERROR response times for %s %f %f %f", request, p90, p95, mean)
					panic("ERROR")
				}
				info.RequestResponseTimes[request] = history.ResponseTimeStats{
					ResponseTimes95Percentile: &p95,
					ResponseTimes90Percentile: &p90,
					ResponseTimesMean:         &mean,
				}
			}

			hash, err := a.Database.Store(a.Workload, currentConfig, info)
			if err != nil {
				log.Panicf("Error while storing run information: %w", err)
			}
			info.HashCode = hash
			log.Println("saved as", hash)
		}
		stackHistory.History = append(stackHistory.History, info)
		isChanged := false
		newSpecs, isChanged, err = a.ConfigurerAgent.Configure(info, currentConfig, configs.GetConfig().TestBed.ServicesToConfigure)

		if _, ok := Validate(newSpecs); !ok {
			log.Println("new config is not valid, breaking out of loop")
			log.Println("stopping load generator")
			if configs.GetConfig().ContinuesRuns {
				a.LoadGenerator.Stop()
			}
			break
		}
		if !isChanged {
			log.Println("is changed is false, breaking out of loop")
			break
		}
		saveHistory(stackHistory)
		fmt.Println("partial results at:", configs.GetConfig().Results.Path+stackHistory.Name+".yml")

		currentConfig = newSpecs

		currentCConfigSimple = make(map[string]swarm.SimpleSpecs)
		for key, value := range currentConfig {
			currentCConfigSimple[key] = swarm.SimpleSpecs{CPU: value.CPULimits, Worker: 1, Replica: value.ReplicaCount}
		}
	}
	saveHistory(stackHistory)
	fmt.Println("final results at:", configs.GetConfig().Results.Path+stackHistory.Name+".yml")
}
