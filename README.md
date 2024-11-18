# pseudo-random-interval-method

[![GoDoc](https://godoc.org/github.com/gektorpestarzt/pseudo-random-interval-method?status.svg)](https://godoc.org/github.com/gektorpestarzt/pseudo-random-interval-method)

Implementation in Go of the pseudo-random interval method for bmp images described in the book by G.F. Konakhovich, A.Yu. Puzirenko, Computer Steganography. Theory and Practice, Moscow: MK-Press, 2006, pages 89-92.

## Installation

```bash
  go install github.com/gektorpestarzt/pseudo-random-interval-method/cmd/prim@latest
  cp $GOPATH/bin/prim ./
```

## Usage

You can run the program with different modes and flags.

### Command-line Flags

- -c: Path to the configuration file.

### Modes

- encode: Encode a message into an image.
- decode: Decode a message from an image.
- compare: Compare two images pixel by pixel.

### Example

```bash
  ./prim encode -c config.yaml
```
  
