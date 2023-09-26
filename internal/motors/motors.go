package motors

import (
	"math"
)

const DefaultDegreesPerStep = 1.8 / 8

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
	fwdRunningError float64
	rvsRunningError float64
}

func New(device Mover, startingAngle float64) *Motor {
	m := &Motor{
		motor:            device,
		currAngleDegrees: startingAngle,
	}

	if m.degreesPerStep == 0 {
		m.degreesPerStep = 1.8 / 8.0
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
		angleDegrees = 0
	}

	degreeDelta := angleDegrees - m.currAngleDegrees
	if math.Abs(degreeDelta) > 180 {
		oldDelta := degreeDelta
		degreeDelta = 360 - math.Abs(degreeDelta)
		if oldDelta > 0 {
			degreeDelta = degreeDelta * -1
		}
	}

	stepDelta := degreeDelta / m.degreesPerStep
	truncatedStepDelta := int(stepDelta)

	m.motor.Enable()
	m.motor.Move(truncatedStepDelta)
	m.motor.Disable()

	if stepDelta < 0 {
		m.rvsRunningError += math.Abs(stepDelta - float64(truncatedStepDelta))
		if math.Abs(m.rvsRunningError) > 1 {
			truncatedRunningError := int(math.Floor(math.Abs(m.rvsRunningError)))

			m.motor.Enable()
			m.motor.Move(-1 * truncatedRunningError)
			m.motor.Disable()

			m.rvsRunningError = m.rvsRunningError - float64(truncatedRunningError)
		}
	} else {
		m.fwdRunningError += stepDelta - float64(truncatedStepDelta)
		if math.Abs(m.fwdRunningError) > 1 {
			truncatedRunningError := int(math.Floor(m.fwdRunningError))

			m.motor.Enable()
			m.motor.Move(truncatedRunningError)
			m.motor.Disable()

			m.fwdRunningError = m.fwdRunningError - float64(truncatedRunningError)
		}
	}

	m.currAngleDegrees = angleDegrees
}
