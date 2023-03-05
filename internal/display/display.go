package display

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	rows = 3
	cols = 18
)

var (
	off = PixelState(color.RGBA{A: 255})
	on  = PixelState(color.RGBA{R: 255, G: 255, B: 255, A: 255})
)

type PixelState color.RGBA

type Device interface {
	SetPixel(x, y int16, c PixelState)
	Display() error
	Clear()
}

type Display struct {
	width, height int
	img           draw.Image
	device        Device
	lines         [][]byte
}

func New(disp Device, width, height int) *Display {
	d := Display{
		width:  width,
		height: height,
		device: disp,
		img:    image.NewRGBA(image.Rect(0, 0, width, height)),
		lines:  getNewLines(),
	}

	return &d
}

func (d *Display) addLabel(x, y int, label string) {
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	drawer := &font.Drawer{
		Dst:  d.img,
		Src:  image.NewUniform(color.RGBA(on)),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	drawer.DrawString(label)
}

func (d *Display) Print(s string) {
	d.pushLines(s)
	d.device.Clear()
	d.update()
}

func (d *Display) PrintAt(line int, s string, clear bool) {
	if clear {
		d.clear()
	}
	d.lines[line] = []byte(s)
	d.update()
}

func (d *Display) display() {
	for j := int16(0); j < int16(d.height); j++ {
		for i := int16(0); i < int16(d.width); i++ {
			r, _, _, _ := d.img.At(int(i), int(j)).RGBA()
			if r == 0 {
				d.device.SetPixel(i, j, off)
			} else {
				d.device.SetPixel(i, j, on)
			}

		}
	}
	_ = d.device.Display()
}

func (d *Display) clear() {
	d.lines = getNewLines()
	d.device.Clear()
}

func (d *Display) update() {
	d.img = image.NewRGBA(image.Rect(0, 0, d.width, d.height))
	d.addLabel(0, 9, string(d.lines[0][:]))
	d.addLabel(0, 21, string(d.lines[1][:]))
	d.addLabel(0, 32, string(d.lines[2][:]))
	d.display()
}

func (d *Display) pushLines(s string) {
	for i := 0; i < len(d.lines)-1; i++ {
		d.lines[i] = d.lines[i+1]
	}
	d.lines[len(d.lines)-1] = []byte(s)
}

func getNewLines() [][]byte {
	var out [][]byte
	for i := 0; i < rows; i++ {
		out = append(out, make([]byte, 0, cols))
	}
	return out
}
