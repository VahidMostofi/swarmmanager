package theory

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/cheggaaa/pb/v3"
	"github.com/montanaflynn/stats"
)

// GoTheory ...
func GoTheory() {
	testCount := 300
	bar := pb.StartNew(testCount)
	log.SetOutput(ioutil.Discard)
	line := "approach,system,steps,min_total_core,max_total_core,resources,classes,sla\n"

	for i := 1; i <= 1; i++ {
		i = 298
		outputs := make(map[string]output)
		fileName := strconv.Itoa(i)
		system := ReadSystem(fileName)

		// ----------------------------------------------------------------------------
		strategy := &strategies.BottleNeckOnlyVersion1{
			StepSize:          4.0,
			Agreements:        []strategies.Agreement{{"ResponseTimesMean", system.SLA}},
			MultiContainer:    true,
			DemandsFilePath:   "./theory/demands/" + fileName + ".yml",
			ConstantInit:      true,
			ConstantInitValue: 1.0,
		}
		strategy.Init()
		approachName := "BNV1-4.0"
		s, o := RunSystemWithStrategy(approachName, system, strategy, false)
		line += s
		outputs[approachName] = o
		// ----------------------------------------------------------------------------
		strategy = &strategies.BottleNeckOnlyVersion1{
			StepSize:          2.0,
			Agreements:        []strategies.Agreement{{"ResponseTimesMean", system.SLA}},
			MultiContainer:    true,
			DemandsFilePath:   "./theory/demands/" + fileName + ".yml",
			ConstantInit:      true,
			ConstantInitValue: 1.0,
		}
		strategy.Init()
		approachName = "BNV1-2.0"
		s, o = RunSystemWithStrategy(approachName, system, strategy, false)
		line += s
		outputs[approachName] = o
		// ----------------------------------------------------------------------------
		strategy = &strategies.BottleNeckOnlyVersion1{
			StepSize:          1.0,
			Agreements:        []strategies.Agreement{{"ResponseTimesMean", system.SLA}},
			MultiContainer:    true,
			DemandsFilePath:   "./theory/demands/" + fileName + ".yml",
			ConstantInit:      true,
			ConstantInitValue: 1.0,
		}
		strategy.Init()
		approachName = "BNV1-1.0"
		s, o = RunSystemWithStrategy(approachName, system, strategy, false)
		line += s
		outputs[approachName] = o
		// ----------------------------------------------------------------------------
		strategy = &strategies.BottleNeckOnlyVersion1{
			StepSize:          0.5,
			Agreements:        []strategies.Agreement{{"ResponseTimesMean", system.SLA}},
			MultiContainer:    true,
			DemandsFilePath:   "./theory/demands/" + fileName + ".yml",
			ConstantInit:      true,
			ConstantInitValue: 1.0,
		}
		strategy.Init()
		approachName = "BNV1-0.5"
		s, o = RunSystemWithStrategy(approachName, system, strategy, false)
		line += s
		outputs[approachName] = o
		// ----------------------------------------------------------------------------
		strategy2 := &strategies.BottleNeckOnlyVersion2{
			StepSize:          2.0,
			MinimumStepSize:   0.25,
			Agreements:        []strategies.Agreement{{"ResponseTimesMean", system.SLA}},
			MultiContainer:    true,
			DemandsFilePath:   "./theory/demands/" + fileName + ".yml",
			ConstantInit:      true,
			ConstantInitValue: 1.0,
		}
		strategy2.Init()
		strategy2.MinimumCPUValue = 1.0
		approachName = "BNV2-2.0"
		s, o = RunSystemWithStrategy(approachName, system, strategy2, false)
		line += s
		outputs[approachName] = o
		// ----------------------------------------------------------------------------
		strategy2 = &strategies.BottleNeckOnlyVersion2{
			StepSize:          4.0,
			MinimumStepSize:   0.25,
			Agreements:        []strategies.Agreement{{"ResponseTimesMean", system.SLA}},
			MultiContainer:    true,
			DemandsFilePath:   "./theory/demands/" + fileName + ".yml",
			ConstantInit:      true,
			ConstantInitValue: 1.0,
		}
		strategy2.Init()
		strategy2.MinimumCPUValue = 1.0
		approachName = "BNV2-4.0"
		s, o = RunSystemWithStrategy(approachName, system, strategy2, false)
		line += s
		outputs[approachName] = o
		// ----------------------------------------------------------------------------
		strategy2 = &strategies.BottleNeckOnlyVersion2{
			StepSize:          1.0,
			MinimumStepSize:   0.25,
			Agreements:        []strategies.Agreement{{"ResponseTimesMean", system.SLA}},
			MultiContainer:    true,
			DemandsFilePath:   "./theory/demands/" + fileName + ".yml",
			ConstantInit:      true,
			ConstantInitValue: 1.0,
		}
		strategy2.Init()
		strategy2.MinimumCPUValue = 1.0
		approachName = "BNV2-1.0"
		s, o = RunSystemWithStrategy(approachName, system, strategy2, false)
		line += s
		outputs[approachName] = o
		// ----------------------------------------------------------------------------
		line += "AMPL," + system.Name + ",0," + strconv.FormatFloat(system.BestObjective, 'f', 2, 64) + ",0,"
		line += strconv.Itoa(len(system.Resources)) + ","
		line += strconv.Itoa(len(system.Classes)) + ","
		line += strconv.FormatFloat(system.SLA, 'f', 1, 64)
		line += "\n"

		bar.Increment()
	}
	bar.Finish()
	// path := "/home/vahid/Dropbox/data/swarm-manager-data/results/theory/model-results-2x-not-that-late-1x-33p.csv"
	// err := ioutil.WriteFile(path, []byte(line), 0777)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("saved to:", path)
}

