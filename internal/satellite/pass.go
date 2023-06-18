package satellite

import (
	"time"

	tinyTime "github.com/pmcanseco/go-sat-tracker/internal/time"
)

type Plan struct {
	Passes []Pass
}

type Pass struct {
	startTime       time.Time
	startLookAngles LookAngles

	midTime       time.Time
	midLookAngles LookAngles

	endTime       time.Time
	endLookAngles LookAngles

	FullPath      []LookAnglesAtTime
	FullPathDelta time.Duration
}

type Timing struct {
	Hour   int
	Minute int
	Second int
}

//func (p *Pass) MarshalJSON() ([]byte, error) {
//	type temp struct {
//		StartTime       time.Time          `json:"start_time"`
//		StartLookAngles LookAngles         `json:"start_look_angles"`
//		MidTime         time.Time          `json:"mid_time"`
//		MidLookAngles   LookAngles         `json:"mid_look_angles"`
//		EndTime         time.Time          `json:"end_time"`
//		EndLookAngles   LookAngles         `json:"end_look_angles"`
//		FullPath        []LookAnglesAtTime `json:"full_path"`
//	}
//	t := temp{
//		StartTime:       p.startTime,
//		StartLookAngles: p.startLookAngles,
//		MidTime:         p.midTime,
//		MidLookAngles:   p.midLookAngles,
//		EndTime:         p.endTime,
//		EndLookAngles:   p.endLookAngles,
//		FullPath:        []LookAnglesAtTime{},
//	}
//	for k, v := range p.FullPath {
//		t.FullPath[k] = v
//	}
//	return json.Marshal(t)
//}

func (p *Pass) CopyPassStartingAt(t time.Time, fullPathDelta time.Duration) Pass {
	pass := Pass{
		startTime:       t,
		startLookAngles: p.startLookAngles,
		midTime:         t.Add(p.midTime.Sub(p.startTime)),
		midLookAngles:   p.midLookAngles,
		endTime:         t.Add(p.endTime.Sub(p.startTime)),
		endLookAngles:   p.endLookAngles,
		FullPath:        []LookAnglesAtTime{},
		FullPathDelta:   p.FullPathDelta,
	}

	for i, v := range p.FullPath {
		v.Time = pass.startTime.Add(fullPathDelta * time.Duration(i))
		pass.FullPath = append(pass.FullPath, v)
	}

	return pass
}

func (p *Pass) IsTimeWithinPass(t time.Time) bool {
	return !t.Before(p.startTime)
}

func (p *Pass) GetStartTime() time.Time {
	return p.startTime
}

func (p *Pass) GetEndTime() time.Time {
	return p.endTime
}

//func (p *Pass) GetDuration() time.Duration {
//	return p.endTime.Sub(p.startTime)
//}

func (p *Pass) GetMaxElevation() int {
	return int(p.midLookAngles.ElevationDegrees)
}

func (p *Pass) GetLookAngle(d time.Duration) *LookAngles {
	//result, ok := p.FullPath[d]
	//if !ok {
	//	return nil
	//}

	for len(p.FullPath) > 0 {
		laat := p.FullPath[0]
		p.FullPath = p.FullPath[1:]

		if time.Until(tinyTime.GetTime().Add(d)) > 0 {
			return &laat.LookAngles
		}
	}

	return nil
}

//func (p *Pass) CoarsePrintPath() {
//	for i := 0; i < len(p.FullPath)-10; i += 10 {
//		fmt.Printf("%s  -- Az: %.1f   El: %.1f\n", p.FullPath[i].Time.Format(time.Layout), p.FullPath[i].AzimuthDegrees, p.FullPath[i].ElevationDegrees)
//	}
//}
