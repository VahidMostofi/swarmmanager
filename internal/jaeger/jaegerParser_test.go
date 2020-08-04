package jaeger

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/montanaflynn/stats"
	"gopkg.in/yaml.v2"
)

func TestJaegerParser(t *testing.T) {
	// TODO this is hardcoded
	formulasPath := "/home/vahid/Desktop/projects/swarmmanager/formulas/nodejs_bookstore.yaml"
	temp := &struct {
		ServiceDetails map[string]jaegerServiceDetail `yaml:"service_details"`
		Formulas       []valueFormula                 `yaml:"formulas"`
	}{}
	b, err := ioutil.ReadFile(formulasPath)
	yaml.Unmarshal(b, temp)
	// TODO this is hardcoded
	jaegerDataFilePath := "/home/vahid/Dropbox/data/swarm-manager-data/jaegers/e6f732ea-7a6d-497a-46fc-1b83f8c920f4.zip"
	r, err := zip.OpenReader(jaegerDataFilePath)
	if err != nil {
		t.Error(err)
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
	j := &Aggregator{
		Values:         make(map[string][]float64),
		ServiceDetails: temp.ServiceDetails,
		Formulas:       temp.Formulas,
	}
	b, err = yaml.Marshal(j)
	ioutil.WriteFile("/home/vahid/Desktop/v.yaml", b, 0777)

	err = j.parseTraces(data.Data)
	if err != nil {
		t.Errorf("error evaluateing formula: %w", err)
	}

	for _, values := range j.Values {
		if len(values) < 1 {
			t.Fail()
		}
		_, err := stats.Mean(values)
		if err != nil {
			t.Error(err)
		}
	}
}
