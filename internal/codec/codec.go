package codec

import (
	"image"
	"math/bits"
	"pseudo-random-interval-method/internal/packer"
)

type Codec struct {
	packer packer.Packer
	entry  int
	key    int
	eot    string
}

func NewCodec(packer packer.Packer, entry int, key int, eot string) *Codec {
	return &Codec{
		packer: packer,
		entry:  entry,
		key:    key,
		eot:    eot,
	}
}

func (c *Codec) Encode(text string) (image.Image, error) {
	index := c.entry
	message := text + c.eot

	for _, b := range []byte(message) {
		for i := 7; i >= 0; i-- {
			err := c.packer.Set(index, (b>>i)&1)
			if err != nil {
				return nil, err
			}
			index = step(index, c.key)
		}
	}

	return c.packer.GetImage(), nil
}

func (c *Codec) Decode() (string, error) {
	index := c.entry
	var result []byte

	for {
		var b byte = 0
		for i := 0; i < 8; i++ {
			bit, err := c.packer.Get(index)
			if err != nil {
				return "", err
			}

			b *= 2
			b += bit

			index = step(index, c.key)
		}

		result = append(result, b)
		if isEOT(result, c.eot) {
			break
		}
	}

	return string(result[:len(result)-len(c.eot)]), nil
}

func step(index int, key int) int {
	return index + key*bits.OnesCount(uint(index))
}

func isEOT(text []byte, eot string) bool {
	var i int
	for i = 0; i < len(text) && i < len(eot); i++ {
		if text[len(text)-i-1] != eot[len(eot)-i-1] {
			return false
		}
	}

	if i == len(eot) {
		return true
	}
	return false
}
