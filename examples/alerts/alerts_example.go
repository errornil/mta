package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/errornil/mta/v3"
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

	msg, err := client.GetFeedMessage(mta.AlertsSubway)
	if err != nil {
		return errors.Wrap(err, "failed to get feed message")
	}

	for _, entity := range msg.GetEntity() {
		log.Println(entity.String())
	}

	return nil
}
