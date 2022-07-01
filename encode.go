package main

import (
	"errors"
	"image/color"
)

func (u *UnderbyteImage) EncodeMessage(rawMessage []byte) error {
	header := buildHeader(rawMessage)
	message := append(header, rawMessage...)

	if len(message) > (u.dimensions.X * u.dimensions.Y) {
		return errors.New("message size > pixel count")
	}

	for i := 0; i < len(message); i++ {
		x, y := u.nthPixelCoordinates(i)

		nrgba := u.colorAtPixel(x, y)

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

		u.image.SetNRGBA(x, y, color)
	}
	return nil
}

func buildHeader(message []byte) (header []byte) {
	total := len(message)
	header = []byte{}

	headerSuffix := []byte{}
	for total > 0 {
		current := total & 255
		headerSuffix = append(headerSuffix, uint8(current))
		total = total >> 8
	}
	headerPrefix := uint8(len(headerSuffix))

	header = append(header, headerPrefix)
	header = append(header, headerSuffix...)

	return
}
