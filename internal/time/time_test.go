package time

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTime(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "time suite")
}

var _ = Describe("time test", func() {

	var futureTime = time.Now().
		Add(5 * time.Minute).
		Add(10 * time.Second)

	It("sets the time", func() {
		unixMillisGetter = func() int64 {
			return 17
		}
		SetTime(futureTime)
		Expect(millisWhenTimeWasSet).To(Equal(int64(17)))
	})

	It("gets the time", func() {
		unixMillisGetter = func() int64 {
			return 9000 // force the passage of time
		}
		Expect(GetTime()).To(BeTemporally("~", futureTime.Add((9000-17)*time.Millisecond), time.Millisecond))
	})
})
