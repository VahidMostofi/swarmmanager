package swarm

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"log"

	"github.com/VahidMostofi/swarmmanager/internal/utils"
)

// TODO add attempts to this too
// RemoveStack ...
func (s *Manager) RemoveStack(attempt int) error {
	log.Println("removing stack")
	cmd := exec.Command("docker", "-H", s.Host, "stack", "remove", s.StackName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if attempt <= 5 {
			var waitTime int64 = 5
			log.Printf("removing stack, attempt %d failed. Wait %d seconds\n", attempt, waitTime)
			time.Sleep(time.Duration(waitTime) * time.Second)
			return s.RemoveStack(attempt + 1)
		}
		return fmt.Errorf("error while removing stack with error: %w; %s", err, string(out))
	}
	s.StackStateCh <- StackStateEmpty
	return nil
}

// DeployStackWithDockerCompose ....
func (s *Manager) DeployStackWithDockerCompose(dockerComposePath string, attempt int, initialConfig map[string]SimpleSpecs) error {
	log.Println("deploying stack")
	cmd := exec.Command("docker", "-H", s.Host, "stack", "deploy", "--compose-file", dockerComposePath, s.StackName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if (strings.Contains(string(out), "not found") || strings.Contains(string(out), "cannot be used with services.")) && strings.Contains(string(out), "network") && attempt <= 25 {
			var waitTime int64 = 5
			log.Printf("deploying stack, attempt %d failed. Wait %d seconds\n", attempt, waitTime)
			time.Sleep(time.Duration(waitTime) * time.Second)
			return s.DeployStackWithDockerCompose(dockerComposePath, attempt+1, initialConfig)
		}
		return fmt.Errorf("deploying stack with docker compose file failed with error: %w; %s", err, string(out))
	}

	s.StackStateCh <- StackStateWaitForServicesToBeDeployed

	go func(s *Manager) {
		for {
			if s.CurrentStackState >= StackStateServicesAreDeployed {
				s.FillDesiredSpecsCurrentSpecs()
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if len(initialConfig) > 0 {
			for serviceName, ss := range initialConfig {
				temp := s.DesiredSpecs[serviceName]
				temp.CPULimits = ss.CPU
				temp.CPUReservation = ss.CPU
				temp.ReplicaCount = ss.Replica
				temp.EnvironmentVariables = utils.UpdateENVWorkerCounts(temp.EnvironmentVariables, ss.Worker)
				s.DesiredSpecs[serviceName] = temp
			}
		}
		s.StackStateCh <- StackStateMustCompare
	}(s)

	return nil
}

// FillDesiredSpecsCurrentSpecs ...
func (s *Manager) FillDesiredSpecsCurrentSpecs() { //TODO update
	s.UpdateCurrentSpecs()
	log.Println("Filling Desired Specs with Current Specs")
	for key := range s.CurrentSpecs {
		envs := make([]string, len(s.CurrentSpecs[key].EnvironmentVariables))
		for i, e := range s.CurrentSpecs[key].EnvironmentVariables {
			envs[i] = e
		}
		s.DesiredSpecs[key] = ServiceSpecs{
			ID:                   s.CurrentSpecs[key].ID,
			Name:                 s.CurrentSpecs[key].Name,
			ImageName:            s.CurrentSpecs[key].ImageName,
			ReplicaCount:         s.CurrentSpecs[key].ReplicaCount,
			EnvironmentVariables: envs,
			StackName:            s.CurrentSpecs[key].StackName,
			CPULimits:            s.CurrentSpecs[key].CPULimits,
			CPUReservation:       s.CurrentSpecs[key].CPUReservation,
			MemoryLimits:         s.CurrentSpecs[key].MemoryLimits,
			MemoryReservations:   s.CurrentSpecs[key].MemoryReservations,
		}
	}
}
