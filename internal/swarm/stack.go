package swarm

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// RemoveStack ...
func (s *Manager) RemoveStack(host string, stackName string) error {
	fmt.Println("removing stack")
	cmd := exec.Command("docker", "-H", host, "stack", "remove", stackName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error while removing stack with error: %w; %s", err, string(out))
	}
	s.StackStateCh <- StackStateEmpty
	return nil
}

// DeployStackWithDockerCompose ....
func (s *Manager) DeployStackWithDockerCompose(dockerComposePath string, host string, stackName string, attempt int) error {
	fmt.Println("deploying stack")
	cmd := exec.Command("docker", "-H", host, "stack", "deploy", "--compose-file", dockerComposePath, stackName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "not found") && strings.Contains(string(out), "network") && attempt <= 25 {
			var waitTime int64 = 5
			fmt.Printf("deploying stack, attempt %d failed. Wait %d seconds\n", attempt, waitTime)
			time.Sleep(time.Duration(waitTime) * time.Second)
			return s.DeployStackWithDockerCompose(dockerComposePath, host, stackName, attempt+1)
		}
		return fmt.Errorf("deploying stack with docker compose file failed with error: %w; %s", err, string(out))
	}
	s.StackStateCh <- StackStateWaitForServicesToBeDeployed

	go func(s *Manager) {
		for {
			if s.CurrentStackState >= StackStateServicesAreDeployed {
				s.FillDesiredSpecsCurrentSpecs()
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}(s)

	return nil
}

// FillDesiredSpecsCurrentSpecs ...
func (s *Manager) FillDesiredSpecsCurrentSpecs() { //TODO update
	s.UpdateCurrentSpecs()
	fmt.Println("Filling Desired Specs with Current Specs")
	for serviceID := range s.CurrentSpecs {
		// containers := make([]string, len(s.CurrentSpecs[serviceID].Containers))
		// for i, c := range s.CurrentSpecs[serviceID].Containers {
		// 	containers[i] = c
		// }

		envs := make([]string, len(s.CurrentSpecs[serviceID].EnvironmentVariables))
		for i, e := range s.CurrentSpecs[serviceID].EnvironmentVariables {
			envs[i] = e
		}
		s.DesiredSpecs[serviceID] = ServiceSpecs{
			ID:                   serviceID,
			Name:                 s.CurrentSpecs[serviceID].Name,
			ImageName:            s.CurrentSpecs[serviceID].ImageName,
			ReplicaCount:         s.CurrentSpecs[serviceID].ReplicaCount,
			EnvironmentVariables: envs,
			StackName:            s.CurrentSpecs[serviceID].StackName,
			CPULimits:            s.CurrentSpecs[serviceID].CPULimits,
			CPUReservation:       s.CurrentSpecs[serviceID].CPUReservation,
			MemoryLimits:         s.CurrentSpecs[serviceID].MemoryLimits,
			MemoryReservations:   s.CurrentSpecs[serviceID].MemoryReservations,
		}
	}
}
