package k8s

import (
	"testing"

	"github.com/VahidMostofi/swarmmanager/internal/swarm"
)

// func TestCPUUsage(t *testing.T) {
// 	configs.FakeInitialize()
// 	configs.GetConfig().TestBed.ServicesToConfigure = []string{"auth", "gateway", "books"}
// 	connector := GetNewConnector("ssh", "136.159.209.204")
// 	res, err := connector.GetCPUUsage()
// 	if err != nil {
// 		panic(err)
// 	}
// 	for k, v := range res {
// 		fmt.Print(k, " ", v)
// 		fmt.Println()
// 	}
// }

// func TestGetPods(t *testing.T) {
// 	connector := GetNewConnector("ssh", "136.159.209.204")
// 	res := connector.GetCurrentPods()
// 	fmt.Println(len(res))
// }

func TestApplyConfig(t *testing.T) {
	connector := GetNewConnector("ssh", "136.159.209.204")
	err := connector.ApplyConfig(map[string]swarm.ServiceSpecs{"auth": {ReplicaCount: 2, CPULimits: 0.31}, "books": {ReplicaCount: 3, CPULimits: 0.12}, "gateway": {ReplicaCount: 1, CPULimits: 0.28}})
	if err != nil {
		panic(err)
	}
}
