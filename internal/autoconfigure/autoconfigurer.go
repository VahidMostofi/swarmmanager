package autoconfigure

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"log"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/caching"
	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/jaeger"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	r2 "github.com/VahidMostofi/swarmmanager/internal/resource"
	resource "github.com/VahidMostofi/swarmmanager/internal/resource/collector"
	"github.com/VahidMostofi/swarmmanager/internal/statutils"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/VahidMostofi/swarmmanager/internal/utils"
	"github.com/VahidMostofi/swarmmanager/internal/workload"
	"github.com/montanaflynn/stats"

	"gopkg.in/yaml.v2"
)

// TimingConfigs ...
type TimingConfigs struct {
	IterationDuration                 time.Duration
	WaitAfterServicesAreReadyDuration time.Duration
	WaitAfterLoadGeneratorStopped     time.Duration
}

// AutoConfigurer ...
type AutoConfigurer struct {
	Workload               string
	LoadGenerator          loadgenerator.LoadGenerator
	ResponseTimeCollector  workload.ResponseTimeCollector
	RequestCountCollector  workload.RequestCountCollector
	ResourceUsageCollector resource.Collector
	ConfigurerAgent        strategies.Configurer
	SwarmManager           *swarm.Manager
	Database               caching.Database
	TimingConfigs
}

// NewAutoConfigurer ...
func NewAutoConfigurer(lg loadgenerator.LoadGenerator, rtc workload.ResponseTimeCollector, rcc workload.RequestCountCollector, ruc resource.Collector, c strategies.Configurer, m *swarm.Manager, workload string, database caching.Database) *AutoConfigurer {
	a := &AutoConfigurer{
		LoadGenerator:          lg,
		ResponseTimeCollector:  rtc,
		RequestCountCollector:  rcc,
		ResourceUsageCollector: ruc,
		ConfigurerAgent:        c,
		SwarmManager:           m,
		TimingConfigs: TimingConfigs{
			IterationDuration:                 time.Duration(configs.GetConfig().Test.Duration),
			WaitAfterServicesAreReadyDuration: 7,
			WaitAfterLoadGeneratorStopped:     time.Duration(configs.GetConfig().Test.WaitAfterLoadGeneratorDone),
		},
		Workload: workload,
		Database: database,
	}
	return a
}

// Validate ...
func Validate(config map[string]swarm.ServiceSpecs) (float64, bool) {
	var sum float64
	for key := range config {
		sum += (float64(config[key].ReplicaCount) * config[key].CPULimits)
	}
	log.Println("there are", sum, "cores required!")
	return sum, sum <= configs.GetConfig().Host.AvailableCPUCount
}

