package strategies

import (
	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
)

// Configurer ...
type Configurer interface {
	Configure(history.Information, map[string]swarm.ServiceSpecs, []string) (map[string]swarm.ServiceSpecs, bool, error)
	GetInitialConfig() (map[string]swarm.SimpleSpecs, error)
	OnFeedbackCallback(map[string]history.ServiceInfo) error
}
