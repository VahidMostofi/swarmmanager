package main

//TODO I NEED TO REstart every single container
import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"flag"

	"github.com/VahidMostofi/swarmmanager/internal/autoconfigure"
	"github.com/VahidMostofi/swarmmanager/internal/jaeger"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	resource "github.com/VahidMostofi/swarmmanager/internal/resource/collector"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/VahidMostofi/swarmmanager/internal/workload"
)

const beforeConfigArgCount = 4

// GetTheResourceUsageCollector ...
func GetTheResourceUsageCollector() resource.Collector {
	//TODO SERVICE COUNT IS HARDCODED!!!!!!!!
	stackName := "bookstore"
	c := resource.GetNewCollector("SingleCollector")
	err := c.Configure(map[string]string{"host": "tcp://136.159.209.204:2375", "stackname": stackName})
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
	authProb, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		panic(err)
	}
	bookProb := 1 - authProb
	sleepDuration, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		panic(err)
	}
	script := loadgenerator.CreateLoadGeneartorScript("/Users/vahid/Desktop/type5.js", vus, duration, authProb, bookProb, 0, sleepDuration)
	l.Prepare(map[string]string{"script": script})
	return l

	// StartLG(l)
	// StopLG(l)
	// FeedbackLG(l)
}

// GetJaegerCollector ...
func GetJaegerCollector() *jaeger.JaegerAggregator {
	//TODO these are hardcoded!
	j := jaeger.NewJaegerAggregator("http://136.159.209.204:16686", []string{"auth_req_login", "books_edit_book", "books_get_book"})
	return j
}

// GetCPUIncreaseConfigurer ...
func GetCPUIncreaseConfigurer() autoconfigure.Configurer {
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
	return &autoconfigure.CPUUsageIncrease{
		Threshold:       *cpuOnlyThreshold,
		ValueToConsider: *cpuOnlyValueName,
	}
}

// GetResponseTimeSimpleIncreaseConfigurer ...
func GetResponseTimeSimpleIncreaseConfigurer() autoconfigure.Configurer {
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
	return &autoconfigure.ResponseTimeSimpleIncrease{
		Agreements: []autoconfigure.Agreement{
			{
				PropertyToConsider: *rtsiValueName,
				Value:              *rtsiThreshold,
			},
		},
	}
}

// GetSwarmManager ...
func GetSwarmManager() *swarm.Manager {
	m, err := swarm.GetNewSwarmManager(map[string]string{"stackname": "bookstore", "host": "tcp://136.159.209.204:2375", "services": "auth,gateway,books"})
	if err != nil {
		log.Panic(err)
	}
	return m
}

func main() {

	var ruc = GetTheResourceUsageCollector()
	var rtc workload.ResponseTimeCollector = GetJaegerCollector()
	var rcc workload.RequestCountCollector = rtc.(workload.RequestCountCollector)
	workloadStr := os.Args[1]
	if !strings.Contains(workloadStr, "_") {
		log.Panic("the first argument must be the workload")
	}
	var lg = GetTheLoadGenerator(os.Args[1])

	if len(os.Args) < beforeConfigArgCount {
		fmt.Println("expect name of test as the first argument, expected 'CPUUsageIncrease' or 'ResponseTimeSimpleIncrease' or 'PredefinedSearch' subcommands")
		os.Exit(1)
	}

	var c autoconfigure.Configurer
	switch os.Args[beforeConfigArgCount-1] {
	case "CPUUsageIncrease":
		c = GetCPUIncreaseConfigurer()
	case "ResponseTimeSimpleIncrease":
		c = GetResponseTimeSimpleIncreaseConfigurer()
	case "PredefinedSearch":
		c = autoconfigure.GetNewPredefinedSearcher()
	case "Single":
		c = &autoconfigure.SingleRun{}
	default:
		log.Println("expected 'Single' or 'CPUUsageIncrease' or 'ResponseTimeSimpleIncrease' or 'PredefinedSearch' subcommands")
		os.Exit(1)
	}
	// var c = GetAnotherConfigurer()
	var m = GetSwarmManager()
	a := autoconfigure.NewAutoConfigurer(lg, rtc, rcc, ruc, c, m, workloadStr)
	log.Println("name of the test is:", os.Args[beforeConfigArgCount-2])
	a.Start(os.Args[beforeConfigArgCount-2])
}
