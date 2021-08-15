package mta

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
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
	client    HTTPClient
	apiKey    string
	userAgent string
}

func NewBusTimeClient(client HTTPClient, apiKey, userAgent string) (*BusTimeClient, error) {
	if client == nil {
		return nil, ErrClientRequired
	}
	if apiKey == "" {
		return nil, ErrAPIKeyRequired
	}
	return &BusTimeClient{
		apiKey:    apiKey,
		userAgent: userAgent,
		client:    client,
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

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new HTTP request")
	}
	req.Header.Add("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
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
