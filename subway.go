package mta

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	gtfs "github.com/chuhlomin/mta/transit_realtime"

	"github.com/golang/protobuf/proto"
)

type Feed int

const (
	Feed123456S Feed = 1  // Red
	FeedACEHS   Feed = 26 // Blue, Franklin Ave. Shuttle
	FeedNQRW    Feed = 16 // Yellow
	FeedBDFM    Feed = 21 // Orange
	FeedL       Feed = 2
	FeedSIR     Feed = 11 // StatenIslandRailway
	FeedG       Feed = 31
	FeedJZ      Feed = 36 // Brown
	Feed7       Feed = 51

	feedURL = "http://datamine.mta.info/mta_esi.php"
)

type SubwayTimeService interface {
	GetFeedMessage(feedID Feed) (*gtfs.FeedMessage, error)
}

type SubwayTimeClient struct {
	apiKey string
	client *http.Client
}

func NewSubwayTimeClient(apiKey string, timeout time.Duration) *SubwayTimeClient {
	return &SubwayTimeClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *SubwayTimeClient) GetFeedMessage(feedID Feed) (*gtfs.FeedMessage, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s?key=%s&feed_id=%d", feedURL, c.apiKey, feedID))
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
