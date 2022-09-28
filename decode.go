package underbyte

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
	"math"
)

// Extract embedded message bytes from supplied image file.
func (u *UnderbyteImage) DecodeMessage(w io.Writer) {
	header := u.decodeHeader()

	var messagePixelCount, startingPixelIndex, endingPixel, pixelsAvailable int

	startingPixelIndex = header.messageOffset()
	pixelsAvailable = u.pixelCount() - startingPixelIndex

	if header.size > pixelsAvailable {
		u.strategy = DoublePackStrategy{}
		messagePixelCount = int(math.Round(float64(header.size) / 2))
	} else {
		u.strategy = SinglePackStrategy{}
		messagePixelCount = header.size
	}

	endingPixel = startingPixelIndex + messagePixelCount

	// Message is 0 bytes long, i.e,
	// empty, so write an empty string.
	if endingPixel == 0 {
		fmt.Fprint(w, "")
		return
	}

	decoded := u.parseBytes(startingPixelIndex, endingPixel+header.size%2)
	message := decoded[:len(decoded)-(header.size%2)]
	fmt.Fprintf(w, "%s", message)
}

func (u *UnderbyteImage) decodeHeader() MessageHeader {
	mh := MessageHeader{strategy: DoublePackStrategy{}}

	headerBytes := mh.strategy.unpack(u, 0, 2)
	headerValue := binary.BigEndian.Uint32(headerBytes)

	mh.size = int(headerValue)
	return mh
}

func (u *UnderbyteImage) parseBytes(first, last int) []byte {
	return u.strategy.unpack(u, first, last)
}

func (sp SinglePackStrategy) unpack(u *UnderbyteImage, pixelStart, pixelEnd int) []byte {
	byteCollection := []byte{}

	revealBytes := func(c color.NRGBA) byte {
		b := (c.R & 3 << 6) | (c.G & 3 << 4) | (c.B & 3 << 2) | (c.A & 3)
		return b
	}

	for i := pixelStart; i < pixelEnd; i++ {
		x, y := u.nthPixelCoordinates(i)

		c := u.colorAtPixel(x, y)
		rb := revealBytes(c)

		byteCollection = append(byteCollection, rb)
	}

	return byteCollection

}

func (dp DoublePackStrategy) unpack(u *UnderbyteImage, pixelStart, pixelEnd int) []byte {
	byteCollection := []byte{}

	revealBytes := func(c color.NRGBA) (byte, byte) {
		firstByte := (c.R&15)<<4 + (c.G & 15)
		secondByte := (c.B&15)<<4 + (c.A & 15)
		return firstByte, secondByte
	}

	for i := pixelStart; i < pixelEnd; i++ {
		x, y := u.nthPixelCoordinates(i)

		c := u.colorAtPixel(x, y)
		firstByte, secondByte := revealBytes(c)

		byteCollection = append(byteCollection, firstByte)
		byteCollection = append(byteCollection, secondByte)

	}

	return byteCollection
}
