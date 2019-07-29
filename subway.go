package mta

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	gtfs "github.com/chuhlomin/mta/transit_realtime"

	"github.com/golang/protobuf/proto"
)

type Line int

const (
	Line123456S Line = 1  // Red
	LineACEHS   Line = 26 // Blue, Franklin Ave. Shuttle
	LineNQRW    Line = 16 // Yellow
	LineBDFM    Line = 21 // Orange
	LineL       Line = 2
	LineSIR     Line = 11 // StatenIslandRailway
	LineG       Line = 31
	LineJZ      Line = 36 // Brown
	Line7       Line = 51

	feedURL = "http://datamine.mta.info/mta_esi.php"
)

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

func (c *SubwayTimeClient) GetFeedMessage(feedID Line) (*gtfs.FeedMessage, error) {
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
