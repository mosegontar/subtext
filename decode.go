package underbyte

import (
	"fmt"
	"image/color"
	"io"
)

// Extract embedded message bytes from supplied image file.
func (u *UnderbyteImage) DecodeMessage(w io.Writer) {
	header := u.decodeHeader()
	start, end := header.messageOffset(), header.messageEnd()

	// Message is 0 bytes long, i.e,
	// empty, so write an empty string.
	if end == 0 {
		fmt.Fprint(w, "")
		return
	}

	var decoded []byte
	for i := start; i <= end; i++ {
		x, y := u.nthPixelCoordinates(int(i))

		c := u.colorAtPixel(x, y)
		val := revealByte(c)

		decoded = append(decoded, byte(val))
	}

	fmt.Fprintf(w, "%s", decoded)
}

func (u *UnderbyteImage) decodeHeader() MessageHeader {
	// Pixel at 0,0 encodes a value indicating the
	// the number of the subsequent bytes that encode
	// the message size. Since we have 1:1 correspondence
	// betewen bytes and pixels, this gives us a way to
	// determine which pixels we need to parse to decode
	// the embedded message.
	c := u.colorAtPixel(0, 0)

	header := MessageHeader{
		data: []byte{revealByte(c)},
	}

	offset := header.messageOffset()

	for i := 1; i < int(offset); i++ {
		x, y := u.nthPixelCoordinates(i)

		pixelColor := u.colorAtPixel(x, y)
		val := revealByte(pixelColor)

		header.data = append(header.data, val)
	}

	return header
}

// Extract the embedded byte from a NRGBA color

// Reconstruct the byte using the first two bits of each color value.
// Shift bits to the left such that the values parsed from c.R become the
// the most significant bits of the byte and the values parsed from c.A
// become the least.
//
// For example, given an NRGBA color of
//	{
//	  R: 255, // 0b11111111
//	  G: 128, // 0b10000000
// 	  B: 5,   // 0b00000101
// 	  A: 1    // 0b00000001
//	}
// we extract the two least significant bits from each channel using a bitwise AND
// 	e.,g, 255 & 3 == 0b11111111
//		        &
//			 0b00000011
//		         ----------
//		       = 0b00000011 == 3
//
//
// Then we shift each value according to the position it should occupy in our reconstructed
// byte. So our example becomes:
// 	0b11 << 6 == 3 << 6 == 0b11000000 +
//	0b00 << 4 == 0 << 4 == 0b00000000 +
// 	0b01 << 2 == 1 << 2 == 0b00000100 +
// 	0b01 << 0 == 1 << 0 == 0b00000001
// 			     --------------
// 			     = 0b11000101 == 197
func revealByte(c color.NRGBA) (val uint8) {
	val += (c.R & 3) << 6
	val += (c.G & 3) << 4
	val += (c.B & 3) << 2
	val += (c.A & 3)

	return
}
