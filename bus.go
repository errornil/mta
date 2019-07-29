package mta

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type DetailLevel string

const (
	stopMonitoringURL = "http://bustime.mta.info/api/siri/stop-monitoring.json"

	minimum DetailLevel = "minimum"
	basic   DetailLevel = "basic"
	normal  DetailLevel = "normal"
	calls   DetailLevel = "calls"
)

type BusTimeClient struct {
	apiKey string
	client *http.Client
}

func NewBusTimeClient(apiKey string, timeout time.Duration) *BusTimeClient {
	return &BusTimeClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *BusTimeClient) GetStopMonitoring(stopID string) (*StopMonitoringResponse, error) {
	return c.GetStopMonitoringWithDetailLevel(stopID, calls)
}

func (c *BusTimeClient) GetStopMonitoringWithDetailLevel(stopID string, detailLevel DetailLevel) (*StopMonitoringResponse, error) {
	v := url.Values{}
	v.Add("key", c.apiKey)
	v.Add("version", "2")
	v.Add("OperatorRef", "MTA")
	v.Add("MonitoringRef", stopID)
	v.Add("StopMonitoringDetailLevel", string(detailLevel))

	url := fmt.Sprintf("%s?%s", stopMonitoringURL, v.Encode())
	log.Println(url)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send GetStopMonitoring request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read GetStopMonitoring response: %v", err)
	}

	response := StopMonitoringResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GetStopMonitoring response: %v, body: %s", err, body)
	}

	return &response, nil
}
