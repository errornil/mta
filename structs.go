package mta

type StopMonitoringResponse struct {
	Siri Siri
}

type Siri struct {
	ServiceDelivery ServiceDelivery
}

type ServiceDelivery struct {
	// The timestamp on the MTA Bus Time server at the time the request was fulfilled
	ResponseTimestamp string

	// SIRI container for VehicleMonitoring response data
	StopMonitoringDelivery []StopMonitoringDelivery
}

type StopMonitoringDelivery struct {
	// Required by the SIRI spec
	ResponseTimestamp string

	// The time until which the response data is valid until
	ValidUntil string

	// SIRI container for data about a particular vehicle service the selected stop
	MonitoredStopVisit []MonitoredStopVisit
}

type MonitoredStopVisit struct {
	// The timestamp of the last real-time update from the particular vehicle
	RecordedAtTime string

	// A complete MonitoredVehicleJourney element
	MonitoredVehicleJourney MonitoredVehicleJourney
}

type MonitoredVehicleJourney struct {
	// The 'fully qualified' route name (GTFS agency ID + route ID) for the trip the vehicle is serving.
	// Not intended to be customer-facing
	LineRef string

	// The GTFS direction for the trip the vehicle is serving
	DirectionRef string

	// A compound element uniquely identifying the trip the vehicle is serving
	FramedVehicleJourneyRef FramedVehicleJourneyRef

	// The GTFS Shape_ID, prefixed by GTFS Agency ID
	JourneyPatternRef string

	// The GTFS route_short_name
	PublishedLineName []string

	// GTFS Agency_ID
	OperatorRef string

	// The GTFS stop ID for the first stop on the trip the vehicle is serving, prefixed by Agency ID
	OriginRef string

	// The GTFS stop ID for the last stop on the trip the vehicle is serving, prefixed by Agency ID
	DestinationRef string

	// The GTFS trip_headsign for the trip the vehicle is serving
	DestinationName []string

	// If a bus has not yet departed, OriginAimedDepartureTime indicates the scheduled departure time of that bus from that terminal in ISO8601 format
	OriginAimedDepartureTime string

	// SituationRef, present only if there is an active service alert covering this call
	SituationRef []SituationRef

	// Always true
	Monitored bool

	// The most recently recorded or inferred coordinates of this vehicle
	VehicleLocation *VehicleLocation

	// Vehicle bearing: 0 is East, increments counter-clockwise
	Bearing float64

	// Indicator of whether the bus is making progress (i.e. moving, generally),
	// not moving (with value noProgress),
	// laying over before beginning a trip (value layover),
	// or serving a trip prior to one which will arrive (prevTrip).
	ProgressRate string

	// Optional indicator of vehicle progress status.
	// Set to "layover" when the bus is in a layover waiting for its next trip to start at a terminal,
	// and/or "prevTrip" when the bus is currently serving the previous trip
	// and the information presented 'wraps around' to the following scheduled trip
	ProgressStatus []string

	// Optional indicator of whether the bus occupancy is deemed to be "full", "seatsAvailable" or "standingAvailable".
	// If bus occupancy information is not available, this indicator is not shown (aka hidden.)
	Occupancy string

	// The vehicle ID, preceded by the GTFS agency ID
	VehicleRef string

	// Depending on the system's level of confidence, the GTFS block_id the bus is serving.
	// Please see "Transparency of Block vs. Trip-Level Assignment" section below
	BlockRef string

	// Call data about a particular stop
	// In StopMonitoring, it is the stop of interest;
	// in VehicleMonitoring it is the next stop the bus will make.
	MonitoredCall *MonitoredCall

	// The collection of calls that a vehicle is going to make
	OnwardCalls *OnwardCalls
}

type FramedVehicleJourneyRef struct {
	// The GTFS service date for the trip the vehicle is serving
	DataFrameRef string

	// The GTFS trip ID for trip the vehicle is serving, prefixed by the GTFS agency ID
	DatedVehicleJourneyRef string
}

type OnwardCalls struct {
	OnwardCalls []OnwardCall
}

type SituationRef struct {
	// SituationRef, present only if there is an active service alert covering this call
	SituationSimpleRef string
}

type VehicleLocation struct {
	Longitude float64
	Latitude  float64
}

type MonitoredCall struct {
	// The GTFS stop ID of the stop prefixed by agency_id
	StopPointRef string

	// The ordinal value of the visit of this vehicle to this stop, always 1 in this implementation
	VisitNumber int

	// Predicted arrival times in ISO8601 format
	ExpectedArrivalTime string

	// Predicted departure times in ISO8601 format
	ExpectedDepartureTime string

	// SIRI container for extensions to the standard
	Extensions Extensions

	ArrivalProximityText string

	DistanceFromStop int

	NumberOfStopsAway int

	StopPointName []string
}

type Extensions struct {
	// The MTA Bus Time extensions to show distance of the vehicle from the stop
	Distances Distances
}

type Distances struct {
	// The distance of the stop from the beginning of the trip/route
	CallDistanceAlongRoute float64

	// The distance from the vehicle to the stop along the route, in meters
	DistanceFromCall float64

	// The distance displayed in the UI, see below for an additional note
	PresentableDistance string

	// The number of stops on the vehicle's current trip until the stop in question, starting from 0
	StopsFromCall int
}

type OnwardCall struct {
	// The GTFS stop ID of the stop
	StopPointRef string

	// The ordinal value of the visit of this vehicle to this stop, always 1 in this implementation
	VisitNumber int

	// The GTFS stop name of the stop
	StopPointName []string

	Extensions Extensions

	ExpectedArrivalTime string

	ArrivalProximityText string

	DistanceFromStop int

	NumberOfStopsAway int
}
