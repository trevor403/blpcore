package blp

import (
	"image"
	"image/color"
	"math/bits"
)

// copied from https://github.com/mewbak/framebuffer/blob/master/bgra.go

type BGRA struct {
	Pix    []byte
	Rect   image.Rectangle
	Stride int
}

func (i *BGRA) Bounds() image.Rectangle { return i.Rect }
func (i *BGRA) ColorModel() color.Model { return color.RGBAModel }

func (i *BGRA) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(i.Rect)) {
		return color.RGBA{}
	}

	n := i.PixOffset(x, y)
	pix := i.Pix[n:]
	return color.RGBA{pix[2], pix[1], pix[0], pix[3]}
}

func (i *BGRA) Set(x, y int, c color.Color) {
	i.SetRGBA(x, y, color.RGBAModel.Convert(c).(color.RGBA))
}

func (i *BGRA) SetRGBA(x, y int, c color.RGBA) {
	if !(image.Point{x, y}.In(i.Rect)) {
		return
	}

	n := i.PixOffset(x, y)
	pix := i.Pix[n:]

	pix[0] = c.B
	pix[1] = c.G
	pix[2] = c.R
	pix[3] = c.A
}

func (i *BGRA) PixOffset(x, y int) int {
	return (y-i.Rect.Min.Y)*i.Stride + (x-i.Rect.Min.X)*4
}

// copied from image/image.go

// mul3NonNeg returns (x * y * z), unless at least one argument is negative or
// if the computation overflows the int type, in which case it returns -1.
func mul3NonNeg(x int, y int, z int) int {
	if (x < 0) || (y < 0) || (z < 0) {
		return -1
	}
	hi, lo := bits.Mul64(uint64(x), uint64(y))
	if hi != 0 {
		return -1
	}
	hi, lo = bits.Mul64(lo, uint64(z))
	if hi != 0 {
		return -1
	}
	a := int(lo)
	if (a < 0) || (uint64(a) != lo) {
		return -1
	}
	return a
}

// pixelBufferLength returns the length of the []uint8 typed Pix slice field
// for the NewXxx functions. Conceptually, this is just (bpp * width * height),
// but this function panics if at least one of those is negative or if the
// computation would overflow the int type.
//
// This panics instead of returning an error because of backwards
// compatibility. The NewXxx functions do not return an error.
func pixelBufferLength(bytesPerPixel int, r image.Rectangle, imageTypeName string) int {
	totalLength := mul3NonNeg(bytesPerPixel, r.Dx(), r.Dy())
	if totalLength < 0 {
		panic("image: New" + imageTypeName + " Rectangle has huge or negative dimensions")
	}
	return totalLength
}

// NewBGRA returns a new BGRA image with the given bounds.
func NewBGRA(r image.Rectangle) *BGRA {
	return &BGRA{
		Pix:    make([]uint8, pixelBufferLength(4, r, "BGRA")),
		Stride: 4 * r.Dx(),
		Rect:   r,
	}
}
