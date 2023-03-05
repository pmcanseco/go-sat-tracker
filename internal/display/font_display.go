package display

type FontDisplay struct {
	width, height int
	rows, cols    int
	font          Font
	device        Device
	lines         [][]byte
	pixels        [][]bool
}

func NewFontDisplay(disp Device, width, height int, font Font) *FontDisplay {
	d := FontDisplay{
		width:  width,
		height: height,
		rows:   height / font.charHeight, // font height
		cols:   width / font.charWidth,   // font width
		font:   font,
		device: disp,
	}
	d.lines = d.getNewLines()
	d.pixels = d.getPixelBuffer()

	return &d
}

func (d *FontDisplay) addLabel(x, y int, label string) {
	spix := d.font.Pixels(label)
	var px, py = x, y
	var fx, fy int

	for fy = 0; fy < len(spix); fy++ {
		if py >= d.height || py < 0 {
			break
		}

		for fx = 0; fx < len(spix[0]); fx++ {
			if px >= d.width {
				continue
			}

			d.pixels[py][px] = spix[fy][fx]
			//fmt.Printf("%d,%d<=%d,%d\n", py, px, fy, fx)
			px++
		}
		px = x
		py++
	}
}

func (d *FontDisplay) Print(s string) {
	d.pushLines(s)
	d.device.Clear()
	d.update()
}

func (d *FontDisplay) PrintAt(line int, s string, clear bool) {
	if clear {
		d.clear()
	}
	d.lines[line] = []byte(s)
	d.update()
}

func (d *FontDisplay) update() {
	d.pixels = d.getPixelBuffer()
	for y, i := 0, 0; i < len(d.lines); y, i = (i+1)*d.font.charHeight, i+1 {
		d.addLabel(0, y, string(d.lines[i][:]))
	}
	d.display()
}

func (d *FontDisplay) display() {
	for j := int16(0); j < int16(d.height); j++ {
		for i := int16(0); i < int16(d.width); i++ {
			isOn := d.pixels[j][i]
			if isOn {
				d.device.SetPixel(i, j, on)
			} else {
				d.device.SetPixel(i, j, off)
			}

		}
	}
	_ = d.device.Display()
}

func (d *FontDisplay) clear() {
	d.lines = d.getNewLines()
	d.device.Clear()
}

func (d *FontDisplay) pushLines(s string) {
	for i := 0; i < len(d.lines)-1; i++ {
		d.lines[i] = d.lines[i+1]
	}
	d.lines[len(d.lines)-1] = []byte(s)
}

func (d *FontDisplay) getNewLines() [][]byte {
	var out [][]byte
	for i := 0; i < d.rows; i++ {
		out = append(out, make([]byte, 0, d.cols))
	}
	return out
}

func (d *FontDisplay) getPixelBuffer() [][]bool {
	var out [][]bool
	for i := 0; i < d.height; i++ {
		out = append(out, make([]bool, d.width))
	}
	return out
}
