package gps

import (
	"context"
	"errors"
	"fmt"
	"time"

	tinygoGPS "tinygo.org/x/drivers/gps"
)

type SentenceGetter func() (string, error)

type GPS struct {
	read      SentenceGetter
	parser    tinygoGPS.Parser
	debugMode bool
	doDebug   func(string)
	quit      bool
	hasFix    bool
	time      time.Time
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
	gps.debugMode = true
	gps.doDebug = d
}

func (gps *GPS) debug(s string) {
	if gps.debugMode {
		gps.doDebug(s)
	}
}

func (gps *GPS) HasFix() bool {
	return gps.hasFix
}

// GetCoordinates returns the time, number of satellites, latitude (degrees), longitude (degrees), and altitude
// (meters) in that order. It returns an error  if a GPS fix has yet to be acquired. In that case, call GetFix
// first.
func (gps *GPS) GetCoordinates() (time.Time, int16, float32, float32, int32, error) {
	if !gps.hasFix {
		return time.Time{}, 0, 0, 0, 0, errors.New("no gps fix")
	}
	return gps.time, gps.numSats, gps.lat, gps.lon, gps.alt, nil
}

func (gps *GPS) GetFix(ctx context.Context) error {
	fixChan := make(chan tinygoGPS.Fix)
	go gps.doGetFix(fixChan)

	for {
		select {
		case <-ctx.Done():
			gps.quit = true // quit the get-fix goroutine
			return ctx.Err()
		case fix := <-fixChan:
			gps.hasFix = true
			gps.time = fix.Time
			gps.numSats = fix.Satellites
			gps.alt = fix.Altitude
			gps.lat = fix.Latitude
			gps.lon = fix.Longitude
			gps.quit = true // quit the get-fix goroutine
			return nil
		}
	}
}

func (gps *GPS) doGetFix(fixChan chan<- tinygoGPS.Fix) {
	for i := 0; ; i++ {
		if gps.quit {
			return
		}

		s, err := gps.read()
		if err != nil {
			// next sentence error
			gps.debug("NXT SNTC ERR")
			time.Sleep(1 * time.Second)
			continue
		}

		fix, parseErr := gps.parser.Parse(s)
		if parseErr != nil {
			// parse error
			gps.debug("PARSE ERR")
			time.Sleep(1 * time.Second)
			continue
		}

		if fix.Valid {
			gps.debug("FIX!")

			// sometimes we have a valid fix but satellites is 0, avoid that case
			if fix.Satellites != 0 {
				gps.debug("SATS>0,DONE")
				fixChan <- fix
				return
			}
		} else {
			// no fix
			gps.debug(fmt.Sprintf("NO FIX %d", i))
			time.Sleep(2 * time.Second)
		}

		time.Sleep(1 * time.Second)
	}
}
