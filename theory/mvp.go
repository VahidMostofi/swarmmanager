package theory

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
)

type System struct {
	Classes           []string
	Resources         []string
	ClassProbs        map[string]float64 `json:"class_probs"`
	Throughput        float64
	Demands           map[string]float64
	BestObjective     float64            `json:"best_objective"`
	BestAlphas        map[string]float64 `json:"best"`
	BestResponseTimes []float64          `json:"best_mrt"`
	SLA               float64            `json:"SLA"`
	Name              string
}

// ReadSystem ...
func ReadSystem(name string) *System {
	b, err := ioutil.ReadFile("./theory/systems/" + name + ".json")
	if err != nil {
		panic(err)
	}
	s := &System{}
	err = json.Unmarshal(b, s)
	if err != nil {
		panic(err)
	}

	rtss := s.GetMeanResponseTimes(s.BestAlphas)

	// check
	for i := range s.Classes {
		if math.Abs(s.BestResponseTimes[i]-rtss[i]) > 1e-5 {
			fmt.Println(rtss[i], s.BestResponseTimes[i])
			panic("These values must be equal!")
		}
	}

	s.Name = "system_" + name
	return s
}

// GetMeanResponseTimes ...
func (s *System) GetMeanResponseTimes(alphas map[string]float64) []float64 {
	responseTimes := make([]float64, len(s.Classes))
	for cIdx, cName := range s.Classes {
		var r float64
		for _, kName := range s.Resources {
			r += (s.Demands[cName+"_"+kName] / alphas[kName]) / (1 - s.GetUtilizations(alphas, kName))
		}
		responseTimes[cIdx] = r * 1000
	}
	return responseTimes
}

// GetUtilizations ...
func (s *System) GetUtilizations(alphas map[string]float64, resource string) float64 {
	var u float64 = 0
	for _, cName := range s.Classes {
		u += s.Throughput * s.ClassProbs[cName] * (s.Demands[cName+"_"+resource] / alphas[resource])
	}
	if u > 0.99999 {
		uStr := strconv.FormatFloat(u, 'f', 3, 64)
		panic("Utilization for " + resource + " is " + uStr)
	}
	return u
}

// MeetSLA ...
func (s *System) MeetSLA(rts []float64) bool {
	for _, r := range rts {
		if r > s.SLA {
			return false
		}
	}
	return true
}
