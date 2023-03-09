package blp

// adapted from https://github.com/norgannon/BLPCore

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	"math"
	"unsafe"
)

type BLPCompressionType uint8

const (
	_ BLPCompressionType = iota
	BLPCompressionTypePalettized
	BLPCompressionTypeDXT
	BLPCompressionTypePlain
)

type BGRAPixel struct {
	b, g, r, a uint8
}

type BLPHeader struct {
	magic           [4]byte
	itype           uint32
	compressionType uint8
	alphaBits       uint8
	alphaType       uint8
	hasMips         uint8
	width           uint32
	height          uint32
	mipmapOffsets   [16]uint32
	mipmapLengths   [16]uint32
	colorPalette    [256]BGRAPixel
}

func isPow2(i int) bool {
	_, rem := math.Modf(math.Log2(float64(i)))
	return rem == 0.0
}

func EncodePlainBLP(img image.Image) ([]byte, error) {
	b := img.Bounds()

	if !isPow2(b.Dx()) {
		return nil, fmt.Errorf("image width must be a power of 2")
	}

	if !isPow2(b.Dy()) {
		return nil, fmt.Errorf("image height must be a power of 2")
	}

	m := NewBGRA(image.Rect(0, 0, b.Dx(), b.Dy()))

	draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
	pixels := m.Pix

	buf := new(bytes.Buffer)

	header := BLPHeader{}
	header.magic = [4]byte{'B', 'L', 'P', '2'}
	header.itype = 1
	header.compressionType = uint8(BLPCompressionTypePlain)
	header.alphaBits = 8
	header.alphaType = 8
	header.hasMips = 0
	header.width = uint32(b.Max.X)
	header.height = uint32(b.Max.Y)
	header.mipmapOffsets[0] = uint32(unsafe.Sizeof(header))
	header.mipmapLengths[0] = uint32(len(pixels))

	for i := 1; i < 16; i++ {
		header.mipmapOffsets[i] = 0
		header.mipmapLengths[i] = 0
	}

	for i := 0; i < 256; i++ {
		blank := BGRAPixel{0, 0, 0, 0}
		header.colorPalette[i] = blank
	}

	err := binary.Write(buf, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(pixels)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
