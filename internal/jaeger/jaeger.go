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

	"github.com/VahidMostofi/swarmmanager"
	uuid "github.com/nu7hatch/gouuid"
	"gopkg.in/yaml.v2"
)

type valueFormula struct {
	RequestName string `yaml:"request_name"`
	ValueName   string `yaml:"value_name"`
	Formula     string `yaml:"formula"`
}

type jaegerServiceDetail struct {
	ServiceName         string `yaml:"service_name"`
	RequestName         string `yaml:"request_name"`
	MinNumberOfSpans    int    `yaml:"min_number_of_spans"`
	UniqueOperationName string `yaml:"unique_operation_name"`
}

// Aggregator ...
type Aggregator struct {
	Host           string
	Values         map[string][]float64
	StorePath      string
	LastStoredFile string
	ServiceDetails map[string]jaegerServiceDetail `yaml:"service_details"` //if changed, change NewAggregator
	Formulas       []valueFormula                 `yaml:"formula"`         //if changed, change NewAggregator
}

// NewAggregator is the constructor for Aggregator
func NewAggregator() *Aggregator {
	temp := &struct {
		ServiceDetails map[string]jaegerServiceDetail `yaml:"service_details"`
		Formulas       []valueFormula                 `yaml:"formulas"`
	}{}

	b, err := ioutil.ReadFile(swarmmanager.GetConfig().JaegerDetailsFilePath)
	if err != nil {
		log.Panic(err)
	}
	yaml.Unmarshal(b, temp)
	j := &Aggregator{
		Host:           swarmmanager.GetConfig().JaegerHost,
		Values:         make(map[string][]float64),
		StorePath:      swarmmanager.GetConfig().JaegerStorePath,
		ServiceDetails: temp.ServiceDetails,
		Formulas:       temp.Formulas,
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

func (j *Aggregator) parseTraces(Data []*trace) error {

	for _, formula := range j.Formulas {
		j.Values[formula.ValueName] = make([]float64, 0)
	}
	incompleteTraceCount := 0
	for _, trace := range Data {
		var service string
		var request string
		var key string

		spans := make(map[string]*span) //TODO this should be a map to an array of span if we want to support a calling same operations multiple times
		for _, span := range trace.Spans {
			span.StartTime /= 1000
			span.Duration /= 1000
			span.EndTime = span.StartTime + span.Duration
			span.OperationName = strings.ReplaceAll(span.OperationName, "-", "_")
			spans[span.OperationName] = span

			for k, details := range j.ServiceDetails {
				if span.OperationName == details.UniqueOperationName {
					service = details.ServiceName
					request = details.RequestName
					key = k
					break
				}
			}
		}

		if len(trace.Spans) < j.ServiceDetails[key].MinNumberOfSpans {
			incompleteTraceCount++
			continue
		} //TODO implement retry

		for _, formula := range j.Formulas {
			if request == formula.RequestName || (strings.HasPrefix(formula.RequestName, "@service:") && formula.RequestName == "@service:"+service) || formula.RequestName == "@any:" {
				value, err := evaluateJaegerFormula(formula.Formula, spans)
				if err != nil {
					panic(err)
				}
				j.Values[formula.ValueName] = append(j.Values[formula.ValueName], value)
			}
		}
	}
	log.Println("warning", "incomplete Trace Count", incompleteTraceCount)
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
func (j *Aggregator) GetRequestCount(name string) (map[string]int, error) {
	log.Println("GetRequestCount:", "called with", name)
	res := make(map[string]int)
	for valueName, values := range j.Values {
		if strings.HasPrefix(valueName, name) {
			res[valueName] = len(values)
		}
	}
	log.Println("GetRequestCount:", "called with", name, "found", len(res), "values")
	return res, nil
}

// GetResponseTimes is to we comply with RequestCountCollector
func (j *Aggregator) GetResponseTimes(name string) (map[string][]float64, error) {
	log.Println("GetResponseTimes:", "called with", name)

	res := make(map[string][]float64)
	for valueName, values := range j.Values {
		if strings.HasPrefix(valueName, name) {
			res[valueName] = values
		}
	}
	log.Println("GetResponseTimes:", "called with", name, "found", len(res), "values")
	return res, nil
}

// Trace struct contains a group of spans
type trace struct {
	Spans     []*span `json:"spans"`
	TraceType string  `json:"-"`
	HasRoot   bool    `json:"-"`
	TraceID   string  `json:"traceID"`
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
}