// Start ...
func (a *AutoConfigurer) Start(name string, command string) {

	stackHistory := &history.ExecutionDetails{
		Name:     name,
		Workload: a.Workload,
		History:  make([]history.Information, 0),
		Command:  command,
		Config:   configs.GetConfig(),
	}

	// Remove the current Stack
	err := a.SwarmManager.RemoveStack(1)
	if err != nil {
		log.Panic(err)
	}

	// Deploy the stack with basic configuration
	// dockerComposePath := "/Users/vahid/workspace/bookstore/docker-compose.yml"
	dockerComposePath := configs.GetConfig().TestBed.DockerComposeFile

	initialConfig, err := a.ConfigurerAgent.GetInitialConfig()
	if err != nil {
		log.Panic(err)
	}

	err = a.SwarmManager.DeployStackWithDockerCompose(dockerComposePath, 1, initialConfig)
	if err != nil {
		log.Panic(err)
	}
	// for {
	// 	if a.SwarmManager.CurrentStackState == swarm.StackStateServicesAreReady {
	// 		break
	// 	}
	// 	time.Sleep(150 * time.Millisecond)
	// }
	time.Sleep(8 * time.Second)
	var start int64
	var end int64
	var iteration int
	for {
		iteration++
		log.Printf("ITERATION %d\n", iteration)
		for {
			if a.SwarmManager.CurrentStackState == swarm.StackStateServicesAreReady {
				log.Println("CurrentStackState is", swarm.GetStateString(a.SwarmManager.CurrentStackState), "so lets break out of the loop")
				a.SwarmManager.UpdateCurrentSpecs()
				break
			}
			time.Sleep(150 * time.Millisecond)
		}
		info, err := a.Database.Retrieve(string(a.Workload), a.SwarmManager.DesiredSpecs)
		if err == nil {
			log.Println("Autoconfigurer: information is found for this configuration/workload")
		} else {
			time.Sleep(a.WaitAfterServicesAreReadyDuration * time.Second)
			go a.LoadGenerator.Start()
			log.Println("load generator started")
			time.Sleep(30 * time.Second)
			a.ResourceUsageCollector = resource.GetTheResourceUsageCollector()
			err = a.ResourceUsageCollector.Start()
			if err != nil {
				log.Panic(err)
			}
			time.Sleep(1 * time.Second)
			start = time.Now().UnixNano() / 1e3
			time.Sleep(a.IterationDuration * time.Second)
			end = time.Now().UnixNano() / 1e3
			log.Println("finished the test")
			a.LoadGenerator.Stop()
			time.Sleep(a.WaitAfterLoadGeneratorStopped * time.Second)
			err = a.ResourceUsageCollector.Stop()
			time.Sleep(10 * time.Second)
			if err != nil {
				log.Panic(err)
			}
			servicesInfo := a.GatherInfo(int64(start), int64(end))
			lgFeedback, err := a.LoadGenerator.GetFeedback()
			if err != nil {
				log.Panic(fmt.Errorf("error while retrieving load generator feedback: %w", err))
			}
			info = history.Information{
				ServicesInfo:          servicesInfo,
				Specs:                 a.SwarmManager.ToHumanReadable(a.SwarmManager.CurrentSpecs),
				JaegerFile:            a.ResponseTimeCollector.(*jaeger.Aggregator).LastStoredFile,
				Workload:              a.Workload,
				RequestResponseTimes:  make(map[string]history.ResponseTimeStats),
				LoadGeneratorFeedback: lgFeedback,
			}
			for _, reqName := range a.RequestCountCollector.GetRequestNames() {
				responseTimes, err := a.ResponseTimeCollector.GetRequestResponseTimes(reqName)
				if err != nil {
					log.Panic(err)
				}
				info.RequestResponseTimes[reqName], err = createStats(responseTimes, []string{"count", "mean", "p90", "p95", "p99", "std", "c90p95"})
				if err != nil {
					log.Panic(err)
				}
			}
			hash, err := a.Database.Store(a.Workload, a.SwarmManager.DesiredSpecs, info)
			if err != nil {
				log.Panicf("Error while storing run information: %w", err)
			}
			info.HashCode = hash
		}
		stackHistory.History = append(stackHistory.History, info)
		newSpecs, isChanged, err := a.ConfigurerAgent.Configure(info, a.SwarmManager.CurrentSpecs, a.SwarmManager.ServicesToManage)
		if err != nil {
			log.Panic(err)
		}
		if _, ok := Validate(newSpecs); !ok {
			log.Println("new config is not valid, breaking out of loop")
			break
		}
		fmt.Println("validated the new specs")
		a.SwarmManager.DesiredSpecs = newSpecs
		a.SwarmManager.StackStateCh <- swarm.StackStateMustCompare
		if !isChanged {
			log.Println("is changed is false, breaking out of loop")
			break
		}
		time.Sleep(5 * time.Second)
		saveHistory(stackHistory)
		fmt.Println("partial results at:", configs.GetConfig().Results.Path+stackHistory.Name+".yml")
	}
	saveHistory(stackHistory)
	fmt.Println("final results at:", configs.GetConfig().Results.Path+stackHistory.Name+".yml")
}

func saveHistory(stackHistory *history.ExecutionDetails) {
	b, err := yaml.Marshal(stackHistory)
	if err != nil {
		log.Panic(err)
	}
	ioutil.WriteFile(configs.GetConfig().Results.Path+stackHistory.Name+".yml", b, os.FileMode(int(0777)))
}

func (a *AutoConfigurer) printRUMap(r map[string]*r2.Utilization) string {
	res := "ruMAP:\n"
	for key, value := range r {
		if len(key) == 25 {
			res += a.SwarmManager.DesiredSpecs[key].Name + " " + strconv.Itoa(len(value.CPUUtilizationsAtTime)) + "\n"
		} else {
			// res += a.ResourceUsageCollector.
		}
	}
	return res
}

