package main

// "fmt"
// "sync"
// "time"

// "github.com/VahidMostofi/swarmmanager/internal/resource/collector"

// "github.com/VahidMostofi/swarmmanager/internal/workload"

// "github.com/montanaflynn/stats"

import (
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
)

func main() {
	// stackName := "bookstore"
	// c := collector.GetNewCollector("SingleCollector")
	// err := c.Configure(map[string]string{"host": "tcp://136.159.209.204:2375", "stackname": stackName})
	// if err != nil {
	// 	panic(err)
	// }
	// err = c.Start()
	// if err != nil {
	// 	panic(err)
	// }

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
	// j := jaeger.NewJaegerAggregator("http://136.159.209.204:16686", []string{"auth_req_login", "books_edit_book", "books_get_book"})
	// fmt.Println(j)
	// j.GetTraces(1593568965858000, 1593568975858000, "gateway")
	// time.Sleep(20 * time.Second)
	// start := time.Now().UnixNano() / 1e3
	// start := int64(1593737202594662)
	// time.Sleep(20 * time.Second)
	// end := time.Now().UnixNano() / 1e3
	// end := int64(1593737212597534)

	// ru := c.GetResourceUtilization()
	// for key, value := range ru {
	// 	count := 0
	// 	for k := range value.CPUUtilizationsAtTime {
	// 		if k/1e3 >= start && k/1e3 <= end {
	// 			count++
	// 		}
	// 	}
	// 	fmt.Println(key[:12], value.GetResourceRecordingRate(), len(value.CPUUtilizationsAtTime), count)

	// 	// value.Print()
	// 	// fmt.Println("======")
	// }

	// j.GetTraces(start, end, "gateway")
	// fmt.Println(start, end)
	// count, err := j.GetRequestCount("auth_req_login")
	// fmt.Println("request count login", count)
	// count, err = j.GetRequestCount("books_edit_book")
	// fmt.Println("request count update book", count)
	// count, err = j.GetRequestCount("books_get_book")
	// fmt.Println("request count get book", count)

	// rt, err := j.GetResponseTimes("auth_req_login")
	// m, err := stats.Mean(rt)
	// fmt.Println("mean response time login", m)
	// rt, err = j.GetResponseTimes("books_edit_book")
	// m, err = stats.Mean(rt)
	// fmt.Println("mean response time update book", m)
	// rt, err = j.GetResponseTimes("books_get_book")
	// m, err = stats.Mean(rt)
	// fmt.Println("mean response time get book", m)
	// if err != nil {
	// 	panic(err)
	// }

	// F2(j)
	// F1(j)
	// script := loadgenerator.CreateLoadGeneartorScript("/Users/vahid/Desktop/type5.js", 15, 3600, 0.2, 0.8, 0, 0.1)
	// l := loadgenerator.NewK6LoadGenerator("http://136.159.209.214:7112")

	// PrepareLG(l, script)
	// StartLG(l)
	// StopLG(l)
	// FeedbackLG(l)
	// ========================================================================
	// m, err := swarm.GetNewSwarmManager(map[string]string{"stackname": "bookstore", "host": "tcp://136.159.209.204:2375"})
	// m.StackStateCh <- swarm.StackStateServicesAreReady
	// m.FillDesiredSpecsCurrentSpecs()
	// if err != nil {
	// 	panic(err)
	// }
	// time.Sleep(5 * time.Second)
	// for serviceID, spec := range m.DesiredSpecs {
	// 	if spec.Name == "gateway" {
	// 		spec.ReplicaCount = 3
	// 		spec.CPULimits = 2
	// 		spec.CPUReservation = 2
	// 		spec.EnvironmentVariables = []string{"JWT_KEY=someKeyIsGoodAndSomeOfThemBNoGEo1ioD!", "WorkerCount=2"}
	// 		m.DesiredSpecs[serviceID] = spec
	// 	}
	// }
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

	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// wg.Wait()

	// c := autoconfigure.CPUUsageIncrease{}
	// c.Configure()

}

func PrepareLG(l loadgenerator.LoadGenerator, script string) {
	err := l.Prepare(map[string]string{"script": script})
	if err != nil {
		panic(err)
	}
}

func StartLG(l loadgenerator.LoadGenerator) {
	err := l.Start(map[string]string{})
	if err != nil {
		panic(err)
	}
}

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
