package loadgenerator

import (
	"fmt"

	"github.com/VahidMostofi/swarmmanager/configs"
)

// LoadGenerator has the general behavior of any load generator
type LoadGenerator interface {
	Prepare() error
	Start() error
	Stop() error
	GetFeedback() (map[string]interface{}, error)
}

// Workload type
type Workload interface {
	GetThroughput() float64
	GetRequestProportion() map[string]float64
}

// GetLoadGenerator ...
func GetLoadGenerator() (LoadGenerator, error) {
	requestedType := configs.GetConfig().LoadGenerator.Type
	if requestedType == "k6" {
		return newK6LoadGenerator()
	}
	return nil, fmt.Errorf("the load generator of type %s is unknown", requestedType)
}
