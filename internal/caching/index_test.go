package caching

import (
	"fmt"
	"testing"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/spf13/viper"
)

func TestIndex(t *testing.T) {
	configs.FakeInitialize()
	configs.GetConfig().Cache.Details = map[string]string{"path": "/Users/vahid/Dropbox/data/swarm-manager-data/cache-v2/"}
	configs.GetConfig().AppName = "bookstore_nodejs"
	configs.GetConfig().Version = "v1"
	configs.GetConfig().Test.Duration = 60
	viper.Set("workloadStr", "100_110_0.33_0.33_0.34")
	fmt.Println(FindBaseConfiguration(250, 3).HashCode)
}
