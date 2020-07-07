package main

//TODO I NEED TO REstart every single container
import (
	"log"

	"github.com/VahidMostofi/swarmmanager/internal/autoconfigure"
	"github.com/VahidMostofi/swarmmanager/internal/jaeger"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	resource "github.com/VahidMostofi/swarmmanager/internal/resource/collector"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/VahidMostofi/swarmmanager/internal/workload"
)

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
	script := loadgenerator.CreateLoadGeneartorScript("/Users/vahid/Desktop/type5.js", 14, 3600, 0.2, 0.8, 0, 0.1)
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

// GetAConfigurer ...
func GetAConfigurer() autoconfigure.Configurer {
	return &autoconfigure.CPUUsageIncrease{
		Threshold:       70,
		ValueToConsider: "CPUUsage90Percentile",
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
	var c = GetAConfigurer()
	// var c = GetAnotherConfigurer()
	var m = GetSwarmManager()
	a := autoconfigure.NewAutoConfigurer(lg, rtc, rcc, ruc, c, m)
	a.Start()
}
