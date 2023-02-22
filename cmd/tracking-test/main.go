package main

import (
	"context"
	"fmt"
	"go-sat-tracker/internal/config"
	"go-sat-tracker/internal/satellite"
	"go-sat-tracker/internal/tracking"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Get()

	sat := satellite.NewSatelliteFromID(25544)
	//sat := satellite.NewSatellite(
	//	"1 25544U 98067A   23032.08288244  .00011898  00000-0  21365-3 0  9993",
	//	"2 25544  51.6434 282.6761 0004766 300.0617 145.9076 15.50324176380688",
	//	gosat.GravityWGS84)
	tracker := tracking.NewTracker(
		sat,
		satellite.Coordinates{
			LatitudeDegrees:  cfg.HomeLatitudeDeg,
			LongitudeDegrees: cfg.HomeLongitudeDeg,
			AltitudeKM:       cfg.HomeAltitudeKM,
		})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go tracker.Track(context.Background())

	for {
		select {
		case <-sigs:
			fmt.Println("Captured OS signal...")
			cancel()
		case <-ctx.Done():
			fmt.Println("Bye!")
			cancel()
			os.Exit(0)
		}
	}
}
