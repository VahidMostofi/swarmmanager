package theory

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/montanaflynn/stats"
)

func GoTheory() {
	log.SetOutput(ioutil.Discard)
	line := "approach, system, steps, total core\n"
	for i := 1; i <= 5; i++ {
		fmt.Println(i)
		fileName := strconv.Itoa(i)
		// fmt.Println(i)
		system := ReadSystem(fileName)

		strategy := &strategies.BottleNeckOnlyVersion1{
			StepSize:          1.0,
			Agreements:        []strategies.Agreement{{"ResponseTimesMean", system.SLA}},
			MultiContainer:    true,
			DemandsFilePath:   "./theory/demands/" + fileName + ".yml",
			ConstantInit:      true,
			ConstantInitValue: 1.0,
		}
		strategy.Init()
		line += RunSystemWithStrategy("BNV1-1.0", system, strategy, false)

		strategy = &strategies.BottleNeckOnlyVersion1{
			StepSize:          0.5,
			Agreements:        []strategies.Agreement{{"ResponseTimesMean", system.SLA}},
			MultiContainer:    true,
			DemandsFilePath:   "./theory/demands/" + fileName + ".yml",
			ConstantInit:      true,
			ConstantInitValue: 1.0,
		}
		strategy.Init()
		line += RunSystemWithStrategy("BNV1-0.5", system, strategy, false)

		strategy = &strategies.BottleNeckOnlyVersion1{
			StepSize:          0.25,
			Agreements:        []strategies.Agreement{{"ResponseTimesMean", system.SLA}},
			MultiContainer:    true,
			DemandsFilePath:   "./theory/demands/" + fileName + ".yml",
			ConstantInit:      true,
			ConstantInitValue: 1.0,
		}
		strategy.Init()
		line += RunSystemWithStrategy("BNV1-0.25", system, strategy, false)

		strategy = &strategies.BottleNeckOnlyVersion1{
			StepSize:          0.1,
			Agreements:        []strategies.Agreement{{"ResponseTimesMean", system.SLA}},
			MultiContainer:    true,
			DemandsFilePath:   "./theory/demands/" + fileName + ".yml",
			ConstantInit:      true,
			ConstantInitValue: 1.0,
		}
		strategy.Init()
		line += RunSystemWithStrategy("BNV1-0.1", system, strategy, false)

		strategy2 := &strategies.BottleNeckOnlyVersion2{
			StepSize:          2.0,
			MinimumStepSize:   0.1,
			Agreements:        []strategies.Agreement{{"ResponseTimesMean", system.SLA}},
			MultiContainer:    true,
			DemandsFilePath:   "./theory/demands/" + fileName + ".yml",
			ConstantInit:      true,
			ConstantInitValue: 1.0,
		}
		strategy2.Init()
		strategy2.MinimumCPUValue = 1.0
		line += RunSystemWithStrategy("BNV2-1.0-0.1", system, strategy2, false)

		line += "BF-0.1, " + system.Name + ", 24300000, " + strconv.FormatFloat(system.BestObjective, 'f', 2, 64) + "\n"
	}
	fmt.Println(line)
}

