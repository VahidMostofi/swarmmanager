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
func (m *Manager) UpdateServices(liveUpdate bool) {

	if liveUpdate {
		m.liveUpdate()
		return
	}

	m.StackStateCh <- StackStateUpdatingSpecs
	for key := range m.DesiredSpecs {
		serviceName := m.DesiredSpecs[key].Name
		serviceID := m.DesiredSpecs[key].ID
		log.Println("updating", serviceName, serviceID)
		areEqual, changes := m.comapeServiceSpecs(serviceName)
		if !ForceAllUpdate {
			if areEqual {
				log.Println(m.DesiredSpecs[key].Name, "is not changed, no update is required (ForceAllUpdate is false)")
				continue
			}
			if len(changes) == 1 && changes[0] == "ReplicaCount" {
				//TODO add all of these to each other and run only one command
				log.Println(m.DesiredSpecs[key].Name, "only replica count is changed, use scale (ForceAllUpdate is false)")
				err := m.ScaleOnlyUpdate(m.DesiredSpecs[key].Name, serviceID, m.DesiredSpecs[key].ReplicaCount)
				if err != nil {
					log.Panic(err)
				}
				continue
			}
		}
		serviceReplicaCount := uint64(m.DesiredSpecs[key].ReplicaCount)
		onlineService, _, err := m.Client.ServiceInspectWithRaw(m.Ctx, serviceID)
		if err != nil {
			log.Panic(err)
		}
		newSpec := onlineService.Spec
		newSpec.TaskTemplate.ContainerSpec.Env = m.DesiredSpecs[key].EnvironmentVariables
		newSpec.TaskTemplate.Resources.Limits.NanoCPUs = int64(m.DesiredSpecs[key].CPULimits * 1e9)
		newSpec.TaskTemplate.Resources.Limits.MemoryBytes = m.DesiredSpecs[key].MemoryLimits
		newSpec.TaskTemplate.Resources.Reservations.NanoCPUs = int64(m.DesiredSpecs[key].CPUReservation * 1e9)
		newSpec.TaskTemplate.Resources.Reservations.MemoryBytes = m.DesiredSpecs[key].MemoryReservations
		newSpec.Mode.Replicated.Replicas = &serviceReplicaCount
		newSpec.TaskTemplate.ForceUpdate++
		log.Println("forcing update on", m.DesiredSpecs[key].Name)

		log.Println("updating service...", m.DesiredSpecs[key].Name)
		serviceUpdateResponse, err := m.Client.ServiceUpdate(m.Ctx, serviceID, dockerswarm.Version{onlineService.Version.Index}, newSpec, types.ServiceUpdateOptions{})
		log.Println("update done", m.DesiredSpecs[key].Name)
		if err != nil {
			log.Panic(err)
		}
		if len(serviceUpdateResponse.Warnings) > 0 {
			log.Println("updating", m.DesiredSpecs[key].Name, "warnings", serviceUpdateResponse.Warnings)
		}
	}
	m.StackStateCh <- StackStateWaitForServicesToBeDeployed
}

// ScaleOnlyUpdate ...
func (m *Manager) ScaleOnlyUpdate(serviceName, serviceID string, count int) error {
	log.Println("scaling", serviceName, "to", count)
	cmd := exec.Command("docker", "-H", m.Host, "service", "scale", serviceName+"="+strconv.Itoa(count))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("scaling %s:%s failed with error: %w; %s", serviceName, serviceID[:12], err, string(out))
	}
	return nil
}
