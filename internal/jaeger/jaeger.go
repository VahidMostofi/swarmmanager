package jaeger

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/VahidMostofi/swarmmanager/configs"
	uuid "github.com/nu7hatch/gouuid"
	"gopkg.in/yaml.v2"
)

type formula struct {
	Key     string `yaml:"name"`
	Formula string `yaml:"value"`
}

type tag struct {
	Key   string
	Value string
}

type requestFormula struct {
	Name         string
	ResponseTime string `yaml:"responseTime"`
	Tags         []tag
}

type serviceFormula struct {
	Name     string
	Formulas []formula
}

type valueFormula struct {
	Requests map[string]requestFormula
	Services map[string]serviceFormula
}

// Aggregator ...
type Aggregator struct {
	Host                  string
	requestsResponseTimes map[string][]float64
	servicesTimeDetails   map[string]map[string]map[string][]float64
	StorePath             string
	LastStoredFile        string
	ValueFormulas         valueFormula
}

// NewAggregator is the constructor for Aggregator
func NewAggregator() *Aggregator {
	temp := valueFormula{}

	b, err := ioutil.ReadFile(configs.GetConfig().Jaeger.DetailsFilePath)
	if err != nil {
		log.Panic(err)
	}
	yaml.Unmarshal(b, &temp)
	j := &Aggregator{
		Host:          configs.GetConfig().Jaeger.Host,
		StorePath:     configs.GetConfig().Jaeger.StorePath,
		ValueFormulas: temp,
	}

	return j
}

func (j *Aggregator) getTraces(start, end int64, service string) ([]*trace, error) {
	body := make([]byte, 0)
	for len(body) < 100 { //TODO WTF with 100?!?!?!?
		url := fmt.Sprintf("%s/api/traces?end=%d&limit=100000&service=%s&start=%d", j.Host, end, service, start)
		method := "GET"
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			panic(fmt.Errorf("error while getting traces: %w", err))
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(fmt.Errorf("error getting traces with http request: %w", err))
			time.Sleep(3 * time.Second)
			continue
		}
		body, err = ioutil.ReadAll(res.Body)
		fmt.Println("len(body)", len(body))
		fmt.Println(url)
		res.Body.Close()
		time.Sleep(3 * time.Second)
	}

	if len(j.StorePath) > 1 {
		id, err := uuid.NewV4()
		if err != nil {
			panic(fmt.Errorf("error in generating uuid %w", err)) //TODO
		}
		j.LastStoredFile = j.StorePath + "/" + id.String() + ".zip"
		log.Println("Jaeger: storing traces at:", j.LastStoredFile)
		newZipFile, err := os.Create(j.LastStoredFile)
		if err != nil {
			panic(err)
		}
		defer newZipFile.Close()
		zipWriter := zip.NewWriter(newZipFile)
		defer zipWriter.Close()
		f, err := zipWriter.Create("jaeger-info.json")
		if err != nil {
			panic(err)
		}
		_, err = f.Write(body)
		if err != nil {
			panic(err)
		}
	} //TODO separate this function into multiple functions

	data := struct {
		Data []*trace `json:"data"`
	}{}
	json.Unmarshal(body, &data)

	return data.Data, nil
}

func (j *Aggregator) identifyRequest(t *trace) string {

	for _, s := range t.Spans {
		for _, request := range j.ValueFormulas.Requests {
			matchingTags := 0
			for _, tag := range s.Tags {
				for _, optionTag := range request.Tags {
					if tag.Key == optionTag.Key && strings.Contains(tag.Value, optionTag.Value) {
						matchingTags++
					}
				}
				if matchingTags == len(request.Tags) {
					return request.Name
				}
			}
		}
	}
	fmt.Println("found no request type for this trace, number of spans are:", len(t.Spans), "marking trace as invalid")
	return ""
}

