package main

import (
	"machine"
	"time"

	"github.com/pmcanseco/go-sat-tracker/internal/steppermotor"
)

func main() {
	dir := machine.GPIO10
	dir.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dir.High()

	step := machine.GPIO9
	step.Configure(machine.PinConfig{Mode: machine.PinOutput})
	step.Low()

	sleep := machine.GPIO8
	sleep.Configure(machine.PinConfig{Mode: machine.PinOutput})
	sleep.Low()

	sm := steppermotor.New(steppermotor.DeviceConfig{
		StepPin:   step,
		DirPin:    dir,
		SleepPin:  sleep,
		StepCount: 200,
		RPM:       50,
	})

	for {
		sm.Enable()
		sm.Move(200)
		sm.Disable()
		time.Sleep(10 * time.Second)

		sm.Enable()
		sm.Move(-200)
		sm.Disable()
		time.Sleep(10 * time.Second)
	}
}
