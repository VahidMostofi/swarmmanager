package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/VahidMostofi/swarmmanager/internal/resource"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// SingleCollector gathers the stats form one single Docker host
type SingleCollector struct {
	Client                  *client.Client
	Ctx                     context.Context
	StatRecordingContext    context.Context
	StatRecordingCancelFunc context.CancelFunc
	CancelFunc              context.CancelFunc
	Containers              []types.Container
	Stackname               string
	ServiceToContainers     map[string][]string // map from service id to containers
	ContainerToService      map[string]string
	Services                map[string]string // map from service id to service name
	ResourceStats           map[string]*resource.Utilization
}

// ToString ...
func (sc *SingleCollector) ToString() string {
	res := "--------------------------------------------------------\n"
	res += "monitoring stats:\n"
	res += "containers:\n"
	for _, c := range sc.Containers {
		res += c.ID[:12] + "_" + c.Names[0] + "for service:" + sc.ContainerToService[c.ID] + "\n"
	}
	res += "services:\n"
	for service, containers := range sc.ServiceToContainers {
		res += service + ": "
		for _, c := range containers {
			res += c[:12] + ","
		}
		res += "\n"
	}
	res += "--------------------------------------------------------"
	return res
}

// Configure the collector, the values needs:
// "host": the host we are connecting to
// "stackname": only containers with com.docker.stack.namespace label equal to stackname would be considered
func (sc *SingleCollector) Configure(values map[string]string) error {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(values["host"], "", nil, defaultHeaders)
	if err != nil {
		return fmt.Errorf("error while creating Docker client in SingleCollector.Configure: %w", err)
	}
	sc.Client = cli

	ctx, cancelFunc := context.WithCancel(context.Background())
	sc.Ctx = ctx
	sc.CancelFunc = cancelFunc

	ctx, cancelFunc = context.WithCancel(context.Background())
	sc.StatRecordingContext = ctx
	sc.StatRecordingCancelFunc = cancelFunc

	sc.Stackname = values["stackname"]
	return nil
}

// Start the collector
func (sc *SingleCollector) Start() error {
	log.Println("SingleCollector:", "Starting SingleCollector...")
	// List the containers

	containers, err := sc.Client.ContainerList(sc.Ctx, types.ContainerListOptions{})
	if err != nil {
		return fmt.Errorf("error while listing containers in SingleCollector.Start: %w", err)
	}

	sc.Containers = make([]types.Container, 0)
	sc.ServiceToContainers = make(map[string][]string)
	sc.Services = make(map[string]string)
	sc.ResourceStats = make(map[string]*resource.Utilization)
	sc.ContainerToService = make(map[string]string)
	log.Printf("SingleCollector: Found %d containers on host\n", len(containers))
	for _, container := range containers {
		if container.Labels["com.docker.stack.namespace"] == sc.Stackname {
			sc.Containers = append(sc.Containers, container)
			// fmt.Println("monitoing stats for these containers:", sc.Containers)
			serviceName := container.Labels["com.docker.swarm.service.name"]
			sc.ServiceToContainers[serviceName] = append(sc.ServiceToContainers[serviceName], container.ID)
			sc.Services[serviceName] = serviceName

			sc.ResourceStats[container.ID] = resource.NewResourceUtilization(container.ID)
			sc.ContainerToService[container.ID] = serviceName
		}
	}

	for name := range sc.Services {
		sc.ResourceStats[name] = resource.NewResourceUtilization(name)
	}

	// fmt.Println(sc.ToString())

	errCh := make(chan error)
	statsCh := make(chan struct {
		string
		float64
		int64
	}, len(sc.Containers))
	for _, container := range sc.Containers {
		go sc.getContainerStats(container.ID, errCh, statsCh)
	}

	go sc.recordStats(statsCh)

	return nil
}

// Stop collecting stats, aggregate and clean the results
func (sc *SingleCollector) Stop() error {
	sc.StatRecordingCancelFunc()
	// fmt.Println("signal to stop!")
	return nil
}

// GetResourceUtilization ... returns stats
func (sc *SingleCollector) GetResourceUtilization() map[string]*resource.Utilization {
	// for key, value := range sc.ResourceStats {
	// 	fmt.Println(key, value.CPUUtilizationsAtTime)
	// }
	return sc.ResourceStats
}

// GetContainerToService ....
func (sc *SingleCollector) GetContainerToService() map[string]string {
	return sc.ContainerToService
}

// GetServiceToContainers ...
func (sc *SingleCollector) GetServiceToContainers() map[string][]string {
	return sc.ServiceToContainers
}

func (sc *SingleCollector) recordStats(statsCh chan struct {
	string
	float64
	int64
}) {
	for {
		select {
		case <-sc.StatRecordingContext.Done():
			return
		case pair := <-statsCh:
			sc.ResourceStats[pair.string].AddCPUUsage(pair.float64, pair.int64)
			sc.ResourceStats[sc.ContainerToService[pair.string]].AddCPUUsage(pair.float64, pair.int64)
		}
	}
}

func (sc *SingleCollector) getContainerStats(containerID string, errorCh chan error, statsChannel chan struct {
	string
	float64
	int64
}) {
	stats, err := sc.Client.ContainerStats(sc.StatRecordingContext, containerID, true)
	if err != nil {
		errorCh <- fmt.Errorf("error while getting container stats for %s: %w", containerID, err)
		return
	}

	decoder := json.NewDecoder(stats.Body)

	var v types.StatsJSON

	for {
		select {
		case <-sc.StatRecordingContext.Done():
			// fmt.Println("got a signal for stop, stopping recording for " + containerID)
			stats.Body.Close()
			return
		default:
			if err := decoder.Decode(&v); err == io.EOF {
				return
			} else if err != nil {
				sc.StatRecordingCancelFunc()
			}
			previousCPU := v.PreCPUStats.CPUUsage.TotalUsage
			previousSystem := v.PreCPUStats.SystemUsage
			cpuPercent := calculateCPUPercent(previousCPU, previousSystem, &v)
			statsChannel <- struct {
				string
				float64
				int64
			}{containerID, cpuPercent, time.Now().UnixNano()}
		}
	}
}

func calculateCPUPercent(previousCPU, previousSystem uint64, v *types.StatsJSON) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.Stats.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.Stats.CPUStats.SystemUsage) - float64(previousSystem)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}

/*
services, err := sc.Client.ServiceList(sc.Ctx, types.ServiceListOptions{})
	if err != nil {
		return fmt.Errorf("error while listing services in SingleCollector.Start: %w", err)
	}

	for _, service := range services {
		detailedService, _, err := sc.Client.ServiceInspectWithRaw(sc.Ctx, service.ID)
		if err != nil {
			return fmt.Errorf("error while inspecting service %s in SingleController.Start: %w", detailedService.ID, err)
		}
	}
*/
