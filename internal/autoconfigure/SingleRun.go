package autoconfigure

import (
	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
)

// SingleRun ...
type SingleRun struct {
}

// Configure ...
func (c *SingleRun) Configure(values map[string]history.ServiceInfo, currentState map[string]swarm.ServiceSpecs, servicesToMonitor []string) (map[string]swarm.ServiceSpecs, bool, error) {
	return nil, false, nil
}
