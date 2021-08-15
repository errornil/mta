package mta

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const LIRRDepartureURL = "https://traintime.lirr.org/api/Departure"

type LIRRService interface {
	Departures(locationCode string) (*DeparturesResponse, error)
}

type LIRRClient struct {
	client    HTTPClient
	userAgent string
}

type DeparturesResponse struct {
	Location string                    `json:"LOC"`    // The three-letter code for that station, e.g.: JAM
	Time     string                    `json:"TIME"`   // The date and time the feed was returned at in mm/dd/yyyy hh:mm:ss format (24-hr time)
	Trains   []DeparturesResponseTrain `json:"TRAINS"` // Countdown items for each arriving train
}

// DeparturesResponseTrain - part of DeparturesResponse
type DeparturesResponseTrain struct {
	ScheduledTime string   `json:"SCHED"`    // The scheduled date and time the train is supposed to arrive at the station in mm/dd/yyyy hh:mm:ss format (24-hr time)
	TrainID       string   `json:"TRAIN_ID"` // The train number. These are typically 1-4 digit train numbers you can find in the timetables or as the train_id’s in the GTFS feeds, though since this feed shows inserts, they can be up to 8 alphanumeric characters long.
	RunDate       string   `json:"RUN_DATE"` // E.g.: "2019-12-14",
	Destination   string   `json:"DEST"`     // The three-letter station code of the final stop on that train.
	Stops         []string `json:"STOPS"`    // The three-letter station codes of all remaining stops that train is supposed to make. If the countdown is showing at the final destination, the field will be blank for “discharge only” trains. NOTE: Hillside stops are not shown in this field.
	Track         string   `json:"TRACK"`    // The track and platform the train is supposed to depart from. May be blank at terminals when the track is not yet posted.
	Direction     string   `json:"DIR"`      // The direction the train is travelling (E = eastbound, W = westbound)
	HSF           bool     `json:"HSF"`      // Indicates whether or not the train will stop at Hillisde (true = train stops, false = train does not stop)
	JAM           *bool    `json:"JAM"`      // Indicates whether or not the train will stop at Jamaica (ture = train stops, false = train will pass through Jamaica but not stop, null = train won’t pass through Jamaica).
	ETA           string   `json:"ETA"`      // The estimated arrival time of that train, updated to account for all reported schedule deviations. Returned in mm/dd/yyyy hh:mm:ss format (24-hr time). The difference between this time and the ScheduledTime is how late or early the train is.
	Countdown     int      `json:"CD"`       // The number of seconds between now and the time the train is supposed to arrive (a countdown field).
}

// NewLIRRClient creates new LIRRClient
func NewLIRRClient(client HTTPClient, userAgent string) (*LIRRClient, error) {
	if client == nil {
		return nil, ErrClientRequired
	}

	return &LIRRClient{
		client:    client,
		userAgent: userAgent,
	}, nil
}

// Departures gets train departures from specified station
// locationCode – three-letter code for the station, for example:
//   "NYK" for NY-Penn Station
//   "ATL" for Brooklyn-Atlantic Term
//   "HVL" for Hicksville
// Full list of codes: https://github.com/errornil/mta/blob/master/lirr.md
func (lc *LIRRClient) Departures(locationCode string) (*DeparturesResponse, error) {
	v := url.Values{}
	v.Add("loc", locationCode)

	url := fmt.Sprintf("%s?%s", LIRRDepartureURL, v.Encode())

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new HTTP request")
	}
	req.Header.Add("User-Agent", lc.userAgent)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send Departures request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non 200 response status: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read Departures response")
	}

	response := DeparturesResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse Departures response body: %s", body)
	}

	return &response, nil
}
