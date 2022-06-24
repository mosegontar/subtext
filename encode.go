package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

// Store each byte of string in a single
// pixel.
func EncodeByteToPixel(filename string, messageBytes []byte) *image.NRGBA {
	newImage, dimensions := drawCopy(filename)

	for i := 0; i < len(messageBytes); i++ {
		byt := messageBytes[i]
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

		x := i / dimensions.Y
		y := i % dimensions.Y

		nrgba := newImage.NRGBAAt(x, y)

		// Replace the last 2 bits of the NRGBA color
		// value with the appropriate "flip" value
		r := ((nrgba.R >> 2) << 2) + rFlip
		g := ((nrgba.G >> 2) << 2) + gFlip
		b := ((nrgba.B >> 2) << 2) + bFlip
		a := ((nrgba.A >> 2) << 2) + aFlip

		color := color.NRGBA{r, g, b, a}
		newImage.SetNRGBA(x, y, color)
	}
	return newImage
}

func EncodeNibbleToPixel(s string) *image.NRGBA {
	nibbles := bitsToNibbles(strToBits(s))
	newImage, dimensions := drawCopy(os.Args[1])

	for i, nibble := range nibbles {
		x := i / dimensions.Y
		y := i % dimensions.Y

		nrgba := newImage.NRGBAAt(x, y)

		r := adjustColor(nibble[0], nrgba.R)
		g := adjustColor(nibble[1], nrgba.G)
		b := adjustColor(nibble[2], nrgba.B)
		a := adjustColor(nibble[3], nrgba.A)

		color := color.NRGBA{r, g, b, a}
		newImage.SetNRGBA(x, y, color)
	}

	return newImage
}

func SaveImage(img *image.NRGBA, filename string) {
	outFile, err := os.Create(filename)
	if err != nil {
		panic(err.Error())
	}
	defer outFile.Close()
	png.Encode(outFile, img)
}

func adjustColor(bit bool, color uint8) uint8 {
	if (bit && color%2 == 0) || (!bit && color%2 != 0) {
		if color == 0 {
			return color + 1
		}
		return color - 1
	}
	return color
}

func bitsToNibbles(bits []bool) (packs [][]bool) {
	for i := 0; i < len(bits); i += 4 {
		packs = append(packs, bits[i:i+4])
	}
	return
}

func strToBits(s string) (bits []bool) {
	s += "\000" // Null terminal str
	byteString := []byte(s)
	for _, byt := range byteString {
		for i := 128; i > 0; i /= 2 {
			bits = append(bits, int(byt)&i > 0)
		}
	}
	return
}
