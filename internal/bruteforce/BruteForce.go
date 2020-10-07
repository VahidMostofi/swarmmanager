package bruteforce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/jaeger"
	"github.com/VahidMostofi/swarmmanager/internal/loadgenerator"
	"github.com/VahidMostofi/swarmmanager/internal/statutils"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/VahidMostofi/swarmmanager/internal/utils"
	"github.com/VahidMostofi/swarmmanager/internal/workload"
	"github.com/montanaflynn/stats"
	"gopkg.in/yaml.v2"
)

// BruteForce ...
type BruteForce struct {
	Workload              string
	LoadGenerator         loadgenerator.LoadGenerator
	ResponseTimeCollector workload.ResponseTimeCollector
	RequestCountCollector workload.RequestCountCollector
	SwarmManager          *swarm.Manager
	LoadGeneratorStarted  bool
	InitialCPUCount       float64
	CPUStepSize           float64
	FinalCPUCount         float64
	Version               string
}

// NewBruteForce ...
func NewBruteForce(lg loadgenerator.LoadGenerator, rtc workload.ResponseTimeCollector, rcc workload.RequestCountCollector, m *swarm.Manager, workload string) *BruteForce {
	b := &BruteForce{
		LoadGenerator:         lg,
		ResponseTimeCollector: rtc,
		RequestCountCollector: rcc,
		SwarmManager:          m,
		Workload:              workload,
		LoadGeneratorStarted:  false,
		InitialCPUCount:       0.20,
		CPUStepSize:           0.20,
		FinalCPUCount:         6,
		Version:               "v1",
	}
	return b
}

func (b *BruteForce) mainLoop() {
	allConfigs := b.generateAllConfigurationsSorted()
	prevConfig := make(map[string]swarm.SimpleSpecs)
	for service := range allConfigs[0] {
		prevConfig[service] = swarm.SimpleSpecs{}
	}

	// var start int64
	// var end int64
	for i := 0; i < len(allConfigs); i++ {

		if !b.isConfigNew(allConfigs[i]) {
			continue
		}
		if getTotalCPUCount(allConfigs[i]) <= 5.22 {
			continue
		}
		currentConfig := allConfigs[i]
		log.Println("working on config", (i + 1), b.hash(currentConfig), "/", len(allConfigs), toString(currentConfig))

		err := b.SwarmManager.FastUpdate(currentConfig, prevConfig)
		if err != nil {
			panic(fmt.Errorf("cant update the services: %w", err))
		}
		time.Sleep(8 * time.Second)
		log.Println("Deployed")
		go b.LoadGenerator.Start()
		prevConfig = currentConfig

		time.Sleep(30 * time.Second)
		b.LoadGenerator.Stop()
		time.Sleep(2 * time.Second)

		data, err := b.LoadGenerator.GetFeedbackRaw()
		if err != nil {
			panic(err)
		}
		lf := &LGFeedback{}
		json.Unmarshal(data, lf)

		isGood := b.isGood(*lf)
		if isGood {
			b.saveLoadGeneratorFeedback(currentConfig)
			b.saveToFile(currentConfig, getGoodDirecotry())
		}
		b.saveToFile(currentConfig, getPartialDirecotry())
	}
}

// Start ...
func (b *BruteForce) Start() {
	b.InitialConfig()

	if err := os.MkdirAll(filepath.Dir(getPartialDirecotry()), 0777); err != nil {
		log.Panic(err)
	}

	if err := os.MkdirAll(filepath.Dir(getGoodDirecotry()), 0777); err != nil {
		log.Panic(err)
	}

	b.mainLoop()

}

// InitialConfig ...
func (b *BruteForce) InitialConfig() {

	// Remove the current Stack
	err := b.SwarmManager.RemoveStack(1)
	if err != nil {
		log.Panic(err)
	}

	dockerComposePath := configs.GetConfig().TestBed.DockerComposeFile

	initialConfig := make(map[string]swarm.SimpleSpecs)
	for _, service := range configs.GetConfig().TestBed.ServicesToConfigure {
		initialConfig[service] = swarm.SimpleSpecs{
			CPU:     b.InitialCPUCount,
			Replica: 1,
			Worker:  1,
		}
	}

	initialConfig = reconfigure(initialConfig)

	err = b.SwarmManager.DeployStackWithDockerCompose(dockerComposePath, 1, initialConfig)
	if err != nil {
		log.Panic(err)
	}
	time.Sleep(10 * time.Second)
	for {
		if b.SwarmManager.CurrentStackState == swarm.StackStateMustCompare {
			log.Println("CurrentStackState is", swarm.GetStateString(b.SwarmManager.CurrentStackState), "so lets break out of the loop")
			b.SwarmManager.UpdateCurrentSpecs()
			break
		}
		time.Sleep(150 * time.Millisecond)
	}
	log.Println("Services are ready")
}

