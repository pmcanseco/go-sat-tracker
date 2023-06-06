package gps

import (
	"errors"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSatellite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GPS Suite")
}

var _ = Describe("gps tests", func() {

	Context("sentence parsing", func() {
		var (
			sentenceQueue  []string
			sentenceGetErr error
			gps            *GPS
			mockReader     SentenceGetter
		)

		BeforeEach(func() {
			sentenceQueue = []string{}
			mockReader = func() (string, error) {
				out := sentenceQueue[0]
				sentenceQueue = sentenceQueue[1:]
				return out, sentenceGetErr
			}
			gps = New(mockReader)
		})

		AfterEach(func() {
			gps.quit = true
		})

		It("gets a fix immediately", func() {
			sentenceQueue = append(sentenceQueue, "$GPGGA,210230,3855.4487,N,09446.0071,W,1,07,1.1,370.5,M,-29.5,M,,*7A")
			gps.GetFix()
		})

		It("gets a fix after a parse error", func() {
			sentenceQueue = append(sentenceQueue,
				"$GPGGA,210230,3855.4487,N,0970.5,M,-29.5,M,,*7A", // bad sentence
				"$GPGGA,210230,3855.4487,N,09446.0071,W,1,07,1.1,370.5,M,-29.5,M,,*7A")
			gps.GetFix()
		})

		It("gets a fix after a next sentence (read) error", func() {
			sentenceQueue = append(sentenceQueue, "", "", "$GPGGA,210230,3855.4487,N,09446.0071,W,1,07,1.1,370.5,M,-29.5,M,,*7A")
			errQueue := []error{errors.New("boom"), errors.New("boom2"), nil}
			mockReader = func() (string, error) {
				outS := sentenceQueue[0]
				sentenceQueue = sentenceQueue[1:]
				err := errQueue[0]
				errQueue = errQueue[1:]
				return outS, err
			}
			gps = New(mockReader)
			gps.GetFix()
		})

		It("parses rmc sentences in order to get the date", func() {
			sentenceQueue = append(sentenceQueue,
				"$GPRMC,201640.00,A,3933.29994,N,10448.37680,W,0.564,,120323,,,A*63",
				"$GPGGA,210230,3855.4487,N,09446.0071,W,1,07,1.1,370.5,M,-29.5,M,,*7A")

			t := gps.Time()

			// from RMC
			Expect(t.Month()).To(Equal(time.Month(3)))
			Expect(t.Day()).To(Equal(12))
			Expect(t.Year()).To(Equal(2023))

			t = gps.Time()

			// date wasn't wiped between checks
			Expect(t.Month()).To(Equal(time.Month(3)))
			Expect(t.Day()).To(Equal(12))
			Expect(t.Year()).To(Equal(2023))

			Expect(t.Hour()).To(Equal(20))
			Expect(t.Minute()).To(Equal(16))
			Expect(t.Second()).To(Equal(40))
		})

		It("gets a fix after seeing a sentence with no fix first", func() {
			sentenceQueue = append(sentenceQueue,
				"$GPGGA,210230,3855.4487,N,09446.0071,W,0,0,1.1,370.5,M,-29.5,M,,*7A",  // no fix
				"$GPGGA,210230,3855.4487,N,09446.0071,W,1,07,1.1,370.5,M,-29.5,M,,*7A") // fix
			errQueue := []error{nil, nil}
			mockReader = func() (string, error) {
				outS := sentenceQueue[0]
				sentenceQueue = sentenceQueue[1:]
				err := errQueue[0]
				errQueue = errQueue[1:]
				return outS, err
			}
			gps = New(mockReader)
			gps.GetFix()
		})
	})
})
