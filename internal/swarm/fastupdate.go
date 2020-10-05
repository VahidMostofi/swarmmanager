package swarm

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
)

// FastUpdate ...
func (m *Manager) FastUpdate(config map[string]SimpleSpecs, prev map[string]SimpleSpecs) error {
	wg := sync.WaitGroup{}
	for serviceName, c := range config {
		wg.Add(1)
		go func(serviceName string, c SimpleSpecs) {
			// log.Println("updating", serviceName)
			err := m.fastUpdateService(serviceName, c, prev[serviceName], config)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(serviceName, c)
	}
	go func() {
		wg.Add(1)
		cmd := exec.Command("docker", "-H", m.Host, "service", "update", "--force", m.StackName+"_jaeger")
		b, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(b))
			panic("couldn't resetart jaeger")
		}
		wg.Done()
	}()
	wg.Wait()
	time.Sleep(5 * time.Second)
	return nil
}

func (m *Manager) fastUpdateService(serviceName string, c, prev SimpleSpecs, config map[string]SimpleSpecs) error {

	if prev.Replica != c.Replica {
		replicaCountStr := strconv.Itoa(c.Replica)
		cmd := exec.Command("docker", "-H", m.Host, "service", "scale", m.StackName+"_"+serviceName+"="+replicaCountStr)
		b, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(b))
			return fmt.Errorf("error while updating scale count: %w", err)
		}
	}

	containers, err := m.Client.ContainerList(m.Ctx, types.ContainerListOptions{})
	if err != nil {
		return fmt.Errorf("error while retrieving containers: %w", err)
	}

	for _, container := range containers {
		if strings.Contains(serviceName, m.findService(container.Names[0], config)) {

			cpuCountStr := strconv.FormatFloat(c.CPU, 'f', 3, 64)
			cmd := exec.Command("docker", "-H", m.Host, "container", "update", "--cpus", cpuCountStr, container.ID)

			b, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(string(b))
				return err
			}
		}
	}

	return nil
}

func (m *Manager) findService(containerName string, cs map[string]SimpleSpecs) string {
	for serviceName := range cs {
		if strings.Contains(containerName, m.StackName+"_"+serviceName) {
			return serviceName
		}
	}
	return "NOSERVICE_TO_MONITOR_FOUND"
}
