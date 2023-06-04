package motors

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockMover struct {
	DoMove func(steps int)
}

func (mm *mockMover) Move(steps int) {
	mm.DoMove(steps)
}

func TestMotors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "motors suite")
}

var _ = Describe("motors test", func() {

	Context("angle commanding", func() {

		var (
			commandedSteps = 0
			mock           = &mockMover{
				DoMove: func(steps int) {
					commandedSteps += steps
				},
			}
			motor = New(mock)
		)

		DescribeTable("moves the correct number of steps relative to the current position",
			func(commandedAngle float64, expectedDegreeDelta float64, expectedStepDelta int, expectedCurrentAngle float64) {
				//  params:
				//  - desired angle
				//  - difference between current and desired angles (positive or negative indicates direction of travel)
				//  - this delta compared to steps keeping the sign (divide previous param by StepsPerDegree, which defaults to 1.8)
				//  - expected current angle after the operation completes, should be the commanded angle mod 360 because multiple revolutions is not desired
				//     (doesn't matter for azimuth but for elevation it will damage the device)

				motor.CommandAngle(commandedAngle)
				Expect(commandedSteps).To(Equal(expectedStepDelta))
				Expect(motor.currAngleDegrees).To(Equal(expectedCurrentAngle))
				commandedSteps = 0
			},
			Entry("move to 30 from 0, running error .667", 30.0, 30.0, 16, 30.0),
			Entry("move to 60 from previously 30, running error .334 "+
				"(after add a step to compensate, making step delta 17", 60.0, 30.0, 17, 60.0),
			Entry("move to 20 from previously 60, running error .112 (.334 - .222)", 20.0, -40.0, -22, 20.0),
			Entry("move to 360 means move to 0 from previously 20, running error .001 (.112 - .111)", 360.0, -20.0, -11, 0.0),
			Entry("move to 365 means move to 5 from previously 0, running error .778 (.001 + 2.777) ", 365.0, 5.0, 2, 5.0),
			Entry("move to 0 from previously 5, running error .001 (.778 - .777)", 0.0, -5.0, -2, 0.0),
			Entry("move to 30.8 from previously 0", 30.8, 30.8, 17, 30.8),
			Entry("move to 900 means move to 180 from previously 30.8", 900.0, 149.2, 82, 180.0),
			Entry("move to 181.5 from previously 180, add a step for error compensation", 181.5, 1.5, 1, 181.5),
		)
	})
})
