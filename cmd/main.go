package main

import (
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"log"
	"os"
	"pseudo-random-interval-method/internal/codec"
	"pseudo-random-interval-method/internal/packer"
)

func main() {
	decode()
}

func encode() {
	file, err := os.Open("samples/sample_640×426.bmp")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := bmp.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	pack, err := packer.NewImagePacker(img, packer.Blue)
	if err != nil {
		log.Fatal(err)
	}

	cdc := codec.NewCodec(pack, 134, 4, "KiHeu,6")
	imgWithText, err := cdc.Encode("I love everyone")
	if err != nil {
		log.Fatal(err)
	}

	outFile, err := os.Create("samples/sample_640x426_encode.bmp")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	err = bmp.Encode(outFile, imgWithText)
	if err != nil {
		log.Fatal(err)
	}

	compareImages(img, imgWithText)

	fmt.Println("Image successfully saved!")
}

func decode() {
	file, err := os.Open("samples/sample_640x426_encode.bmp")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := bmp.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	pack, err := packer.NewImagePacker(img, packer.Blue)
	if err != nil {
		log.Fatal(err)
	}

	cdc := codec.NewCodec(pack, 134, 4, "KiHeu,6")
	str, err := cdc.Decode()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded string: \"%s\"\n", str)
}

func compareImages(img1, img2 image.Image) bool {
	if img1.Bounds() != img2.Bounds() {
		fmt.Println("Images have different sizes!")
		return false
	}

	result := true

	for y := img1.Bounds().Min.Y; y < img1.Bounds().Max.Y; y++ {
		for x := img1.Bounds().Min.X; x < img1.Bounds().Max.X; x++ {
			color1 := img1.At(x, y)
			color2 := img2.At(x, y)

			if color1 != color2 {
				fmt.Printf("%d x %d\n", x, y)
				result = false
			}
		}
	}

	return result
}
