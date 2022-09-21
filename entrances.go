package mta

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

const EntrancesURL = "http://web.mta.info/developers/data/nyct/subway/StationEntrances.csv"

type EntranceType int

const (
	EntranceTypeUnknown EntranceType = iota
	EntranceTypeStair
	EntranceTypeEasement
	EntranceTypeDoor
	EntranceTypeElevator
	EntranceTypeEscalator
	EntranceTypeRamp
	EntranceTypeWalkway
)

type Staffing int

const (
	StaffingUnknown Staffing = iota
	StaffingFull
	StaffingNone
	StaffingPart
	StaffingSpecialEvent
)

type Entrance struct {
	Division         string // Specifies which division the station belongs to; either BMT, IND or IRT. Example: "BMT"
	Line             string
	StationName      string // Name of the station as shown on subway handout map. Example: "14 St-Union Sq"
	StationLatitude  float64
	StationLongitude float64
	Routes           []string     // Routes that serves the station (or complex, if applicable) during weekday service
	EntranceType     EntranceType // Describes the physical entry way indicated by the coordinates and is usually one of the following: Door or Doorway, Stair, Escalator, Elevator, Ramp or Easement (which only specifies that the entrance is within other property and does not indicate if it is a stair, etc)
	Entry            bool         // Indicates whether this 'Entrance' has turnstiles that allow entry (as opposed to exit only - see below)
	ExitOnly         bool         // Indicates if this 'Entrance' is actually an exit only location and prohibits entry
	Vending          bool
	Staffing         Staffing // Indicates the current level of staffing that is available at this entrance; FULL indicates a booth agent is available 24/7, NONE means the area is not staffed, and PART indicates that an agent is on duty at certain times but not 24/7
	StaffHours       string
	ADA              bool    // If a station is ADA accessible value is TRUE; however note this does not indicate that the actual entrance is ADA accessible
	ADANotes         string  // Any special circumstances regarding ADA accessibility at this station (or complex, if applicable)
	FreeCrossover    bool    // If a passenger can switch route directions without exiting and paying another fare at this station (note this is station level information, not entrance level)
	NorthSouthStreet string  // Primary vertical street adjacent to station entrance for wayfinding purposes
	EastWestStreet   string  // Primary horizontal street adjacent to station entrance for wayfinding purposes
	Corner           string  // Directional corner from the intersection of the vertical and horizontal streets
	Latitude         float64 // Latitude degrees of the station entrance
	Longitude        float64 // Longitude degrees of the station entrance
}

func ParseEntrancesCSV(path string) ([]Entrance, error) {
	// open file by path and pass reader to ParseCSVReader
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

	return ParseCSVReader(file)
}

func ParseCSVReader(r io.Reader) ([]Entrance, error) {
	reader := csv.NewReader(r)
	header, err := reader.Read()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read CSV")
	}

	line := 0
	results := []Entrance{}
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to read CSV")
		}

		line++

		entrance, err := parseEntrance(header, rec)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse record at line %d", line)
		}

		results = append(results, entrance)
	}

	return results, nil
}

func parseEntrance(header, rec []string) (Entrance, error) {
	entrance := Entrance{}

	entrance.Routes = []string{}

	for i, v := range rec {
		switch header[i] {
		case "Division":
			entrance.Division = v

		case "Line":
			entrance.Line = v

		case "Station_Name":
			entrance.StationName = v

		case "Station_Latitude":
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return entrance, errors.Wrapf(err, "failed to parse float (Station_Latitude, %q)", v)
			}
			entrance.StationLatitude = f

		case "Station_Longitude":
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return entrance, errors.Wrapf(err, "failed to parse float (Station_Longitude, %q)", v)
			}
			entrance.StationLongitude = f

		case "Route_1", "Route_2", "Route_3", "Route_4", "Route_5",
			"Route_6", "Route_7", "Route_8", "Route_9", "Route_10", "Route_11":
			if v != "" {
				entrance.Routes = append(entrance.Routes, v)
			}

		case "Entrance_Type":
			switch v {
			case "Stair":
				entrance.EntranceType = EntranceTypeStair
			case "Easement":
				entrance.EntranceType = EntranceTypeEasement
			case "Door":
				entrance.EntranceType = EntranceTypeDoor
			case "Elevator":
				entrance.EntranceType = EntranceTypeElevator
			case "Escalator":
				entrance.EntranceType = EntranceTypeEscalator
			case "Ramp":
				entrance.EntranceType = EntranceTypeRamp
			case "Walkway":
				entrance.EntranceType = EntranceTypeWalkway
			}

		case "Entry":
			b := false
			if v == "YES" {
				b = true
			}
			entrance.Entry = b

		case "Exit Only":
			b, err := strconv.ParseBool(v)
			if err != nil {
				return entrance, errors.Wrapf(err, "failed to parse bool (Vending, %q)", v)
			}
			entrance.ExitOnly = b

		case "Vending":
			b := false
			if v == "YES" {
				b = true
			}
			entrance.Vending = b

		case "Staffing":
			switch v {
			case "FULL":
				entrance.Staffing = StaffingFull
			case "NONE":
				entrance.Staffing = StaffingNone
			case "PART":
				entrance.Staffing = StaffingPart
			case "Spc Ev":
				entrance.Staffing = StaffingSpecialEvent
			default:
				entrance.Staffing = StaffingUnknown
			}

		case "Staff_Hours":
			entrance.StaffHours = v

		case "ADA":
			b, err := strconv.ParseBool(v)
			if err != nil {
				return entrance, errors.Wrapf(err, "failed to parse bool (ADA, %q)", v)
			}
			entrance.ADA = b

		case "ADA_Notes":
			entrance.ADANotes = v

		case "Free_Crossover":
			b, err := strconv.ParseBool(v)
			if err != nil {
				return entrance, errors.Wrapf(err, "failed to parse bool (Free_Crossover, %q)", v)
			}
			entrance.FreeCrossover = b

		case "North_South_Street":
			entrance.NorthSouthStreet = v

		case "East_West_Street":
			entrance.EastWestStreet = v

		case "Corner":
			entrance.Corner = v

		case "Latitude":
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return entrance, errors.Wrapf(err, "failed to parse float (Latitude, %q)", v)
			}
			entrance.Latitude = f

		case "Longitude":
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return entrance, errors.Wrapf(err, "failed to parse float (Longitude, %q)", v)
			}
			entrance.Longitude = f
		}
	}

	return entrance, nil
}
