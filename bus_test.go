package mta

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBusErrAPIKeyRequired(t *testing.T) {
	_, err := NewBusTimeClient("", nil)
	require.Error(t, err, ErrAPIKeyRequired)
}

func TestBusErrClientRequired(t *testing.T) {
	_, err := NewBusTimeClient("apiKey", nil)
	require.Error(t, err, ErrClientRequired)
}

func TestGetStopMonitoring(t *testing.T) {
	GetFunc = func(url string) (*http.Response, error) {
		json := `{
			"Siri": {
				"ServiceDelivery": {
					"ResponseTimestamp": "2006-01-02T15:04:05.000-07:00",
					"StopMonitoringDelivery": [
						{
							"ResponseTimestamp": "2006-01-02T15:04:05.000-07:00",
							"MonitoredStopVisit": [
								{
									"MonitoredVehicleJourney": {
										"LineRef": "MTA NYCT_M20",
										"FramedVehicleJourneyRef": {
											"DataFrameRef": "M20",
											"DatedVehicleJourneyRef": "M20"
										},
										"OperatorRef": "MTA",
										"OriginRef": "404847",
										"PublishedLineName": [
											"M20"
										],
										"DestinationName": [
											"LINCOLN CENTER 66 ST via 8 AV"
										],
										"MonitoredCall": {
											"ExpectedArrivalTime": "2021-08-14T23:40:57.787-04:00"
										}
									}
								}
							]
						}
					],
					"SituationExchangeDelivery": []
				}
			}
		}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	}
	c, err := NewBusTimeClient("apiKey", mockClient{})
	require.NoError(t, err)

	resp, err := c.GetStopMonitoring("404847")
	require.NoError(t, err)

	expectedResp := &StopMonitoringResponse{
		Siri: Siri{
			ServiceDelivery: ServiceDelivery{
				ResponseTimestamp: "2006-01-02T15:04:05.000-07:00",
				StopMonitoringDelivery: []StopMonitoringDelivery{
					{
						ResponseTimestamp: "2006-01-02T15:04:05.000-07:00",
						MonitoredStopVisit: []MonitoredStopVisit{
							{
								MonitoredVehicleJourney: MonitoredVehicleJourney{
									LineRef: "MTA NYCT_M20",
									FramedVehicleJourneyRef: FramedVehicleJourneyRef{
										DataFrameRef:           "M20",
										DatedVehicleJourneyRef: "M20",
									},
									OperatorRef:       "MTA",
									OriginRef:         "404847",
									PublishedLineName: []string{"M20"},
									DestinationName:   []string{"LINCOLN CENTER 66 ST via 8 AV"},
									MonitoredCall: &MonitoredCall{
										ExpectedArrivalTime: "2021-08-14T23:40:57.787-04:00",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	require.Equal(t, expectedResp, resp)
}

func TestGetStopMonitoringErrAPIKeyNotAuthorized(t *testing.T) {
	GetFunc = func(url string) (*http.Response, error) {
		json := `{
			"Siri": {
				"ServiceDelivery": {
					"ResponseTimestamp": "2006-01-02T15:04:05.000-07:00",
					"VehicleMonitoringDelivery": [
						{
							"ResponseTimestamp": "2006-01-02T15:04:05.000-07:00",
							"ErrorCondition": {
								"OtherError": {
									"ErrorText": "API key is not authorized."
								},
								"Description": "API key is not authorized."
							}
						}
					]
				}
			}
		}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	}
	c, _ := NewBusTimeClient("apiKey", mockClient{})

	_, err := c.GetStopMonitoring("404847")
	require.Error(t, err, ErrAPIKeyNotAuthorized)
}

func TestGetStopMonitoringErrAPIKeyRequired2(t *testing.T) {
	GetFunc = func(url string) (*http.Response, error) {
		json := `{
			"Siri": {
				"ServiceDelivery": {
					"ResponseTimestamp": "2006-01-02T15:04:05.000-07:00",
					"VehicleMonitoringDelivery": [
						{
							"ResponseTimestamp": "2006-01-02T15:04:05.000-07:00",
							"ErrorCondition": {
								"OtherError": {
									"ErrorText": "API key required."
								},
								"Description": "API key required."
							}
						}
					]
				}
			}
		}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	}
	c, _ := NewBusTimeClient("apiKey", mockClient{})

	_, err := c.GetStopMonitoring("404847")
	require.Error(t, err, "API key required.")
}

func TestGetStopMonitoringErrRequestSend(t *testing.T) {
	GetFunc = func(url string) (*http.Response, error) {
		return nil, net.UnknownNetworkError("...")
	}
	c, _ := NewBusTimeClient("apiKey", mockClient{})

	_, err := c.GetStopMonitoring("404847")
	require.Error(t, err, "failed to send GetStopMonitoring request: unknown network ...")
}

func TestGetStopMonitoringErrReadBody(t *testing.T) {
	mockReadCloser := mockReadCloser{}
	mockReadCloser.On("Read", mock.AnythingOfType("[]uint8")).Return(0, fmt.Errorf("error reading"))
	mockReadCloser.On("Close").Return(nil)

	GetFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       &mockReadCloser,
		}, nil
	}
	c, _ := NewBusTimeClient("apiKey", mockClient{})

	_, err := c.GetStopMonitoring("404847")
	require.Error(t, err, "failed to parse GetStopMonitoring response: error reading")
}

func TestGetStopMonitoringErrBadResponse(t *testing.T) {
	GetFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("not JSON"))),
		}, nil
	}
	c, _ := NewBusTimeClient("apiKey", mockClient{})

	_, err := c.GetStopMonitoring("404847")
	require.Error(t, err, "failed to parse GetStopMonitoring response: invalid character 'o' in literal null (expecting 'u'), body: not JSON")
}
