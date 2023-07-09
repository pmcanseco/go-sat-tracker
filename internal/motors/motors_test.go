package motors

import (
	"math"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockMover struct {
	DoMove func(steps int)
	//DoEnable  func()
	//DoDisable func()
}

func (mm *mockMover) Move(steps int) {
	mm.DoMove(steps)
}

func (mm *mockMover) Enable() {
	//mm.DoEnable()
}

func (mm *mockMover) Disable() {
	//mm.DoDisable()
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
			motor = New(mock, 0)
		)

		DescribeTable("moves the correct number of steps relative to the current position",
			func(commandedAngle float64, expectedDegreeDelta float64, expectedStepDelta int, expectedCurrentAngle float64, expectedRunningError float64) {
				//  params:
				//  - desired angle
				//  - difference between current and desired angles (positive or negative indicates direction of travel)
				//  - this delta compared to steps keeping the sign (divide previous param by StepsPerDegree, which defaults to 1.8)
				//  - expected current angle after the operation completes, should be the commanded angle mod 360 because multiple revolutions is not desired
				//     (doesn't matter for azimuth but for elevation it will damage the device)
				motor.CommandAngle(commandedAngle)
				Expect(commandedSteps).To(Equal(expectedStepDelta))
				Expect(motor.currAngleDegrees).To(Equal(expectedCurrentAngle))
				if commandedSteps < 0 {
					Expect(motor.rvsRunningError).To(BeNumerically("~", expectedRunningError, 0.001))
				} else {
					Expect(motor.fwdRunningError).To(BeNumerically("~", expectedRunningError, 0.001))
				}
				commandedSteps = 0
			},
			Entry("move to 30 from 0, fwd running error .333", 30.0, 30.0, int(math.Floor(30.0/DefaultDegreesPerStep)), 30.0, 0.333),
			Entry("move to 60 from previously 30, fwd running error .666", 60.0, 30.0, int(math.Floor(30.0/DefaultDegreesPerStep)), 60.0, 0.666),
			Entry("move to 20 from previously 60, rvs running error .777", 20.0, -40.0, -1*int(math.Floor(40.0/DefaultDegreesPerStep)), 20.0, 0.777),
			Entry("move to 360 means move to 0 from previously 20, rvs running error .666 (.777 + .888) add error step", 360.0, -20.0, -1*(int(math.Floor(20.0/DefaultDegreesPerStep))+1), 0.0, .666),
			Entry("move to 365 means move to 5 from previously 0, fwd running error .888 (.666 + .222) ", 365.0, 5.0, int(math.Floor(5.0/DefaultDegreesPerStep)), 5.0, .888),
			Entry("move to 10 from previously 5, fwd running error .111 (.888 + .222) add error step", 10.0, 5.0, int(math.Floor(5.0/DefaultDegreesPerStep))+1, 10.0, .111),
			Entry("move to 30.8 from previously 10, fwd running error .555 (.111 + .444)", 30.8, 20.8, int(math.Floor(20.8/DefaultDegreesPerStep)), 30.8, 0.555),
			Entry("move to 900 means move to 180 from previously 30.8, fwd running error .666 (.555 + .111)", 900.0, 149.2, int(math.Floor(149.2/DefaultDegreesPerStep)), 180.0, .666),
			Entry("move to 358 from previously 180, fwd running error .777 (.666 + .111)", 358.0, 178.0, int(math.Floor(178/DefaultDegreesPerStep)), 358.0, .777),
			Entry("move to 2 from previously 358, fwd running error .555 (.777 + .777) add error step", 2.0, 4.0, int(math.Floor(4.0/DefaultDegreesPerStep))+1, 2.0, .555),
		)
	})

	FContext("passing the 360/0 boundary", func() {

		var (
			commandedSteps = 0
			mock           = &mockMover{
				DoMove: func(steps int) {
					commandedSteps += steps
				},
			}
			motor *Motor
		)

		BeforeEach(func() {
			motor = New(mock, 0)
		})

		It("smoothly passes the 360 degree boundary going forward", func() {
			motor.CommandAngle(358.0)
			commandedSteps = 0

			motor.CommandAngle(2.0)
			Expect(commandedSteps).To(BeNumerically("~", int(math.Floor(4/DefaultDegreesPerStep)), 1),
				"the degree delta should be 4, not -356")
		})

		It("smoothly passes the 360 degree boundary going reverse", func() {
			motor.CommandAngle(2.0)
			commandedSteps = 0

			motor.CommandAngle(358.0)
			Expect(commandedSteps).To(BeNumerically("~", int(math.Floor(-4/DefaultDegreesPerStep)), 1),
				"the degree delta should be -4, not 356")
		})
	})

	Context("realistic angles", func() {

		var (
			commandedSteps = 0
			mock           = &mockMover{
				DoMove: func(steps int) {
					commandedSteps += steps
				},
			}
			motor = New(mock, 0)
		)

		PDescribeTable("generates rotation steps accurately with a realistic set of azimuth angle inputs",
			func(commandedAngle float64, expectedDegreeDelta float64, expectedStepDelta int, expectedCurrentAngle float64) {
				motor.CommandAngle(commandedAngle)
				Expect(commandedSteps).To(Equal(expectedStepDelta))
				Expect(motor.currAngleDegrees).To(Equal(expectedCurrentAngle))
				//fmt.Println("commanded angle: ", commandedAngle, ", commanded steps: ", commandedSteps, ", running error: ", motor.fwdRunningError)
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
