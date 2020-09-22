package loadgenerator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/spf13/viper"
)

// K6Workload ...
type K6Workload struct {
	CmdStr            string
	requestProportion map[string]float64
	SleepTime         float64
	args              map[string]string
	VirtualUsersCount float64
	Duration          float64
	BaseURL           string
	Architecture      string
}

// GetRequestProportion ...
func (k *K6Workload) GetRequestProportion() map[string]float64 {
	return k.requestProportion
}

// GetThroughput ...
func (k *K6Workload) GetThroughput() float64 {
	return k.VirtualUsersCount / k.SleepTime
}

// K6 connector to work with a K6 Wrapper
type K6 struct {
	Host     string
	Workload *K6Workload
	Script   string
}

// NewK6LoadGenerator is constructor
func newK6LoadGenerator() (*K6, error) {
	w := &K6Workload{
		CmdStr: viper.GetString("workloadStr"),
		args:   make(map[string]string),
	}

	// VUS
	vus, err := strconv.ParseFloat(strings.Split(w.CmdStr, "_")[0], 64)
	if err != nil {
		return nil, fmt.Errorf("Cant parse number of VUS in workload: %s", strings.Split(strings.Split(w.CmdStr, "_")[0], "_")[0])
	}
	w.VirtualUsersCount = vus
	w.args["ARG_VUS"] = strings.Split(w.CmdStr, "_")[0]
	// -------------------------------------------------------------------

	// DURATION
	duration, err := strconv.ParseFloat(strings.Split(w.CmdStr, "_")[1], 64)
	if err != nil {
		return nil, fmt.Errorf("Cant parse Duration in workload: %s", strings.Split(strings.Split(w.CmdStr, "_")[1], "_")[0])
	}
	if duration < float64(configs.GetConfig().Test.Duration) {
		panic("for now these two values should be equal or duration should be more than TestDuration!")
	}
	w.Duration = duration
	w.args["ARG_DURATION"] = strings.Split(w.CmdStr, "_")[1]
	// -------------------------------------------------------------------

	// SLEEP DURATION
	sleepTime, err := strconv.ParseFloat(configs.GetConfig().LoadGenerator.Args["ARG_SLEEP_DURATION"], 64)
	if err != nil {
		return nil, fmt.Errorf("Cant parse SleepDuration in workload: %s", strings.Split(configs.GetConfig().LoadGenerator.Args["ARG_SLEEP_DURATION"], "_")[3])
	}
	w.SleepTime = sleepTime
	w.args["ARG_SLEEP_DURATION"] = configs.GetConfig().LoadGenerator.Args["ARG_SLEEP_DURATION"]
	// -------------------------------------------------------------------

	// BASE URL
	baseURL := configs.GetConfig().LoadGenerator.Args["ARG_BASE_URL"]
	w.BaseURL = baseURL
	w.args["ARG_BASE_URL"] = configs.GetConfig().LoadGenerator.Args["ARG_BASE_URL"]
	// -------------------------------------------------------------------

	// ARCHITECTURE
	architecture := configs.GetConfig().LoadGenerator.Args["ARG_ARCHITECTURE"]
	w.Architecture = architecture
	w.args["ARG_ARCHITECTURE"] = configs.GetConfig().LoadGenerator.Args["ARG_ARCHITECTURE"]
	// -------------------------------------------------------------------

	// REQUESTS
	reqNames := make([]string, 0)
	for _, reqName := range strings.Split(configs.GetConfig().LoadGenerator.Args["REQUEST_NAMES"], ",") {
		reqNames = append(reqNames, reqName)
	}
	reqProbs := make(map[string]float64)
	for idx, reqProbStr := range strings.Split(w.CmdStr, "_")[2:] {
		reqProb, err := strconv.ParseFloat(reqProbStr, 64)
		if err != nil {
			return nil, fmt.Errorf("cant parse request prob number %d %s", idx, reqProbStr)
		}
		reqProbs[reqNames[idx]] = reqProb
		w.args["ARG_"+reqNames[idx]] = reqProbStr
	}
	w.requestProportion = reqProbs
	lg := &K6{
		Host:     configs.GetConfig().LoadGenerator.Details["host"],
		Workload: w,
	}
	lg.Script = createLoadGeneartorScript(configs.GetConfig().LoadGenerator.Details["script"], w.args)

	return lg, nil
}

// GetWorkload ...
func (k *K6) GetWorkload() Workload {
	return k.Workload
}

// Prepare the load generator
func (k *K6) Prepare() error {
	url := k.Host + "/prepare"
	method := "POST"

	b, err := json.Marshal(struct {
		Script string `json:"script"`
	}{Script: k.Script})
	if err != nil {
		return fmt.Errorf("error while convert k6 prepare input to json: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(b))

	if err != nil {
		return fmt.Errorf("error while creating prepare request to k6 wrapper server: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error from sending prepare request to k6 wrapper server: %w", err)
	}
	defer res.Body.Close()
	return nil
}

// Start ...
func (k *K6) Start() error {
	url := k.Host + "/start"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return fmt.Errorf("error while creating start request to k6 wrapper server: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error from sending start request to k6 wrapper server: %w", err)
	}
	defer res.Body.Close()

	return nil
}

// Stop ...
func (k *K6) Stop() error {
	url := k.Host + "/stop"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return fmt.Errorf("error while creating stop request to k6 wrapper server: %w", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error from sending stop request to k6 wrapper server: %w", err)
	}
	defer res.Body.Close()
	return nil
}

// GetFeedback ...
func (k *K6) GetFeedback() (map[string]interface{}, error) {
	url := k.Host + "/feedback"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, fmt.Errorf("error while creating feedback request to k6 wrapper server: %w", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error from sending feedback request to k6 wrapper server: %w", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	feedbackRes := make(map[string]interface{})
	if err = json.Unmarshal(body, &feedbackRes); err != nil {
		return nil, fmt.Errorf("error while parsing json in feedback of k6 wrapper: %w", err)
	}

	return feedbackRes, nil
}

// createLoadGeneartorScript ...
func createLoadGeneartorScript(scriptPath string, args map[string]string) string {
	res := readLoadGeneratorScript(scriptPath)
	for key, value := range args {
		res = strings.ReplaceAll(res, key, value)
	}
	return res
}

// ReadLoadGeneratorScript ...
func readLoadGeneratorScript(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(fmt.Errorf("cant load k6 load geneartor script at: %s; %w", path, err))
	}
	return string(b)
}