// RunSystemWithStrategy ...
func RunSystemWithStrategy(name string, system *System, strategy strategies.Configurer, debug bool) string {

	t := theoryWorkload{Throughput: system.Throughput, ClassProbs: system.ClassProbs}
	currentConfig, err := strategy.GetInitialConfig(t)

	if err != nil {
		panic(err)
	}
	if len(currentConfig) == 0 {
		for _, service := range system.Resources {
			currentConfig[service] = swarm.SimpleSpecs{CPU: 1, Replica: 1, Worker: 1}
		}
	}

	iterations := make([]iterationInfo, 0)
	currentState := make(map[string]swarm.ServiceSpecs)
	stepCount := 1
	for {
		itr := iterationInfo{make(map[string]float64), make(map[string]float64), make(map[string]float64)}
		alphas := make(map[string]float64)
		for serviceName, simpleConfig := range currentConfig {
			currentState[serviceName] = swarm.ServiceSpecs{
				ReplicaCount: simpleConfig.Replica,
				CPULimits:    simpleConfig.CPU,
			}
			alphas[serviceName] = simpleConfig.CPU * float64(simpleConfig.Replica)
			itr.Alphas[serviceName] = alphas[serviceName]
		}

		info := history.Information{}
		info.RequestResponseTimes = make(map[string]history.ResponseTimeStats)
		meanResponseTimes := system.GetMeanResponseTimes(alphas)
		for requestIdx, requestName := range system.Classes {
			mrt := meanResponseTimes[requestIdx]
			if mrt <= 0 {
				fmt.Println(mrt, requestIdx, requestName)
				panic("mrt is less than 0!!!!!!!!!!")
			}

			info.RequestResponseTimes[requestName] = history.ResponseTimeStats{ResponseTimesMean: &mrt}
			itr.ResponseTimes[requestName] = mrt
		}

		info.ServicesInfo = make(map[string]history.ServiceInfo)
		for _, serviceName := range system.Resources {
			info.ServicesInfo[serviceName] = history.ServiceInfo{
				CPUUsageMean: system.GetUtilizations(alphas, serviceName),
			}
			itr.Utilizations[serviceName] = system.GetUtilizations(alphas, serviceName)
		}
		newState, isChanged, err := strategy.Configure(info, currentState, system.Resources)
		if err != nil {
			panic(err)
		}
		for serviceName, serviceState := range newState {
			ss := swarm.SimpleSpecs{}
			ss.CPU = serviceState.CPULimits
			ss.Replica = serviceState.ReplicaCount
			currentConfig[serviceName] = ss
		}
		iterations = append(iterations, itr)
		if debug {
			printIteration(itr, system)
		}
		if !isChanged {
			break
		}
		stepCount++
		// if stepCount == 100 {
		// 	break
		// }
	}

	bestIteration := GetBestIteration(iterations, system.SLA)

	row := name + ", "
	row += system.Name + ", "
	row += strconv.Itoa(stepCount) + ", "
	row += strconv.FormatFloat(bestIteration.GetSum(), 'f', 1, 64)
	row += "\n"
	return row
}

type iterationInfo struct {
	ResponseTimes map[string]float64
	Utilizations  map[string]float64
	Alphas        map[string]float64
}

func (i iterationInfo) GetSum() float64 {
	var s float64
	for _, a := range i.Alphas {
		s += a
	}
	return s
}

type theoryWorkload struct {
	Throughput float64
	ClassProbs map[string]float64
}

func (t theoryWorkload) GetThroughput() float64 {
	return t.Throughput
}

func (t theoryWorkload) GetRequestProportion() map[string]float64 {
	return t.ClassProbs
}

// GetBestIteration ...
func GetBestIteration(its []iterationInfo, SLA float64) iterationInfo {
	var minCPUCount float64 = 100000
	var best iterationInfo
	for _, i := range its {
		meets := true
		for _, r := range i.ResponseTimes {
			if r > SLA {
				meets = false
				break
			}
		}
		if meets && i.GetSum() < minCPUCount {
			minCPUCount = i.GetSum()
			best = i
		}
	}
	return best
}

func printIteration(i iterationInfo, system *System) {
	for _, rq := range system.Classes {
		rts := i.ResponseTimes[rq]
		fmt.Printf("%6.2f ", rts)
	}
	fmt.Print(", ")

	utils := make([]float64, 0)
	for _, s := range system.Resources {
		u := i.Utilizations[s]
		fmt.Printf("%4.2f ", u)
		utils = append(utils, u)
	}
	std, _ := stats.StandardDeviation(utils)
	fmt.Printf("%4.3f ", std)
	fmt.Print(", ")
	var sum float64 = 0
	for _, s := range system.Resources {
		a := i.Alphas[s]
		sum += a
		fmt.Printf("%2.2f ", a)
	}
	fmt.Printf("%5.2f\n", sum)
}
