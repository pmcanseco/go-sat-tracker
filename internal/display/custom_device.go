package display

import "image/color"

var _ Device = (*CustomDevice)(nil)

// CustomDevice is a customizable implementation of the Device interface.
// Set the functions to whatever you'd like to keep the complexity of the
// hardware outside of this package.
type CustomDevice struct {
	PixelSetter func(int16, int16, color.RGBA)
	Displayer   func() error
	Clearer     func()
}

func (m CustomDevice) SetPixel(x, y int16, ps PixelState) {
	m.PixelSetter(x, y, color.RGBA(ps))
}

func (m CustomDevice) Display() error {
	return m.Displayer()
}

func (m CustomDevice) Clear() {
	m.Clearer()
}
