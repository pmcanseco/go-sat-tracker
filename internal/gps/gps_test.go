package gps

import (
	"errors"
	"testing"

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

		PIt("parses rmc sentences in order to get the date", func() {
			sentenceQueue = append(sentenceQueue,
				"$GPRMC,201640.00,A,3933.29994,N,10448.37680,W,0.564,,120323,,,A*63",
				"$GPGGA,210230,3855.4487,N,09446.0071,W,1,07,1.1,370.5,M,-29.5,M,,*7A")
			gps.GetFix()
			Expect(gps.lastFix.Time.Month()).To(Equal(3))
			Expect(gps.lastFix.Time.Day()).To(Equal(12))
			Expect(gps.lastFix.Time.Year()).To(Equal(2023))
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
