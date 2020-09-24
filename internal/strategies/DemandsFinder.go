package strategies

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"gopkg.in/yaml.v3"
)

// DemandsFinder ...
type DemandsFinder struct {
	ResultPath string
}

// GetInitialConfig ...
func (c *DemandsFinder) GetInitialConfig(workload loadgenerator.Workload) (map[string]swarm.SimpleSpecs, error) {
	return make(map[string]swarm.SimpleSpecs), nil
}

// Configure ....
func (c *DemandsFinder) Configure(info history.Information, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	demands := make(map[string]map[string]float64)

	for serviceName, serviceInfo := range info.ServicesInfo {
		demands[serviceName] = make(map[string]float64)
		for request, timesDetails := range serviceInfo.TimesDetails {
			if _, contains := timesDetails["service_time"]; contains {
				demands[serviceName][request] = *timesDetails["service_time"].ResponseTimesMean
			} else {
				demands[serviceName][request] = 0
			}
		}
	}

	b, err := yaml.Marshal(demands)
	if err != nil {
		return nil, false, fmt.Errorf("error while marshaling demands to yaml: %w", err)
	}

	err = ioutil.WriteFile(c.ResultPath, b, 0777)
	if err != nil {
		return nil, false, fmt.Errorf("error while saving demands to %s: %w", c.ResultPath, err)
	}

	log.Println("demands wrote to:", c.ResultPath)

	return currentState, false, nil
}

// OnFeedbackCallback ...
func (c *DemandsFinder) OnFeedbackCallback(map[string]history.ServiceInfo) error { return nil }
