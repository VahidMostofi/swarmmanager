package strategies

import (
	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
)

// SingleRun ...
type SingleRun struct {
}

// OnFeedbackCallback ...
func (c *SingleRun) OnFeedbackCallback(map[string]history.ServiceInfo) error {
	return nil
}

// Configure ...
func (c *SingleRun) Configure(values map[string]history.ServiceInfo, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	return nil, false, nil
}

// GetInitialConfig ...
func (c *SingleRun) GetInitialConfig() (map[string]swarm.SimpleSpecs, error) {
	config := make(map[string]swarm.SimpleSpecs)
	config["auth"] = swarm.SimpleSpecs{
		CPU:     2,
		Replica: 1,
		Worker:  2,
	}
	config["books"] = swarm.SimpleSpecs{
		CPU:     1,
		Replica: 1,
		Worker:  1,
	}
	config["gateway"] = swarm.SimpleSpecs{
		CPU:     2,
		Replica: 1,
		Worker:  2,
	}
	return config, nil
}
