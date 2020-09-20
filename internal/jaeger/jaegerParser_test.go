package jaeger

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestJaegerParser(t *testing.T) {
	// TODO this is hardcoded
	formulasPath := "/home/vahid/Desktop/projects/swarmmanager/configurations/formulas/muck_general.yaml"
	temp := valueFormula{}
	b, err := ioutil.ReadFile(formulasPath)
	yaml.Unmarshal(b, &temp)
	if err != nil {
		t.Error(err)
	}

	// TODO this is hardcoded
	jaegerDataFilePath := "/home/vahid/Dropbox/data/swarm-manager-data/jaegers/64c12667-8dcb-4581-412a-ff505d925c05.zip"
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
		ValueFormulas: temp,
	}
	b, err = yaml.Marshal(j)
	ioutil.WriteFile("/home/vahid/Desktop/v.yaml", b, 0777)
	fmt.Println(len(data.Data))
	err = j.parseTraces(data.Data)
	if err != nil {
		t.Errorf("error evaluateing formula: %w", err)
	}
	for _, t := range data.Data {
		fmt.Println(t.Valid)
	}
	for r := range j.ValueFormulas.Requests {
		v, _ := j.GetRequestResponseTimes(r)
		fmt.Println(len(v))
	}

	// for service, sValue := range j.servicesTimeDetails {
	// 	for request, rValue := range sValue {
	// 		for name, values := range rValue {
	// 			fmt.Println(service, request, name, len(values))
	// 		}
	// 	}
	// }
}
