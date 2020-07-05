package autoconfigure

//TODO I NEED TO REstart every single container
//TODO check validity of a configuration
import (
	"fmt"
	"strconv"
	"time"

	"github.com/VahidMostofi/swarmmanager/internal/jaeger"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	r2 "github.com/VahidMostofi/swarmmanager/internal/resource"
	resource "github.com/VahidMostofi/swarmmanager/internal/resource/collector"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/VahidMostofi/swarmmanager/internal/workload"
	"github.com/montanaflynn/stats"
)

// AutoConfigurer ...
type AutoConfigurer struct {
	LoadGenerator          loadgenerator.LoadGenerator
	ResponseTimeCollector  workload.ResponseTimeCollector
	RequestCountCollector  workload.RequestCountCollector
	ResourceUsageCollector resource.Collector
	ConfigurerAgent        Configurer
	SwarmManager           *swarm.Manager
}

// GetTheResourceUsageCollector ...
func GetTheResourceUsageCollector() resource.Collector {
	//TODO SERVICE COUNT IS HARDCODED!!!!!!!!
	stackName := "bookstore"
	c := resource.GetNewCollector("SingleCollector")
	err := c.Configure(map[string]string{"host": "tcp://136.159.209.204:2375", "stackname": stackName})
	if err != nil {
		panic(err)
	}

	return c
}

// NewAutoConfigurer ...
func NewAutoConfigurer(lg loadgenerator.LoadGenerator, rtc workload.ResponseTimeCollector, rcc workload.RequestCountCollector, ruc resource.Collector, c Configurer, m *swarm.Manager) *AutoConfigurer {
	a := &AutoConfigurer{
		LoadGenerator:          lg,
		ResponseTimeCollector:  rtc,
		RequestCountCollector:  rcc,
		ResourceUsageCollector: ruc,
		ConfigurerAgent:        c,
		SwarmManager:           m,
	}
	return a
}

