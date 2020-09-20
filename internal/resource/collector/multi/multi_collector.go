package multi

import (
	"fmt"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/resource"
	"github.com/VahidMostofi/swarmmanager/internal/resource/collector/single"
)

// Collector collects stat usage from multiple hosts using multiple single collectors
type Collector struct {
	SingleCollectors []*single.Collector
	Hosts            []string
}

// Configure ...
func (mc *Collector) Configure(map[string]string) error {
	mc.SingleCollectors = make([]*single.Collector, 0)
	mc.Hosts = make([]string, 0)
	for _, s := range configs.GetConfig().UsageCollector.Details["hosts"].([]interface{}) {
		host := s.(string)
		mc.Hosts = append(mc.Hosts, host)

		sc := &single.Collector{}
		err := sc.Configure(map[string]string{"host": host})
		if err != nil {
			return fmt.Errorf("error while configuring single collector for host %s: %w", host, err)
		}
		mc.SingleCollectors = append(mc.SingleCollectors, sc)
	}
	return nil
}

// Start ...
func (mc *Collector) Start() error {
	for i, sc := range mc.SingleCollectors {
		err := sc.Start()
		if err != nil {
			return fmt.Errorf("error while starting SingleCollector(%s): %w", mc.Hosts[i], err)
		}
	}
	return nil
}

// Stop ...
func (mc *Collector) Stop() error {
	for i, sc := range mc.SingleCollectors {
		err := sc.Stop()
		if err != nil {
			return fmt.Errorf("error while stopping SingleCollector(%s): %w", mc.Hosts[i], err)
		}
	}
	return nil
}

// GetResourceUtilization ...
func (mc *Collector) GetResourceUtilization() map[string]*resource.Utilization {
	combined := make(map[string]*resource.Utilization)
	for _, sc := range mc.SingleCollectors {
		singleResourceUtilization := sc.GetResourceUtilization()
		for resourceID := range singleResourceUtilization {
			if _, ok := combined[resourceID]; !ok {
				ru := &resource.Utilization{
					ResourceID:            resourceID,
					CPUUtilizationsAtTime: make(map[int64]float64),
				}
				combined[resourceID] = ru
			}
			for timeStamp, value := range singleResourceUtilization[resourceID].CPUUtilizationsAtTime {
				combined[resourceID].CPUUtilizationsAtTime[timeStamp] = value
			}
		}
	}
	return combined
}

// GetContainerToService ContainerID to ServiceName
func (mc *Collector) GetContainerToService() map[string]string {
	combined := make(map[string]string)
	for _, sc := range mc.SingleCollectors {
		for container, service := range sc.GetContainerToService() {
			combined[container] = service
		}
	}
	return combined
}

// GetServiceToContainers ServiceName to ContainerIDs
func (mc *Collector) GetServiceToContainers() map[string][]string {
	combined := make(map[string][]string)
	for _, sc := range mc.SingleCollectors {
		for service, containers := range sc.GetServiceToContainers() {
			if existingList, ok := combined[service]; ok {
				combined[service] = append(existingList, containers...)
			} else {
				combined[service] = containers[:]
			}
		}
	}
	return combined
}
