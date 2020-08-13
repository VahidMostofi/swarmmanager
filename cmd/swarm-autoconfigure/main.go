package main

//TODO I NEED TO REstart every single container
import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

// GetPerPathEUBasedScaling ...
func GetPerPathEUBasedScaling(workloadStr string) strategies.Configurer {
	workload, err := loadgenerator.WorkloadFromString(workloadStr)
	if err != nil {
		log.Panic("error while parsing workload", workloadStr, err)
	}
	ppeusCmd := flag.NewFlagSet("PerPathEUBasedScaling", flag.ExitOnError)
	ppeusValueName := ppeusCmd.String("property", "", "Which property of a run to consider? CPUUsageMean,CPUUsage90Percentile 70-95, 99")
	ppeusThreshold := ppeusCmd.Float64("value", 0, "what is the threshold")
	ppeusStepSize := ppeusCmd.Float64("step", -1, "how much core to add at each step to each path")
	ppeusContainerStrategy := ppeusCmd.Bool("mc", false, "run it with multiple containers or not")
	ppeusCmd.Parse(os.Args[beforeConfigArgCount:])
	if *ppeusStepSize < 0 {
		log.Panic("invalid value for stepSize")
		os.Exit(1)
	}

	c := &strategies.PerPathEUBasedScaling{ //todo this is hardcoded!
		Path2Service2EUtilization: map[string]map[string]float64{
			"auth": map[string]float64{
				"gateway": workload.Throughput * 8 * workload.PathProportion["auth"] / 1000.0,
				"auth":    workload.Throughput * 74 * workload.PathProportion["auth"] / 1000.0,
			},
			"books": map[string]float64{
				"gateway": workload.Throughput * 42 * workload.PathProportion["books"] / 1000.0,
				"books":   workload.Throughput * 62 * workload.PathProportion["books"] / 1000.0,
			},
		},
		NormalizedPath2Service2EUtilization: map[string]map[string]float64{
			"auth": map[string]float64{
				"gateway": 0.10,
				"auth":    0.90,
			},
			"books": map[string]float64{
				"gateway": 0.40,
				"books":   0.60,
			},
		},
		MultiContainer: *ppeusContainerStrategy,
		StepSize:       *ppeusStepSize,
		Agreements:     []strategies.Agreement{strategies.Agreement{PropertyToConsider: *ppeusValueName, Value: *ppeusThreshold}},
	}

	err = c.Init()
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	return c
}

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
func GetJaegerCollector() *jaeger.Aggregator {
	j := jaeger.NewAggregator()
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

// GetAddDifferentFractionalCPUcores ...
func GetAddDifferentFractionalCPUcores(workload string) strategies.Configurer {
	adfccCmd := flag.NewFlagSet("AddDifferentFractionalCPUcores", flag.ExitOnError)
	adfccValueName := adfccCmd.String("property", "", "Which property of a run to consider? CPUUsageMean,CPUUsage90Percentile 70-95, 99")
	adfccThreshold := adfccCmd.Float64("value", 0, "what is the threshold")
	adfccAmount := adfccCmd.Float64("amount", -1, "how much core to add at each step")
	adfccIndicator := adfccCmd.String("indicator", "", "what is the indicator? Demand/Utilization")
	adfccContainerStrategy := adfccCmd.Bool("mc", false, "run it with multiple containers")

	adfccCmd.Parse(os.Args[beforeConfigArgCount:])
	adfccCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[beforeConfigArgCount-1])
		adfccCmd.PrintDefaults()
	}

	if *adfccValueName == "" {
		adfccCmd.Usage()
		os.Exit(1)
	}

	if *adfccThreshold == 0 {
		adfccCmd.Usage()
		os.Exit(1)
	}

	if *adfccAmount < 0 {
		adfccCmd.Usage()
		os.Exit(1)
	}

	if *adfccIndicator == "" {
		adfccCmd.Usage()
		os.Exit(1)
	}

	values, maxIncrease, err := strategies.GetFractionalCPUIncreaseValues(workload, *adfccIndicator, *adfccAmount)
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
	log.Println("values for Fractional Increase:", values)
	log.Println("values for Max Increase:", maxIncrease)

	log.Println("Configuring AddFractionalCPUcores with Value:", *adfccThreshold, "and property of", *adfccValueName, " and core amount of", *adfccAmount, "and indicator=", *adfccIndicator)
	return &strategies.AddDifferentFractionalCPUcores{
		ServiceToAmount:   values,
		MaxServiceIncease: maxIncrease,
		Agreements: []strategies.Agreement{
			{
				PropertyToConsider: *adfccValueName,
				Value:              *adfccThreshold,
			},
		},
		MultiContainer: *adfccContainerStrategy,
	}
}

