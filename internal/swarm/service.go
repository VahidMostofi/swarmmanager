package swarm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"log"

	"github.com/docker/docker/api/types"
	dockerswarm "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// ServiceSpecs ...
type ServiceSpecs struct {
	ID                   string   `yaml:"-"`
	Name                 string   `yaml:"-"`
	ImageName            string   `yaml:"-"`
	ReplicaCount         int      `yaml:"replicaCount"`
	EnvironmentVariables []string `yaml:"envs"`
	StackName            string   `yaml:"-"`
	CPULimits            float64  `yaml:"CPULimits"`
	CPUReservation       float64  `yaml:"CPUReservation"`
	MemoryLimits         int64    `yaml:"-"`
	MemoryReservations   int64    `yaml:"-"`
	Containers           []string `yaml:"-"`
}

// StackSpecs ...
type StackSpecs map[string]ServiceSpecs

// Manager manages the swarm cluster
type Manager struct {
	Client            *client.Client
	Host              string
	Ctx               context.Context
	CtxCancelFunc     context.CancelFunc
	StackName         string
	DesiredSpecs      map[string]ServiceSpecs
	CurrentSpecs      map[string]ServiceSpecs
	StackStateCh      chan int
	CurrentStackState int
	ServicesToManage  []string
}

// ToHumanReadable ...
func (s *Manager) ToHumanReadable(m map[string]ServiceSpecs) map[string]ServiceSpecs {
	m2 := make(map[string]ServiceSpecs)
	for _, value := range m {
		flag := false
		for _, str := range s.ServicesToManage {
			if str == value.Name {
				flag = true
				break
			}
		}
		if flag {
			m2[value.Name] = value
		}
	}
	return m2
}

// StackStates ...
const (
	StackStateEmpty                       = 0
	StackStateWaitForServicesToBeDeployed = 5
	StackStateServicesAreDeployed         = 10
	StackStateWaitForServicesToBeReady    = 15
	StackStateServicesAreReady            = 20
	StackStateMustCompare                 = 22
	StackStateUpdatingSpecs               = 25
)

// GetNewSwarmManager is constructor
// required fields: host, stackname, services: one string with comma separated services' names
func GetNewSwarmManager(values map[string]string) (*Manager, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(values["host"], "", nil, defaultHeaders)
	if err != nil {
		return nil, fmt.Errorf("error while creating Docker client in Manager.GetNewSwarmManager: %w", err)
	}

	if _, ok := values["stackname"]; !ok {
		return nil, fmt.Errorf("no stackname is provided in the values map")
	}

	if _, ok := values["services"]; !ok {
		return nil, fmt.Errorf("no 'services' field is provided in the value map")
	}

	servicesToMonitor := strings.Split(values["services"], ",")
	for i, s := range servicesToMonitor {
		servicesToMonitor[i] = strings.Trim(s, " ")
	}
	log.Println("Services to monitor:", servicesToMonitor)
	ctx, cancelFunc := context.WithCancel(context.Background())

	m := &Manager{
		Client:           cli,
		Host:             values["host"],
		Ctx:              ctx,
		CtxCancelFunc:    cancelFunc,
		StackName:        values["stackname"],
		CurrentSpecs:     make(map[string]ServiceSpecs),
		DesiredSpecs:     make(map[string]ServiceSpecs),
		StackStateCh:     make(chan int),
		ServicesToManage: servicesToMonitor,
	}

	go m.monitorState()
	go m.manageState()
	go m.monitorSpecs()
	return m, nil
}

