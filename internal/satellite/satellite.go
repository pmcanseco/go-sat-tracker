package satellite

import (
	"fmt"
	"math"
	"os"
	"time"

	gosat "github.com/pmcanseco/go-satellite"
)

type Satellite struct {
	noradID uint64
	sat     gosat.Satellite
}

type Coordinates struct {
	LatitudeDegrees  float64
	LongitudeDegrees float64
	AltitudeKM       float64
}

type LookAngles struct {
	AzimuthDegrees   float64 `json:"az"`
	ElevationDegrees float64 `json:"el"`
}

type LookAnglesAtTime struct {
	LookAngles
	Time time.Time
}

func NewSatellite(line1, line2 string, gravityConstant gosat.Gravity) *Satellite {
	sat, err := gosat.TLEToSat(
		line1,
		line2,
		gravityConstant)

	if err != nil {
		panic(err)
	}

	return &Satellite{sat: *sat}
}

func NewSatelliteFromID(noradID uint64) *Satellite {
	sat := getSatelliteFromSpaceTrack(noradID)
	fmt.Printf("Fetched TLE:\n  %s\n  %s\n", sat.sat.Line1, sat.sat.Line2)
	return &sat
}

func getSatelliteFromSpaceTrack(noradID uint64) Satellite {
	username := os.Getenv("SPACETRACK_USERNAME")
	password := os.Getenv("SPACETRACK_PASSWORD")
	spacetrack := gosat.NewSpacetrack(username, password)

	sat, err := spacetrack.GetTLE(noradID, time.Now(), gosat.GravityWGS84)
	if err != nil {
		panic(err)
	}

	return Satellite{
		noradID: noradID,
		sat:     sat,
	}
}

// getTime returns the individual components of a time.Time for use in the go-satellite library
func getTime(t time.Time) (int, int, int, int, int, int) {
	return t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second()
}

// deg2rad takes in degrees and returns the same in radians
func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// rad2deg takes in radians and returns the same in degrees
func rad2deg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

func (s *Satellite) GetLookAnglesAt(observer Coordinates, time time.Time) LookAngles {
	year, month, day, hour, minute, second := getTime(time.UTC())

	// get the satellite position
	position, _ := gosat.Propagate(s.sat, year, month, day, hour, minute, second)

	// declare my current location, altitude
	location := gosat.LatLong{
		Latitude:  deg2rad(observer.LatitudeDegrees),
		Longitude: deg2rad(observer.LongitudeDegrees),
	}

	// get my observation angles in radian
	obs := gosat.ECIToLookAngles(
		position, location, observer.AltitudeKM,
		// get Julian date
		gosat.JDay(
			year, month, day,
			hour, minute, second,
		),
	)

	return LookAngles{
		AzimuthDegrees:   rad2deg(obs.Az),
		ElevationDegrees: rad2deg(obs.El),
	}
}

func (s *Satellite) Plan(location Coordinates, minMinElevation, minMaxElevation float64, startTime, endTime time.Time, delta time.Duration) []Pass {
	currTime := startTime

	var (
		passes   []Pass
		inPass   bool
		currPass = Pass{FullPath: []LookAnglesAtTime{}}
	)

	for currTime.Before(endTime) {
		lookAngles := s.GetLookAnglesAt(location, currTime)

		if !inPass &&
			lookAngles.ElevationDegrees > minMinElevation {

			//fmt.Println("starting pass!")

			inPass = true
			currPass = Pass{
				startTime:       currTime,
				startLookAngles: lookAngles,
				midTime:         currTime,
				midLookAngles:   lookAngles,
				FullPath:        []LookAnglesAtTime{},
				FullPathDelta:   delta,
			}
		}

		if inPass {
			currPass.FullPath = append(currPass.FullPath, LookAnglesAtTime{
				LookAngles: lookAngles,
				Time:       currTime,
			})

			if currPass.midLookAngles.ElevationDegrees < lookAngles.ElevationDegrees {
				//fmt.Printf("found midpoint at %s, El: %.2f\n", currTime.Format(time.Stamp), lookAngles.ElevationDegrees)

				currPass.midLookAngles = lookAngles
				currPass.midTime = currTime
			}
		}

		if inPass &&
			lookAngles.ElevationDegrees < minMinElevation {
			//fmt.Println("ending pass")
			inPass = false

			if currPass.midLookAngles.ElevationDegrees > minMaxElevation {
				//fmt.Printf("pass meets minMaxElevation %.0f, saving to pass list\n", minMaxElevation)

				//fmt.Printf("Start: %s, %.0f, %.0f \t Max: %s, %.0f, %.0f \t End: %s, %.0f, %.0f\n",
				//	currPass.startTime.Format(time.RFC822), currPass.startLookAngles.AzimuthDegrees, currPass.startLookAngles.ElevationDegrees,
				//	currPass.midTime.Format(time.RFC822), currPass.midLookAngles.AzimuthDegrees, currPass.midLookAngles.ElevationDegrees,
				//	currTime.Format(time.RFC822), lookAngles.AzimuthDegrees, lookAngles.ElevationDegrees)

				currPass.endTime = currTime
				currPass.endLookAngles = lookAngles
				passes = append(passes, currPass)
			}
		}

		currTime = currTime.Add(delta)
	}

	return passes
}
