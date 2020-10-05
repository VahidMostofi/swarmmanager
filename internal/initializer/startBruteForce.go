package initializer

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/bruteforce"
	"github.com/VahidMostofi/swarmmanager/internal/workload"
	"github.com/spf13/viper"
)

// StartBruteForce ...
func StartBruteForce() {

	if strings.Contains(configs.GetConfig().Results.Path, "$STRATEGY") {
		configs.GetConfig().Results.Path = strings.Replace(configs.GetConfig().Results.Path, "$STRATEGY", "brute-force", 1)
		log.Println("Updating result path to", configs.GetConfig().Results.Path)
	}

	// creating directories for ResultDirectoryPath
	if err := os.MkdirAll(filepath.Dir(configs.GetConfig().Results.Path), 0777); err != nil {
		log.Panic(err)
	}

	workloadStr := viper.GetString("workloadStr")

	var rtc workload.ResponseTimeCollector = getJaegerCollector()
	var rcc workload.RequestCountCollector = rtc.(workload.RequestCountCollector)
	var lg = GetTheLoadGenerator(workloadStr)
	var m = GetSwarmManager(false)

	b := bruteforce.NewBruteForce(lg, rtc, rcc, m, workloadStr)
	b.Start()
	log.Println("name of the test is:", viper.GetString("testName"))
}
