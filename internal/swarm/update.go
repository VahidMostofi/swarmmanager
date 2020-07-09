package swarm

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/docker/docker/api/types"
	dockerswarm "github.com/docker/docker/api/types/swarm"
)

var ForceAllUpdate = true

// UpdateServices based on desired specs
func (m *Manager) UpdateServices() {

	m.StackStateCh <- StackStateUpdatingSpecs
	for serviceID := range m.DesiredSpecs {

		areEqual, changes := m.comapeServiceSpecs(serviceID)
		if !ForceAllUpdate {
			if areEqual {
				log.Println(m.DesiredSpecs[serviceID].Name, "is not changed, no update is required (ForceAllUpdate is false)")
				continue
			}
			if len(changes) == 1 && changes[0] == "ReplicaCount" {
				//TODO add all of these to each other and run only one command
				log.Println(m.DesiredSpecs[serviceID].Name, "only replica count is changed, use scale (ForceAllUpdate is false)")
				err := m.ScaleOnlyUpdate(m.DesiredSpecs[serviceID].Name, serviceID, m.DesiredSpecs[serviceID].ReplicaCount)
				if err != nil {
					log.Panic(err)
				}
				continue
			}
		}

		serviceReplicaCount := uint64(m.DesiredSpecs[serviceID].ReplicaCount)
		onlineService, _, err := m.Client.ServiceInspectWithRaw(m.Ctx, serviceID)
		if err != nil {
			log.Panic(err)
		}
		newSpec := onlineService.Spec
		newSpec.TaskTemplate.ContainerSpec.Env = m.DesiredSpecs[serviceID].EnvironmentVariables
		newSpec.TaskTemplate.Resources.Limits.NanoCPUs = int64(m.DesiredSpecs[serviceID].CPULimits * 1e9)
		newSpec.TaskTemplate.Resources.Limits.MemoryBytes = m.DesiredSpecs[serviceID].MemoryLimits
		newSpec.TaskTemplate.Resources.Reservations.NanoCPUs = int64(m.DesiredSpecs[serviceID].CPUReservation * 1e9)
		newSpec.TaskTemplate.Resources.Reservations.MemoryBytes = m.DesiredSpecs[serviceID].MemoryReservations
		newSpec.Mode.Replicated.Replicas = &serviceReplicaCount
		newSpec.TaskTemplate.ForceUpdate++
		log.Println("forcing update on", m.DesiredSpecs[serviceID].Name)

		log.Println("updating service...", m.DesiredSpecs[serviceID].Name)
		serviceUpdateResponse, err := m.Client.ServiceUpdate(m.Ctx, serviceID, dockerswarm.Version{onlineService.Version.Index}, newSpec, types.ServiceUpdateOptions{})
		log.Println("update done", m.DesiredSpecs[serviceID].Name)
		if err != nil {
			log.Panic(err)
		}
		if len(serviceUpdateResponse.Warnings) > 0 {
			log.Println("updating", m.DesiredSpecs[serviceID].Name, "warnings", serviceUpdateResponse.Warnings)
		}
	}
	m.StackStateCh <- StackStateWaitForServicesToBeDeployed
}

// ScaleOnlyUpdate ...
func (m *Manager) ScaleOnlyUpdate(serviceName, serviceID string, count int) error {
	log.Println("scaling", serviceName, "to", count)
	cmd := exec.Command("docker", "-H", m.Host, "service", "scale", serviceID+"="+strconv.Itoa(count))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("scaling %s:%s failed with error: %w; %s", serviceName, serviceID[:12], err, string(out))
	}
	return nil
}
