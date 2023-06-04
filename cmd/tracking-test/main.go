package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	gosat "github.com/pmcanseco/go-satellite"

	"github.com/pmcanseco/go-sat-tracker/internal/config"
	"github.com/pmcanseco/go-sat-tracker/internal/satellite"
	tinyTime "github.com/pmcanseco/go-sat-tracker/internal/time"
	"github.com/pmcanseco/go-sat-tracker/internal/tracking"
)

func main() {
	cfg := config.Get()
	tinyTime.SetTime(time.Now())

	sat := satellite.NewSatellite(
		"1 25544U 98067A   23144.86253841  .00014539  00000-0  26082-3 0  9994",
		"2 25544  51.6411  84.1532 0005431  16.7594  53.2286 15.50149056398178",
		gosat.GravityWGS84)
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
