package strategies

import (
	"math"

	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
)

// SingleRun ...
type SingleRun struct {
	Config map[string]swarm.SimpleSpecs
}

// OnFeedbackCallback ...
func (c *SingleRun) OnFeedbackCallback(map[string]history.ServiceInfo) error {
	return nil
}

// Configure ...
func (c *SingleRun) Configure(info history.Information, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	return nil, false, nil
}

// GetInitialConfig ...
func (c *SingleRun) GetInitialConfig(loadgenerator.Workload) (map[string]swarm.SimpleSpecs, error) {
	return c.getReconfiguredConfiguration(c.Config), nil
}

func (c *SingleRun) getReconfiguredConfiguration(inputConfig map[string]swarm.SimpleSpecs) map[string]swarm.SimpleSpecs {
	service2totalResource := make(map[string]float64)
	for serviceName, c := range inputConfig {
		service2totalResource[serviceName] = c.CPU * float64(c.Replica)
	}
	reconfiguredSpecs := make(map[string]swarm.SimpleSpecs)

	for service, totalCPU := range service2totalResource {
		reconfiguredSpecs[service] = swarm.SimpleSpecs{
			CPU:     totalCPU,
			Replica: 1,
			Worker:  int(math.Ceil(totalCPU)),
		}
	}

	return reconfiguredSpecs
}
