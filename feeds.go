package mta

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	gtfs "github.com/errornil/mta/v3/transit_realtime"
)

type Feed string

const (
	Feed123456S Feed = "nyct/gtfs"      // Red
	FeedACEHS   Feed = "nyct/gtfs-ace"  // Blue, Franklin Ave. Shuttle
	FeedNQRW    Feed = "nyct/gtfs-nqrw" // Yellow
	FeedBDFM    Feed = "nyct/gtfs-bdfm" // Orange
	FeedL       Feed = "nyct/gtfs-l"
	FeedSIR     Feed = "nyct/gtfs-si" // StatenIslandRailway
	FeedG       Feed = "nyct/gtfs-g"
	FeedJZ      Feed = "nyct/gtfs-jz" // Brown
	Feed7       Feed = "nyct/gtfs-7"
	FeedLIRR    Feed = "lirr/gtfs-lirr" // Long Island Rail Road
	FeedMNR     Feed = "mnr/gtfs-mnr"   // Metro-North Railroad

	AlertsAll    Feed = "camsys/all-alerts"    // All Service Alerts
	AlertsSubway Feed = "camsys/subway-alerts" // Subway Alerts
	AlertsBus    Feed = "camsys/bus-alerts"    // Bus Alerts
	AlertsLIRR   Feed = "camsys/lirr-alerts"   // Long Island Rail Road Alerts
	AlertsMNR    Feed = "camsys/mnr-alerts"    // Metro-North Railroad Alerts

	FeedURL = "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/"

	FeedBusTripUpdates      Feed = "tripUpdates"
	FeedBusVehiclePositions Feed = "vehiclePositions"
	FeedBusAlerts           Feed = "alerts"

	FeedURLBus = "http://gtfsrt.prod.obanyc.com/"
)

var (
	SubwayFeeds []Feed = []Feed{
		Feed123456S,
		FeedACEHS,
		FeedNQRW,
		FeedBDFM,
		FeedL,
		FeedSIR,
		FeedG,
		FeedJZ,
		Feed7,
	}

	BusFeeds []Feed = []Feed{
		FeedBusTripUpdates,
		FeedBusVehiclePositions,
		FeedBusAlerts,
	}

	AllFeeds []Feed = append(SubwayFeeds, FeedLIRR, FeedMNR)

	AllAlerts []Feed = []Feed{
		AlertsAll,
		AlertsSubway,
		AlertsBus,
		AlertsLIRR,
		AlertsMNR,
	}
)

type FeedsService interface {
	GetFeedMessage(feedID Feed) (*gtfs.FeedMessage, error)
}

// FeedsClient provides MTA GTFS-Realtime data
// Implements FeedsService interface.
type FeedsClient struct {
	client      HTTPClient
	feedsApiKey string
	busApiKey   string
	userAgent   string
}

// NewFeedsClient creates new FeedsClient
func NewFeedsClient(client HTTPClient, feedsApiKey, busApiKey, userAgent string) (*FeedsClient, error) {
	if client == nil {
		return nil, ErrClientRequired
	}
	return &FeedsClient{
		client:      client,
		feedsApiKey: feedsApiKey,
		busApiKey:   busApiKey,
		userAgent:   userAgent,
	}, nil
}

// GetFeedMessage sends request to MTA server to get latest GTFS-Realtime data from specified feed
func (f *FeedsClient) GetFeedMessage(feedID Feed) (*gtfs.FeedMessage, error) {
	feedURL, key, err := f.buildURL(feedID)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new HTTP request")
	}

	req.Header.Add("x-api-key", key)

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send GET request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non 200 response status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	feed := &gtfs.FeedMessage{}
	err = proto.Unmarshal(body, feed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshall GTFS Realtime Feed Message")
	}

	return feed, nil
}

func (f *FeedsClient) buildURL(feedID Feed) (feedURL string, key string, err error) {
	if isBusFeed(feedID) {
		feedURL = fmt.Sprintf("%s%s", FeedURLBus, feedID)
		key = f.busApiKey
	} else {
		feedURL = fmt.Sprintf("%s%s", FeedURL, url.PathEscape(string(feedID)))
		key = f.feedsApiKey
	}

	if key == "" {
		err = ErrAPIKeyRequired
	}
	return
}

func isBusFeed(feedID Feed) bool {
	return feedID == FeedBusTripUpdates ||
		feedID == FeedBusVehiclePositions ||
		feedID == FeedBusAlerts
}
