package strategies

import (
	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
)

// Configurer ...
type Configurer interface {
	Configure(map[string]history.ServiceInfo, map[string]swarm.ServiceSpecs, []string) (map[string]swarm.ServiceSpecs, bool, error)
}
