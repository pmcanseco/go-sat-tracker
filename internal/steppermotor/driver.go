package steppermotor

import (
	"machine"
	"time"
)

// Device holds the pins and the delay between steps
type Device struct {
	stepPin   machine.Pin
	dirPin    machine.Pin
	sleepPin  machine.Pin
	stepDelay time.Duration
}

// DeviceConfig contains the configuration data for a single easystepper driver
type DeviceConfig struct {
	StepPin  machine.Pin
	DirPin   machine.Pin
	SleepPin machine.Pin

	// StepCount is the number of steps required to perform a full revolution of the stepper motor
	StepCount uint
	// RPM determines the speed of the stepper motor in 'Revolutions per Minute'
	RPM uint
}

// New returns a new stepper driver given a DeviceConfig
func New(config DeviceConfig) *Device {
	if config.StepCount == 0 || config.RPM == 0 {
		panic("config.StepCount and config.RPM must be > 0")
	}
	return &Device{
		stepPin:   config.StepPin,
		dirPin:    config.DirPin,
		sleepPin:  config.SleepPin,
		stepDelay: time.Second * 60 / time.Duration(config.StepCount*config.RPM),
	}
}

// Move rotates the motor the number of given steps
// (negative steps will rotate it the opposite direction)
func (d *Device) Move(steps int) {
	direction := steps > 0
	if steps < 0 {
		steps = -steps
	}

	if direction {
		d.dirPin.High()
	} else {
		d.dirPin.Low()
	}
	time.Sleep(1000 * time.Microsecond)

	for s := 0; s < steps; s++ {
		d.step()
		time.Sleep(d.stepDelay)
	}
}

func (d *Device) step() {
	d.stepPin.High()
	time.Sleep(2 * time.Microsecond)
	d.stepPin.Low()
}

func (d *Device) Disable() {
	d.sleepPin.Low()
}

func (d *Device) Enable() {
	d.sleepPin.High()
}
