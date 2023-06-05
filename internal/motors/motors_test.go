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

	Context("realistic angles", func() {

		var (
			commandedSteps = 0
			mock           = &mockMover{
				DoMove: func(steps int) {
					commandedSteps += steps
				},
			}
			motor = New(mock)
		)

		DescribeTable("generates rotation steps accurately with a realistic set of azimuth angle inputs",
			func(commandedAngle float64, expectedDegreeDelta float64, expectedStepDelta int, expectedCurrentAngle float64) {
				motor.CommandAngle(commandedAngle)
				Expect(commandedSteps).To(Equal(expectedStepDelta))
				Expect(motor.currAngleDegrees).To(Equal(expectedCurrentAngle))
				//fmt.Println("commanded angle: ", commandedAngle, ", commanded steps: ", commandedSteps, ", running error: ", motor.runningError)
				commandedSteps = 0
			},
			Entry("", 247.7, 247.7, 137, 247.7), // .611 error
			Entry("", 247.7, 0.0, 0, 247.7),
			Entry("", 247.8, 0.1, 0, 247.8), // .666
			Entry("", 247.9, 0.1, 0, 247.9), // .722
			Entry("", 247.9, 0.0, 0, 247.9),
			Entry("", 248.0, 0.1, 0, 248.0), // .777
			Entry("", 248.1, 0.1, 0, 248.1), // .833
			Entry("", 248.2, 0.1, 0, 248.2), // .888
			Entry("", 248.2, 0.0, 0, 248.2),
			Entry("", 248.3, 0.1, 0, 248.3), // .944
			Entry("", 248.4, 0.1, 1, 248.4), // 1.00000000001 or something, error step here, pattern repeats adding 0.555~ per 0.1 degree delta
			Entry("", 248.4, 0.0, 0, 248.4),
			Entry("", 248.5, 0.1, 0, 248.5),
			Entry("", 248.6, 0.1, 0, 248.6),
			Entry("", 248.7, 0.1, 0, 248.7),
			Entry("", 248.7, 0.0, 0, 248.7),
			Entry("", 248.8, 0.1, 0, 248.8),
			Entry("", 248.9, 0.1, 0, 248.9),
			Entry("", 249.0, 0.1, 0, 249.0),
			Entry("", 249.1, 0.1, 0, 249.1),
			Entry("", 249.1, 0.0, 0, 249.1),
			Entry("", 249.2, 0.1, 0, 249.2),
			Entry("", 249.3, 0.1, 0, 249.3),
			Entry("", 249.4, 0.1, 0, 249.4),
			Entry("", 249.5, 0.1, 0, 249.5),
			Entry("", 249.6, 0.1, 0, 249.6),
			Entry("", 249.6, 0.0, 0, 249.6),
			Entry("", 249.7, 0.1, 0, 249.7),
			Entry("", 249.8, 0.1, 0, 249.8),
			Entry("", 249.9, 0.1, 0, 249.9),
			Entry("", 250.0, 0.1, 0, 250.0),
			Entry("", 250.1, 0.1, 0, 250.1),
			Entry("", 250.2, 0.1, 1, 250.2),
		)
	})
})