// Start ...
func (a *AutoConfigurer) Start() {

	// Remove the current Stack
	err := a.SwarmManager.RemoveStack()
	if err != nil {
		panic(err)
	}

	// Deploy the stack with basic configuration
	dockerComposePath := "/Users/vahid/workspace/bookstore/docker-compose.yml" //TODO this is hard coded!
	err = a.SwarmManager.DeployStackWithDockerCompose(dockerComposePath, 1)
	if err != nil {
		panic(err)
	}
	for {
		if a.SwarmManager.CurrentStackState == swarm.StackStateServicesAreReady {
			break
		}
		time.Sleep(150 * time.Millisecond)
	}
	time.Sleep(15 * time.Second)
	for {
		for {
			// fmt.Println(a.SwarmManager.CurrentStackState)
			if a.SwarmManager.CurrentStackState == swarm.StackStateServicesAreReady {
				fmt.Println(time.Now().UnixNano(), "a.SwarmManager.CurrentStackState is", a.SwarmManager.CurrentStackState, "so lets break out of the loop")
				break
			}
			time.Sleep(150 * time.Millisecond)
		}
		fmt.Println(time.Now().UnixNano(), a.SwarmManager.CurrentStackState, "write after the loop")
		time.Sleep(15 * time.Second)
		fmt.Println(time.Now().UnixNano(), a.SwarmManager.CurrentStackState, "15 seconds after the after the loop")
		go a.LoadGenerator.Start(make(map[string]string))
		fmt.Println("load generator started")
		time.Sleep(15 * time.Second)
		a.ResourceUsageCollector = GetTheResourceUsageCollector()
		err := a.ResourceUsageCollector.Start()
		if err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)
		fmt.Println(time.Now().UnixNano(), "NEW ITERATION")
		start := time.Now().UnixNano() / 1e3
		time.Sleep(15 * time.Second)
		end := time.Now().UnixNano() / 1e3
		a.LoadGenerator.Stop(make(map[string]string))
		time.Sleep(2 * time.Second)
		err = a.ResourceUsageCollector.Stop()
		if err != nil {
			panic(err)
		}
		info := a.GatherInfo(int64(start), int64(end))
		newSpecs, isChanged, err := a.ConfigurerAgent.Configure(info, a.SwarmManager.CurrentSpecs)
		a.SwarmManager.StackStateCh <- swarm.StackStateMustCompare
		if err != nil {
			panic(err)
		}
		a.SwarmManager.DesiredSpecs = newSpecs
		if !isChanged {
			fmt.Println("is changed is false, breaking out of loop")
		}
		// if a.SwarmManager.CompareSpecs() {
		// fmt.Println("specs are the same break out of the loop")
		// break
		// }
		time.Sleep(5 * time.Second)
	}
	fmt.Println(a.SwarmManager.DesiredSpecs)
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
func (a *AutoConfigurer) GatherInfo(start, end int64) map[string]ServiceInfo {
	ruMap := a.ResourceUsageCollector.GetResourceUtilization()
	a.RequestCountCollector.(*jaeger.JaegerAggregator).GetTraces(start, end, "gateway") //TODO this is hardcoded!
	info := make(map[string]ServiceInfo)
	for serviceID := range a.SwarmManager.CurrentSpecs {
		serviceName := a.SwarmManager.CurrentSpecs[serviceID].Name
		if !(serviceName == "books" || serviceName == "auth" || serviceName == "gateway") {
			continue
		} //TODO
		fmt.Println("gathering info about", serviceName, serviceID)
		serviceInfo := ServiceInfo{
			Start:        start,
			End:          end,
			ReplicaCount: a.SwarmManager.CurrentSpecs[serviceID].ReplicaCount,
		}

		// CPU usage
		cpuUsages := make([]float64, 0)
		serviceInfo.NumberOfCores = a.SwarmManager.CurrentSpecs[serviceID].CPULimits
		fmt.Println(a.printRUMap(ruMap))
		for timestamp, usage := range ruMap[serviceID].CPUUtilizationsAtTime {
			if timestamp >= start*1e3 && timestamp <= end*1e3 {
				cpuUsages = append(cpuUsages, usage/serviceInfo.NumberOfCores)
			}
		}
		v, err := stats.Mean(cpuUsages)
		if err != nil {
			panic(err)
		}
		serviceInfo.CPUUsageMean = v

		v, err = stats.Percentile(cpuUsages, 90)
		if err != nil {
			panic(err)
		}
		serviceInfo.CPUUsage95Percentile = v

		v, err = stats.Percentile(cpuUsages, 95)
		if err != nil {
			panic(err)
		}
		serviceInfo.CPUUsage99Percentile = v

		v, err = stats.Percentile(cpuUsages, 99)
		if err != nil {
			panic(err)
		}
		serviceInfo.CPUUsage90Percentile = v

		// Request Count
		if serviceName == "books" {
			c, e := a.RequestCountCollector.GetRequestCount("books_edit_book")
			if e != nil {
				panic(e)
			}
			serviceInfo.RequestCount = c
			c, e = a.RequestCountCollector.GetRequestCount("books_get_book")
			if e != nil {
				panic(e)
			}
			serviceInfo.RequestCount += c
		} else if serviceName == "auth" {
			c, e := a.RequestCountCollector.GetRequestCount("auth_req_login")
			if e != nil {
				panic(e)
			}
			serviceInfo.RequestCount = c
		} else if serviceName == "gateway" {
			c, e := a.RequestCountCollector.GetRequestCount("auth_req_login")
			if e != nil {
				panic(e)
			}
			serviceInfo.RequestCount = c
			c, e = a.RequestCountCollector.GetRequestCount("books_edit_book")
			if e != nil {
				panic(e)
			}
			serviceInfo.RequestCount += c
			c, e = a.RequestCountCollector.GetRequestCount("books_get_book")
			if e != nil {
				panic(e)
			}
			serviceInfo.RequestCount += c
		}

		// Response Times
		responseTimes := make([]float64, 0)
		if serviceName == "books" {
			rts, e := a.ResponseTimeCollector.GetResponseTimes("books_edit_book")
			if e != nil {
				panic(e)
			}
			responseTimes = append(responseTimes, rts...)
			rts, e = a.ResponseTimeCollector.GetResponseTimes("books_get_book")
			if e != nil {
				panic(e)
			}
			responseTimes = append(responseTimes, rts...)
		} else if serviceName == "auth" {
			rts, e := a.ResponseTimeCollector.GetResponseTimes("auth_req_login")
			if e != nil {
				panic(e)
			}
			responseTimes = append(responseTimes, rts...)
		} else if serviceName == "gateway" {
			rts, e := a.ResponseTimeCollector.GetResponseTimes("auth_req_login")
			if e != nil {
				panic(e)
			}
			responseTimes = append(responseTimes, rts...)
			rts, e = a.ResponseTimeCollector.GetResponseTimes("books_edit_book")
			if e != nil {
				panic(e)
			}
			responseTimes = append(responseTimes, rts...)
			rts, e = a.ResponseTimeCollector.GetResponseTimes("books_get_book")
			if e != nil {
				panic(e)
			}
			responseTimes = append(responseTimes, rts...)
		}
		m, err := stats.Mean(responseTimes)
		if err != nil {
			panic(err)
		}
		serviceInfo.ResponseTimesMean = m

		m, err = stats.Percentile(responseTimes, 90)
		if err != nil {
			panic(err)
		}
		serviceInfo.ResponseTimes90Percentile = m

		m, err = stats.Percentile(responseTimes, 95)
		if err != nil {
			panic(err)
		}
		serviceInfo.ResponseTimes95Percentile = m

		m, err = stats.Percentile(responseTimes, 99)
		if err != nil {
			panic(err)
		}
		serviceInfo.ResponseTimes99Percentile = m

		info[serviceName] = serviceInfo
	}
	return info
}
