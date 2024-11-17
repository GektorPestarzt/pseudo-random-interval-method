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
	case Blue, Red, Green:
		err = ip.setColor(index, bit)
	}

	return err
}

func (ip *ImagePacker) setColor(index int, bit uint8) error {
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
		B: originalColor.B,
		A: originalColor.A,
	}

	switch ip.op {
	case Blue:
		newColor.B = setLastBit(originalColor.B, bit)
	case Red:
		newColor.R = setLastBit(originalColor.R, bit)
	case Green:
		newColor.G = setLastBit(originalColor.G, bit)
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
	case Blue, Red, Green:
		bit, err = ip.getColor(index)
	}

	return bit, err
}

func (ip *ImagePacker) getColor(index int) (uint8, error) {
	width := ip.img.Bounds().Max.X - ip.img.Bounds().Min.X
	height := ip.img.Bounds().Max.Y - ip.img.Bounds().Min.Y

	x := ip.img.Bounds().Min.X + index%width
	y := ip.img.Bounds().Min.Y + index/width

	if y > height {
		return 2, fmt.Errorf("image index %d is too large", index)
	}

	clr := ip.img.At(x, y).(color.RGBA)
	var bit uint8 = 2

	switch ip.op {
	case Blue:
		bit = clr.B & 0x01
	case Red:
		bit = clr.R & 0x01
	case Green:
		bit = clr.G & 0x01
	}

	if bit == 2 {
		return bit, fmt.Errorf("wrong packer op")
	}
	return bit, nil
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

func OpMap(text string) uint {
	opMap := map[string]uint{
		"blue":  Blue,
		"red":   Red,
		"green": Green,
	}

	return opMap[text]
}
