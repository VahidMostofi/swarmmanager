package swarm

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"
)

func (m *Manager) liveUpdate() {
	m.StackStateCh <- StackStateUpdatingSpecs

	for key := range m.DesiredSpecs {
		serviceName := m.DesiredSpecs[key].Name
		serviceID := m.DesiredSpecs[key].ID
		log.Println("live updating", serviceName, serviceID)
		areEqual, _ := m.comapeServiceSpecs(serviceName)
		if areEqual {
			log.Println(m.DesiredSpecs[key].Name, "is not changed, no update is required (ForceAllUpdate is false)")
			continue
		}

		serviceReplicaCount := uint64(m.DesiredSpecs[key].ReplicaCount)

		log.Println("updating service...", m.DesiredSpecs[key].Name)
		err := m.liveUpdateRunCommand(m.DesiredSpecs[key].CPULimits, serviceReplicaCount, m.DesiredSpecs[key].Name)

		log.Println("update done", m.DesiredSpecs[key].Name)
		if err != nil {
			log.Panic(err)
		}
	}
	time.Sleep(1 * time.Second)
	m.StackStateCh <- StackStateServicesAreReady
}

func (m *Manager) liveUpdateRunCommand(cpuCount float64, replicaCount uint64, serviceName string) error {
	cpuCountStr := strconv.FormatFloat(cpuCount, 'f', 3, 64)
	replicaCountStr := strconv.FormatUint(replicaCount, 10)
	cmd := exec.Command("docker", "-H", m.Host, "service", "update", "--update-order", "start-first", "--limit-cpu", cpuCountStr, "--reserve-cpu", cpuCountStr, "--replicas="+replicaCountStr, "--update-parallelism", "1", "--update-delay", "15s", m.StackName+"_"+serviceName)
	b, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(b))
		return err
	}
	return nil
}
