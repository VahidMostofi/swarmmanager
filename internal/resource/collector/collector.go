package collector

import (
	"log"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/resource"
	"github.com/VahidMostofi/swarmmanager/internal/resource/collector/multi"
	"github.com/VahidMostofi/swarmmanager/internal/resource/collector/single"
)

// Collector is the interface of any collection of tools and techniques which collect and aggregate resource utilization
type Collector interface {
	Configure(map[string]string) error
	Start() error
	Stop() error
	GetResourceUtilization() map[string]*resource.Utilization
	GetContainerToService() map[string]string    //ContainerID to ServiceName
	GetServiceToContainers() map[string][]string //ServiceName to ContainerIDs
}

// GetNewCollector is the factory method for constructing a new Collector
func GetNewCollector(kind string) Collector {
	if kind == "SingleCollector" {
		return &single.Collector{}
	} else if kind == "MultiHostCollector" {
		return &multi.Collector{}
	}
	return nil
}

// GetTheResourceUsageCollector ...
func GetTheResourceUsageCollector() Collector {
	c := GetNewCollector(configs.GetConfig().UsageCollector.Type)
	err := c.Configure(map[string]string{})
	if err != nil {
		log.Panic(err)
	}

	return c
}
