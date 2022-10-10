package main

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/pkg/errors"

	"github.com/errornil/mta/v3"
)

func main() {
	if err := run(); err != nil {
		log.Printf("ERROR: %v", err)
	}
}

func run() error {
	path := flag.String("path", "", "Path to CSV file")
	flag.Parse()
	if *path == "" {
		return errors.New("missing path to CSV file, pass it with -path flag")
	}

	entrances, err := mta.ParseEntrancesCSV(*path)
	if err != nil {
		return err
	}

	for _, entrance := range entrances {
		// json encode entrance
		b, _ := json.MarshalIndent(entrance, "", "  ")
		log.Printf("%s", b)
		// break
	}

	return nil
}
