package initializer

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	l := loadgenerator.NewK6LoadGenerator("http://136.159.209.214:7112")
	log.Println("workload string:", workloadStr)
	parts := strings.Split(workloadStr, "_")
	vus, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	duration, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	if duration < configs.GetConfig().Test.Duration {
		panic("for now these two values should be equal or duration should be more than TestDuration!")
	}
	authProb, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		panic(err)
	}
	bookProb := 1 - authProb
	sleepDuration, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		panic(err)
	}
	script := loadgenerator.CreateLoadGeneartorScript(configs.GetConfig().LoadGenerator.Details["script"], vus, duration, authProb, bookProb, 0, sleepDuration)
	l.Prepare(map[string]string{"script": script})
	return l
}

// GetSwarmManager ...
func GetSwarmManager() *swarm.Manager {
	m, err := swarm.GetNewSwarmManager(map[string]string{"stackname": configs.GetConfig().TestBed.StackName, "host": configs.GetConfig().Host.Host})
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

	if strings.Contains(configs.GetConfig().Results.Path, "$STRATEGY") {
		configs.GetConfig().Results.Path = strings.Replace(configs.GetConfig().Results.Path, "$STRATEGY", strategyName, 1)
		log.Println("Updating result path to", configs.GetConfig().Results.Path)
	}

	workloadStr := viper.GetString("workloadStr")

	var ruc = resource.GetTheResourceUsageCollector()
	var rtc workload.ResponseTimeCollector = getJaegerCollector()
	var rcc workload.RequestCountCollector = rtc.(workload.RequestCountCollector)
	var lg = GetTheLoadGenerator(workloadStr)

	// var c strategies.Configurer
	// switch os.Args[beforeConfigArgCount-1] {
	// case "CPUUsageIncrease":
	// 	c = GetCPUIncreaseConfigurer()
	// case "ResponseTimeSimpleIncrease":
	// 	c = GetResponseTimeSimpleIncreaseConfigurer()
	// case "CPUUtil_RT_Hybrid":
	// 	c = GetCPUUtilRTHybridConfigurer()
	// case "PredefinedSearch":
	// 	c = strategies.GetNewPredefinedSearcher()
	// case "MOBO":
	// 	c = GetMOBOConfigurer()
	// case "AddFractionalCPUcores":
	// 	c = GetAddFractionalCPUcoresConfigurer()
	// case "Single":
	// 	c = &strategies.SingleRun{}
	// case "AddDifferentFractionalCPUcores":
	// 	c = GetAddDifferentFractionalCPUcores(workloadStr)
	// case "PerPathEUBasedScaling":
	// 	c = GetPerPathEUBasedScaling(workloadStr)
	// default:
	// 	log.Println("expected 'Single' or 'CPUUsageIncrease' or 'ResponseTimeSimpleIncrease' or 'PredefinedSearch' subcommands but got", os.Args[beforeConfigArgCount-1])
	// 	os.Exit(1)
	// }
	// var c = GetAnotherConfigurer()
	var m = GetSwarmManager()
	db := GetNewDatabase()
	a := autoconfigure.NewAutoConfigurer(lg, rtc, rcc, ruc, strategy, m, workloadStr, db)
	log.Println("name of the test is:", viper.GetString("testName"))
	a.Start(viper.GetString("testName"), strings.Join(os.Args, " "))
}
