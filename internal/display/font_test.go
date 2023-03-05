package display

import (
	"fmt"
	"image/color"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Font Tests", func() {

	Context("Consolas 8pt", func() {

		It("grabs a letter", func() {
			bitmap := Consolas8pt.Bitmap('@')
			Expect(bitmap).To(Equal([]uint8{
				// @352 '@' (7 pixels wide)
				0x00, //
				0x38, //   ###
				0x4C, //  #  ##
				0x44, //  #   #
				0x94, // #  # #
				0xB4, // # ## #
				0xB4, // # ## #
				0xB8, // # ###
				0x80, // #
				0x40, //  #
				0x78, //  ####
			}))
		})

		It("renders a string", func() {
			result := Consolas8pt.StringBitmap("hello")
			Expect(result).To(HaveLen(Consolas8pt.charHeight))
			for _, row := range result {
				Expect(row).To(HaveLen(len("hello")))
			}
			Expect(result[:2]).To(Equal([][]uint8{
				{0x00, 0x00, 0x00, 0x00, 0x00},
				{0x40, 0x00, 0x70, 0x70, 0x00},
			}))
		})

		It("renders pixels", func() {
			result := Consolas8pt.Pixels("hello")
			Expect(result).To(HaveLen(Consolas8pt.charHeight))
			for _, row := range result {
				Expect(row).To(HaveLen(Consolas8pt.charWidth * len("hello")))
			}

			// test the top two rows, there's enough there to know if we're good or not.
			Expect(result[:2]).To(Equal([][]bool{
				{
					false, false, false, false, false, false, false, // h (0 - 0b0000000)
					false, false, false, false, false, false, false, // e (0 - 0b0000000)
					false, false, false, false, false, false, false, // l (0 - 0b0000000)
					false, false, false, false, false, false, false, // l (0 - 0b0000000)
					false, false, false, false, false, false, false, // o (0 - 0b0000000)
				},
				{
					false, true, false, false, false, false, false, // h (1 - 0b0100000)
					false, false, false, false, false, false, false, // e (0 - 0b000000)
					false, true, true, true, false, false, false, // l (0 - 0b0111000)
					false, true, true, true, false, false, false, // l (0 - 0b0111000)
					false, false, false, false, false, false, false, // o (0 - 0b000000)
				},
			}))
		})

		It("displays on the pixel buffer", func() {
			fd := NewFontDisplay(CustomDevice{
				PixelSetter: func(i int16, i2 int16, rgba color.RGBA) {},
				Displayer:   func() error { return nil },
				Clearer:     func() {},
			}, 128, 32, Consolas8pt)

			fd.addLabel(0, 0, "hello!yy")
			fd.addLabel(0, 11, "world? f")

			for _, row := range fd.pixels {
				for _, col := range row {
					fmt.Printf("%s", btoa(col))
				}
				fmt.Println()
			}
		})

	})

	Context("Consolas 7pt", func() {

		It("grabs a letter", func() {
			bitmap := Consolas7pt.Bitmap('@')
			Expect(bitmap).To(Equal([]uint8{
				// @279 '@' (5 pixels wide)
				0x00, //
				0x30, //   ##
				0x48, //  #  #
				0xB8, // # ###
				0xB8, // # ###
				0xB8, // # ###
				0xA8, // # # #
				0x80, // #
				0x70, //  ###
			}))
		})

		It("renders pixels", func() {
			result := Consolas7pt.Pixels("hello")
			Expect(result).To(HaveLen(Consolas7pt.charHeight))
			for _, row := range result {
				Expect(row).To(HaveLen(Consolas7pt.charWidth * len("hello")))
			}

			// todo - check the booleans
		})

		It("renders a string", func() {
			result := Consolas7pt.StringBitmap("hello")
			Expect(result).To(HaveLen(Consolas7pt.charHeight))
			for _, row := range result {
				Expect(row).To(HaveLen(len("hello")))
			}
			Expect(result).To(Equal([][]uint8{
				//  h     e     l     l     o
				{0x00, 0x00, 0x00, 0x00, 0x00}, // 0 top row
				{0x40, 0x00, 0x70, 0x70, 0x00}, // 1 second row
				{0x40, 0x00, 0x10, 0x10, 0x00}, // 2 ... etc
				{0x78, 0x38, 0x10, 0x10, 0x60},
				{0x48, 0x78, 0x10, 0x10, 0x90},
				{0x48, 0x40, 0x10, 0x10, 0x90},
				{0x48, 0x38, 0x78, 0x78, 0x60},
				{0x00, 0x00, 0x00, 0x00, 0x00}, // 7 second-to-last row
				{0x00, 0x00, 0x00, 0x00, 0x00}, // 8 bottom row
			}))
		})

		It("displays on the pixel buffer", func() {
			fd := NewFontDisplay(CustomDevice{
				PixelSetter: func(i int16, i2 int16, rgba color.RGBA) {},
				Displayer:   func() error { return nil },
				Clearer:     func() {},
			}, 128, 32, Consolas7pt)

			fd.addLabel(0, 0, "why, hello there overflow test")
			fd.addLabel(0, 9, "GOLANg world!")
			fd.addLabel(0, 18, "I'm PABLO")
			fd.addLabel(0, 27, "vertical overflow")

			for _, row := range fd.pixels {
				for _, col := range row {
					fmt.Printf("%s", btoa(col))
				}
				fmt.Println()
			}
		})
	})
})
