package main

//TODO I NEED TO REstart every single container
import (
	"fmt"
	"log"
	"os"

	"flag"

	"github.com/VahidMostofi/swarmmanager/internal/autoconfigure"
	"github.com/VahidMostofi/swarmmanager/internal/jaeger"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	resource "github.com/VahidMostofi/swarmmanager/internal/resource/collector"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/VahidMostofi/swarmmanager/internal/workload"
)

const beforeConfigArgCount = 3

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
func GetTheLoadGenerator() loadgenerator.LoadGenerator {
	l := loadgenerator.NewK6LoadGenerator("http://136.159.209.214:7112")
	//TODO: what about the duration of generated load
	//TODO: this is hard coded
	script := loadgenerator.CreateLoadGeneartorScript("/Users/vahid/Desktop/type5.js", 20, 80, 0.2, 0.8, 0, 0.1)
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
	fmt.Println("Configuring CPUUsageIncreaseConfigurer with Threshold:", *cpuOnlyThreshold, "and property of", *cpuOnlyValueName)
	return &autoconfigure.CPUUsageIncrease{
		Threshold:       *cpuOnlyThreshold,
		ValueToConsider: *cpuOnlyValueName,
	}
}

// GetAnotherConfigurer ...
func GetAnotherConfigurer() autoconfigure.Configurer {
	return &autoconfigure.ResponseTimeIncrease{
		Agreements: []autoconfigure.Agreement{
			{
				PropertyToConsider: "ResponseTimes95Percentile",
				Value:              200,
			},
			{
				PropertyToConsider: "ResponseTimes99Percentile",
				Value:              400,
			},
		},
	}
}

// GetSwarmManager ...
func GetSwarmManager() *swarm.Manager {
	m, err := swarm.GetNewSwarmManager(map[string]string{"stackname": "bookstore", "host": "tcp://136.159.209.204:2375"})
	if err != nil {
		log.Panic(err)
	}
	return m
}

func main() {

	var ruc = GetTheResourceUsageCollector()
	var rtc workload.ResponseTimeCollector = GetJaegerCollector()
	var rcc workload.RequestCountCollector = rtc.(workload.RequestCountCollector)
	var lg = GetTheLoadGenerator()

	if len(os.Args) < beforeConfigArgCount {
		fmt.Println("expect name of test as the first argument, expected 'CPUUsageIncrease' or '' subcommands")
		os.Exit(1)
	}

	var c autoconfigure.Configurer
	switch os.Args[beforeConfigArgCount-1] {
	case "CPUUsageIncrease":
		c = GetCPUIncreaseConfigurer()
	default:
		log.Println("expected 'CPUUsageIncrease' or '' subcommands")
		os.Exit(1)
	}
	// var c = GetAnotherConfigurer()
	var m = GetSwarmManager()
	a := autoconfigure.NewAutoConfigurer(lg, rtc, rcc, ruc, c, m)
	fmt.Println("name of the test is:", os.Args[beforeConfigArgCount-2])
	a.Start(os.Args[beforeConfigArgCount-2])
}
