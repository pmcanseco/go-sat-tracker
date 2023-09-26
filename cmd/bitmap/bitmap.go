package main

import (
	"fmt"
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 128, 32))
	addLabel(img, 0, 13, "Pablo\nCanseco")

	//f, err := os.Create("hello-go.bmp")
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//if err := bmp.Encode(f, img); err != nil {
	//	panic(err)
	//}

	for j := 0; j < 32; j++ {
		fmt.Printf("%2d: ", j)
		for i := 0; i < 128; i++ {
			r, _, _, _ := img.At(i, j).RGBA()
			c := " "
			if r == 0 {
				c = " "
			} else {
				c = "*"
			}
			fmt.Printf("%s", c)

			if i == 127 {
				fmt.Println(" ", i)
			}
		}
	}
}
