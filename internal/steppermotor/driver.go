package steppermotor

import (
	"machine"
	"time"
)

const (
	pulseLength = 2 * time.Microsecond
	wakeTime    = 1000 * time.Microsecond
)

// Device holds the pins and the delay between steps
type Device struct {
	stepPin        machine.Pin
	dirPin         machine.Pin
	sleepPin       machine.Pin
	stepDelay      time.Duration
	enableSleeping bool
	reverse        bool // invert the direction of movement
}

// DeviceConfig contains the configuration data for a single easystepper driver
type DeviceConfig struct {
	StepPin  machine.Pin
	DirPin   machine.Pin
	SleepPin machine.Pin

	// EnableSleeping determines whether to keep the motor energized between commands.
	// This is useful for the elevation motor to hold up the dish, otherwise it falls at low-ish angles.
	// This is less useful for azimuth since the assembly should hold its position (in theory).
	EnableSleeping bool

	// StepCount is the number of steps required to perform a full revolution of the stepper motor
	StepCount uint
	// RPM determines the speed of the stepper motor in 'Revolutions per Minute'
	RPM uint
}

// New returns a new stepper driver given a DeviceConfig
func New(config DeviceConfig, reverse bool) *Device {
	if config.StepCount == 0 || config.RPM == 0 {
		panic("StepCount and RPM must be > 0")
	}
	return &Device{
		stepPin:        config.StepPin,
		dirPin:         config.DirPin,
		reverse:        reverse,
		enableSleeping: config.EnableSleeping,
		sleepPin:       config.SleepPin,
		stepDelay:      time.Second * 60 / time.Duration(config.StepCount*config.RPM),
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
		if d.reverse {
			d.dirPin.Low()
		} else {
			d.dirPin.High()
		}
	} else {
		if d.reverse {
			d.dirPin.High()
		} else {
			d.dirPin.Low()
		}
	}
	time.Sleep(wakeTime)

	for s := 0; s < steps; s++ {
		d.step()
		time.Sleep(d.stepDelay)
	}
}

func (d *Device) step() {
	d.stepPin.High()
	time.Sleep(pulseLength)
	d.stepPin.Low()
}

func (d *Device) Disable() {
	if d.enableSleeping {
		d.sleepPin.Low()
	}
}

func (d *Device) Enable() {
	d.sleepPin.High()
}
