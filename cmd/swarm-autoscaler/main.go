package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/VahidMostofi/swarmmanager/internal/resource/collector"
)

func main() {
	c := collector.GetNewCollector("SingleCollector")
	err := c.Configure(map[string]string{"host": "tcp://136.159.209.204:2375", "stackname": "bookstore"})
	if err != nil {
		panic(err)
	}
	err = c.Start()
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	time.AfterFunc(time.Second*5, func() { c.Stop(); wg.Done() })
	wg.Wait()
	time.Sleep(3 * time.Second)
	ru := c.GetResourceUtilization()
	for _, value := range ru {
		fmt.Println(value.GetResourceRecordingRate())
		value.Print()
		fmt.Println("======")
	}
}
