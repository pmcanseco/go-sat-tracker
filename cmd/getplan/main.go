package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	gosat "github.com/pmcanseco/go-satellite"

	"github.com/pmcanseco/go-sat-tracker/internal/config"
	"github.com/pmcanseco/go-sat-tracker/internal/satellite"
	tinyTime "github.com/pmcanseco/go-sat-tracker/internal/time"
)

func main() {
	cfg := config.Get()
	tinyTime.SetTime(time.Now())

	sat := satellite.NewSatellite(
		"1 25544U 98067A   23169.18024747  .00013351  00000-0  24027-3 0  9995",
		"2 25544  51.6405 323.6092 0004372  52.6861 119.0308 15.50144016401949",
		gosat.GravityWGS84)

	passes := sat.Plan(satellite.Coordinates{
		LatitudeDegrees:  cfg.HomeLatitudeDeg,
		LongitudeDegrees: cfg.HomeLongitudeDeg,
		AltitudeKM:       cfg.HomeAltitudeKM,
	}, 10, 45, time.Now(), time.Now().Add(2*24*time.Hour), time.Second*2)

	plan := satellite.Plan{
		Passes: passes,
	}

	planJSON, err := json.Marshal(plan)
	if err != nil {
		panic(fmt.Errorf("failed to marshal plan: %v", err.Error()))
	}
	fileName := fmt.Sprintf("plan-%d.json", time.Now().Unix())
	writeErr := os.WriteFile(fileName, planJSON, 0777)
	if writeErr != nil {
		panic("failed to write file" + writeErr.Error())
	}
}
