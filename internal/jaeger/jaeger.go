package jaeger

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VahidMostofi/swarmmanager"
	uuid "github.com/nu7hatch/gouuid"
)

type JaegerAggregator struct {
	Host           string
	Keys           []string // important spans that we want to consider
	Spans          map[string][]float64
	StorePath      string
	LastStoredFile string
}

// NewJaegerAggregator is the constructor for JaegerAggregator
func NewJaegerAggregator(host string, keys []string) *JaegerAggregator {
	j := &JaegerAggregator{
		Host:      host,
		Keys:      keys,
		Spans:     make(map[string][]float64),
		StorePath: swarmmanager.GetConfig().JaegerStorePath,
	}

	for _, key := range j.Keys {
		j.Spans[key] = make([]float64, 0)
	}
	return j
}

// GetTraces retrieves traces form Jaeger instance
func (j *JaegerAggregator) GetTraces(start, end int64, service string) {
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

	for _, key := range j.Keys {
		j.Spans[key] = make([]float64, 0)
	}

	for _, trace := range data.Data {
		var service string
		var request string

		spans := make(map[string]*span)
		for _, span := range trace.Spans {
			if len(trace.Spans) < 6 {
				log.Println("warning", "len(trace.Spans) is", len(trace.Spans))
				continue
			} //TODO implement retry
			span.EndTime = span.StartTime + span.Duration
			spans[span.OperationName] = span

			if span.OperationName == "auth_req_login" {
				service = "auth"
				request = "login"
			}
			if span.OperationName == "books_get_book" {
				service = "books"
				request = "get_book"
			}
			if span.OperationName == "books_edit_book" {
				service = "books"
				request = "edit_book"
			}
		}

		var sup float64
		var sub float64
		if request == "login" {
			sup = spans["auth_req_login"].Duration
			sub = spans["auth"].EndTime - spans["auth_connect"].EndTime
		} else if request == "edit_book" {
			sup = spans["books_edit_book"].Duration
			sub = spans["books"].EndTime - spans["books_connect"].EndTime
		} else if request == "get_book" {
			sup = spans["books_get_book"].Duration
			sub = spans["books"].EndTime - spans["books_connect"].EndTime
		} else {
			continue
		}
		sub /= 1000
		sup /= 1000
		if service == "auth" {
			j.Spans["auth"] = append(j.Spans["auth"], sup)
			j.Spans["auth_total"] = append(j.Spans["auth_total"], sup)
			j.Spans["auth_gateway"] = append(j.Spans["auth_gateway"], sup-sub)
			j.Spans["auth_sub"] = append(j.Spans["auth_sub"], sub)
		} else {
			j.Spans["books"] = append(j.Spans["books"], sup)
			j.Spans["books_total"] = append(j.Spans["books_total"], sup)
			j.Spans["books_gateway"] = append(j.Spans["books_gateway"], sup-sub)
			j.Spans["books_sub"] = append(j.Spans["books_sub"], sub)
		}
		j.Spans["gateway"] = append(j.Spans["gateway"], sup)
	}
	fmt.Println("these are the keys")
	for key := range j.Spans {
		fmt.Println(key)
	}
}

// GetRequestCount is to we comply with ResponseTimeCollector
func (j *JaegerAggregator) GetRequestCount(name string) (int, error) {
	log.Println("GetRequestCount:", "called with", name)
	if _, ok := j.Spans[name]; ok {
		return len(j.Spans[name]), nil
	}
	return 0, fmt.Errorf("No such key found: %s", name)
}

// GetResponseTimes is to we comply with RequestCountCollector
func (j *JaegerAggregator) GetResponseTimes(name string) ([]float64, error) {
	log.Println("GetResponseTimes:", "called with", name)
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
	StartTime     float64       `json:"startTime"`
	EndTime       float64       `json:"endTime"`
	Duration      float64       `json:"duration"`
	OperationName string        `json:"operationName"`
	SpanID        string        `json:"spanID"`
	TraceID       string        `json:"traceID"`
	IsRoot        bool          `json:"-"`
	References    []interface{} `json:"references"`
}
