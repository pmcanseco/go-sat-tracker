package motors

import (
	"math"
)

const DefaultDegreesPerStep = 1.8

// this package takes in two motors and makes sure they are controlled
// such that the overall device points to the given azimuth and elevation

type Mover interface {
	Move(steps int)
	Enable()
	Disable()
}

type Motor struct {
	motor            Mover
	currAngleDegrees float64
	degreesPerStep   float64

	// this tracks what we've discarded in the float to int conversion. When it's absolute value is greater than 1 it's
	// converted to an integer and included in the number of steps to command the motor to rotate, to prevent this error
	// from adding up over time and deviating the real angle of the device.
	runningError float64
}

func New(device Mover) *Motor {
	m := &Motor{
		motor: device,
	}

	if m.degreesPerStep == 0 {
		m.degreesPerStep = DefaultDegreesPerStep
	}

	if m.motor == nil {
		m.motor = device
	}

	return m
}

func (m *Motor) CommandAngle(angleDegrees float64) {
	if angleDegrees >= 360 {
		angleDegrees = math.Mod(angleDegrees, 360.0)
	}
	if angleDegrees < 0 {
		angleDegrees = 0 // negative angles don't exist
	}
	degreeDelta := angleDegrees - m.currAngleDegrees
	stepDelta := degreeDelta / m.degreesPerStep
	truncatedStepDelta := int(stepDelta)

	m.motor.Enable()
	m.motor.Move(truncatedStepDelta)
	m.motor.Disable()

	m.runningError += stepDelta - float64(truncatedStepDelta)

	if math.Abs(m.runningError) > 1 {
		truncatedRunningError := int(m.runningError)

		m.motor.Enable()
		m.motor.Move(truncatedRunningError)
		m.motor.Disable()

		m.runningError = m.runningError - float64(truncatedRunningError)
	}

	m.currAngleDegrees = angleDegrees
}
