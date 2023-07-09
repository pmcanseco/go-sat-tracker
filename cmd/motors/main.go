package main

import (
	"machine"
	"time"

	"github.com/pmcanseco/go-sat-tracker/internal/motors"
	"github.com/pmcanseco/go-sat-tracker/internal/steppermotor"
	"github.com/pmcanseco/go-sat-tracker/internal/tracking"
)

// wire the motor cable from the top down: (right side of driver, pot on the bottom left)
// - blue
// - red
// - green
// - black
// where blue is immediately adjacent to gnd on the top
// and there's a one-pin gap between the gnd on the bottom

const (
	ElevationSleepPin     = machine.GPIO8  // WIRING - Pin 6 white (immediately below pins 10 and 9)
	ElevationStepPin      = machine.GPIO7  // WIRING - Pin 5 green
	ElevationDirectionPin = machine.GPIO10 // WIRING - blue
	AzimuthSleepPin       = machine.GPIO26 // WIRING - A0 white
	AzimuthStepPin        = machine.GPIO27 // WIRING - A1 green
	AzimuthDirectionPin   = machine.GPIO28 // WIRING - A2 blue
)

func main() {
	//dir := machine.GPIO10
	//dir.Configure(machine.PinConfig{Mode: machine.PinOutput})
	//dir.High()
	//
	//step := machine.GPIO9
	//step.Configure(machine.PinConfig{Mode: machine.PinOutput})
	//step.Low()
	//
	//sleep := machine.GPIO8
	//sleep.Configure(machine.PinConfig{Mode: machine.PinOutput})
	//sleep.Low()
	//
	//sm := steppermotor.New(steppermotor.DeviceConfig{
	//	StepPin:   step,
	//	DirPin:    dir,
	//	SleepPin:  sleep,
	//	StepCount: 200,
	//	RPM:       40,
	//})
	//
	//for {
	//	sm.Enable()
	//	sm.Move(200)
	//	sm.Disable()
	//	time.Sleep(5 * time.Second)
	//
	//	sm.Enable()
	//	sm.Move(-200)
	//	sm.Disable()
	//	time.Sleep(5 * time.Second)
	//}

	elm := getElevationMotor()
	azm := getAzimuthMotor()

	for {
		azm.CommandAngle(340)
		elm.CommandAngle(90)
		time.Sleep(15 * time.Second)
		azm.CommandAngle(45)
		elm.CommandAngle(45)
		time.Sleep(15 * time.Second)
		azm.CommandAngle(270)
		elm.CommandAngle(15)
		time.Sleep(15 * time.Second)
		azm.CommandAngle(0)
		time.Sleep(15 * time.Second)
	}
}

func getElevationMotor() tracking.Angler {
	dir := ElevationDirectionPin
	dir.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dir.High()

	step := ElevationStepPin
	step.Configure(machine.PinConfig{Mode: machine.PinOutput})
	step.Low()

	sleep := ElevationSleepPin
	sleep.Configure(machine.PinConfig{Mode: machine.PinOutput})
	sleep.Low()

	sm := steppermotor.New(steppermotor.DeviceConfig{
		StepPin:        step,
		DirPin:         dir,
		SleepPin:       sleep,
		EnableSleeping: false,
		StepCount:      200,
		RPM:            40,
	}, false)

	return motors.New(sm, 15)
}

func getAzimuthMotor() tracking.Angler {
	dir := AzimuthDirectionPin
	dir.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dir.High()

	step := AzimuthStepPin
	step.Configure(machine.PinConfig{Mode: machine.PinOutput})
	step.Low()

	sleep := AzimuthSleepPin
	sleep.Configure(machine.PinConfig{Mode: machine.PinOutput})
	sleep.Low()

	sm := steppermotor.New(steppermotor.DeviceConfig{
		StepPin:        step,
		DirPin:         dir,
		SleepPin:       sleep,
		EnableSleeping: false,
		StepCount:      200,
		RPM:            60,
	}, false)
	return motors.New(sm, 0)
}
