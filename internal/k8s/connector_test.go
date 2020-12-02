package k8s

import (
	"fmt"
	"testing"

	"github.com/VahidMostofi/swarmmanager/configs"
)

func TestCPUUsage(t *testing.T) {
	configs.FakeInitialize()
	configs.GetConfig().TestBed.ServicesToConfigure = []string{"auth", "gateway", "books"}
	connector := GetNewConnector("ssh", "136.159.209.204")
	res, err := connector.GetCPUUsage()
	if err != nil {
		panic(err)
	}
	for k, v := range res {
		fmt.Print(k, " ", v)
		fmt.Println()
	}
}

func TestGetPods(t *testing.T) {
	connector := GetNewConnector("ssh", "136.159.209.204")
	res := connector.GetCurrentPods()
	fmt.Println(len(res))
}
