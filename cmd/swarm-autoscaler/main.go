package main

// "fmt"
// "sync"
// "time"

// "github.com/VahidMostofi/swarmmanager/internal/resource/collector"

// "github.com/VahidMostofi/swarmmanager/internal/workload"

// "github.com/montanaflynn/stats"

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

}

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
