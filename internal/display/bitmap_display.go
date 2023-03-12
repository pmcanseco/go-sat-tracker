package display

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var (
	off = PixelState(color.RGBA{A: 255})
	on  = PixelState(color.RGBA{R: 255, G: 255, B: 255, A: 255})
)

type PixelState color.RGBA

type BitmapDisplay struct {
	width, height int
	rows, cols    int
	img           draw.Image
	device        Device
	lines         [][]byte
}

// compile-time interface satisfaction check
var _ Printer = (*BitmapDisplay)(nil)

func New(disp Device, width, height int) *BitmapDisplay {
	d := BitmapDisplay{
		width:  width,
		height: height,
		rows:   height / 10, // font height
		cols:   width / 7,   // font width
		device: disp,
		img:    image.NewRGBA(image.Rect(0, 0, width, height)),
	}
	d.lines = d.getNewLines()

	return &d
}

func (d *BitmapDisplay) addLabel(x, y int, label string) {
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	drawer := &font.Drawer{
		Dst:  d.img,
		Src:  image.NewUniform(color.RGBA(on)),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	drawer.DrawString(label)
}

func (d *BitmapDisplay) Print(s string) {
	d.pushLines(s)
	d.device.Clear()
	d.update()
}

func (d *BitmapDisplay) PrintAt(line int, s string, clear bool) {
	if clear {
		d.clear()
	}
	d.lines[line] = []byte(s)
	d.update()
}

func (d *BitmapDisplay) display() {
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

func (d *BitmapDisplay) clear() {
	d.lines = d.getNewLines()
	d.device.Clear()
}

func (d *BitmapDisplay) update() {
	d.img = image.NewRGBA(image.Rect(0, 0, d.width, d.height))
	for y, i := d.height, len(d.lines)-1; y > 0; y, i = y-(10+i), i-1 { // 10 = font height
		d.addLabel(0, y, string(d.lines[i][:]))
	}
	d.display()
}

func (d *BitmapDisplay) pushLines(s string) {
	for i := 0; i < len(d.lines)-1; i++ {
		d.lines[i] = d.lines[i+1]
	}
	d.lines[len(d.lines)-1] = []byte(s)
}

func (d *BitmapDisplay) getNewLines() [][]byte {
	var out [][]byte
	for i := 0; i < d.rows; i++ {
		out = append(out, make([]byte, 0, d.cols))
	}
	return out
}