// GatherInfo ...
func (a *AutoConfigurer) GatherInfo(start, end int64) map[string]history.ServiceInfo {
	ruMap := a.ResourceUsageCollector.GetResourceUtilization()
	a.RequestCountCollector.(*jaeger.Aggregator).GetTraces(start, end, configs.GetConfig().Jaeger.RootService)
	info := make(map[string]history.ServiceInfo)
	for key := range a.SwarmManager.CurrentSpecs {
		serviceName := a.SwarmManager.CurrentSpecs[key].Name
		if !utils.ContainsString(a.SwarmManager.ServicesToManage, serviceName) {
			continue
		}
		serviceInfo := history.ServiceInfo{
			Start:        start,
			End:          end,
			ReplicaCount: a.SwarmManager.CurrentSpecs[key].ReplicaCount,
		}

		// CPU usage
		cpuUsages := make([]float64, 0)
		serviceInfo.NumberOfCores = a.SwarmManager.CurrentSpecs[key].CPULimits
		for timestamp, usage := range ruMap[a.SwarmManager.StackName+"_"+serviceName].CPUUtilizationsAtTime {
			if timestamp >= start*1e3 && timestamp <= end*1e3 {
				cpuUsages = append(cpuUsages, usage/serviceInfo.NumberOfCores)
			}
		}
		// fmt.Println(serviceName, cpuUsages)
		v, err := stats.Mean(cpuUsages)
		if err != nil {
			log.Panic(err)
		}
		serviceInfo.CPUUsageMean = v

		v, err = stats.Percentile(cpuUsages, 70)
		if err != nil {
			log.Panic(err)
		}
		serviceInfo.CPUUsage70Percentile = v

		v, err = stats.Percentile(cpuUsages, 75)
		if err != nil {
			log.Panic(err)
		}
		serviceInfo.CPUUsage75Percentile = v

		v, err = stats.Percentile(cpuUsages, 80)
		if err != nil {
			log.Panic(err)
		}
		serviceInfo.CPUUsage80Percentile = v
		//--------------------------------------------------
		v, err = stats.Percentile(cpuUsages, 85)
		if err != nil {
			log.Panic(err)
		}
		serviceInfo.CPUUsage85Percentile = v
		//--------------------------------------------------
		v, err = stats.Percentile(cpuUsages, 90)
		if err != nil {
			log.Panic(err)
		}
		serviceInfo.CPUUsage90Percentile = v
		//--------------------------------------------------
		v, err = stats.Percentile(cpuUsages, 95)
		if err != nil {
			log.Panic(err)
		}
		serviceInfo.CPUUsage95Percentile = v
		//--------------------------------------------------
		v, err = stats.Percentile(cpuUsages, 99)
		if err != nil {
			log.Panic(err)
		}
		serviceInfo.CPUUsage99Percentile = v
		//--------------------------------------------------
		// Request Count
		c, e := a.RequestCountCollector.GetRequestCount(serviceName)
		if e != nil {
			log.Panic(e)
		}
		serviceInfo.RequestCount = c
		//--------------------------------------------------
		// TimeDetails
		requestToValueNameToValues, e := a.ResponseTimeCollector.GetServiceDetails(serviceName)
		if e != nil {
			log.Panic(e)
		}

		serviceInfo.TimesDetails = make(map[string]map[string]history.ResponseTimeStats)
		for req := range requestToValueNameToValues {
			serviceInfo.TimesDetails[req] = make(map[string]history.ResponseTimeStats)
			for valueName := range requestToValueNameToValues[req] {
				rts, err := createStats(requestToValueNameToValues[req][valueName], []string{"mean", "count"})
				if err != nil {
					log.Panic(err)
				}
				serviceInfo.TimesDetails[req][valueName] = rts
			}
		}
		info[serviceName] = serviceInfo
	}

	return info
}

func createStats(values []float64, names []string) (history.ResponseTimeStats, error) {
	if len(values) == 0 {
		values = append(values, 0)
	}
	rts := history.ResponseTimeStats{}

	if utils.ContainsString(names, "mean") {
		mean, err := stats.Mean(values)
		if err != nil {
			log.Panic(err)
		}
		rts.ResponseTimesMean = &mean
	}
	//--------------------------------------------------
	if utils.ContainsString(names, "p90") {
		p90, err := stats.Percentile(values, 90)
		if err != nil {
			log.Panic(err)
		}
		rts.ResponseTimes90Percentile = &p90
	}

	//--------------------------------------------------
	if utils.ContainsString(names, "p95") {
		p95, err := stats.Percentile(values, 95)
		if err != nil {
			log.Panic(err)
		}
		rts.ResponseTimes95Percentile = &p95
	}

	//--------------------------------------------------
	if utils.ContainsString(names, "p99") {
		p99, err := stats.Percentile(values, 99)
		if err != nil {
			log.Panic(err)
		}
		rts.ResponseTimes99Percentile = &p99
	}

	//--------------------------------------------------
	if utils.ContainsString(names, "std") {
		std, err := stats.StandardDeviation(values)
		if err != nil {
			log.Panic(err)
		}
		rts.ResponseTimesSTD = &std
	}

	//--------------------------------------------------
	if utils.ContainsString(names, "count") {
		c := len(values)
		rts.Count = &c
	}

	//--------------------------------------------------
	if utils.ContainsString(names, "c90p95") && len(values) > 10 {
		_, uti, err := statutils.ComputeToleranceIntervalNonParametric(values, 0.90, 0.95)
		if err != nil {
			log.Println("response times:", values)
			log.Panic(err)
		}
		rts.RTToleranceIntervalUBoundConfidence90p95 = &uti
	}
	return rts, nil
}
