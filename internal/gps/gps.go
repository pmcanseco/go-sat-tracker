package gps

import (
	"sync"
	"time"

	tinygoGPS "tinygo.org/x/drivers/gps"
)

type SentenceGetter func() (string, error)

type GPS struct {
	read      SentenceGetter
	parser    tinygoGPS.Parser
	doDebug   func(string)
	quit      bool
	isRunning bool
	hasFix    bool
	lastFix   tinygoGPS.Fix
	fixes     <-chan tinygoGPS.Fix
	time      time.Time
	timeSet   bool
	dateSet   bool
	timeMutex sync.Mutex
	numSats   int16
	lat       float32
	lon       float32
	alt       int32 // meters
}

func New(reader SentenceGetter) *GPS {
	return &GPS{
		read:   reader,
		parser: tinygoGPS.NewParser(),
	}
}

func (gps *GPS) SetDebug(d func(string)) {
	gps.doDebug = d
}

func (gps *GPS) debug(s string) {
	if gps.doDebug != nil {
		gps.doDebug(s)
	}
}

func (gps *GPS) HasFix() bool {
	return gps.hasFix
}

// GetCoordinates returns the time, number of satellites, latitude (degrees), longitude (degrees), and altitude
// (meters) in that order. It returns an error  if a GPS fix has yet to be acquired. In that case, call GetFix
// first.
func (gps *GPS) GetCoordinates() (time.Time, int16, float32, float32, int32) {
	gps.GetFix()
	return gps.lastFix.Time,
		gps.lastFix.Satellites,
		gps.lastFix.Latitude,
		gps.lastFix.Longitude,
		gps.lastFix.Altitude
}

func (gps *GPS) Time() time.Time {
	gps.GetFix()
	for !gps.dateSet || !gps.timeSet {
	}

	gps.timeMutex.Lock()
	t := gps.time
	gps.timeMutex.Unlock()
	return t
}

func (gps *GPS) FixChan() <-chan tinygoGPS.Fix {
	return gps.fixes
}

func (gps *GPS) GetFix() {
	if !gps.isRunning {
		fixChan := make(chan tinygoGPS.Fix)
		gps.fixes = fixChan
		go gps.doGetFix(fixChan)
	}

	for {
		select {
		case <-gps.fixes:
			gps.hasFix = true
			return
		}
	}
}

func (gps *GPS) doGetFix(fixChan chan<- tinygoGPS.Fix) {
	gps.isRunning = true
	for i := 0; ; i++ {
		if gps.quit {
			gps.isRunning = false
			return
		}

		s, err := gps.read()
		if err != nil {
			// next sentence error
			gps.debug("NXT SNTC ERR")
			time.Sleep(1 * time.Millisecond)
			continue
		}

		fix, parseErr := gps.parser.Parse(s)
		if parseErr != nil {
			// parse error
			gps.debug("PARSE ERR")
			time.Sleep(1 * time.Millisecond)
			continue
		}

		if fix.Valid {
			gps.debug("FIX!")

			gps.alt = fix.Altitude
			gps.lat = fix.Latitude
			gps.lon = fix.Longitude
			gps.setDate(fix)
			gps.setTime(fix)

			if fix.Satellites != 0 {
				gps.numSats = fix.Satellites
			}

			gps.lastFix = fix
			fixChan <- fix
		} else {
			// no fix
			gps.debug("NO FIX")
		}

		time.Sleep(199 * time.Millisecond)
	}
}

func (gps *GPS) setDate(fix tinygoGPS.Fix) {
	if gps.time.Year() == 1 {
		gps.timeMutex.Lock()
		gps.time = time.Date(fix.Time.Year(), fix.Time.Month(), fix.Time.Day(), 0, 0, 0, 0, time.UTC)
		gps.dateSet = true
		gps.timeMutex.Unlock()
	}
}

func (gps *GPS) setTime(fix tinygoGPS.Fix) {
	gps.timeMutex.Lock()
	gps.timeSet = true
	gps.time = time.Date(gps.time.Year(), gps.time.Month(), gps.time.Day(), fix.Time.Hour(), fix.Time.Minute(), fix.Time.Second(), 0, time.UTC)
	gps.timeMutex.Unlock()
}
