package main

import (
	"context"
	"fmt"
	"go-sat-tracker/internal/satellite"
	"go-sat-tracker/internal/tracking"

	gosat "github.com/pmcanseco/go-satellite"
)

func main() {
	sat := satellite.NewSatellite(
		"1 25544U 98067A   23032.08288244  .00011898  00000-0  21365-3 0  9993",
		"2 25544  51.6434 282.6761 0004766 300.0617 145.9076 15.50324176380688",
		gosat.GravityWGS84)

	tracker := tracking.NewTracker(
		sat,
		satellite.Coordinates{
			// todo - fill this in from gps
			LatitudeDegrees:  39.0,
			LongitudeDegrees: -104.0,
			AltitudeKM:       1.77,
		})

	fmt.Printf("hello! I have a satellite! %+v", sat)

	// todo - see if you can pass the look angles through a channel and if it still compiles with tinygo
	go tracker.Track(context.Background())

	// todo - get the look angles here and configure the antenna to use them.
	for {
		//select {
		//case <-sigs:
		//	fmt.Println("Captured OS signal...")
		//	cancel()
		//case <-ctx.Done():
		//	fmt.Println("Bye!")
		//	cancel()
		//	os.Exit(0)
		//}
	}
}
