package collector

import (
	"time"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/k8s"
)

// K8sResourceCollector ...
type K8sResourceCollector struct {
	cpus    map[string][]float64
	closeCh chan bool
}

// Start ...
func (k *K8sResourceCollector) Start() {
	ticker := time.NewTicker(time.Second)

	k.cpus = make(map[string][]float64)
	for _, name := range configs.GetConfig().TestBed.ServicesToConfigure {
		k.cpus[name] = make([]float64, 0)
	}
	done := make(chan bool)
	k.closeCh = done
	go func() {
		connector := k8s.GetNewConnector("ssh", configs.GetConfig().Host.Host)
		for {
			select {
			case <-ticker.C:
				v, err := connector.GetCPUUsage()
				if err != nil {
					panic(err)
				}
				for key, value := range v {
					k.cpus[key] = append(k.cpus[key], value)
				}
			case <-done:
				return
			}
		}
	}()
}

// Stop ...
func (k *K8sResourceCollector) Stop() {
	k.closeCh <- true
}

// GetCPUValues ...
func (k *K8sResourceCollector) GetCPUValues() map[string][]float64 {
	return k.cpus
}
