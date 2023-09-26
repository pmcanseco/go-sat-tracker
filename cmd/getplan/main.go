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
		"1 25544U 98067A   23185.67588067  .00125576  00000-0  21210-2 0  9990",
		"2 25544  51.6410 241.8291 0004659  96.9441  45.8138 15.50681602404509",
		gosat.GravityWGS84)

	passes := sat.Plan(satellite.Coordinates{
		LatitudeDegrees:  cfg.HomeLatitudeDeg,
		LongitudeDegrees: cfg.HomeLongitudeDeg,
		AltitudeKM:       cfg.HomeAltitudeKM,
	}, 20, 45, time.Now(), time.Now().Add(2*24*time.Hour), time.Second*1)

	plan := satellite.Plan{
		Passes: passes,
	}

	fakePass := plan.Passes[0].CopyPassStartingAt(tinyTime.GetTime().Add(10*time.Second), 1*time.Second)
	plan.Passes = []satellite.Pass{fakePass}

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
