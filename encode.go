package underbyte

import (
	"errors"
	"image/color"
)

func (u *UnderbyteImage) EncodeMessage(message []byte) error {
	header := newHeader(message)

	pixelsAvailable := u.pixelCount() - header.messageOffset()

	if len(message) > 2*pixelsAvailable {
		return errors.New("message too large for pixel count")
	}

	if len(message) > pixelsAvailable {
		u.strategy = DoublePackStrategy{}
	} else {
		u.strategy = SinglePackStrategy{}
	}

	err := u.embedBytes(&header, message)
	if err != nil {
		return err
	}

	return nil
}

func (u *UnderbyteImage) embedBytes(header *MessageHeader, message []byte) error {
	offset := header.strategy.pack(u, header.Bytes(), 0)
	u.strategy.pack(u, message, offset)

	return nil
}

func (dp DoublePackStrategy) pack(u *UnderbyteImage, message []byte, offset int) (nPixel int) {

	nthPixel := func(offset, index int) int {
		return (index / 2) + offset
	}

	var i int

	for i = 0; i < len(message); i += 2 {
		var rFlip, gFlip, bFlip, aFlip uint8

		bytesToPack := bytesInRange(message, i, 2)

		rFlip, gFlip = splitByteInTwo(bytesToPack[0])

		if len(bytesToPack) == 2 {
			bFlip, aFlip = splitByteInTwo(bytesToPack[1])
		}

		x, y := u.nthPixelCoordinates(nthPixel(offset, i))
		nrgba := u.colorAtPixel(x, y)

		// Replace the last 4 bits of the NRGBA color
		// value with the appropriate "flip" value
		r := ((nrgba.R >> 4) << 4) + rFlip
		g := ((nrgba.G >> 4) << 4) + gFlip
		b := ((nrgba.B >> 4) << 4) + bFlip
		a := ((nrgba.A >> 4) << 4) + aFlip

		color := color.NRGBA{r, g, b, a}

		u.image.SetNRGBA(x, y, color)
	}

	return nthPixel(offset, i)
}

func (sp SinglePackStrategy) pack(u *UnderbyteImage, message []byte, offset int) (nPixel int) {

	nthPixel := func(offset, index int) int {
		return offset + index
	}

	var i int

	for i = 0; i < len(message); i++ {
		messageByte := message[i]

		rFlip := (messageByte & 192) >> 6
		gFlip := (messageByte & 48) >> 4
		bFlip := (messageByte & 12) >> 2
		aFlip := (messageByte & 3)

		x, y := u.nthPixelCoordinates(nthPixel(offset, i))
		nrgba := u.colorAtPixel(x, y)

		// Replace the last 4 bits of the NRGBA color
		// value with the appropriate "flip" value
		r := ((nrgba.R >> 2) << 2) + rFlip
		g := ((nrgba.G >> 2) << 2) + gFlip
		b := ((nrgba.B >> 2) << 2) + bFlip
		a := ((nrgba.A >> 2) << 2) + aFlip

		color := color.NRGBA{r, g, b, a}

		u.image.SetNRGBA(x, y, color)
	}

	return nthPixel(offset, i)
}

func bytesInRange(message []byte, index, count int) []byte {
	if index+2 > len(message) {
		return message[index:]
	} else {
		return message[index : index+2]
	}
}

/* Splits byte in half, resulting in two bytes, the first
of which represents the most significant digits, and the second
represents the least significant digits. */
func splitByteInTwo(b byte) (uint8, uint8) {
	msd := b & 240 >> 4
	lsd := b & 15
	return msd, lsd
}