func (b *BruteForce) generateAllConfigurationsSorted() []map[string]swarm.SimpleSpecs {
	all := make([]map[string]swarm.SimpleSpecs, 0)
	if len(configs.GetConfig().TestBed.ServicesToConfigure) == 3 {
		for i := b.InitialCPUCount; i < b.FinalCPUCount; i += b.CPUStepSize {
			for j := b.InitialCPUCount; j < b.FinalCPUCount; j += b.CPUStepSize {
				for k := b.InitialCPUCount; k < b.FinalCPUCount; k += b.CPUStepSize {
					config := make(map[string]swarm.SimpleSpecs)
					config[configs.GetConfig().TestBed.ServicesToConfigure[0]] = swarm.SimpleSpecs{
						CPU: i, Worker: 1, Replica: 1,
					}
					config[configs.GetConfig().TestBed.ServicesToConfigure[1]] = swarm.SimpleSpecs{
						CPU: j, Worker: 1, Replica: 1,
					}
					config[configs.GetConfig().TestBed.ServicesToConfigure[2]] = swarm.SimpleSpecs{
						CPU: k, Worker: 1, Replica: 1,
					}
					config = reconfigure(config)
					all = append(all, config)
				}
			}
		}
		sort.SliceStable(all, func(i, j int) bool {
			if getTotalCPUCount(all[i]) == getTotalCPUCount(all[j]) {
				return getMulCPUCount(all[i]) > getMulCPUCount(all[j])
			}
			return getTotalCPUCount(all[i]) < getTotalCPUCount(all[j])
		})
	} else {
		log.Panic("HEY!")
	}
	return all
}

func (b *BruteForce) isConfigNew(config map[string]swarm.SimpleSpecs) bool {
	return !fileExists(getPartialDirecotry() + b.hash(config) + ".yml")
}

func round(value float64) float64 {
	return math.Floor(value*1000) / 1000
}

func reconfigure(config map[string]swarm.SimpleSpecs) map[string]swarm.SimpleSpecs {
	newConfig := make(map[string]swarm.SimpleSpecs)

	for service := range config {
		totalCPU := config[service].CPU * float64(config[service].Replica)
		replicaCount := int(math.Ceil(totalCPU))
		eachContainerCPU := round(totalCPU / float64(replicaCount))
		c := swarm.SimpleSpecs{
			CPU:     eachContainerCPU,
			Replica: replicaCount,
			Worker:  1,
		}
		newConfig[service] = c
	}

	return newConfig
}

func getTotalCPUCount(config map[string]swarm.SimpleSpecs) float64 {
	var total float64
	for _, c := range config {
		total += c.CPU * float64(c.Replica)
	}
	return total
}

func getMulCPUCount(config map[string]swarm.SimpleSpecs) float64 {
	var total float64
	for _, c := range config {
		total *= c.CPU * float64(c.Replica)
	}
	return total
}

