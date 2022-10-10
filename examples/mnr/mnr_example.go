package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/errornil/mta/v3"
	"github.com/errornil/mta/v3/proto/mnr"
)

func main() {
	if err := run(); err != nil {
		log.Printf("ERROR: %v", err)
	}
}

func run() error {
	apiKey := flag.String("key", "", "API key")
	flag.Parse()
	if *apiKey == "" {
		return errors.New("missing API key, pass it with -key flag")
	}

	client, err := mta.NewFeedsClient(
		&http.Client{
			Timeout: 30 * time.Second,
		},
		*apiKey,
		"",
		"github.com/errornil/mta:v2.0",
	)
	if err != nil {
		return errors.Wrap(err, "failed to get feed message")
	}

	msg, err := client.GetFeedMessage(mta.FeedMNR)
	if err != nil {
		return errors.Wrap(err, "failed to get feed message")
	}

	for _, entity := range msg.GetEntity() {
		for _, stopTimeUpdate := range entity.GetTripUpdate().GetStopTimeUpdate() {
			// optional MnrStopTimeUpdate
			var mnrStopTimeUpdate *mnr.MnrStopTimeUpdate
			mnrStopTimeUpdate = proto.GetExtension(stopTimeUpdate, mnr.E_MnrStopTimeUpdate).(*mnr.MnrStopTimeUpdate)
			b, err := json.MarshalIndent(mnrStopTimeUpdate, "", "  ")
			if err != nil {
				return errors.Wrap(err, "failed to marshal entity")
			}
			log.Println(string(b))
		}

		b, err := json.MarshalIndent(entity, "", "  ")
		if err != nil {
			return errors.Wrap(err, "failed to marshal entity")
		}

		log.Println(string(b))
	}

	return nil
}
