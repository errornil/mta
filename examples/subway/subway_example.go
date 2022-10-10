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
	"github.com/errornil/mta/v3/proto/subway"
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

	msg, err := client.GetFeedMessage(mta.Feed123456S)
	if err != nil {
		return errors.Wrap(err, "failed to get feed message")
	}

	// optional NyctFeedHeader
	var nyctFeedHeader *subway.NyctFeedHeader
	nyctFeedHeader = proto.GetExtension(msg.GetHeader(), subway.E_NyctFeedHeader).(*subway.NyctFeedHeader)
	b, err := json.MarshalIndent(nyctFeedHeader, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal entity")
	}
	log.Println(string(b))

	for _, entity := range msg.GetEntity() {
		// optional NyctTripDescriptor
		var nyctTripDescriptor *subway.NyctTripDescriptor
		nyctTripDescriptor = proto.GetExtension(entity.GetTripUpdate().GetTrip(), subway.E_NyctTripDescriptor).(*subway.NyctTripDescriptor)
		b, err := json.MarshalIndent(nyctTripDescriptor, "", "  ")
		if err != nil {
			return errors.Wrap(err, "failed to marshal entity")
		}
		log.Println(string(b))

		b, err = json.MarshalIndent(entity, "", "  ")
		if err != nil {
			return errors.Wrap(err, "failed to marshal entity")
		}

		log.Println(string(b))
	}

	return nil
}
