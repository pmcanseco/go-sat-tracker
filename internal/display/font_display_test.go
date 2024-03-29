package display

import (
	"image/color"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("font display tests", func() {
	Context("Line stack", func() {

		It("prints are pushed onto the line stack", func() {
			p := NewFontDisplay(CustomDevice{
				PixelSetter: func(i int16, i2 int16, rgba color.RGBA) {},
				Displayer:   func() error { return nil },
				Clearer:     func() {},
			}, 128, 32, Consolas7pt)
			d := p.(*FontDisplay)

			hello := "hello"
			d.Print(hello)
			Expect(d.lines[2]).To(Equal([]byte(hello)))

			worlds := "worlds"
			d.Print(worlds)
			Expect(d.lines[1]).To(Equal([]byte(hello)))
			Expect(d.lines[2]).To(Equal([]byte(worlds)))

			bye := "bye"
			d.Print(bye)
			Expect(d.lines[0]).To(Equal([]byte(hello)))
			Expect(d.lines[1]).To(Equal([]byte(worlds)))
			Expect(d.lines[2]).To(Equal([]byte(bye)))

			overflow := "overflow"
			d.Print(overflow)
			Expect(d.lines[0]).To(Equal([]byte(worlds)))
			Expect(d.lines[1]).To(Equal([]byte(bye)))
			Expect(d.lines[2]).To(Equal([]byte(overflow)))

			at := "at"
			d.PrintAt(1, at, false)
			Expect(d.lines[0]).To(Equal([]byte(worlds)))
			Expect(d.lines[1]).To(Equal([]byte(at)))
			Expect(d.lines[2]).To(Equal([]byte(overflow)))

			cleared := "cleared!"
			d.PrintAt(0, cleared, true)
			Expect(d.lines[0]).To(Equal([]byte(cleared)))
			Expect(d.lines[1]).To(Equal(make([]byte, 0, d.cols)))
			Expect(d.lines[2]).To(Equal(make([]byte, 0, d.cols)))
		})
	})
})
