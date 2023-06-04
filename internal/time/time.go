package time

import "time"

// This package allows for tracking time based on the monotonic clock. When the time is set, the unix milliseconds are
// saved. When the time is retrieved, the milliseconds since the time was set is added to the previously-set time. The
// reason for this package is that this codebase is ran on a microcontroller which doesn't have a wall time. The intent
// is to set the time from a GPS module and then use the monotonic clock to retrieve the current time after that.

var (
	millisWhenTimeWasSet int64
	wallClock            time.Time
	unixMillisGetter     = func() int64 { return time.Now().UnixMilli() }
)

func SetTime(t time.Time) {
	millisWhenTimeWasSet = unixMillisGetter()
	wallClock = t
}

func GetTime() time.Time {
	elapsedMillis := unixMillisGetter() - millisWhenTimeWasSet
	return wallClock.Add(time.Millisecond * time.Duration(elapsedMillis))
}
