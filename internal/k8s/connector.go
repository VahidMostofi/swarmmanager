package k8s

import (
	"strconv"
	"strings"
	"time"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/VahidMostofi/swarmmanager/internal/utils"
	"github.com/montanaflynn/stats"
	"gopkg.in/yaml.v2"
)

// Connector ...
type Connector interface {
	GetCPUUsage() (map[string]float64, error)
	GetCurrentConfiguration() (map[string]swarm.SimpleSpecs, error)
	GetCurrentPods() []string
	AreAllPodsRunning() bool
}

// GetNewConnector ...
func GetNewConnector(kind string, host string) Connector {
	if kind == "ssh" {
		return &sshConnector{
			executor: sshExecuter{Host: host},
		}
	}
	panic("unknown kind")
}

type sshConnector struct {
	executor sshExecuter
}

// AreAllPodsRunning ...
func (s *sshConnector) AreAllPodsRunning() bool {
	sRes := s.executor.executeCommand("kubectl get pods -o yaml")
	b := []byte(sRes)
	data := make(map[string]interface{})
	yaml.Unmarshal(b, data)

	for _, item := range data["items"].([]interface{}) {
		phase := item.(map[interface{}]interface{})["status"].(map[interface{}]interface{})["phase"]

		if phase != "Running" {
			return false
		}

	}
	return true
}

func (s *sshConnector) GetCurrentPods() []string {
	sRes := s.executor.executeCommand("kubectl get pods -o yaml")
	b := []byte(sRes)
	data := make(map[string]interface{})
	yaml.Unmarshal(b, data)

	pods := make([]string, 0)

	for _, item := range data["items"].([]interface{}) {
		name := item.(map[interface{}]interface{})["metadata"].(map[interface{}]interface{})["name"].(string)

		phase := item.(map[interface{}]interface{})["status"].(map[interface{}]interface{})["phase"]

		if phase != "Running" {
			panic("one of the pods is not running!")
		}

		pods = append(pods, name)
	}
	return pods
}

func (s *sshConnector) GetCurrentConfiguration() (map[string]swarm.SimpleSpecs, error) {
	for !s.AreAllPodsRunning() {
		time.Sleep(1 * time.Second)
	}

	sRes := s.executor.executeCommand("kubectl get deployment -o yaml")
	b := []byte(sRes)
	data := make(map[string]interface{})
	yaml.Unmarshal(b, data)
	configs := make(map[string]swarm.SimpleSpecs)

	for _, item := range data["items"].([]interface{}) {
		name := item.(map[interface{}]interface{})["metadata"].(map[interface{}]interface{})["name"].(string)

		if _, ok := configs[name]; ok {
			continue
		}

		replicas := item.(map[interface{}]interface{})["spec"].(map[interface{}]interface{})["replicas"].(int)

		containers := item.(map[interface{}]interface{})["spec"].(map[interface{}]interface{})["template"].(map[interface{}]interface{})["spec"].(map[interface{}]interface{})["containers"].([]interface{})
		container := containers[0]

		cpuStr := container.(map[interface{}]interface{})["resources"].(map[interface{}]interface{})["limits"].(map[interface{}]interface{})["cpu"].(string)
		cpu, err := strconv.ParseFloat(cpuStr[:len(cpuStr)-1], 64)
		if err != nil {
			panic(err)
		}
		cpu /= 1000

		configs[name] = swarm.SimpleSpecs{CPU: cpu, Worker: 1, Replica: replicas}
	}

	return configs, nil
}

func (s *sshConnector) GetCPUUsage() (map[string]float64, error) {
	currentSpecs, err := s.GetCurrentConfiguration()
	if err != nil {
		panic(err)
	}
	values := make(map[string][]float64)
	for serviceName := range currentSpecs {
		values[serviceName] = make([]float64, 0)
	}
	sRes := s.executor.executeCommand("kubectl top pods")
	for iL, line := range strings.Split(sRes, "\n") {
		if iL == 0 {
			continue
		}
		line = strings.Trim(line, " ")
		if len(line) < 2 {
			continue
		}
		for strings.Contains(line, "  ") {
			line = strings.ReplaceAll(line, "  ", " ")
		}
		parts := strings.Split(line, " ")
		name := strings.Split(parts[0], "-")[0]
		if !utils.ContainsString(configs.GetConfig().TestBed.ServicesToConfigure, name) {
			continue
		}
		cpuF, err := strconv.ParseFloat(parts[1][:len(parts[1])-1], 64)
		if err != nil {
			panic(err)
		}
		cpuF /= 1000
		utilization := cpuF / currentSpecs[name].CPU
		values[name] = append(values[name], utilization)
	}
	res := make(map[string]float64)
	for serviceName, utils := range values {	
		res[serviceName], err = stats.Mean(utils)
		if err != nil {
			panic(err)
		}
	}
	return res, nil
}
