package main

import (
	"fmt"
	"math"
	"time"

	"github.com/pmcanseco/go-sat-tracker/internal/config"

	gosat "github.com/pmcanseco/go-satellite"
)

const (
	SatNoradID = 44747 // rando starlink satellite
	ISSNoradID = 25544 // International Space Station
)

func main() {
	cfg := config.Get()

	//sat, err := spacetrack.GetTLE(ISSNoradID, time.Now(), gosat.GravityWGS84)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("%s\n%s\n", sat.Line1, sat.Line2)

	t, err := time.Parse(time.RFC822, "31 Jan 23 19:09 MST")
	if err != nil {
		panic(err)
	}

	// declare current time
	year, month, day, hour, minute, second := getTime(t.UTC())

	// initialize satellite
	sat, _ := gosat.TLEToSat(
		"1 25544U 98067A   23030.92270735  .00012277  00000-0  22041-3 0  9994",
		"2 25544  51.6436 288.4282 0004934 297.6548 148.9640 15.50296141380504",
		gosat.GravityWGS84)

	// get the satellite position
	position, _ := gosat.Propagate(
		*sat, year, month, day,
		hour, minute, second,
	)

	// declare my current location, altitude
	location := gosat.LatLong{
		Latitude:  deg2rad(cfg.HomeLatitudeDeg),
		Longitude: deg2rad(cfg.HomeLongitudeDeg),
	}

	// get my observation angles in radian
	obs := gosat.ECIToLookAngles(
		position, location, cfg.HomeAltitudeKM,
		// get Julian date
		gosat.JDay(
			year, month, day,
			hour, minute, second,
		),
	)

	// print my observation azimuth in angle
	fmt.Printf("AzimuthDegrees:  %.2f\n", rad2deg(obs.Az))
	fmt.Printf("ElevationDegrees %.2f\n", rad2deg(obs.El))
}

func rad2deg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

func getTime(t time.Time) (int, int, int, int, int, int) {
	return t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second()
}
