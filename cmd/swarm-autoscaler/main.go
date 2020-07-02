package main

// "fmt"
// "sync"
// "time"

// "github.com/VahidMostofi/swarmmanager/internal/resource/collector"

// "github.com/VahidMostofi/swarmmanager/internal/workload"

// "github.com/montanaflynn/stats"

import (
	"sync"
	"time"

	"github.com/VahidMostofi/swarmmanager/internal/swarm"
)

func main() {
	// c := collector.GetNewCollector("SingleCollector")
	// err := c.Configure(map[string]string{"host": "tcp://136.159.209.204:2375", "stackname": "bookstore"})
	// if err != nil {
	// 	panic(err)
	// }
	// err = c.Start()
	// if err != nil {
	// 	panic(err)
	// }

	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// time.AfterFunc(time.Second*5, func() { c.Stop(); wg.Done() })
	// wg.Wait()
	// time.Sleep(3 * time.Second)
	// ru := c.GetResourceUtilization()
	// for _, value := range ru {
	// 	fmt.Println(value.GetResourceRecordingRate())
	// 	value.Print()
	// 	fmt.Println("======")
	// }

	// "getone", "update","login"
	// j := jaeger.NewJaegerAggregator("http://136.159.209.204:16686", []string{"auth_req_login", "update_book", "get_book"})
	// j.GetTraces(1593568965858000, 1593568975858000, "gateway")

	// F2(j)
	// F1(j)
	// script := "import{check,sleep}from'k6';import{execute_random_login,execute_get_book,execute_edit_book}from'./bookstore_content/bookstore_units.js';export let options={vus:3,duration:'5s',userAgent:'MyK6UserAgentString/1.0',};const SLEEP_DURATION=0.1;export function setup(){}\nexport default function(data){const auth_token=execute_random_login();sleep(SLEEP_DURATION);const book=execute_get_book(auth_token);sleep(SLEEP_DURATION);execute_edit_book(auth_token,book);};export function teardown(data){}"
	// l := loadgenerator.NewK6LoadGenerator("http://136.159.209.214:7112")

	// PrepareLG(l, script)
	// StartLG(l)
	// StopLG(l)
	// FeedbackLG(l)
	// ========================================================================
	m, err := swarm.GetNewSwarmManager(map[string]string{"stackname": "bookstore", "host": "tcp://136.159.209.204:2375"})
	m.StackStateCh <- swarm.StackStateServicesAreReady
	m.FillDesiredSpecsCurrentSpecs()
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
	for serviceID, spec := range m.DesiredSpecs {
		if spec.Name == "gateway" {
			spec.ReplicaCount = 3
			spec.CPULimits = 2
			spec.CPUReservation = 2
			spec.EnvironmentVariables = []string{"JWT_KEY=someKeyIsGoodAndSomeOfThemBNoGEo1ioD!", "WorkerCount=2"}
			m.DesiredSpecs[serviceID] = spec
		}
	}
	// go m.Temp()

	// err = m.RemoveStack("tcp://136.159.209.204:2375", "bookstore")
	// if err != nil {
	// 	panic(err)
	// }

	// dockerComposePath := "/Users/vahid/workspace/bookstore/docker-compose.yml"
	// err = m.DeployStackWithDockerCompose(dockerComposePath, "tcp://136.159.209.204:2375", "bookstore", 1)
	// if err != nil {
	// 	panic(err)
	// }

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}

// func PrepareLG(l loadgenerator.LoadGenerator, script string) {
// 	err := l.Prepare(map[string]string{"script": script})
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func StartLG(l loadgenerator.LoadGenerator) {
// 	err := l.Start(map[string]string{})
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func StopLG(l loadgenerator.LoadGenerator) {
// 	err := l.Stop(map[string]string{})
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func FeedbackLG(l loadgenerator.LoadGenerator) {
// 	f, err := l.GetFeedback(map[string]string{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	for k, _ := range f {
// 		fmt.Println(k)
// 	}
// }

// func F1(rtc workload.ResponseTimeCollector) {
// 	rt, err := rtc.GetResponseTimes("auth_req_login")
// 	fmt.Println(stats.Mean(rt))
// 	rt, err = rtc.GetResponseTimes("update_book")
// 	fmt.Println(stats.Mean(rt))
// 	rt, err = rtc.GetResponseTimes("get_book")
// 	fmt.Println(stats.Mean(rt))
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func F2(rcc workload.RequestCountCollector) {
// 	fmt.Println(rcc.GetRequestCount("auth_req_login"))
// 	fmt.Println(rcc.GetRequestCount("update_book"))
// 	fmt.Println(rcc.GetRequestCount("get_book"))
// }
