package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"os"
)

type UnderbyteImage struct {
	image      *image.NRGBA
	dimensions image.Point
}

func NewUnderbyteImage(sourceFilename string) *UnderbyteImage {
	original := imageData(sourceFilename)
	originalBounds := original.Bounds()

	rectangle := image.Rect(0, 0, originalBounds.Dx(), originalBounds.Dy())
	newImage := image.NewNRGBA(rectangle)

	draw.Draw(newImage, newImage.Bounds(), original, originalBounds.Min, draw.Src)

	return &UnderbyteImage{image: newImage, dimensions: originalBounds.Size()}
}

func (s *UnderbyteImage) EncodeMessage(message []byte) {
	for i := 0; i < len(message); i++ {
		x := i / s.dimensions.Y
		y := i % s.dimensions.Y

		if x > s.dimensions.X || y > s.dimensions.Y {
			return
		}

		nrgba := s.image.NRGBAAt(x, y)

		byt := message[i]
		// Interpret r,g,b,a values of a pixel as
		// corresponding to bits in a byte, where
		// 'r' represents the most significant bits
		// and 'a' represents the least.
		//
		// For each color value, the ones place value
		// is used to determine the value of the bits
		// in the byte. As a decimal, this value will be
		// between 0 and 3, which corresponds to the following
		// binary representations:
		// 11, 10, 01, 00
		//
		// For example, given RGBA values of {121, 255, 28, 4}
		// and the character 'l', which has an ASCII code of 108,
		// and 8 bit binary representation of 01101100, we modify
		// the RGBA values to be 121, 254, 31, 4, which have the following
		// 8 bit representations:
		// 	121 - 011110<01>
		//	254 - 111111<10>
		//       31 - 001111<11>
		//        4 - 000001<00>
		// Extracting the least significant bits from each number and shifting
		// according to the color's designated position in our scheme, we get
		// 	01101100
		// which is 108, ASCII character 'l'.

		// Use bitwise AND to determine value of the two least significant bits
		rFlip := (byt >> 6) & 3
		gFlip := (byt >> 4) & 3
		bFlip := (byt >> 2) & 3
		aFlip := (byt >> 0) & 3

		// Replace the last 2 bits of the NRGBA color
		// value with the appropriate "flip" value
		r := ((nrgba.R >> 2) << 2) + rFlip
		g := ((nrgba.G >> 2) << 2) + gFlip
		b := ((nrgba.B >> 2) << 2) + bFlip
		a := ((nrgba.A >> 2) << 2) + aFlip

		color := color.NRGBA{r, g, b, a}

		s.image.SetNRGBA(x, y, color)
	}
}

func (s *UnderbyteImage) WriteImage(w io.Writer) {
	png.Encode(w, s.image)
}

func (img *UnderbyteImage) DecodeMessage(w io.Writer) {
	var decoded string

	for x := 0; x < img.dimensions.X; x++ {
		message, endOfMessage := processColumn(x, img.image, img.dimensions.Y)
		decoded += message

		if endOfMessage {
			break
		}
	}

	fmt.Fprintf(w, "%s", decoded)
}

func processColumn(column int, img *image.NRGBA, maxRow int) (string, bool) {
	var message string

	for y := 0; y < maxRow; y++ {
		nrgba := img.NRGBAAt(column, y)
		r, g, b, a := nrgba.R, nrgba.G, nrgba.B, nrgba.A

		var val uint8
		// Reconstruct the byte using the
		// first two bits of each color value.
		// Shift bits to the left such that
		// the values parsed from r become the
		// the most significant bits of the byte
		// and the values parsed from a become the least.
		val += (r & 3) << 6
		val += (g & 3) << 4
		val += (b & 3) << 2
		val += (a & 3)

		if string(val) == "\000" {
			// Return true to indicate that this substring is the final one in the encoded message.
			return message, true
		} else {
			message += string(val)
		}

	}
	return message, false
}

func openImage(filepath string) *os.File {
	imgfile, err := os.Open(filepath)
	if err != nil {
		panic(err.Error())
	}
	return imgfile
}

func decodeImage(f *os.File) image.Image {
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err.Error())
	}

	return img
}

func imageData(filepath string) image.Image {
	imgfile := openImage(filepath)
	defer imgfile.Close()
	img := decodeImage(imgfile)

	return img
}
