package jaeger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type JaegerAggregator struct {
	Host  string
	Keys  []string // important spans that we want to consider
	Spans map[string][]float64
}

// NewJaegerAggregator is the constructor for JaegerAggregator
func NewJaegerAggregator(host string, keys []string) *JaegerAggregator {
	j := &JaegerAggregator{
		Host:  host,
		Keys:  keys,
		Spans: make(map[string][]float64),
	}

	for _, key := range j.Keys {
		j.Spans[key] = make([]float64, 0)
	}
	return j
}

// GetTraces retrieves traces form Jaeger instance
func (j *JaegerAggregator) GetTraces(start, end int64, service string) {
	url := fmt.Sprintf("%s/api/traces?end=%d&limit=100000&service=%s&start=%d", j.Host, end, service, start)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	data := struct {
		Data []*trace `json:"data"`
	}{}
	json.Unmarshal(body, &data)

	for _, key := range j.Keys {
		j.Spans[key] = make([]float64, 0)
	}

	for _, trace := range data.Data {
		for _, span := range trace.Spans {
			if j.isSpanKey(span) {
				j.Spans[span.OperationName] = append(j.Spans[span.OperationName], span.Duration/1e3)
			}
		}
	}
}

// GetRequestCount is to we comply with ResponseTimeCollector
func (j *JaegerAggregator) GetRequestCount(name string) (int, error) {
	if _, ok := j.Spans[name]; ok {
		return len(j.Spans[name]), nil
	}
	return 0, fmt.Errorf("No such key found: %s", name)
}

// GetResponseTimes is to we comply with RequestCountCollector
func (j *JaegerAggregator) GetResponseTimes(name string) ([]float64, error) {
	if _, ok := j.Spans[name]; ok {
		return j.Spans[name], nil
	}
	return nil, fmt.Errorf("No such key found: %s", name)
}

func (j *JaegerAggregator) isSpanKey(s *span) bool {
	for _, key := range j.Keys {
		if key == s.OperationName {
			return true
		}
	}
	return false
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
	StartTime     float64 `json:"startTime"`
	Duration      float64 `json:"duration"`
	OperationName string  `json:"operationName"`
	SpanID        string  `json:"spanID"`
	TraceID       string  `json:"traceID"`
	IsRoot        bool    `json:"-"`
}
