package tracking

import (
	"context"
	"fmt"
	"time"

	"github.com/pmcanseco/go-sat-tracker/internal/satellite"
	tinyTime "github.com/pmcanseco/go-sat-tracker/internal/time"
)

type Tracker interface {
	Track(context.Context)
}

type Track struct {
	satellite        *satellite.Satellite
	location         satellite.Coordinates
	plan             []satellite.Pass
	populatedPlanLen int
	latestPlanTime   time.Time
	mode             mode
	currentPass      satellite.Pass
}

type mode int

const (
	idle             mode = iota
	awaitingPass     mode = iota
	tracking         mode = iota
	trackingComplete mode = iota
)

const (
	timeLayout = "02 Jan 15:04:05 MST"
)

func NewTracker(sat *satellite.Satellite, observer satellite.Coordinates) Tracker {
	t := &Track{
		satellite:      sat,
		location:       observer,
		latestPlanTime: tinyTime.GetTime(),
		plan: sat.Plan(observer, 10, 45,
			tinyTime.GetTime(), tinyTime.GetTime().Add(7*24*time.Hour), time.Second),
		mode: idle,
	}

	// make a fake pass that starts in 5 seconds and put it at the beginning for easier testing
	fakePass := t.plan[0].CopyPassStartingAt(tinyTime.GetTime().Add(5*time.Second), time.Second)
	newPlan := []satellite.Pass{fakePass}
	newPlan = append(newPlan, t.plan...)
	t.plan = newPlan
	t.populatedPlanLen = len(t.plan)

	fmt.Printf("Populated %d passes:\n", t.populatedPlanLen)
	for _, p := range t.plan {
		fmt.Printf("  Start: %s, Max Elevation: %d\n",
			p.GetStartTime().Format(timeLayout),
			p.GetMaxElevation())
	}

	return t
}

// if the plan has fewer items than the last time it was populated, re-populate it with another day's worth of passes
func (t *Track) populatePlan() {
	if len(t.plan) < t.populatedPlanLen {
		t.plan = append(t.plan, t.satellite.Plan(t.location, 0, 55,
			t.latestPlanTime, t.latestPlanTime.Add(24*time.Hour), time.Second)...)

		t.latestPlanTime = t.latestPlanTime.Add(24 * time.Hour)
		t.populatedPlanLen = len(t.plan)
	}
}

// grab a pass from our plan
func (t *Track) dequeuePass() *satellite.Pass {
	if len(t.plan) > 0 {
		p := t.plan[0]
		t.plan = t.plan[1:]
		return &p
	}
	return nil
}

// Track grab the current time, figure out where we should be looking at
func (t *Track) Track(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case tick := <-ticker.C:

			switch t.mode {
			case idle:
				fmt.Printf("%s \t Mode: %d - Idle\n", tick.Format(timeLayout), t.mode)
				p := t.dequeuePass()
				if p != nil {
					t.currentPass = *p
					t.mode = awaitingPass
				}

			case awaitingPass:
				fmt.Printf("%s \t Mode: %d - Awaiting Pass: %s to %s\n",
					tick.Format(timeLayout), t.mode, t.currentPass.GetStartTime().Format(timeLayout), t.currentPass.GetEndTime().Format(timeLayout))
				d := time.Second
				if time.Until(t.currentPass.GetStartTime()) > 2*time.Minute {
					d = time.Until(t.currentPass.GetStartTime().Add(-1 * time.Minute))
				}
				fmt.Printf("waiting %s ...\n", d.String())
				time.Sleep(d)
				if t.currentPass.IsTimeWithinPass(tinyTime.GetTime()) {
					t.mode = tracking
					fmt.Println("switched to tracking mode")
				}

			case tracking:
				now := getTiming(tinyTime.GetTime())
				la := t.currentPass.GetLookAngle(time.Since(t.currentPass.GetStartTime())) //.Round(time.Second))
				if la != nil {
					fmt.Printf("Tracking - %02d:%02d:%02d \t Az: %.1f \t El: %.1f \n", now.Hour, now.Minute, now.Second, la.AzimuthDegrees, la.ElevationDegrees)
				}

				if len(t.currentPass.FullPath) == 0 || time.Since(t.currentPass.GetEndTime()) > 0 {
					t.mode = trackingComplete
					fmt.Println("switch to trackingComplete")
				}

			case trackingComplete:
				fmt.Printf("%s \t Mode: %d Tracking Complete\n", tick.Format(timeLayout), t.mode)
				if len(t.plan) < 3 {
					fmt.Println("populating plan to fetch more passes...")
					t.populatePlan()
				} else {
					fmt.Println("not populating plan as there's still 3 or more passes coming up")
				}
				t.mode = idle
			}
		}
	}
}

func getTiming(t time.Time) satellite.Timing {
	return satellite.Timing{
		Hour:   t.Hour(),
		Minute: t.Minute(),
		Second: t.Second(),
	}
}
