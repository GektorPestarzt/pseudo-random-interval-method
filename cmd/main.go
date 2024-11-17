package main

import (
	"flag"
	"fmt"
	"golang.org/x/image/bmp"
	"gopkg.in/yaml.v3"
	"image"
	"image/color"
	"log"
	"os"
	"pseudo-random-interval-method/internal/codec"
	"pseudo-random-interval-method/internal/packer"
)

const DEFAULT_CONFIG_FILE = "config.yaml"

type Config struct {
	Encode  Encode  `yaml:"encode"`
	Decode  Decode  `yaml:"decode"`
	Compare Compare `yaml:"compare"`

	EOT   string `yaml:"eot"`
	Entry int    `yaml:"entry"`
	Key   int    `yaml:"key"`
	Op    string `yaml:"op"`
	op    uint
}

type Encode struct {
	Input  string `yaml:"input"`
	Output string `yaml:"output"`
	Text   string `yaml:"text"`
}

type Decode struct {
	Input string `yaml:"input"`
}

type Compare struct {
	First  string `yaml:"first"`
	Second string `yaml:"second"`
	Output string `yaml:"output"`
}

func main() {
	mode, configFile, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	config, err := parseConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	switch mode {
	case "encode":
		err = encode(config)
		if err == nil {
			fmt.Println("Image successfully saved!")
		}
	case "decode":
		text, err := decode(config)
		if err == nil {
			fmt.Printf("Decoded text: \"%s\"\n", text)
		}
	case "compare":
		err = compare(config)

	default:
		log.Fatalf("Unknown mode: %s", mode)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func parseArgs() (string, string, error) {
	configFile := flag.String("c", DEFAULT_CONFIG_FILE, "Config file")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		return "", "", fmt.Errorf("Usage: program <mode>")
	}

	return args[0], *configFile, nil
}

func parseConfig(configFile string) (*Config, error) {
	bytes, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file: %v\n", err)
	}

	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("Failed to parse config file: %v\n", err)
	}

	config.op = packer.OpMap(config.Op)

	return &config, nil
}

func encode(config *Config) error {
	file, err := os.Open(config.Encode.Input)
	if err != nil {
		return fmt.Errorf("Failed to open input file: %v\n", err)
	}
	defer file.Close()

	img, err := bmp.Decode(file)
	if err != nil {
		return fmt.Errorf("Failed to decode bmp image: %v\n", err)
	}

	pack, err := packer.NewImagePacker(img, config.op)
	if err != nil {
		return fmt.Errorf("Failed to create packer: %v\n", err)
	}

	cdc := codec.NewCodec(pack, config.Entry, config.Key, config.EOT)
	imgWithText, err := cdc.Encode(config.Encode.Text)
	if err != nil {
		return fmt.Errorf("Failed to encode text to image: %v\n", err)
	}

	outFile, err := os.Create(config.Encode.Output)
	if err != nil {
		return fmt.Errorf("Failed to open output file: %v\n", err)
	}
	defer outFile.Close()

	err = bmp.Encode(outFile, imgWithText)
	if err != nil {
		return fmt.Errorf("Failed to encode bmp image: %v\n", err)
	}

	return nil
}

func decode(config *Config) (string, error) {
	file, err := os.Open(config.Decode.Input)
	if err != nil {
		return "", fmt.Errorf("Failed to open input file: %v\n", err)
	}
	defer file.Close()

	img, err := bmp.Decode(file)
	if err != nil {
		return "", fmt.Errorf("Failed to decode bmp image: %v\n", err)
	}

	pack, err := packer.NewImagePacker(img, config.op)
	if err != nil {
		return "", fmt.Errorf("Failed to create packer: %v\n", err)
	}

	cdc := codec.NewCodec(pack, config.Entry, config.Key, config.EOT)
	str, err := cdc.Decode()
	if err != nil {
		return "", fmt.Errorf("Failed to decode text from image: %v\n", err)
	}

	return str, nil
}

func compare(config *Config) error {
	file1, err := os.Open(config.Compare.First)
	if err != nil {
		return fmt.Errorf("Failed to open first file: %v\n", err)
	}
	defer file1.Close()

	img1, err := bmp.Decode(file1)
	if err != nil {
		return fmt.Errorf("Failed to decode first bmp image: %v\n", err)
	}

	file2, err := os.Open(config.Compare.Second)
	if err != nil {
		return fmt.Errorf("Failed to open second file: %v\n", err)
	}
	defer file2.Close()

	img2, err := bmp.Decode(file2)
	if err != nil {
		return fmt.Errorf("Failed to decode second bmp image: %v\n", err)
	}

	img, err := compareImages(img1, img2)
	if err != nil {
		return fmt.Errorf("Failed to compare images: %v\n", err)
	}

	outFile, err := os.Create(config.Compare.Output)
	if err != nil {
		return fmt.Errorf("Failed to open output file: %v\n", err)
	}
	defer outFile.Close()

	err = bmp.Encode(outFile, img)
	if err != nil {
		return fmt.Errorf("Failed to encode bmp image: %v\n", err)
	}

	return nil
}

func compareImages(image1, image2 image.Image) (image.Image, error) {
	bounds1 := image1.Bounds()
	bounds2 := image2.Bounds()

	if bounds1 != bounds2 {
		return nil, fmt.Errorf("Image bounds do not match")
	}

	diffImage := image.NewRGBA(bounds1)

	for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
		for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
			c1 := image1.At(x, y)
			c2 := image2.At(x, y)

			if c1 != c2 {
				diffImage.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 255})
			} else {
				diffImage.Set(x, y, color.Black)
			}
		}
	}

	return diffImage, nil
}