// RunSystemWithStrategy ...
func RunSystemWithStrategy(name string, system *System, strategy strategies.Configurer, debug bool) (string, output) {
	// fmt.Println(name)
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

		// ss := time.Now().Nanosecond()
		info := history.Information{}
		info.RequestResponseTimes = make(map[string]history.ResponseTimeStats)
		meanResponseTimes := system.GetMeanResponseTimes(alphas)
		for requestIdx, requestName := range system.Classes {
			mrt := meanResponseTimes[requestIdx]
			// if stepCount == 1 {
			// 	fmt.Println(mrt)
			// }
			if mrt <= 0 {
				fmt.Println(mrt, requestIdx, requestName)
				panic("mrt is less than 0!!!!!!!!!!")
			}

			info.RequestResponseTimes[requestName] = history.ResponseTimeStats{ResponseTimesMean: &mrt}
			itr.ResponseTimes[requestName] = mrt
		}
		// ee := time.Now().Nanosecond()
		// start := time.Now().Nanosecond()
		info.ServicesInfo = make(map[string]history.ServiceInfo)
		for _, serviceName := range system.Resources {
			info.ServicesInfo[serviceName] = history.ServiceInfo{
				CPUUsageMean: system.GetUtilizations(alphas, serviceName),
			}
			itr.Utilizations[serviceName] = system.GetUtilizations(alphas, serviceName)
		}
		// enda := time.Now().Nanosecond()
		// fmt.Println(enda-start, ee-ss)
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
		// if stepCount%100 == 0 {
		// 	fmt.Println(stepCount)
		// }
		if stepCount == 25000 {
			break
		}
	}

	for _, iter := range iterations {
		fmt.Println(iter.GetSum())
	}

	bestIteration := GetBestIteration(iterations, system.SLA)
	worstIterationWhichMeets := GetWorstIterationWhichMeetsSLA(iterations, system.SLA)

	row := name + ","
	row += system.Name + ","
	row += strconv.Itoa(stepCount) + ","
	row += strconv.FormatFloat(bestIteration.GetSum(), 'f', 1, 64) + ","
	row += strconv.FormatFloat(worstIterationWhichMeets.GetSum(), 'f', 1, 64) + ","
	row += strconv.Itoa(len(system.Resources)) + ","
	row += strconv.Itoa(len(system.Classes)) + ","
	row += strconv.FormatFloat(system.SLA, 'f', 1, 64)
	row += "\n"
	// fmt.Println(row)

	o := output{
		Name:       name,
		SystemName: system.Name,
		Steps:      stepCount,
		CPUs:       bestIteration.GetSum(),
	}
	return row, o
}

type output struct {
	Name       string
	SystemName string
	Steps      int
	CPUs       float64
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
	if s == 0 {
		return 100000
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

// GetWorstIterationWhichMeetsSLA ...
func GetWorstIterationWhichMeetsSLA(its []iterationInfo, SLA float64) iterationInfo {
	var maxCPUCount float64 = 0
	var worst iterationInfo
	for _, i := range its {
		meets := true
		for _, r := range i.ResponseTimes {
			if r > SLA {
				meets = false
				break
			}
		}
		if meets && i.GetSum() > maxCPUCount {
			maxCPUCount = i.GetSum()
			worst = i
		}
	}
	return worst
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
