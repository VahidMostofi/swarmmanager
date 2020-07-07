package autoconfigure

import "github.com/VahidMostofi/swarmmanager/internal/swarm"

// Configurer ...
type Configurer interface {
	Configure(map[string]ServiceInfo, map[string]swarm.ServiceSpecs, []string) (map[string]swarm.ServiceSpecs, bool, error)
}
