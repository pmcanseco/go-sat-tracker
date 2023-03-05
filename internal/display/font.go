package display

type FontBitmaps []uint8
type FontCharInfo [2]int

// Font contains font information for Consolas 8pt
type Font struct {
	charHeight  int
	charWidth   int
	startChar   rune
	endChar     rune
	descriptors []FontCharInfo
	bitmaps     FontBitmaps
}

func (f *Font) Bitmap(c rune) []uint8 {
	bitmapOffset := f.GetBitmapOffset(c)
	return f.bitmaps[bitmapOffset : bitmapOffset+f.charHeight]
}

func (f *Font) RowN(c rune, y int) uint8 {
	bitmapOffset := f.GetBitmapOffset(c)
	i := bitmapOffset + y
	return f.bitmaps[i]
}

func (f *Font) GetBitmapOffset(c rune) int {
	offset := c - f.startChar
	return f.descriptors[offset][1]
}

func (f *Font) StringBitmap(s string) [][]uint8 {
	var out [][]uint8
	for row := 0; row < f.charHeight; row++ {
		var r []uint8
		//fmt.Printf("%d: ", row)
		for _, c := range s {
			r = append(r, f.RowN(c, row))
			//fmt.Printf("%s(0x%x) ", string(c), f.RowN(c, row))
		}
		out = append(out, r)
		//fmt.Println()
	}
	return out
}

func (f *Font) Pixels(s string) [][]bool {
	var out [][]bool

	for row := 0; row < f.charHeight; row++ {
		var r []bool
		//fmt.Printf("\n%2d-   ", row)
		for _, c := range s {
			n := f.RowN(c, row)
			var cr []bool
			//fmt.Printf("%c: 0b", c)
			//fmt.Printf("%c: ", c)
			//for b := 7; b >= 0; b-- {
			for b := 0; b < 7; b++ {
				cr = append(cr, valueOfKthBit(n, 7-b))
				//fmt.Printf("%s", btoa(valueOfKthBit(n, b)))
				//fmt.Printf("%d ", (i*8)+b)
			}
			r = append(r, cr[:f.charWidth]...)
			//fmt.Print("  ")
		}
		out = append(out, r)
	}
	return out
}

func valueOfKthBit(n uint8, k int) bool {
	//fmt.Printf("n=%d, k=%d\n", n, k)

	return (n & (1 << k)) > 0
}

func btoa(b bool) string {
	if !b {
		return " "
	}
	return "*"
}
