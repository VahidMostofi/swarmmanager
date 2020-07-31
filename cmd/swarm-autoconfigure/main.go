package main

//TODO I NEED TO REstart every single container
import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"flag"

	"github.com/VahidMostofi/swarmmanager"
	"github.com/VahidMostofi/swarmmanager/internal/autoconfigure"
	"github.com/VahidMostofi/swarmmanager/internal/caching"
	"github.com/VahidMostofi/swarmmanager/internal/jaeger"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	resource "github.com/VahidMostofi/swarmmanager/internal/resource/collector"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/VahidMostofi/swarmmanager/internal/workload"
)

const beforeConfigArgCount = 4

// GetTheResourceUsageCollector ...
func GetTheResourceUsageCollector() resource.Collector {
	stackName := swarmmanager.GetConfig().StackName
	c := resource.GetNewCollector("SingleCollector")
	err := c.Configure(map[string]string{"host": swarmmanager.GetConfig().Host, "stackname": stackName})
	if err != nil {
		log.Panic(err)
	}

	return c
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
	if duration < swarmmanager.GetConfig().TestDuration {
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
	script := loadgenerator.CreateLoadGeneartorScript(swarmmanager.GetConfig().K6Script, vus, duration, authProb, bookProb, 0, sleepDuration)
	l.Prepare(map[string]string{"script": script})
	return l

	// StartLG(l)
	// StopLG(l)
	// FeedbackLG(l)
}

// GetJaegerCollector ...
func GetJaegerCollector() *jaeger.JaegerAggregator {
	//TODO these are hardcoded!
	j := jaeger.NewJaegerAggregator(swarmmanager.GetConfig().JaegerHost, []string{"gateway", "auth", "books", "gateway", "auth_total", "auth_gateway", "auth_sub", "books_total", "books_gateway", "books_sub"})
	return j
}

// GetMOBOConfigurer ...
func GetMOBOConfigurer() strategies.Configurer {
	config := make(map[string]int)
	for i := 0; i < len(os.Args[beforeConfigArgCount:])/2; i++ {
		serviceName := strings.Trim(os.Args[beforeConfigArgCount+2*i], " ")
		count, err := strconv.Atoi(strings.Trim(os.Args[beforeConfigArgCount+2*i+1], " "))
		if err != nil {
			log.Panic("GetMOBOConfigurer: invalid input as count: %w", err)
		}
		config[serviceName] = count
	}

	return strategies.GetnewMOBOConfigurer(config)
}

// GetCPUIncreaseConfigurer ...
func GetCPUIncreaseConfigurer() strategies.Configurer {
	cpuOnlyCmd := flag.NewFlagSet("CPUUsageIncrease", flag.ExitOnError)
	cpuOnlyValueName := cpuOnlyCmd.String("property", "", "Which property of a run to consider? CPUUsageMean,CPUUsage90Percentile 70-95, 99")
	cpuOnlyThreshold := cpuOnlyCmd.Float64("threshold", 0, "what is the threshold")
	cpuOnlyCmd.Parse(os.Args[beforeConfigArgCount:])
	cpuOnlyCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[beforeConfigArgCount-1])
		cpuOnlyCmd.PrintDefaults()
	}

	if *cpuOnlyValueName == "" {
		cpuOnlyCmd.Usage()
		os.Exit(1)
	}

	if *cpuOnlyThreshold == 0 {
		cpuOnlyCmd.Usage()
		os.Exit(1)
	}
	log.Println("Configuring CPUUsageIncreaseConfigurer with Threshold:", *cpuOnlyThreshold, "and property of", *cpuOnlyValueName)
	return &strategies.CPUUsageIncrease{
		Threshold:       *cpuOnlyThreshold,
		ValueToConsider: *cpuOnlyValueName,
	}
}

// GetCPUUtilRTHybridConfigurer ... //TODO add input values for CPU util to this too
func GetCPUUtilRTHybridConfigurer() strategies.Configurer {
	rtsiCmd := flag.NewFlagSet("CPUUtil_RT_Hybrid", flag.ExitOnError)
	rtsiValueName := rtsiCmd.String("property", "", "Which property of a run to consider? CPUUsageMean,CPUUsage90Percentile 70-95, 99")
	rtsiThreshold := rtsiCmd.Float64("value", 0, "what is the threshold")
	rtsiCmd.Parse(os.Args[beforeConfigArgCount:])
	rtsiCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[beforeConfigArgCount-1])
		rtsiCmd.PrintDefaults()
	}

	if *rtsiValueName == "" {
		rtsiCmd.Usage()
		os.Exit(1)
	}

	if *rtsiThreshold == 0 {
		rtsiCmd.Usage()
		os.Exit(1)
	}
	log.Println("Configuring CPUUtil_RT_Hybrid with Value:", *rtsiThreshold, "and property of", *rtsiValueName)
	return &strategies.HybridCPUUtilResponseTimeSimpleIncrease{ //TODO add input values for CPU util to this too
		Agreements: []strategies.Agreement{
			{
				PropertyToConsider: *rtsiValueName,
				Value:              *rtsiThreshold,
			},
		},
	}
}

