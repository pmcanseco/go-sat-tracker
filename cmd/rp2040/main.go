package main

// This is the most minimal blinky example and should run almost everywhere.

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/drivers/ws2812"

	"github.com/pmcanseco/go-sat-tracker/internal/display"
)

func main() {
	np := machine.GPIO16
	np.Configure(machine.PinConfig{Mode: machine.PinOutput})
	neopixel := ws2812.New(machine.GPIO16)
	_ = neopixel.WriteColors([]color.RGBA{{0, 0, 0, 255}})

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()
	time.Sleep(1 * time.Second)
	led.Low()
	time.Sleep(1 * time.Second)
	led.High()

	//for {
	//	_ = neopixel.WriteColors([]color.RGBA{{255, 0, 0, 0}})
	//
	//	led.Low()
	//	time.Sleep(time.Millisecond * 500)
	//
	//	_ = neopixel.WriteColors([]color.RGBA{{0, 255, 0, 0}})
	//
	//	led.High()
	//	time.Sleep(time.Millisecond * 500)
	//}

	time.Sleep(1 * time.Second)
	i2c := machine.I2C1
	err := i2c.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
		SDA:       machine.I2C1_SDA_PIN,
		SCL:       machine.I2C1_SCL_PIN,
	})
	if err != nil {
		panic(err)
	}

	oled := ssd1306.NewI2C(i2c)
	oled.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  32,
	})
	oled.ClearDisplay()

	d := getDevice(oled)
	screen := display.New(d, 128, 32)

	screen.Print("IDLE")
	time.Sleep(2 * time.Second)
	screen.Print("WAITPASS")
	time.Sleep(2 * time.Second)
	screen.Print("TRACK AZ:150 EL:30")
	time.Sleep(2 * time.Second)
	screen.Print("TRACKING COMPLETE")

	for {
		time.Sleep(1 * time.Millisecond)
	}
}

func getDevice(device ssd1306.Device) display.Device {
	return display.CustomDevice{
		PixelSetter: device.SetPixel,
		Displayer:   device.Display,
		Clearer:     device.ClearDisplay,
	}
}
