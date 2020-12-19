package jaeger

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/montanaflynn/stats"
)

var SCRIPT_TEMPLATE = `
import pandas as pd
import numpy as np
import scipy.stats as stats


def medianCI(data, ci, p):
	if type(data) is pd.Series or type(data) is pd.DataFrame:
		#transfer data into np.array
		data = data.values

	#flat to one dimension array
	data = data.reshape(-1)
	data = np.sort(data)
	N = data.shape[0]
	
	lowCount, upCount = stats.binom.interval(ci, N, p, loc=0)
	#given this: https://onlinecourses.science.psu.edu/stat414/node/316
	#lowCount and upCount both refers to  W's value, W follows binomial Dis.
	#lowCount need to change to lowCount-1, upCount no need to change in python indexing
	lowCount -= 1
	# print lowCount, upCount
	return data[int(lowCount)], data[int(upCount)]

if __name__ == '__main__':
	data = np.array([DATA_GOES_HERE])
	data = pd.Series(data)
	print(medianCI(data, 0.1, 0.95))
`

func TestOne(t *testing.T) {
	pythonPath := "/home/vahid/.virtualenvs/with-data/bin/python"
	log.SetOutput(ioutil.Discard)

	// appName := "bookstore_nodejs"
	// for _, w := range []int{75, 100, 125, 150} {
	// 	for _, lastPart := range []string{"_110_0.33_0.33_0.34/bnv2/bnv2-250-2.0-mc-c-1.0.yml",
	// 		"_110_0.33_0.33_0.34/bnv2/bnv2-250-1.0-mc-c-1.0.yml",
	// 		"_110_0.33_0.33_0.34/bnv1/bnv1-250-2.0-mc-c-1.0.yml",
	// 		"_110_0.33_0.33_0.34/bnv1/bnv1-250-1.0-mc-c-1.0.yml",
	// 		"_110_0.33_0.33_0.34/mobo/mobo.yml"} {
	appName := "muck_star-small"
	for _, w := range []int{10, 20, 40} {
		for _, lastPart := range []string{"_250_0.45_0.25_0.15_0.15/bnv2/bnv2-300-2.0-mc-c-0.5.yml",
			"_250_0.45_0.25_0.15_0.15/bnv2/bnv2-300-1.0-mc-c-0.5.yml",
			"_250_0.45_0.25_0.15_0.15/bnv1/bnv1-300-2.0-mc-c-0.5.yml",
			"_250_0.45_0.25_0.15_0.15/bnv1/bnv1-300-1.0-mc-c-0.5.yml",
			// "_110_0.45_0.25_0.15_0.15/mobo/mobo.yml"
		} {
			var formulasPath string
			if appName == "muck_star-small" {
				formulasPath = "/home/vahid/Desktop/projects/swarmmanager/configurations/formulas/muck_general.yaml"
			} else if appName == "bookstore_nodejs" {
				formulasPath = "/home/vahid/Desktop/projects/swarmmanager/configurations/formulas/" + appName + ".yaml"
			}

			wStr := strconv.Itoa(w)
			// lastPart := "_110_0.33_0.33_0.34/bnv2/bnv2-250-2.0-mc-c-1.0.yml"
			inputFilePath := "/home/vahid/Dropbox/data/swarm-manager-data/results/" + appName + "/" + wStr + lastPart

			b, err := ioutil.ReadFile(inputFilePath)
			if err != nil {
				panic(err)
			}
			stackHistory := &history.ExecutionDetails{}
			err = yaml.Unmarshal(b, stackHistory)
			if err != nil {
				panic(err)
			}
			cache := make(map[string]bool)

			configs.FakeInitialize()
			configs.GetConfig().Jaeger.DetailsFilePath = formulasPath
			// fmt.Println(inputFilePath, len(stackHistory.History))
			// continue
			for i := 0; i < len(stackHistory.History); i++ {
				if _, ok := cache[stackHistory.History[i].JaegerFile]; ok {
					continue
				}
				cache[stackHistory.History[i].JaegerFile] = true
				r, err := zip.OpenReader(stackHistory.History[i].JaegerFile)
				if err != nil {
					r, err = zip.OpenReader(strings.Replace(stackHistory.History[i].JaegerFile, "/home/vahid/Dropbox/data/swarm-manager-data/jaegers/", "/run/media/vahid/DISK_VAHID/other_hard/jaegers/", 1))
					// r, err := zip.OpenReader(strings.R)
					if err != nil {
						fmt.Println("This is bad!")
						panic("ERR!")
					}
				}
				defer r.Close()
				data := &struct {
					Data []*trace `json:"data"`
				}{}
				if len(r.File) > 1 {
					t.Error(fmt.Errorf("more than 1 file in zip file."))
				}
				for _, f := range r.File {
					r, err := f.Open()
					if err != nil {
						t.Error(err)
					}
					b, err := ioutil.ReadAll(r)
					if err != nil {
						t.Error(err)
					}
					json.Unmarshal(b, data)
					break
				}
				j := NewAggregator()
				b, err = yaml.Marshal(j)
				// ioutil.WriteFile("/home/vahid/Desktop/v.yaml", b, 0777)
				// fmt.Println(len(data.Data))
				if len(data.Data) < 100 {
					fmt.Println("len(data.Data)", len(data.Data))
					panic("")
				}
				err = j.ParseTraces(data.Data)
				if err != nil {
					t.Errorf("error evaluateing formula: %w", err)
				}

				for r := range j.ValueFormulas.Requests {
					v, _ := j.GetRequestResponseTimes(r)
					p95, err := stats.Percentile(v, 95)
					if err != nil {
						panic(err)
					}
					str := ""
					for i, f := range v {
						str += strconv.FormatFloat(f, 'f', 3, 64)
						if i != len(v)-1 {
							str += ", "
						}
					}
					script := strings.Replace(SCRIPT_TEMPLATE, "DATA_GOES_HERE", str, 1)
					err = ioutil.WriteFile("main.py", []byte(script), 0644)
					if err != nil {
						fmt.Println("error writing file")
						panic(err)
					}
					cmd := exec.Command(pythonPath, "main.py")
					out, err := cmd.Output()
					if err != nil {
						fmt.Println(out)
						panic(err)
					}
					res := string(out)
					res = strings.Replace(res, "(", "", 1)
					res = strings.Replace(res, ")", "", 1)
					res = strings.Replace(res, "\n", "", 1)

					fmt.Println(appName+",", res, ",", p95, ",", len(v))
				}
			}
		}
	}
}
