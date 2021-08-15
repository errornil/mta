package main

import (
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/errornil/mta/v2"
)

func main() {
	if err := run(); err != nil {
		log.Printf("ERROR: %v", err)
	}
}

func run() error {
	client, err := mta.NewFeedsClient(
		&http.Client{
			Timeout: 30 * time.Second,
		},
		"lKTSZpn9bX58Nmg11rHhX1dsKaBpoFakmMSuqeh0",
		"github.com/errornil/mta:v2.0",
	)
	if err != nil {
		return errors.Wrap(err, "failed to get feed message")
	}

	msg, err := client.GetFeedMessage(mta.FeedLIRR)
	if err != nil {
		return errors.Wrap(err, "failed to get feed message")
	}

	for _, entity := range msg.GetEntity() {
		log.Println(entity.String())
	}

	return nil
}
