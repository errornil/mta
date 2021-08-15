package mta

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
)

type DetailLevel string

const (
	StopMonitoringURL = "http://bustime.mta.info/api/siri/stop-monitoring.json"

	minimum DetailLevel = "minimum"
	basic   DetailLevel = "basic"
	normal  DetailLevel = "normal"
	calls   DetailLevel = "calls"
)

type BusTimeService interface {
	GetStopMonitoring(stopID string) (*StopMonitoringResponse, error)
	GetStopMonitoringWithDetailLevel(stopID string, detailLevel DetailLevel) (*StopMonitoringResponse, error)
}

type BusTimeClient struct {
	apiKey string
	client HTTPClient
}

func NewBusTimeClient(apiKey string, client HTTPClient) (*BusTimeClient, error) {
	if apiKey == "" {
		return nil, ErrAPIKeyRequired
	}
	if client == nil {
		return nil, ErrClientRequired
	}
	return &BusTimeClient{
		apiKey: apiKey,
		client: client,
	}, nil
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

	url := fmt.Sprintf("%s?%s", StopMonitoringURL, v.Encode())

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

	if len(response.Siri.ServiceDelivery.VehicleMonitoringDelivery) > 0 {
		message := response.Siri.ServiceDelivery.VehicleMonitoringDelivery[0].ErrorCondition.Description
		if message == "API key is not authorized." {
			return &response, ErrAPIKeyNotAuthorized
		}
		return &response, errors.New(message)
	}

	return &response, nil
}
