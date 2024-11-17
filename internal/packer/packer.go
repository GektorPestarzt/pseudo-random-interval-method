package packer

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

const (
	Blue = iota
	Red
	Green
)

type Packer interface {
	Set(index int, bit uint8) error
	Get(index int) (uint8, error)
	GetImage() image.Image
}

type ImagePacker struct {
	img     image.Image
	newRGBA *image.RGBA
	op      uint
}

func NewImagePacker(img image.Image, op uint) (*ImagePacker, error) {
	if op > Green { // ???
		return nil, fmt.Errorf("package type %d is not supported", op)
	}

	bounds := img.Bounds()
	newRGBA := image.NewRGBA(bounds)
	draw.Draw(newRGBA, bounds, img, bounds.Min, draw.Src)

	packer := &ImagePacker{
		op:      op,
		img:     img,
		newRGBA: newRGBA,
	}

	return packer, nil
}

func (ip *ImagePacker) Set(index int, bit uint8) error {
	if index < 0 {
		return fmt.Errorf("index out of range")
	}

	var err error = nil

	switch ip.op {
	case Blue:
		err = ip.setBlue(index, bit)
	case Red:
		err = fmt.Errorf("dont work")
	case Green:
		err = fmt.Errorf("dont work")
	}

	return err
}

func (ip *ImagePacker) setBlue(index int, bit uint8) error {
	width := ip.img.Bounds().Max.X - ip.img.Bounds().Min.X
	height := ip.img.Bounds().Max.Y - ip.img.Bounds().Min.Y

	x := ip.img.Bounds().Min.X + index%width
	y := ip.img.Bounds().Min.Y + index/width

	if y > height {
		return fmt.Errorf("image index %d is too large", index)
	}

	originalColor := ip.img.At(x, y).(color.RGBA)
	newColor := color.RGBA{
		R: originalColor.R,
		G: originalColor.G,
		B: setLastBit(originalColor.B, bit),
		A: originalColor.A,
	}
	ip.newRGBA.Set(x, y, newColor)

	return nil
}

func (ip *ImagePacker) Get(index int) (uint8, error) {
	if index < 0 {
		return 2, fmt.Errorf("index out of range")
	}

	var err error = nil
	var bit uint8 = 2

	switch ip.op {
	case Blue:
		bit, err = ip.getBlue(index)
	case Red:
	case Green:
	}

	return bit, err
}

func (ip *ImagePacker) getBlue(index int) (uint8, error) {
	width := ip.img.Bounds().Max.X - ip.img.Bounds().Min.X
	height := ip.img.Bounds().Max.Y - ip.img.Bounds().Min.Y

	x := ip.img.Bounds().Min.X + index%width
	y := ip.img.Bounds().Min.Y + index/width

	if y > height {
		return 2, fmt.Errorf("image index %d is too large", index)
	}

	clr := ip.img.At(x, y).(color.RGBA)
	if clr.B&0x01 == 0x01 {
		return 1, nil
	}
	return 0, nil
}

func setLastBit(value uint8, bit uint8) uint8 {
	if bit == 0 {
		return value & 0xFE
	}
	return value | 0x01
}

func (ip *ImagePacker) GetImage() image.Image {
	return ip.newRGBA
}
