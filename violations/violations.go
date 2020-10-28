package violations

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/jaeger"
	"github.com/montanaflynn/stats"
	"gopkg.in/yaml.v3"

	"github.com/cheggaaa/pb/v3"
)

// Run ...
func Run() {
	testName := "single_" + os.Args[2]
	appName := "bookstore_nodejs"
	workload := os.Args[3]
	resultsDirectory := "/home/vahid/Dropbox/data/swarm-manager-data/results/"
	var sla float64 = 250

	resultFile := resultsDirectory + appName + "/" + workload + "/single/" + testName + ".yml"

	b, err := ioutil.ReadFile(resultFile)
	if err != nil {
		panic(err)
	}

	info := &history.ExecutionDetails{}
	err = yaml.Unmarshal(b, info)
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
	fmt.Println(info.History[0].JaegerFile)

	configs.FakeInitialize()
	configs.GetConfig().Jaeger.DetailsFilePath = "/home/vahid/Desktop/projects/swarmmanager/configurations/formulas/bookstore_nodejs.yaml"

	j := jaeger.NewAggregator()

	startTime := info.History[0].ServicesInfo["entry"].Start
	endTime := info.History[0].ServicesInfo["entry"].End

	var duration int64 = 60 * 1e6
	meetsStats := make([]bool, 0)
	falseCount := 0
	var i int64
	log.SetOutput(ioutil.Discard)

	responseTimes := make(map[string][]float64)
	counts := make(map[string][]int)
	for key := range info.History[0].RequestResponseTimes {
		responseTimes[key] = make([]float64, 0)
		counts[key] = make([]int, 0)
	}

	bar := pb.StartNew(int((endTime - startTime) / duration))
	doLog := true
	for i = 0; i < (endTime-startTime)/duration; i++ {
		fmt.Println(i)
		s := float64(startTime + i*duration)
		e := float64(startTime + (i+1)*duration)

		j.ParseTraceFile(info.History[0].JaegerFile, s, e)

		meets := true
		for reqName := range info.History[0].RequestResponseTimes {
			rs, err := j.GetRequestResponseTimes(reqName)
			if err != nil {
				panic(err)
			}
			p95, err := stats.Percentile(rs, 95)
			if doLog {
				fmt.Print(p95, ", ")
			}
			meets = meets && p95 <= sla
			if !(p95 <= sla) {
				fmt.Print(reqName + ", ")
			}
			c, err := j.GetRequestCount(reqName)
			if err != nil {
				panic(err)
			}
			responseTimes[reqName] = append(responseTimes[reqName], p95)
			counts[reqName] = append(counts[reqName], c)
		}
		if doLog {
			fmt.Println(meets)
		}
		meetsStats = append(meetsStats, meets)
		if !meets {
			falseCount++
		}
		bar.Increment()
	}
	fmt.Println("didn't meet:", falseCount, "out of", len(meetsStats))
	for req, rts := range responseTimes {
		fmt.Print(req, " ")
		for _, v := range rts {
			fmt.Print(v, ",")
		}
		fmt.Println()
	}

	for req, cs := range counts {
		fmt.Print(req, " ")
		for _, v := range cs {
			fmt.Print(v, ",")
		}
		fmt.Println()
	}
}
