package dxt

import (
	"image"
	"image/color"
)

func parseDXT1(b *Buffer, image *image.RGBA) {
	getTables := func() ([]color.RGBA, []byte) {
		C := b.Read(UINT64).(uint64)

		colors := []color.RGBA{}
		c0 := uint16(C & 0xFFFF)
		colors = append(colors, Rgb565toargb8888(c0))
		C >>= 16

		c1 := uint16(C & 0xFFFF)
		colors = append(colors, Rgb565toargb8888(c1))
		C >>= 16

		if c0 > c1 {
			colors = append(colors, color.RGBA{A: 255, R: uint8((2*uint16(colors[0].R) + 1*uint16(colors[1].R)) / 3), G: uint8((2*uint16(colors[0].G) + 1*uint16(colors[1].G)) / 3), B: uint8((2*uint16(colors[0].B) + 1*uint16(colors[1].B)) / 3)})
			colors = append(colors, color.RGBA{A: 255, R: uint8((1*uint16(colors[0].R) + 2*uint16(colors[1].R)) / 3), G: uint8((1*uint16(colors[0].G) + 2*uint16(colors[1].G)) / 3), B: uint8((1*uint16(colors[0].B) + 2*uint16(colors[1].B)) / 3)})
		} else {
			colors = append(colors, color.RGBA{A: 255, R: uint8((uint16(colors[0].R) + uint16(colors[1].R)) / 2), G: uint8((uint16(colors[0].G) + uint16(colors[1].G)) / 2), B: uint8((uint16(colors[0].B) + uint16(colors[1].B)) / 2)})
			colors = append(colors, color.RGBA{})
		}

		clookup := make([]byte, 16)
		for i := 0; i < 16; i++ {
			clookup[i] = byte(C & 0x03)
			C >>= 2
		}

		return colors, clookup
	}

	size := image.Rect.Size()
	dx, dy := size.X/4, size.Y/4

	for iy := 0; iy < dy; iy++ { // dy
		for jx := 0; jx < dx; jx++ { // dx
			colors, clookup := getTables()
			for k := 0; k < 16; k++ {
				x := jx*4 + k%4
				y := 4*iy + k/4
				c := colors[clookup[k]]

				image.SetRGBA(x, y, c)
			}
		}
	}
}
