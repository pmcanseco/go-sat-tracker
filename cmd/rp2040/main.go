package main

import (
	"context"
	"machine"
	"time"

	gosat "github.com/pmcanseco/go-satellite"
	tinygoGPS "tinygo.org/x/drivers/gps"
	"tinygo.org/x/drivers/ssd1306"

	"github.com/pmcanseco/go-sat-tracker/internal/display"
	gpsDevice "github.com/pmcanseco/go-sat-tracker/internal/gps"
	"github.com/pmcanseco/go-sat-tracker/internal/satellite"
	"github.com/pmcanseco/go-sat-tracker/internal/steppermotor"
	tinyTime "github.com/pmcanseco/go-sat-tracker/internal/time"
	"github.com/pmcanseco/go-sat-tracker/internal/tracking"
)

//var (
//	npRed   = []color.RGBA{{255, 0, 0, 127}}
//	npGreen = []color.RGBA{{0, 255, 0, 127}}
//	npBlue  = []color.RGBA{{0, 0, 255, 127}}
//)

func main() {
	//rgb := machine.GPIO16
	//rgb.Configure(machine.PinConfig{Mode: machine.PinOutput})
	//neopixel := ws2812.New(machine.GPIO16)
	//_ = neopixel.WriteColors([]color.RGBA{{0, 0, 0, 255}})

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()
	time.Sleep(500 * time.Millisecond)
	led.Low()
	time.Sleep(500 * time.Millisecond)
	led.High()
	time.Sleep(500 * time.Millisecond)
	led.Low()
	time.Sleep(500 * time.Millisecond)
	led.High()
	time.Sleep(500 * time.Millisecond)

	println("BOOTIN")
	for i := 0; i < 5; i++ {
		print(".")
		time.Sleep(3 * time.Second)
	}

	println("\n\n\nhello world!!!")

	//defer func() {
	//	if r := recover(); r != nil {
	//		println("panicked!")
	//	}
	//}()

	time.Sleep(1 * time.Second)

	sat := satellite.NewSatellite(
		"1 25544U 98067A   23071.22950734  .00021411  00000-0  39277-3 0  9995",
		"2 25544  51.6409  88.8414 0005771  75.2083  23.5161 15.49204123386753",
		gosat.GravityWGS84)
	println("MADE SATELLITE!")

	//i2c := machine.I2C1
	//err := i2c.Configure(machine.I2CConfig{
	//	Frequency: 400 * machine.KHz,
	//	SDA:       machine.I2C1_SDA_PIN,
	//	SCL:       machine.I2C1_SCL_PIN,
	//})
	//if err != nil {
	//	_ = neopixel.WriteColors(npRed)
	//}

	//oled := ssd1306.NewI2C(i2c)
	//oled.Configure(ssd1306.Config{
	//	Address: 0x3C,
	//	Width:   128,
	//	Height:  32,
	//})
	//oled.ClearDisplay()

	//d := getDevice(oled)
	//wipeAnimation := display.NewWipeAnimation(d)
	//wipeAnimation.Run()

	//screen := display.NewFontDisplay(d, 128, 32, display.Consolas7pt)

	//screen.Print("IDLE")
	//time.Sleep(1 * time.Second)
	//screen.Print("WAITPASS")
	//time.Sleep(1 * time.Second)
	//screen.Print("TRACK AZ:150 EL:30")
	//time.Sleep(1 * time.Second)
	//screen.Print("TRACKING COMPLETE")
	//time.Sleep(1 * time.Second)

	_ = machine.UART1.Configure(machine.UARTConfig{
		BaudRate: 9600,
		RX:       machine.GPIO9, // WIRING - WHITE GPS WIRE GOES HERE
	})

	ublox := tinygoGPS.NewUART(machine.UART1)
	gps := gpsDevice.New(func() (string, error) {
		s, err := ublox.NextSentence()
		//println(s)
		return s, err
	})

	println("GETTING FIX")

	//gps.SetDebug(screen.Print)
	gps.GetFix()
	println("GOT FIX!")
	gps.SetDebug(nil)
	now, _, lat, lon, alt := gps.GetCoordinates()
	println("GOT COORDS! Lat: ", lat, "\tLon: ", lon, "\t Alt: ", alt)
	for alt == 0 && now.Year() != 1 {
		now, _, lat, lon, alt = gps.GetCoordinates()
		println("GOT COORDS! Lat: ", lat, "\tLon: ", lon, "\t Alt: ", alt)
		time.Sleep(300 * time.Millisecond)
	}
	now = gps.Time()
	gps.Quit()
	time.Sleep(2 * time.Second)
	tinyTime.SetTime(now)
	println("SET TIME! ", tinyTime.GetTime().Format(time.RFC822))
	//elMotor := motors.New(getElevationMotor())
	//println("MADE ELEV MOTOR!")

	tracker := tracking.NewTracker(
		sat,
		satellite.Coordinates{
			LatitudeDegrees:  float64(lat),
			LongitudeDegrees: float64(lon),
			AltitudeKM:       float64(alt / 1000),
		})
	println("MADE TRACKER!")

	tracker.Track(context.Background())
	//tracker = tracker

	println("TRACKING ENDED!")
}

func getDevice(device ssd1306.Device) display.Device {
	return display.CustomDevice{
		PixelSetter: device.SetPixel,
		Displayer:   device.Display,
		Clearer:     device.ClearDisplay,
	}
}

func printGPS(printer display.Printer, ts time.Time, numSats int16, lat, lon float32, alt int32) {
	//println(fmt.Sprintf("%s lat %.2f lon %.2f alt %d sats %d",
	//	ts.Format(time.RFC1123), lat, lon, alt, numSats))
	//printer.PrintAt(
	//	0,
	//	fmt.Sprintf("%s SATS:%d",
	//		ts.Format("01/02 15:04:05"),
	//		numSats),
	//	false)
	//printer.PrintAt(
	//	1,
	//	fmt.Sprintf("LON:%.3f",
	//		lon),
	//	false)
	//printer.PrintAt(
	//	2,
	//	fmt.Sprintf("LAT:%.3f @%dM",
	//		lat,
	//		alt),
	//	false)
}

func getElevationMotor() *steppermotor.Device {
	dir := machine.GPIO10
	dir.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dir.High()

	step := machine.GPIO9
	step.Configure(machine.PinConfig{Mode: machine.PinOutput})
	step.Low()

	sleep := machine.GPIO8
	sleep.Configure(machine.PinConfig{Mode: machine.PinOutput})
	sleep.Low()

	sm := steppermotor.New(steppermotor.DeviceConfig{
		StepPin:   step,
		DirPin:    dir,
		SleepPin:  sleep,
		StepCount: 200,
		RPM:       50,
	})
	return sm
}