func (j *Aggregator) parseTraces(Data []*trace) error {

	for i := range Data {
		Data[i].RequestType = j.identifyRequest(Data[i])
		if Data[i].RequestType == "" {
			Data[i].Valid = false
			continue
		} else {
			Data[i].Valid = true
		}
		Data[i].Services = make(map[string]string)
		for j := range Data[i].Spans {
			Data[i].Spans[j].ServiceName = Data[i].Processes[Data[i].Spans[j].ProcessID].ServiceName
			Data[i].Services[Data[i].Spans[j].ServiceName] = ""
		}
	}

	j.requestsResponseTimes = make(map[string][]float64)
	for _, request := range j.ValueFormulas.Requests {
		j.requestsResponseTimes[request.Name] = make([]float64, 0)
	}

	failedEvaluations := 0

	for _, trace := range Data {
		if !trace.Valid {
			continue
		}
		value, err := evaluateJaegerFormula(j.ValueFormulas.Requests[trace.RequestType].ResponseTime, trace)
		if err != nil {
			// log.Panic(fmt.Errorf("error while evaluating Jaeger Formula: %w", err))
			failedEvaluations++
		} else {
			j.requestsResponseTimes[trace.RequestType] = append(j.requestsResponseTimes[trace.RequestType], value)
			// fmt.Println(value)
		}
	}
	log.Println("Jaeger: failed evaluations count", failedEvaluations)
	j.servicesTimeDetails = make(map[string]map[string]map[string][]float64)
	for _, service := range j.ValueFormulas.Services {
		j.servicesTimeDetails[service.Name] = make(map[string]map[string][]float64)
		for _, request := range j.ValueFormulas.Requests {
			j.servicesTimeDetails[service.Name][request.Name] = make(map[string][]float64)
		}
	}

	for _, t := range Data {
		if !t.Valid {
			continue
		}
		for serviceInTrace := range t.Services {
			for _, f := range j.ValueFormulas.Services[serviceInTrace].Formulas {
				value, err := evaluateJaegerFormula(f.Formula, t)
				if value == 0 {
					fmt.Println("this cant be 0! (I guess)")
					return nil
				}
				if err != nil {
					return err
				}

				if _, contains := j.servicesTimeDetails[serviceInTrace][t.RequestType][f.Key]; !contains {
					j.servicesTimeDetails[serviceInTrace][t.RequestType][f.Key] = make([]float64, 0)
				}
				j.servicesTimeDetails[serviceInTrace][t.RequestType][f.Key] = append(j.servicesTimeDetails[serviceInTrace][t.RequestType][f.Key], value)
			}
		}
	}

	return nil
}

// GetTraces retrieves traces form Jaeger instance
func (j *Aggregator) GetTraces(start, end int64, service string) {

	Data, err := j.getTraces(start, end, service)
	if err != nil {
		panic(err)
	}

	err = j.parseTraces(Data)
	if err != nil {
		panic(err)
	}

}

// GetRequestCount is to we comply with ResponseTimeCollector
func (j *Aggregator) GetRequestCount(name string) (int, error) {
	log.Println("GetRequestCount:", "called with", name)
	return len(j.requestsResponseTimes[name]), nil
}

// GetRequestResponseTimes ....
func (j *Aggregator) GetRequestResponseTimes(name string) ([]float64, error) {
	log.Println("GetResponseTimes:", "called with", name)
	res := j.requestsResponseTimes[name]
	log.Println("GetResponseTimes:", "called with", name, "found", len(res), "values")
	return res, nil
}

// GetServiceDetails ...
func (j *Aggregator) GetServiceDetails(name string) (map[string]map[string][]float64, error) {
	log.Println("GetServiceDetails:", "called with", name)
	return j.servicesTimeDetails[name], nil
}

// GetRequestNames ...
func (j *Aggregator) GetRequestNames() []string {
	names := make([]string, 0)
	for n := range j.ValueFormulas.Requests {
		names = append(names, n)
	}
	return names
}

// Trace struct contains a group of spans
type trace struct {
	Spans       []*span            `json:"spans"`
	RequestType string             `json:"-"`
	Valid       bool               `json:"-"`
	TraceID     string             `json:"traceID"`
	Processes   map[string]process `json:"processes"`
	Services    map[string]string  `json:"-"`
}

type process struct {
	ServiceName string `json:"serviceName"`
}

// Span struct contains information about each span
type span struct {
	StartTime     float64       `json:"startTime"`
	EndTime       float64       `json:"endTime"`
	Duration      float64       `json:"duration"`
	OperationName string        `json:"operationName"`
	SpanID        string        `json:"spanID"`
	TraceID       string        `json:"traceID"`
	IsRoot        bool          `json:"-"`
	References    []interface{} `json:"references"`
	Tags          []tag         `json:"tags"`
	ProcessID     string        `json:"processID"`
	ServiceName   string        `json:"-"`
}