func (s *Manager) monitorSpecs() {
	waitTime := 10
	log.Printf("monitoring specs every %d seconds\n", waitTime)
	for {
		if s.CurrentStackState >= StackStateServicesAreReady || s.CurrentStackState == StackStateMustCompare {
			err := s.UpdateCurrentSpecs()
			if err != nil {
				log.Panic(err)
			}
			comparision := s.CompareSpecs()
			if !comparision {
				log.Println("specs comparision is:", comparision)
			}
			if !comparision {
				s.UpdateServices()
			} else {
				//todo remove this
				// specs := s.DesiredSpecs["books"]
				// specs.ReplicaCount = 3
				// s.DesiredSpecs["books"] = specs
				// fmt.Println(s.CurrentSpecs["books"])
			}
		}
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}

func (s *Manager) manageState() {
	for {
		if s.CurrentStackState == StackStateEmpty {
			//
		} else if s.CurrentStackState <= StackStateWaitForServicesToBeDeployed {
			s.CheckServicedDeployment(5)
		} else if s.CurrentStackState <= StackStateWaitForServicesToBeReady {
			s.CheckforServicesReadiness()
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *Manager) monitorState() {
	for {
		select {
		case newState := <-s.StackStateCh:
			if newState != s.CurrentStackState {
				log.Println("changed state to:", GetStateString(newState))
				s.CurrentStackState = newState
			}
		}
	}
}

func (s *Manager) removeStackName(serviceName string) string {
	return strings.Replace(serviceName, s.StackName+"_", "", 1)
}

// UpdateCurrentSpecs ...
func (s *Manager) UpdateCurrentSpecs() error {
	services, err := s.Client.ServiceList(s.Ctx, types.ServiceListOptions{})
	if err != nil {
		return fmt.Errorf("error in manager.UpdateCurrentSpecs(): %w", err)
	}
	s.CurrentSpecs = make(map[string]ServiceSpecs)
	for _, service := range services {
		serviceName := s.removeStackName(service.Spec.Name)
		serviceID := service.ID

		s.CurrentSpecs[serviceID] = ServiceSpecs{}
		tempState := s.CurrentSpecs[serviceID]
		tempState.ID = service.ID
		tempState.ImageName = service.Spec.Labels["com.docker.stack.image"]
		tempState.EnvironmentVariables = service.Spec.TaskTemplate.ContainerSpec.Env
		tempState.Name = serviceName
		tempState.ReplicaCount = int(*service.Spec.Mode.Replicated.Replicas)
		tempState.StackName = service.Spec.Labels["com.docker.stack.namespace"]
		tempState.CPULimits = float64(service.Spec.TaskTemplate.Resources.Limits.NanoCPUs) / 1e9
		tempState.CPUReservation = float64(service.Spec.TaskTemplate.Resources.Reservations.NanoCPUs) / 1e9
		tempState.MemoryLimits = service.Spec.TaskTemplate.Resources.Limits.MemoryBytes
		tempState.MemoryReservations = service.Spec.TaskTemplate.Resources.Reservations.MemoryBytes
		s.CurrentSpecs[serviceID] = tempState
	}
	return nil
}

// UpdateServicesSpecs ...
func (s *Manager) UpdateServicesSpecs() error {

	// check running containers for services
	for serviceID := range s.CurrentSpecs {
		temp := s.CurrentSpecs[serviceID]
		temp.Containers = make([]string, 0)
		s.CurrentSpecs[serviceID] = temp
	}
	tasks, err := s.Client.TaskList(s.Ctx, types.TaskListOptions{})
	if err != nil {
		return fmt.Errorf("error while retrieving tasks: %w", err)
	}
	// fmt.Println("tasks")
	for _, t := range tasks {
		// fmt.Println(t)
		if t.Status.State == "running" && (time.Now().UnixNano()-t.Status.Timestamp.UnixNano())/1e9 > 10 {
			temp := s.CurrentSpecs[t.ServiceID]
			temp.Containers = append(temp.Containers, t.Status.ContainerStatus.ContainerID)
			s.CurrentSpecs[t.ServiceID] = temp
		}
	}
	return nil
}

// comapeServiceSpecs ... returns true if they are equal
func (s *Manager) comapeServiceSpecs(serviceID string) (bool, []string) {
	changes := []string{}
	if s.CurrentSpecs[serviceID].ImageName != s.DesiredSpecs[serviceID].ImageName {
		log.Println("CompareSpecs: ImageName is changed")
		changes = append(changes, "ImageName")
	}
	if s.CurrentSpecs[serviceID].ReplicaCount != s.DesiredSpecs[serviceID].ReplicaCount {
		log.Println("CompareSpecs: ReplicaCount is changed")
		changes = append(changes, "ReplicaCount")
	}
	if s.CurrentSpecs[serviceID].CPULimits != s.DesiredSpecs[serviceID].CPULimits {
		log.Println("CompareSpecs: CPULimits is changed")
		changes = append(changes, "CPULimits")
	}
	if s.CurrentSpecs[serviceID].CPUReservation != s.DesiredSpecs[serviceID].CPUReservation {
		log.Println("CompareSpecs: CPUReservation is changed")
		changes = append(changes, "CPUReservation")
	}
	if s.CurrentSpecs[serviceID].MemoryLimits != s.DesiredSpecs[serviceID].MemoryLimits {
		log.Println("CompareSpecs: MemoryLimits is changed")
		changes = append(changes, "MemoryLimits")
	}
	if s.CurrentSpecs[serviceID].MemoryReservations != s.DesiredSpecs[serviceID].MemoryReservations {
		log.Println("CompareSpecs: MemoryReservations is changed")
		changes = append(changes, "MemoryReservations")
	}

	if !Equal(s.CurrentSpecs[serviceID].EnvironmentVariables, s.DesiredSpecs[serviceID].EnvironmentVariables) {
		log.Println("CompareSpecs: EnvironmentVariables is changed")
		changes = append(changes, "EnvironmentVariables")
	}
	return len(changes) == 0, changes
}

// CompareSpecs ...
func (s *Manager) CompareSpecs() bool {
	for serviceID := range s.CurrentSpecs {
		flag, _ := s.comapeServiceSpecs(serviceID)
		if !flag {
			return false
		}
	}
	return true
}

// Equal ...
func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, ai := range a {
		flag := false
		for _, bi := range b {
			if ai == bi {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

func listOfContainersToString(cs []string) string {
	res := ""
	for _, c := range cs {
		res += c[:10] + ", "
	}
	return res
}

// IsServiceReady ...
func (s *Manager) IsServiceReady(serviceID string) bool {
	if len(s.CurrentSpecs[serviceID].Containers) == s.DesiredSpecs[serviceID].ReplicaCount {
		return true
	}
	// fmt.Println(s.CurrentSpecs[serviceID].Name, listOfContainersToString(s.CurrentSpecs[serviceID].Containers), s.DesiredSpecs[serviceID].ReplicaCount)
	return false
}

// CheckforServicesReadiness ...
func (s *Manager) CheckforServicesReadiness() {
	s.StackStateCh <- StackStateWaitForServicesToBeReady
	flag := true
	err := s.UpdateServicesSpecs()
	if err != nil {
		log.Panic(err)
	}

	for serviceID := range s.CurrentSpecs {
		if !s.IsServiceReady(serviceID) {
			// fmt.Println(s.CurrentSpecs[serviceID].Name, "is not ready")
			flag = false
			break
		}
	}

	if flag {
		s.StackStateCh <- StackStateServicesAreReady
	}
}

// CheckServicedDeployment ...
func (s *Manager) CheckServicedDeployment(numberOfServices int) {
	var err error
	var services []dockerswarm.Service
	s.StackStateCh <- StackStateWaitForServicesToBeDeployed

	services, err = s.Client.ServiceList(s.Ctx, types.ServiceListOptions{})
	if err != nil {
		log.Panic(err)
	}
	if len(services) == numberOfServices {
		s.StackStateCh <- StackStateServicesAreDeployed
	}
}

// GetStateString ...
func GetStateString(stateValue int) string {
	if stateValue == StackStateEmpty {
		return "StackState Empty"
	}
	if stateValue == StackStateWaitForServicesToBeDeployed {
		return "StackState WaitForServicesToBeDeployed"
	}
	if stateValue == StackStateServicesAreDeployed {
		return "StackState ServicesAreDeployed"
	}
	if stateValue == StackStateWaitForServicesToBeReady {
		return "StackState WaitForServicesToBeReady"
	}
	if stateValue == StackStateServicesAreReady {
		return "StackState ServicesAreReady"
	}
	if stateValue == StackStateUpdatingSpecs {
		return "StackState UpdatingSpecs"
	}
	if stateValue == StackStateMustCompare {
		return "StackState MustCompare"
	}
	return "unknown state"
}
