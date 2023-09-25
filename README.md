# go-sat-tracker

Welcome! This repo holds the code to use the `github.com/joshuaferrara/go-satellite` library in order to get the look
angles for satellites in an application that controls a miniature toy dish antenna with two stepper motors.

The code is in a very "experimental" phase, written quickly to see how/if things work and should not be used for
anything that matters. Feel free to fork and do whatever you'd like. Pull requests are welcome as well.

## Prerequisites
- TinyGo 0.30

## Construction
The 3D model (paid) can be found here: https://pinshape.com/items/27103-3d-printed-desktop-satellite-antenna

Some assembly will be required with glue, bearings, and machine screws. Some wiring instructions can be found in the 
main `cmd/rp2040` application.

### Electronics:
- Adafruit Feather RP2040 ([link](https://www.adafruit.com/product/4884))
- Adafruit FeatherWing OLED ([link](https://www.adafruit.com/product/2900))
- Adafruit Perma-Proto Small Mint Tin Size Breadboard PCB ([link](https://www.adafruit.com/product/1214))
- Adafruit 6-wire Slip Ring ([link](https://www.amazon.com/gp/product/B00QSHPIHE/ref=ppx_od_dt_b_asin_title_s00?ie=UTF8&psc=1))
- Pololu DRV8834 Low-Voltage Stepper Motor Driver Carrier ([link](https://www.pololu.com/product/2134))
- 2x NEMA 11 Stepper Motors ([link](https://www.amazon.com/gp/product/B00PNEPF4O/ref=ppx_yo_dt_b_asin_title_o00_s00?ie=UTF8&psc=1))
- ublox 7M-based GPS Module ([link](https://www.amazon.com/Satellite-Positioning-Antenna-Microcontroller-Receiver/dp/B09F3JHW77/ref=sr_1_6?crid=2H35U20QFXLZ7&keywords=NEO-7M+gps&qid=1695617632&s=electronics&sprefix=neo-7m+gp%2Celectronics%2C255&sr=1-6))

## Compile + Flash

```
$ tinygo [build|flash] -target feather-rp2040 cmd/rp2040/main.go
```

## Structure
The `cmd/` folder contains various apps that were quickly whipped up to test various functionalities. The main one
powering the dish is `cmd/rp2040`. 

- `cmd/getplan`
  - an application to retrieve the next satellite pass that meets the hard-coded criteria. The pass will be emitted to
  a JSON file in the same directory.
- `cmd/instantaneous-look-angles-test`
  - a super quick test app to get the look angles for a satellite at a given time, location, altitude, and TLE
- `cmd/motors`
  - used to test the two motors when the `motors` and `steppermotor` packages were being written
- `cmd/rp2040`
  - the application demo'd at Gophercon 2023, it awaits the pass defined in the hard-coded JSON, does the tracking, and
  terminates. This is essentially the main function for this repo. 
- `cmd/tracking-test`
  - a test app to help verifying the functionality of the `tracking` package as it was being written

In the `internal/` folder are various packages for which each part of the toy dish antenna's functionality are 
encapsulated. The original iteration was pretty ambitious and included self-location via hardware GPS, as well as an
OLED display screen to show the current angles, which satellite was being tracked, when the next pass was going to be,
etc. The code to do all that is still present in the `gps/` and `display/` packages but are currently unused. A lot of
the code in `display/` was written before I was aware of the TinyFont library, which I strongly recommend using instead
if the display functionality is being revived uafter upgrading to a more powerful microcontroller. 

Where you see a`*_test.go` file, you can just run `go test .` in that direcgtory, the tests have no dependencies. 