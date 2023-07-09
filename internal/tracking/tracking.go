package tracking

import (
	"context"
	"time"

	"github.com/mailru/easyjson"

	"github.com/pmcanseco/go-sat-tracker/internal/satellite"
	tinyTime "github.com/pmcanseco/go-sat-tracker/internal/time"
)

type Angler interface {
	CommandAngle(angleDegrees float64)
}

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
	//currentPass      satellite.Pass
	azimuthMotor   Angler
	elevationMotor Angler
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
		mode:           idle,
	}

	println("MADE TRACKER! PLANNING ...")

	t.plan = sat.Plan(observer, 10, 45,
		tinyTime.GetTime(), tinyTime.GetTime().Add(2*24*time.Hour), 3*time.Second)

	generateFakePassStartingSoon(t, 5*time.Second)
	t.populatedPlanLen = len(t.plan)

	println("Populated ", t.populatedPlanLen, " passes:")
	for _, p := range t.plan {
		println("  Start: ", p.GetStartTime().Format(timeLayout), "Max Elevation: ", p.GetMaxElevation())
	}

	return t
}

// generateFakePassStartingSoon makes a fake pass at the top of the plan that starts in duration d for easier testing
func generateFakePassStartingSoon(t *Track, d time.Duration) {
	fakePass := t.plan[0].CopyPassStartingAt(tinyTime.GetTime().Add(d), 3*time.Second)
	newPlan := []satellite.Pass{fakePass}
	newPlan = append(newPlan, t.plan...)
	t.plan = newPlan
}

func NewTrackerWithPlan(planJSON []byte) *Track {
	var plan satellite.Plan
	err := easyjson.Unmarshal(planJSON, &plan)
	if err != nil {
		println("failed to unmarshal plan json")
	}

	t := &Track{
		latestPlanTime: tinyTime.GetTime(),
		mode:           idle,
		plan:           []satellite.Pass{plan.Passes[0]},
	}

	time.Sleep(2 * time.Second)

	//generateFakePassStartingSoon(t, 5*time.Second)

	t.populatedPlanLen = len(t.plan)

	return t
}

func WithElevationMotor(el Angler) func(*Track) {
	return func(t *Track) {
		t.elevationMotor = el
	}
}

func WithAzimuthMotor(az Angler) func(*Track) {
	return func(t *Track) {
		t.azimuthMotor = az
	}
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
		case <-ticker.C:

			switch t.mode {
			case idle:
				println(" \t Mode: ", t.mode, " - Idle")
				//p := t.dequeuePass()
				//if p != nil {
				//	t.currentPass = *p
				//	t.mode = awaitingPass
				//}
				//t.currentPass = t.plan[0]
				t.mode = awaitingPass

			case awaitingPass:
				s := t.plan[0].GetStartTime()
				now := tinyTime.GetTime()
				println(now.Format(timeLayout), " \t Mode: ", t.mode, " - Awaiting Pass: ", s.Format(timeLayout)) //, " to ", t.plan[0].GetEndTime().Format(timeLayout))
				d := time.Second
				if s.Sub(now) > 2*time.Minute {
					d = s.Add(-1 * time.Minute).Sub(now)
				}
				println("waiting ", d.String(), "...")
				time.Sleep(d)
				if t.plan[0].IsTimeWithinPass(tinyTime.GetTime()) {
					t.mode = tracking
					println("switched to tracking mode")
				}

			case tracking:
				now := getTiming(tinyTime.GetTime())
				la := t.plan[0].GetLookAngle(time.Since(t.plan[0].GetStartTime())) //.Round(time.Second))
				if la != nil {
					//dualLog("Tracking - %02d:%02d:%02d \t Az: %.1f \t El: %.1f \n", now.Hour, now.Minute, now.Second, la.AzimuthDegrees, la.ElevationDegrees)
					println("tracking - ", now.Hour, ":", now.Minute, ":", now.Second, "\t Az:", la.AzimuthDegrees, "\t El:", la.ElevationDegrees)
					//println("Tracking")
					if t.azimuthMotor != nil {
						t.azimuthMotor.CommandAngle(la.AzimuthDegrees)
					}
					if t.elevationMotor != nil {
						t.elevationMotor.CommandAngle(la.ElevationDegrees)
					}
				}

				if len(t.plan[0].FullPath) == 0 || time.Since(t.plan[0].GetEndTime()) > 0 {
					t.mode = trackingComplete
					println("switch to trackingComplete")
				}

			case trackingComplete:
				//dualLog("%s \t Mode: %d Tracking Complete\n", tick.Format(timeLayout), t.mode)
				println("tracking complete")
				//if len(t.plan) < 3 {
				//	println("populating plan to fetch more passes...")
				//	t.populatePlan()
				//} else {
				//	println("not populating plan as there's still 3 or more passes coming up")
				//}
				//t.mode = idle
				t.azimuthMotor.CommandAngle(0)
				t.elevationMotor.CommandAngle(15)
				return
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

//func dualLog(s string, args ...interface{}) {
//fmt.Printf(s, args...)
//print(fmt.Sprintf(s, args...))
//}
