package main

import (
	"log"
	"net/http"
	"time"

	mta "github.com/chuhlomin/mta/v2"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Printf("ERROR: %v", err)
	}
}

func run() error {
	client, err := mta.NewFeedsClient(
		"lKTSZpn9bX58Nmg11rHhX1dsKaBpoFakmMSuqeh0",
		&http.Client{
			Timeout: 30 * time.Second,
		},
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
