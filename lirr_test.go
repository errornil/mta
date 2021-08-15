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

func b(v bool) *bool {
	return &v
}

func TestLIRRErrClientRequired(t *testing.T) {
	_, err := NewLIRRClient(nil, "")
	require.Error(t, err, ErrClientRequired)
}

func TestDepartures(t *testing.T) {
	DoFunc = func(req *http.Request) (*http.Response, error) {
		json := `{
			"LOC": "NYK",
			"TIME": "08/15/2021 14:48:01",
			"TRAINS": [
				{
					"SCHED": "01/02/2006 15:04:05",
					"TRAIN_ID": "6112",
					"RUN_DATE": "2006-01-02",
					"DEST": "BTA",
					"STOPS": [
						"FHL",
						"KGN"
					],
					"TRACK": "16",
					"DIR": "E",
					"HSF": false,
					"JAM": true,
					"ETA": "01/02/2006 15:04:05",
					"CD": 419
				}
			]
		}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	}

	c, err := NewLIRRClient(mockClient{}, "")
	require.NoError(t, err)

	d, err := c.Departures("NYK")
	require.NoError(t, err)

	expected := DeparturesResponse{
		Location: "NYK",
		Time:     "08/15/2021 14:48:01",
		Trains: []DeparturesResponseTrain{
			{
				ScheduledTime: "01/02/2006 15:04:05",
				TrainID:       "6112",
				RunDate:       "2006-01-02",
				Destination:   "BTA",
				Stops: []string{
					"FHL",
					"KGN",
				},
				Track:     "16",
				Direction: "E",
				HSF:       false,
				JAM:       b(true),
				ETA:       "01/02/2006 15:04:05",
				Countdown: 419,
			},
		},
	}
	require.Equal(t, &expected, d)
}

func TestDeparturesErrRequestSend(t *testing.T) {
	DoFunc = func(req *http.Request) (*http.Response, error) {
		return nil, net.UnknownNetworkError("...")
	}

	c, err := NewLIRRClient(mockClient{}, "")
	require.NoError(t, err)

	_, err = c.Departures("NYK")
	require.Error(t, err, "failed to send Departures request: unknown network ...")
}

func TestDeparturesErrReadBody(t *testing.T) {
	mockReadCloser := mockReadCloser{}
	mockReadCloser.On("Read", mock.AnythingOfType("[]uint8")).Return(0, fmt.Errorf("error reading"))
	mockReadCloser.On("Close").Return(nil)

	DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       &mockReadCloser,
		}, nil
	}

	c, err := NewLIRRClient(mockClient{}, "")
	require.NoError(t, err)

	_, err = c.Departures("NYK")
	require.Error(t, err, "read response body: error reading")
}

func TestDeparturesErrBadResponse(t *testing.T) {
	DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("not JSON"))),
		}, nil
	}

	c, err := NewLIRRClient(mockClient{}, "")
	require.NoError(t, err)

	_, err = c.Departures("NYK")
	require.Error(t, err, "failed to parse Departures response body: invalid character 'o' in literal null (expecting 'u'), body: not JSON")
}
