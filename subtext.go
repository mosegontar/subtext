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
	"strings"
	"sync"
)

type SubtextImage struct {
	image      *image.NRGBA
	dimensions image.Point
}

func NewSubtextImage(sourceFilename string) *SubtextImage {
	original := imageData(sourceFilename)
	originalBounds := original.Bounds()

	rectangle := image.Rect(0, 0, originalBounds.Dx(), originalBounds.Dy())
	newImage := image.NewNRGBA(rectangle)

	draw.Draw(newImage, newImage.Bounds(), original, originalBounds.Min, draw.Src)

	return &SubtextImage{image: newImage, dimensions: originalBounds.Size()}
}

func (s *SubtextImage) EncodeMessage(message []byte) {
	for i := 0; i < len(message); i++ {
		x := i / s.dimensions.Y
		y := i % s.dimensions.Y

		nrgba := s.image.NRGBAAt(x, y)

		byt := message[i]
		// Interpret r,g,b,a values of pixel as
		// corresponding to bits in a byte, where
		// r represents the most significant bits
		// and a represents the least.
		//
		// For each color value, the ones place value
		// is used to determine the value of the bits
		// in the byte. As a decimal, this value will be
		// between 0 and 3, which corresponds to the following
		// binary representations:
		// 11, 10, 01, 00
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

func (s *SubtextImage) WriteImage(w io.Writer) {
	png.Encode(w, s.image)
}

func (img *SubtextImage) DecodeMessage(w io.Writer) {
	var wg sync.WaitGroup

	columns := make([]string, img.dimensions.X)

	var lastMessageRow int

	for x := 0; x < img.dimensions.X; x++ {
		wg.Add(1)

		go func(column int) {
			message, terminalColumn := processRow(&wg, column, img.image, img.dimensions.Y)
			columns[column] = message

			if terminalColumn {
				lastMessageRow = column
			}
		}(x)
	}

	wg.Wait()

	message := strings.Join(columns[:lastMessageRow+1], "")
	fmt.Fprintf(w, "%s", message)
}

func processRow(wg *sync.WaitGroup, column int, img *image.NRGBA, maxRow int) (string, bool) {
	defer wg.Done()

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
			message += string(val)
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
