package loadgenerator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// K6 connector to work with a K6 Wrapper
type K6 struct {
	Host string
}

// NewK6LoadGenerator is constructor
func NewK6LoadGenerator(host string) *K6 {
	return &K6{
		Host: host,
	}
}

// Prepare the load generator
func (k *K6) Prepare(values map[string]string) error {
	if _, ok := values["script"]; !ok {
		return fmt.Errorf("the k6 load generator needs script in the prepare method")
	}

	url := k.Host + "/prepare"
	method := "POST"

	b, err := json.Marshal(struct {
		Script string `json:"script"`
	}{Script: values["script"]})
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

func (k *K6) Start(values map[string]string) error {
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

func (k *K6) Stop(values map[string]string) error {
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

func (k *K6) GetFeedback(values map[string]string) (map[string]interface{}, error) {
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
