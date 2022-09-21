# mta

[![main](https://github.com/errornil/mta/actions/workflows/main.yml/badge.svg)](https://github.com/errornil/mta/actions/workflows/main.yml)

`mta` is the library that provides an interface to MTA GTFS Real-Time Data Feeds.

## Subway

Read and Accept [agreement](https://api.mta.info/#/DataFeedAgreement).

Grab the API key from https://api.mta.info/#/AccessKey.

Example:

```go
import "github.com/errornil/mta/v3"

client, err := mta.NewFeedsClient(
    &http.Client{
        Timeout: 30 * time.Second,
    },
    "53b2c13dbc574e8cb4bf964dd2a215e253b2c13d", // API Key (this is a fake one)
    "", // API Key for Bus feeds (optional)
    "github.com/errornil/mta:v3.0",
)

resp, err := client.GetFeedMessage(mta.Line123456S)
```

`resp` has type [FeedMessage](https://github.com/errornil/mta/blob/master/transit_realtime/gtfs-realtime.pb.go#L488-L506) (generated).

### ProtoBuf

MTA uses realtime-GTFS with their own [extension for subway feeds](http://datamine.mta.info/sites/all/files/pdfs/nyct-subway.proto.txt).
To re-regenerate generatated code run following command with [protoc](https://github.com/protocolbuffers/protobuf):

```bash
cd proto
protoc --go_out=../transit_realtime gtfs-realtime.proto nyct-subway.proto
```

## Bus

Request API key [here](https://register.developer.obanyc.com).

Example:

```go
import "github.com/errornil/mta/v3"

client, err := mta.NewFeedsClient(
    &http.Client{
        Timeout: 30 * time.Second,
    },
    "", // API Key for Subway feeds (optional)
    "bc574e8cb4bf964dd2a215e253b2c13d", // API Key for Bus feeds (this is a fake one)
    "github.com/errornil/mta:v3.0",
)

resp, err := client.GetFeedMessage(mta.FeedBusTripUpdates)
```

## Legal

This repository is not endorsed by, directly affiliated with, maintained, authorized, or sponsored by MTA. All product and company names are the registered trademarks of their original owners. The use of any trade name or trademark is for identification and reference purposes only and does not imply any association with the trademark owner.
