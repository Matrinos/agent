// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package edgex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/mainflux/mainflux/logger"

	model "github.com/edgexfoundry/go-mod-core-contracts/models"
)

type Client interface {

	// PushOperation - pushes operation to EdgeX components
	PushOperation([]string) (string, error)

	// FetchConfig - fetches config from EdgeX components
	FetchConfig([]string) (string, error)

	// FetchMetrics - fetches metrics from EdgeX components
	FetchMetrics(cmdArr []string) (string, error)

	ControlDevice(cmdArr []string) (string, error)

	// Ping - ping EdgeX SMA
	Ping() (string, error)
}

type edgexClient struct {
	url    string
	logger log.Logger
}

// NewClient - Creates ne EdgeX client
func NewClient(edgexURL string, logger log.Logger) Client {
	return &edgexClient{
		url:    edgexURL,
		logger: logger,
	}
}

// ControlDevice - control device
func (ec *edgexClient) ControlDevice(cmdArr []string) (string, error) {

	url := "http://192.168.124.150:59882/api/v2/device/name/Nader100/SwitchCommand"
	fmt.Printf("EdgeX Url: %v\n", url)
	data, err := json.Marshal(map[string]string{
		"Switch": "65280",
	})
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer req.Body.Close()
	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// PushOperation - pushes operation to EdgeX components
func (ec *edgexClient) PushOperation(cmdArr []string) (string, error) {
	url := ec.url + "operation"

	m := model.Operation{
		Action:   cmdArr[0],
		Services: cmdArr[1:],
	}
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// FetchConfig - fetches config from EdgeX components
func (ec *edgexClient) FetchConfig(cmdArr []string) (string, error) {
	cmdStr := strings.Replace(strings.Join(cmdArr, ","), " ", "", -1)
	url := ec.url + "config/" + cmdStr

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// FetchMetrics - fetches metrics from EdgeX components
func (ec *edgexClient) FetchMetrics(cmdArr []string) (string, error) {
	cmdStr := strings.Replace(strings.Join(cmdArr, ","), " ", "", -1)
	url := ec.url + "metrics/" + cmdStr

	resp, err := http.Get(url)
	if err != nil {

		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Ping - ping EdgeX SMA
func (ec *edgexClient) Ping() (string, error) {
	url := ec.url + "ping"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