// GetResponseTimeSimpleIncreaseConfigurer ...
func GetResponseTimeSimpleIncreaseConfigurer() strategies.Configurer {
	rtsiCmd := flag.NewFlagSet("ResponseTimeSimpleIncrease", flag.ExitOnError)
	rtsiValueName := rtsiCmd.String("property", "", "Which property of a run to consider? CPUUsageMean,CPUUsage90Percentile 70-95, 99")
	rtsiThreshold := rtsiCmd.Float64("value", 0, "what is the threshold")
	rtsiCmd.Parse(os.Args[beforeConfigArgCount:])
	rtsiCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[beforeConfigArgCount-1])
		rtsiCmd.PrintDefaults()
	}

	if *rtsiValueName == "" {
		rtsiCmd.Usage()
		os.Exit(1)
	}

	if *rtsiThreshold == 0 {
		rtsiCmd.Usage()
		os.Exit(1)
	}
	log.Println("Configuring ResponseTimeSimpleIncrease with Value:", *rtsiThreshold, "and property of", *rtsiValueName)
	return &strategies.ResponseTimeSimpleIncrease{
		Agreements: []strategies.Agreement{
			{
				PropertyToConsider: *rtsiValueName,
				Value:              *rtsiThreshold,
			},
		},
	}
}

// GetSwarmManager ...
func GetSwarmManager() *swarm.Manager {
	m, err := swarm.GetNewSwarmManager(map[string]string{"stackname": swarmmanager.GetConfig().StackName, "host": swarmmanager.GetConfig().Host, "services": swarmmanager.GetConfig().ServicesToMonitor})
	if err != nil {
		log.Panic(err)
	}
	return m
}

// GetNewDatabase ...
func GetNewDatabase() caching.Database {
	db, err := caching.GetNewDropboxDatabase(swarmmanager.GetConfig().DropboxPath)
	if err != nil {
		panic(fmt.Errorf("error in getting mongo database for caching: %w", err))
	}
	return db
}

func main() {
	workloadStr := os.Args[1]
	if !strings.Contains(workloadStr, "_") {
		log.Panic("the first argument must be the workload")
	}

	if strings.Contains(swarmmanager.GetConfig().ResultsDirectoryPath, "$WORKLOAD") {
		swarmmanager.GetConfig().ResultsDirectoryPath = strings.Replace(swarmmanager.GetConfig().ResultsDirectoryPath, "$WORKLOAD", workloadStr, 1)
		log.Println("Updating result path to", swarmmanager.GetConfig().ResultsDirectoryPath)
	}

	var ruc = GetTheResourceUsageCollector()
	var rtc workload.ResponseTimeCollector = GetJaegerCollector()
	var rcc workload.RequestCountCollector = rtc.(workload.RequestCountCollector)
	var lg = GetTheLoadGenerator(workloadStr)

	if len(os.Args) < beforeConfigArgCount {
		fmt.Println("expect name of test as the first argument, expected 'CPUUsageIncrease' or 'ResponseTimeSimpleIncrease' or 'PredefinedSearch' subcommands")
		os.Exit(1)
	}

	var c strategies.Configurer
	switch os.Args[beforeConfigArgCount-1] {
	case "CPUUsageIncrease":
		c = GetCPUIncreaseConfigurer()
	case "ResponseTimeSimpleIncrease":
		c = GetResponseTimeSimpleIncreaseConfigurer()
	case "CPUUtil_RT_Hybrid":
		c = GetCPUUtilRTHybridConfigurer()
	case "PredefinedSearch":
		c = strategies.GetNewPredefinedSearcher()
	case "MOBO":
		c = GetMOBOConfigurer()
	case "Single":
		c = &strategies.SingleRun{}
	default:
		log.Println("expected 'Single' or 'CPUUsageIncrease' or 'ResponseTimeSimpleIncrease' or 'PredefinedSearch' subcommands but got", os.Args[beforeConfigArgCount-1])
		os.Exit(1)
	}
	// var c = GetAnotherConfigurer()
	var m = GetSwarmManager()
	db := GetNewDatabase()
	a := autoconfigure.NewAutoConfigurer(lg, rtc, rcc, ruc, c, m, workloadStr, db)
	log.Println("name of the test is:", os.Args[beforeConfigArgCount-2])
	a.Start(os.Args[beforeConfigArgCount-2])
}
