package dxt

import (
	"image/color"
)

func Rgb565toargb8888(packed uint16) color.RGBA {
	c := color.RGBA{A: 255}

	c.R = uint8((packed >> 11) & 0x1F)
	c.G = uint8((packed >> 5) & 0x3F)
	c.B = uint8((packed) & 0x1F)

	c.R = (c.R << 3) | (c.R >> 2)
	c.G = (c.G << 2) | (c.G >> 4)
	c.B = (c.B << 3) | (c.B >> 2)
	return c
}
