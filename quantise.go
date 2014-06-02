package quantise

import (
	"image"
	"image/color"
)

type Quantiser interface {
	Quantise(in image.Image) color.Palette
}

func Quantise(in image.Image, q Quantiser) image.Image {
	bounds := in.Bounds()

	out := image.NewPaletted(bounds, q.Quantise(in))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			out.Set(x, y, in.At(x, y))
		}
	}

	return out
}
