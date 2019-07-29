# mta

The library that provides an interface to MTA Real-Time Data Feeds.

### Accept [agreement](http://web.mta.info/developers/developer-data-terms.html)

## Subway (GTFS-realtime feeds)

### [Register to get API key](http://datamine.mta.info/user/register)

### Example

```go
import "github.com/chuhlomin/mta"

client := mta.NewSubwayClient(
    "53b2c13dbc574e8cb4bf964dd2a215e2", // API Key (this is a fake one)
    10*time.Seconds, // HTTP client timeout
)

resp, err := client.GetFeedMessage(mta.Line123456S)
// check err
```

`resp` has type [FeedMessage](https://github.com/chuhlomin/mta/blob/master/transit_realtime/gtfs-realtime.pb.go#L488-L506) (generated).

### ProtoBuf

MTA uses realtime-GTFS with their own extension for subway feeds.
To re-regenerate generatated code run following command with [protoc](https://github.com/protocolbuffers/protobuf):

```bash
cd proto
protoc --go_out=../transit_realtime gtfs-realtime.proto nyct-subway.proto
```

## Bus Times

### [Request API key](http://spreadsheets.google.com/viewform?hl=en&formkey=dG9kcGIxRFpSS0NhQWM4UjA0V0VkNGc6MQ#gid=0)

### Example

```go
import "github.com/chuhlomin/mta"

client := mta.NewBusTimeClient(
    "fa05aa30-3c71-4953-91c8-65b46c6e5f78", // API Key (this is a fake one)
    10*time.Seconds, // HTTP client timeout
)

resp, err := client.GetStopMonitoring(400933) // 400933 is the stop ID for "AV OF THE AMERICANS/W 34 ST" bus stop
// check err
```

`resp` has type [StopMonitoringResponse](https://github.com/chuhlomin/mta/blob/master/structs.go#L3-L5).

## Legal

This repository is not endorsed by, directly affiliated with, maintained, authorized, or sponsored by MTA. All product and company names are the registered trademarks of their original owners. The use of any trade name or trademark is for identification and reference purposes only and does not imply any association with the trademark owner.
