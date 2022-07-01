package main

import (
	"fmt"
	"image/color"
	"io"
)

func (u *UnderbyteImage) DecodeMessage(w io.Writer) {
	start, end := u.messageStartAndEnd()

	var decoded []byte
	for i := start; i <= int(end); i++ {
		x, y := u.nthPixelCoordinates(i)

		c := u.colorAtPixel(x, y)
		val := revealByte(c)

		decoded = append(decoded, byte(val))
	}

	fmt.Fprintf(w, "%s", decoded)
}

func (u *UnderbyteImage) messageStartAndEnd() (int, int) {
	// Pixel at 0,0 encodes a value indicating the
	// the number of the subsequent bytes that encode
	// the message size. Since we 1:1 correspondence
	// betewen bytes and pixels, this gives us a way to
	// determine which pixels we need to parse to decode
	// the embedded message.
	c := u.colorAtPixel(0, 0)
	headerPrefix := int(revealByte(c))

	// If there's no prefix value, it means
	// message size is 0
	if headerPrefix == 0 {
		return 0, -1
	}

	// Decode message size (which is stored in the "headerSuffix")
	var headerSuffix []byte
	for i := 1; i <= headerPrefix; i++ {
		x, y := u.nthPixelCoordinates(i)

		pixelColor := u.colorAtPixel(x, y)
		val := revealByte(pixelColor)

		headerSuffix = append(headerSuffix, val)
	}

	// Message begins at the start-th pixel and ends
	// at the end-th pixel.
	start := int(1 + headerPrefix)
	end := headerPrefix + bytesToInt(headerSuffix)

	return start, end
}

func bytesToInt(b []byte) (total int) {
	// bytes are in little endian order
	for i := len(b) - 1; i >= 0; i-- {
		shift := i * 8
		total += int(b[i]) << shift
	}
	return
}

func revealByte(c color.NRGBA) uint8 {
	var val uint8
	// Reconstruct the byte using the
	// first two bits of each color value.
	// Shift bits to the left such that
	// the values parsed from r become the
	// the most significant bits of the byte
	// and the values parsed from a become the least.
	val += (c.R & 3) << 6
	val += (c.G & 3) << 4
	val += (c.B & 3) << 2
	val += (c.A & 3)
	return val
}
