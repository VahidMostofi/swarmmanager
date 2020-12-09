package initializer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/autoconfigure"
	"github.com/VahidMostofi/swarmmanager/internal/caching"
	"github.com/VahidMostofi/swarmmanager/internal/jaeger"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	resource "github.com/VahidMostofi/swarmmanager/internal/resource/collector"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/VahidMostofi/swarmmanager/internal/workload"
	"github.com/spf13/viper"
)

// getJaegerCollector ...
func getJaegerCollector() *jaeger.Aggregator {
	j := jaeger.NewAggregator()
	return j
}

// GetTheLoadGenerator ...
func GetTheLoadGenerator(workloadStr string) loadgenerator.LoadGenerator {
	l, err := loadgenerator.GetLoadGenerator()
	if err != nil {
		log.Panic(err)
	}
	err = l.Prepare()
	if err != nil {
		log.Panic(fmt.Errorf("error while preparing load generator: %w", err))
	}
	return l
}

// GetSwarmManager ...
func GetSwarmManager(monitorSpecs bool) *swarm.Manager {
	m, err := swarm.GetNewSwarmManager(map[string]string{"stackname": configs.GetConfig().TestBed.StackName, "host": configs.GetConfig().Host.Host}, monitorSpecs)
	if err != nil {
		log.Panic(err)
	}
	return m
}

// GetNewDatabase ...
func GetNewDatabase() caching.Database {
	db, err := caching.GetNewDropboxDatabase(configs.GetConfig().Cache.Details["path"])
	if err != nil {
		panic(fmt.Errorf("error in getting mongo database for caching: %w", err))
	}
	return db
}

// StartAutoconfig starts the auto configuration based on provided strategy
func StartAutoconfig(strategy strategies.Configurer, strategyName string) {
	isK8s := false
	if isK8s == true {
		StartK8sAutoConfigurer(strategy, strategyName)
		return
	}

	if strings.Contains(configs.GetConfig().Results.Path, "$STRATEGY") {
		configs.GetConfig().Results.Path = strings.Replace(configs.GetConfig().Results.Path, "$STRATEGY", strategyName, 1)
		log.Println("Updating result path to", configs.GetConfig().Results.Path)
	}

	// creating directories for ResultDirectoryPath
	if err := os.MkdirAll(filepath.Dir(configs.GetConfig().Results.Path), 0777); err != nil {
		log.Panic(err)
	}

	workloadStr := viper.GetString("workloadStr")

	var ruc = resource.GetTheResourceUsageCollector()
	var rtc workload.ResponseTimeCollector = getJaegerCollector()
	var rcc workload.RequestCountCollector = rtc.(workload.RequestCountCollector)
	var lg = GetTheLoadGenerator(workloadStr)

	var m = GetSwarmManager(true)
	db := GetNewDatabase()
	a := autoconfigure.NewAutoConfigurer(lg, rtc, rcc, ruc, strategy, m, workloadStr, db)
	log.Println("name of the test is:", viper.GetString("testName"))
	a.Start(viper.GetString("testName"), strings.Join(os.Args, " "))
}
