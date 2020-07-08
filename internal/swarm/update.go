package swarm

import (
	"log"

	"github.com/docker/docker/api/types"
	dockerswarm "github.com/docker/docker/api/types/swarm"
)

// UpdateServices based on desired specs
func (m *Manager) UpdateServices() {

	m.StackStateCh <- StackStateUpdatingSpecs
	for serviceID := range m.DesiredSpecs {
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
