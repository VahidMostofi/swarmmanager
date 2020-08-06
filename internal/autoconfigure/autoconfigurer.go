package autoconfigure

//TODO I NEED TO REstart every single container
//TODO check validity of a configuration
import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"log"

	"github.com/VahidMostofi/swarmmanager"
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

// GetTheResourceUsageCollector ...
func GetTheResourceUsageCollector() resource.Collector {
	//TODO SERVICE COUNT IS HARDCODED!!!!!!!!
	stackName := swarmmanager.GetConfig().StackName
	c := resource.GetNewCollector("SingleCollector")
	err := c.Configure(map[string]string{"host": swarmmanager.GetConfig().Host, "stackname": stackName})
	if err != nil {
		log.Panic(err)
	}

	return c
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
			IterationDuration:                 time.Duration(swarmmanager.GetConfig().TestDuration),
			WaitAfterServicesAreReadyDuration: 7,
			WaitAfterLoadGeneratorStopped:     time.Duration(swarmmanager.GetConfig().WaitAfterLoadGenerator),
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
	return sum, sum <= swarmmanager.GetConfig().AvailabeCPUCount
}

// Start ...
func (a *AutoConfigurer) Start(name string, command string) {

	stackHistory := &history.ExecutionDetails{
		Name:     name,
		Workload: a.Workload,
		History:  make([]history.Information, 0),
		Command:  command,
		Config:   swarmmanager.GetConfig(),
	}

	// Remove the current Stack
	err := a.SwarmManager.RemoveStack(1)
	if err != nil {
		log.Panic(err)
	}

	// Deploy the stack with basic configuration
	// dockerComposePath := "/Users/vahid/workspace/bookstore/docker-compose.yml"
	dockerComposePath := swarmmanager.GetConfig().DockerComposeFile

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
			// a.SwarmManager.CurrentSpecs = make(map[string]swarm.ServiceSpecs)

			// services, err := a.SwarmManager.Client.ServiceList(a.SwarmManager.Ctx, types.ServiceListOptions{})
			// if err != nil {
			// 	panic(err)
			// }
			// for _, srv := range services {
			// 	serviceName := strings.Split(srv.Spec.Name, "_")[1]
			// 	temp := info.Specs[serviceName]
			// 	temp.Name = serviceName
			// 	temp.ID = srv.ID
			// 	a.SwarmManager.CurrentSpecs[serviceName] = temp
			// }
		} else {
			time.Sleep(a.WaitAfterServicesAreReadyDuration * time.Second)
			go a.LoadGenerator.Start(make(map[string]string))
			log.Println("load generator started")
			time.Sleep(30 * time.Second)
			a.ResourceUsageCollector = GetTheResourceUsageCollector()
			err = a.ResourceUsageCollector.Start()
			if err != nil {
				log.Panic(err)
			}
			time.Sleep(1 * time.Second)
			start = time.Now().UnixNano() / 1e3
			time.Sleep(a.IterationDuration * time.Second)
			end = time.Now().UnixNano() / 1e3
			log.Println("finished the test")
			a.LoadGenerator.Stop(make(map[string]string))
			time.Sleep(a.WaitAfterLoadGeneratorStopped * time.Second)
			err = a.ResourceUsageCollector.Stop()
			if err != nil {
				log.Panic(err)
			}
			servicesInfo := a.GatherInfo(int64(start), int64(end))
			info = history.Information{
				ServicesInfo: servicesInfo,
				Specs:        a.SwarmManager.ToHumanReadable(a.SwarmManager.CurrentSpecs),
				JaegerFile:   a.ResponseTimeCollector.(*jaeger.Aggregator).LastStoredFile,
				Workload:     a.Workload,
			}
			hash, err := a.Database.Store(a.Workload, a.SwarmManager.DesiredSpecs, info)
			if err != nil {
				log.Panicf("Error while storing run information: %w", err)
			}
			info.HashCode = hash
		}
		stackHistory.History = append(stackHistory.History, info)
		newSpecs, isChanged, err := a.ConfigurerAgent.Configure(info.ServicesInfo, a.SwarmManager.CurrentSpecs, a.SwarmManager.ServicesToManage)
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
		fmt.Println("partial results at:", swarmmanager.GetConfig().ResultsDirectoryPath+stackHistory.Name+".yml")
	}
	saveHistory(stackHistory)
	fmt.Println("final results at:", swarmmanager.GetConfig().ResultsDirectoryPath+stackHistory.Name+".yml")
}

func saveHistory(stackHistory *history.ExecutionDetails) {
	b, err := yaml.Marshal(stackHistory)
	if err != nil {
		log.Panic(err)
	}
	ioutil.WriteFile(swarmmanager.GetConfig().ResultsDirectoryPath+stackHistory.Name+".yml", b, os.FileMode(int(0777)))
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
	a.RequestCountCollector.(*jaeger.Aggregator).GetTraces(start, end, swarmmanager.GetConfig().JaegerRootService) //TODO this is hardcoded!
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
		serviceInfo.RequestCount = c[serviceName+"_total"]
		//--------------------------------------------------
		// Response Times
		valueNameToResponseTimes, e := a.ResponseTimeCollector.GetResponseTimes(serviceName)
		if e != nil {
			log.Panic(e)
		}

		serviceInfo.ResponseTimes = make(map[string]history.ResponseTimeStats)
		for valueName, responseTimes := range valueNameToResponseTimes {

			mean, std, p90, p95, p99 := getDifferentResponseTimes(responseTimes)

			responseTimesStats := history.ResponseTimeStats{
				ResponseTimesMean:         &mean,
				ResponseTimesSTD:          &std,
				ResponseTimes90Percentile: &p90,
				ResponseTimes95Percentile: &p95,
				ResponseTimes99Percentile: &p99,
			}

			//TODO the 90 and 95 values and the fact that which values should be computed should come form config file
			_, uti, err := statutils.ComputeToleranceIntervalNonParametric(responseTimes, 0.90, 0.95)
			if err != nil {
				log.Println("response times:", responseTimes)
				log.Panic(err)
			}
			responseTimesStats.RTToleranceIntervalUBoundConfidence90p95 = &uti
			serviceInfo.ResponseTimes[strings.Split(valueName, "_")[1]] = responseTimesStats

		}

		info[serviceName] = serviceInfo
	}
	return info
}

func getDifferentResponseTimes(responseTimes []float64) (float64, float64, float64, float64, float64) {
	if len(responseTimes) == 0 {
		responseTimes = append(responseTimes, 0)
	}
	mean, err := stats.Mean(responseTimes)
	if err != nil {
		log.Panic(err)
	}
	//--------------------------------------------------
	p90, err := stats.Percentile(responseTimes, 90)
	if err != nil {
		log.Panic(err)
	}

	//--------------------------------------------------
	p95, err := stats.Percentile(responseTimes, 95)
	if err != nil {
		log.Panic(err)
	}

	//--------------------------------------------------
	p99, err := stats.Percentile(responseTimes, 99)
	if err != nil {
		log.Panic(err)
	}

	//--------------------------------------------------
	std, err := stats.StandardDeviation(responseTimes)
	if err != nil {
		log.Panic(err)
	}

	return mean, std, p90, p95, p99
}
