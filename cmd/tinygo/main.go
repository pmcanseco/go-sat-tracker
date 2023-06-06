package main

import (
	"context"
	"fmt"
	"image/color"
	"machine"
	"time"

	"github.com/pmcanseco/go-sat-tracker/internal/satellite"
	"github.com/pmcanseco/go-sat-tracker/internal/tracking"

	gosat "github.com/pmcanseco/go-satellite"
	"tinygo.org/x/drivers/ssd1306"
)

func main() {
	sat := satellite.NewSatellite(
		"1 25544U 98067A   23032.08288244  .00011898  00000-0  21365-3 0  9993",
		"2 25544  51.6434 282.6761 0004766 300.0617 145.9076 15.50324176380688",
		gosat.GravityWGS84)

	tracker := tracking.NewTracker(
		sat,
		satellite.Coordinates{
			// todo - fill this in from gps
			LatitudeDegrees:  39.0,
			LongitudeDegrees: -104.0,
			AltitudeKM:       1.77,
		})

	fmt.Printf("hello! I have a satellite! %+v", sat)

	// todo - see if you can pass the look angles through a channel and if it still compiles with tinygo
	go tracker.Track(context.Background())

	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
		SDA:       machine.I2C0_SDA_PIN,
		SCL:       machine.I2C0_SCL_PIN,
	})
	if err != nil {
		panic(err)
	}

	display := ssd1306.NewI2C(i2c)

	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  32,
	})

	display.ClearDisplay()

	x := int16(0)
	y := int16(0)
	deltaX := int16(1)
	deltaY := int16(1)
	for {
		pixel := display.GetPixel(x, y)
		c := color.RGBA{R: 255, G: 255, B: 255, A: 255}
		if pixel {
			c = color.RGBA{A: 255}
		}
		display.SetPixel(x, y, c)
		_ = display.Display() // implementation never returns error

		x += deltaX
		y += deltaY

		if x == 0 || x == 127 {
			deltaX = -deltaX
		}

		if y == 0 || y == 31 {
			deltaY = -deltaY
		}
		time.Sleep(1 * time.Millisecond)
	}

	// todo - get the look angles here and configure the antenna to use them.
	for {
		//select {
		//case <-sigs:
		//	fmt.Println("Captured OS signal...")
		//	cancel()
		//case <-ctx.Done():
		//	fmt.Println("Bye!")
		//	cancel()
		//	os.Exit(0)
		//}
	}
}