func (b *BruteForce) hash(config map[string]swarm.SimpleSpecs) string {
	services := make([]string, 0)
	for service := range config {
		services = append(services, service)
	}
	sort.Strings(services)
	res := ""
	for _, service := range services {
		serviceCPU := config[service].CPU * float64(config[service].Replica)
		res += strconv.FormatFloat(serviceCPU, 'f', 3, 64) + "_"
	}
	res += b.Version
	return res
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getPartialDirecotry() string {
	return configs.GetConfig().Results.Path + "partials/"
}

func getGoodDirecotry() string {
	return configs.GetConfig().Results.Path + "good/"
}

func toString(config map[string]swarm.SimpleSpecs) string {
	res := "| "
	for name, c := range config {
		count := float64(c.Replica) * c.CPU
		cStr := strconv.FormatFloat(count, 'f', 3, 64)
		res += name + ": " + cStr + " | "
	}
	return res
}

// GatherResponeTimes ...
func (b *BruteForce) GatherResponeTimes(start, end int64) map[string]history.ResponseTimeStats {

	b.RequestCountCollector.(*jaeger.Aggregator).GetTraces(start, end, configs.GetConfig().Jaeger.RootService, false)
	res := make(map[string]history.ResponseTimeStats)

	for _, req := range []string{"login", "get_book", "edit_book"} {
		values, err := b.ResponseTimeCollector.GetRequestResponseTimes(req)
		if err != nil {
			log.Panic(err)
		}
		responseTimeDetails, err := createStats(values, []string{"count", "mean", "p90", "p95", "p99", "std", "c90p95"})
		if err != nil {
			log.Panic(err)
		}
		res[req] = responseTimeDetails
	}
	return res
}

func createStats(values []float64, names []string) (history.ResponseTimeStats, error) {
	if len(values) == 0 {
		values = append(values, 0)
	}
	rts := history.ResponseTimeStats{}

	if utils.ContainsString(names, "mean") {
		mean, err := stats.Mean(values)
		if err != nil {
			log.Panic(err)
		}
		rts.ResponseTimesMean = &mean
	}
	//--------------------------------------------------
	if utils.ContainsString(names, "p90") {
		p90, err := stats.Percentile(values, 90)
		if err != nil {
			log.Panic(err)
		}
		rts.ResponseTimes90Percentile = &p90
	}

	//--------------------------------------------------
	if utils.ContainsString(names, "p95") {
		p95, err := stats.Percentile(values, 95)
		if err != nil {
			log.Panic(err)
		}
		rts.ResponseTimes95Percentile = &p95
	}

	//--------------------------------------------------
	if utils.ContainsString(names, "p99") {
		p99, err := stats.Percentile(values, 99)
		if err != nil {
			log.Panic(err)
		}
		rts.ResponseTimes99Percentile = &p99
	}

	//--------------------------------------------------
	if utils.ContainsString(names, "std") {
		std, err := stats.StandardDeviation(values)
		if err != nil {
			log.Panic(err)
		}
		rts.ResponseTimesSTD = &std
	}

	//--------------------------------------------------
	if utils.ContainsString(names, "count") {
		c := len(values)
		rts.Count = &c
	}

	//--------------------------------------------------
	if utils.ContainsString(names, "c90p95") && len(values) > 10 {
		_, uti, err := statutils.ComputeToleranceIntervalNonParametric(values, 0.90, 0.95)
		if err != nil {
			log.Println("response times:", values)
			log.Panic(err)
		}
		rts.RTToleranceIntervalUBoundConfidence90p95 = &uti
	}
	return rts, nil
}

func (b *BruteForce) isGood(lf LGFeedback) bool {

	lf.Metrics.LoginDurationTrend.P95 = round(lf.Metrics.LoginDurationTrend.P95)
	lf.Metrics.EditBookDurationTrend.P95 = round(lf.Metrics.EditBookDurationTrend.P95)
	lf.Metrics.GetBookDurationTrend.P95 = round(lf.Metrics.GetBookDurationTrend.P95)
	fmt.Print("login", " ", lf.Metrics.LoginDurationTrend.P95, " | ", "edit_book ", lf.Metrics.EditBookDurationTrend.P95, " | ", "get_booK ", lf.Metrics.GetBookDurationTrend.P95)
	if lf.Metrics.LoginDurationTrend.P95 > 250 {
		fmt.Println(" ")
		return false
	}
	if lf.Metrics.EditBookDurationTrend.P95 > 250 {
		fmt.Println(" ")
		return false
	}
	if lf.Metrics.GetBookDurationTrend.P95 > 250 {
		fmt.Println(" ")
		return false
	}
	fmt.Println(" meets SLA")
	return true
}

func (b *BruteForce) saveToFile(config map[string]swarm.SimpleSpecs, directory string) {
	data, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(directory+b.hash(config)+".yml", data, 0777)
}

func (b *BruteForce) saveLoadGeneratorFeedback(config map[string]swarm.SimpleSpecs) {
	feedback, err := b.LoadGenerator.GetFeedback()
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(feedback, "", "    ")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(getGoodDirecotry()+b.hash(config)+".json", data, 0777)

}

type DurationTrend struct {
	Avg float64 `json:"avg"`
	P95 float64 `json:"p(95)"`
}

type LGFeedback struct {
	Metrics Metrics `json:"metrics"`
}

type Metrics struct {
	EditBookDurationTrend DurationTrend `json:"edit_book_duration_trend"`
	LoginDurationTrend    DurationTrend `json:"get_book_duration_trend"`
	GetBookDurationTrend  DurationTrend `json:"login_duration_trend"`
}
