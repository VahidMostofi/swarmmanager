package caching

import (
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func round1(value float64) float64 {
	return math.Round(value*10) / 10
}

func getPartialWorkload(workload string) (string, int) {
	parts := strings.Split(workload, "_")
	w := ""
	for i, p := range parts {
		if i > 0 {
			w += p
			if i < len(parts) {
				w += "_"
			}
		}
	}
	l, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	return w, l
}

// FindBaseConfiguration ...
func FindBaseConfiguration(sla float64, currentLIndex int) *history.Information {
	cachePath := configs.GetConfig().Cache.Details["path"]

	files, err := ioutil.ReadDir(cachePath)
	if err != nil {
		panic(err)
	}

	theFiles := make([]string, 0)
	partialWorkload, l := getPartialWorkload(viper.GetString("workloadStr"))
	largestL := l

	for _, f := range files {
		if strings.Contains(f.Name(), configs.GetConfig().AppName) {
			if strings.Contains(f.Name(), configs.GetConfig().Version) {
				if strings.Contains(f.Name(), partialWorkload) {
					currentL, err := strconv.Atoi(strings.Split(f.Name(), "_")[currentLIndex])

					if err != nil {
						panic(err)
					}

					if currentL < largestL {
						largestL = currentL
					}
					if strings.Contains(f.Name(), "_"+strconv.Itoa(configs.GetConfig().Test.Duration)+"_") {
						theFiles = append(theFiles, f.Name())
					}
				}
			}
		}
	}

	theWorkload := "_" + strconv.Itoa(largestL) + "_" + partialWorkload

	var minCPU float64 = 10000000
	var bestFile = ""
	for _, f := range files {
		meets := true
		if strings.Contains(f.Name(), theWorkload) {
			b, err := ioutil.ReadFile(configs.GetConfig().Cache.Details["path"] + f.Name() + "/" + "info.yml")
			if err != nil {
				panic(err)
			}
			info := &history.Information{}

			yaml.Unmarshal(b, info)

			for _, rts := range info.RequestResponseTimes {
				if *rts.RTToleranceIntervalUBoundConfidence90p95 > sla {
					meets = false
					break
				}
			}
			if meets {
				var cpuCount float64 = 0
				for _, sr := range info.Specs {
					cpuCount += round1(sr.CPULimits * float64(sr.ReplicaCount))
				}
				// for _, rts := range info.RequestResponseTimes {
				// 	fmt.Println(*rts.RTToleranceIntervalUBoundConfidence90p95)
				// }
				// fmt.Println(cpuCount)
				// fmt.Println("==========")
				if cpuCount < minCPU {
					minCPU = cpuCount
					bestFile = f.Name()
				}
			}
		}
	}

	b, err := ioutil.ReadFile(configs.GetConfig().Cache.Details["path"] + bestFile + "/" + "info.yml")
	if err != nil {
		panic(err)
	}
	info := &history.Information{}

	yaml.Unmarshal(b, info)
	return info
}
