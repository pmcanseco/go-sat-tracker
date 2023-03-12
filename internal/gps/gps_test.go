package gps

import (
	"context"
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

		It("gets a fix immediately", func() {
			sentenceQueue = append(sentenceQueue, "$GPGGA,210230,3855.4487,N,09446.0071,W,1,07,1.1,370.5,M,-29.5,M,,*7A")
			err := gps.GetFix(context.Background())
			Expect(err).ToNot(HaveOccurred())
		})

		It("gets a fix after a parse error", func() {
			sentenceQueue = append(sentenceQueue,
				"$GPGGA,210230,3855.4487,N,0970.5,M,-29.5,M,,*7A", // bad sentence
				"$GPGGA,210230,3855.4487,N,09446.0071,W,1,07,1.1,370.5,M,-29.5,M,,*7A")
			err := gps.GetFix(context.Background())
			Expect(err).ToNot(HaveOccurred())
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
			err := gps.GetFix(context.Background())
			Expect(err).ToNot(HaveOccurred())
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
			err := gps.GetFix(context.Background())
			Expect(err).ToNot(HaveOccurred())
		})

		It("errors out after the context is cancelled or expires", func() {
			// intentionally not populate the sentence queue so it will hang forever reading from it
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			err := gps.GetFix(ctx)
			Expect(err).To(HaveOccurred())
		})
	})
})
