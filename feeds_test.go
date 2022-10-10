package mta

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/errornil/mta/v3/proto/transit_realtime"
)

func str(v string) *string {
	return &v
}

func TestFeedsErrClientRequired(t *testing.T) {
	_, err := NewFeedsClient(nil, "", "", "")
	require.Error(t, err)
	require.Equal(t, ErrClientRequired, err)
}

func TestFeedsErrAPIKeyRequired(t *testing.T) {
	fc, err := NewFeedsClient(&mockClient{}, "", "", "")
	require.NoError(t, err)

	_, err = fc.GetFeedMessage(Feed(FeedBusTripUpdates))
	require.Error(t, err)
	require.Equal(t, ErrAPIKeyRequired, err)
}

func TestGetFeedMessage(t *testing.T) {
	msg := transit_realtime.FeedMessage{
		Header: &transit_realtime.FeedHeader{
			GtfsRealtimeVersion: str("1.0"),
		},
		Entity: []*transit_realtime.FeedEntity{
			{
				Id: str("1"),
				TripUpdate: &transit_realtime.TripUpdate{
					Trip: &transit_realtime.TripDescriptor{
						StartTime: str("20160101T000000"),
					},
					StopTimeUpdate: []*transit_realtime.TripUpdate_StopTimeUpdate{
						{
							StopId: str("1"),
							Departure: &transit_realtime.TripUpdate_StopTimeEvent{
								Time: proto.Int64(1),
							},
						},
					},
				},
			},
		},
	}

	DoFunc = func(req *http.Request) (*http.Response, error) {
		require.Equal(t, "apiKey", req.Header.Get("x-api-key"))
		b, _ := proto.Marshal(&msg)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(b)),
		}, nil
	}
	c, err := NewFeedsClient(&mockClient{}, "apiKey", "", "")
	require.NoError(t, err)

	feedMessage, err := c.GetFeedMessage(Feed123456S)
	require.NoError(t, err)
	require.True(t, proto.Equal(&msg, feedMessage))
}

func TestGetFeedMessageErrRequestSend(t *testing.T) {
	DoFunc = func(req *http.Request) (*http.Response, error) {
		return nil, net.UnknownNetworkError("...")
	}
	c, _ := NewFeedsClient(&mockClient{}, "apiKey", "", "")

	_, err := c.GetFeedMessage(Feed123456S)
	require.Error(t, err)
	require.Equal(t, "failed to send GET request: unknown network ...", err.Error())
}

func TestGetFeedMessageErrNon200(t *testing.T) {
	DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			Status:     "500 Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewReader([]byte("..."))),
		}, nil
	}
	c, _ := NewFeedsClient(&mockClient{}, "apiKey", "", "")

	_, err := c.GetFeedMessage(Feed123456S)
	require.Error(t, err)
	require.Equal(t, "non 200 response status: 500 Internal Server Error", err.Error())
}

func TestGetFeedMessageErrReadBody(t *testing.T) {
	mockReadCloser := mockReadCloser{}
	mockReadCloser.On("Read", mock.AnythingOfType("[]uint8")).Return(0, fmt.Errorf("error reading"))
	mockReadCloser.On("Close").Return(nil)

	DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       &mockReadCloser,
		}, nil
	}
	c, _ := NewFeedsClient(&mockClient{}, "apiKey", "", "")

	_, err := c.GetFeedMessage(Feed123456S)
	require.Error(t, err)
	require.Equal(t, "failed to read response body: error reading", err.Error())
}

func TestGetFeedMessageErrBadResponse(t *testing.T) {
	DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte("not Protobuf"))),
		}, nil
	}
	c, _ := NewFeedsClient(&mockClient{}, "apiKey", "", "")

	_, err := c.GetFeedMessage(Feed123456S)
	require.Error(t, err)
	require.Equal(t, "failed to unmarshall GTFS Realtime Feed Message: proto: cannot parse invalid wire-format data", err.Error())
}
