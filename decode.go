package underbyte

import (
	"encoding/binary"
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

	decoded := u.parseBytes(start, end)
	message := decoded[:len(decoded)-(header.size%2)]
	fmt.Fprintf(w, "%s", message)
}

func (u *UnderbyteImage) decodeHeader() MessageHeader {
	headerBytes := u.parseBytes(0, 8)
	headerValue := binary.BigEndian.Uint32(headerBytes)

	header := MessageHeader{
		size:           int(headerValue),
		pixelByteRatio: 0.5,
	}

	return header
}

func (u *UnderbyteImage) parseBytes(first, last int) []byte {
	collection := []byte{}

	for i := first; i < last; i++ {
		x, y := u.nthPixelCoordinates(i)

		c := u.colorAtPixel(x, y)
		firstByte, secondByte := revealBytes(c)

		collection = append(collection, firstByte)
		collection = append(collection, secondByte)

	}
	return collection
}

func revealBytes(c color.NRGBA) (firstByte, secondByte byte) {
	firstByte = (c.R&15)<<4 + (c.G & 15)
	secondByte = (c.B&15)<<4 + (c.A & 15)
	return
}