// GetAddFractionalCPUcoresConfigurer ...
func GetAddFractionalCPUcoresConfigurer() strategies.Configurer {
	afccCmd := flag.NewFlagSet("AddFractionalCPUcoresConfigurer", flag.ExitOnError)
	afccValueName := afccCmd.String("property", "", "Which property of a run to consider? CPUUsageMean,CPUUsage90Percentile 70-95, 99")
	afccThreshold := afccCmd.Float64("value", 0, "what is the threshold")
	afccAmount := afccCmd.Float64("amount", -1, "how much core to add at each step")

	afccCmd.Parse(os.Args[beforeConfigArgCount:])
	afccCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[beforeConfigArgCount-1])
		afccCmd.PrintDefaults()
	}

	if *afccValueName == "" {
		afccCmd.Usage()
		os.Exit(1)
	}

	if *afccThreshold == 0 {
		afccCmd.Usage()
		os.Exit(1)
	}

	if *afccAmount < 0 {
		afccCmd.Usage()
		os.Exit(1)
	}

	log.Println("Configuring AddFractionalCPUcores with Value:", *afccThreshold, "and property of", *afccValueName, " and core amount of", *afccAmount)
	return &strategies.AddFractionalCPUcores{
		EachStepIncrease: *afccAmount,
		Agreements: []strategies.Agreement{
			{
				PropertyToConsider: *afccValueName,
				Value:              *afccThreshold,
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

	if strings.Contains(swarmmanager.GetConfig().ResultsDirectoryPath, "$STRATEGY") {
		swarmmanager.GetConfig().ResultsDirectoryPath = strings.Replace(swarmmanager.GetConfig().ResultsDirectoryPath, "$STRATEGY", os.Args[beforeConfigArgCount-1], 1)
		log.Println("Updating result path to", swarmmanager.GetConfig().ResultsDirectoryPath)

	}

	if strings.Contains(swarmmanager.GetConfig().ResultsDirectoryPath, "$SYSTEM_NAME") {
		swarmmanager.GetConfig().ResultsDirectoryPath = strings.Replace(swarmmanager.GetConfig().ResultsDirectoryPath, "$SYSTEM_NAME", swarmmanager.GetConfig().SystemName, 1)
		log.Println("Updating result path to", swarmmanager.GetConfig().ResultsDirectoryPath)
	}

	// creating directories for ResultDirectoryPath
	if err := os.MkdirAll(filepath.Dir(swarmmanager.GetConfig().ResultsDirectoryPath), 0770); err != nil {
		log.Panic(err)
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
	case "AddFractionalCPUcores":
		c = GetAddFractionalCPUcoresConfigurer()
	case "Single":
		c = &strategies.SingleRun{}
	case "AddDifferentFractionalCPUcores":
		c = GetAddDifferentFractionalCPUcores(workloadStr)
	case "PerPathEUBasedScaling":
		c = GetPerPathEUBasedScaling(workloadStr)
	default:
		log.Println("expected 'Single' or 'CPUUsageIncrease' or 'ResponseTimeSimpleIncrease' or 'PredefinedSearch' subcommands but got", os.Args[beforeConfigArgCount-1])
		os.Exit(1)
	}
	// var c = GetAnotherConfigurer()
	var m = GetSwarmManager()
	db := GetNewDatabase()
	a := autoconfigure.NewAutoConfigurer(lg, rtc, rcc, ruc, c, m, workloadStr, db)
	log.Println("name of the test is:", os.Args[beforeConfigArgCount-2])
	a.Start(os.Args[beforeConfigArgCount-2], strings.Join(os.Args, " "))
}
