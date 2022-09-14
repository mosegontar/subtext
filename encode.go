package underbyte

import (
	"errors"
	"image/color"
)

func (u *UnderbyteImage) EncodeMessage(rawMessage []byte) error {
	header := newHeader(rawMessage, 0.5)
	message := append(header.Bytes(), rawMessage...)

	if len(message) > 2*(u.dimensions.X*u.dimensions.Y) {
		return errors.New("message size > pixel count")
	}

	for i := 0; i < len(message); i += 2 {
		var rFlip, gFlip, bFlip, aFlip uint8

		bytesToPack := bytesInRange(message, i, 2)

		rFlip, gFlip = splitByte(bytesToPack[0])

		if len(bytesToPack) == 2 {
			bFlip, aFlip = splitByte(bytesToPack[1])
		}

		x, y := u.nthPixelCoordinates(i / 2)
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
	return nil
}

func bytesInRange(message []byte, index, count int) []byte {
	if index+2 > len(message) {
		return message[index:]
	} else {
		return message[index : index+2]
	}
}

/* Splits byte in half, resulting in two bytes, the first
of which represents the most significant digit, and the second
represents the least significant digit. */
func splitByte(b byte) (uint8, uint8) {
	msd := b & 240 >> 4
	lsd := b & 15
	return msd, lsd
}
