package mta

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	gtfs "github.com/chuhlomin/mta/v2/transit_realtime"
	"github.com/pkg/errors"

	"github.com/golang/protobuf/proto"
)

type Feed string

const (
	Feed123456S Feed = "nyct%2Fgtfs"      // Red
	FeedACEHS   Feed = "nyct%2Fgtfs-ace"  // Blue, Franklin Ave. Shuttle
	FeedNQRW    Feed = "nyct%2Fgtfs-nqrw" // Yellow
	FeedBDFM    Feed = "nyct%2Fgtfs-bdfm" // Orange
	FeedL       Feed = "nyct%2Fgtfs-l"
	FeedSIR     Feed = "nyct%2Fgtfs-si" // StatenIslandRailway
	FeedG       Feed = "nyct%2Fgtfs-g"
	FeedJZ      Feed = "nyct%2Fgtfs-jz" // Brown
	Feed7       Feed = "nyct%2Fgtfs-7"
	FeedLIRR    Feed = "lirr%2Fgtfs-lirr" // Long Island Rail Road
	FeedMNR     Feed = "mnr%2Fgtfs-mnr"   // Metro-North Railroad

	AlertsAll    Feed = "camsys%2Fall-alerts"    // All Service Alerts
	AlertsSubway Feed = "camsys%2Fsubway-alerts" // Subway Alerts
	AlertsBus    Feed = "camsys%2Fbus-alerts"    // Bus Alerts
	AlertsLIRR   Feed = "camsys%2Flirr-alerts"   // Long Island Rail Road Alerts
	AlertsMNR    Feed = "camsys%2Fmnr-alerts"    // Metro-North Railroad Alerts

	feedURL = "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/"
)

type FeedsService interface {
	GetFeedMessage(feedID Feed) (*gtfs.FeedMessage, error)
}

// FeedsClient provides MTA GTFS-Realtime data
// Implements FeedsService interface.
type FeedsClient struct {
	apiKey string
	client *http.Client
}

// NewFeedsClient creates new FeedsClient
func NewFeedsClient(apiKey string, timeout time.Duration) *FeedsClient {
	return &FeedsClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// GetFeedMessage sends request to MTA server to get latest GTFS-Realtime data from specified feed
func (f *FeedsClient) GetFeedMessage(feedID Feed) (*gtfs.FeedMessage, error) {
	url := fmt.Sprintf("%s%s", feedURL, feedID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new HTTP request")
	}

	req.Header.Add("x-api-key", f.apiKey)

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	feed := &gtfs.FeedMessage{}
	err = proto.Unmarshal(body, feed)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall GTFS Realtime Feed Message: %v", err)
	}

	return feed, nil
}
