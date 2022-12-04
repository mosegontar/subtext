package underbyte

import (
	"errors"
	"image/color"
)

type packer func([]byte, *PixelCursor)

func (u *UnderbyteImage) Encode(message []byte) error {
	header := messageHeader(message)
	totalLength := len(message) + len(header)

	var pack packer

	if totalLength > 2*u.pixelCount() {
		return errors.New("message too large for pixel count")
	} else if totalLength > u.pixelCount() {
		pack = u.doublePack
	} else {
		pack = u.singlePack
	}

	cursor := NewPixelCursor(totalLength, 0)
	u.singlePack(header, cursor)
	pack(message, cursor)

	return nil
}

func (u *UnderbyteImage) singlePack(message []byte, cursor *PixelCursor) {
	for i := 0; i < len(message); i++ {
		messageByte := message[i]

		rFlip := (messageByte & 192) >> 6
		gFlip := (messageByte & 48) >> 4
		bFlip := (messageByte & 12) >> 2
		aFlip := (messageByte & 3)

		nthPixel, ok := cursor.next()
		if !ok {
			panic("cursor cannot return any new positions")
		}
		x, y := u.nthPixelCoordinates(nthPixel)
		nrgba := u.colorAtPixel(x, y)

		// Replace the last 2 bits of the NRGBA color
		// value with the appropriate "flip" value
		r := ((nrgba.R >> 2) << 2) + rFlip
		g := ((nrgba.G >> 2) << 2) + gFlip
		b := ((nrgba.B >> 2) << 2) + bFlip
		a := ((nrgba.A >> 2) << 2) + aFlip

		color := color.NRGBA{r, g, b, a}
		u.SetNRGBA(x, y, color)
	}
}

func (u *UnderbyteImage) doublePack(message []byte, cursor *PixelCursor) {
	for i := 0; i < len(message); i += 2 {
		nthPixel, ok := cursor.next()
		if !ok {
			panic("cursor cannot return any new positions")
		}
		x, y := u.nthPixelCoordinates(nthPixel)
		nrgba := u.colorAtPixel(x, y)

		bytesToPack := bytesInRange(message, i)
		rFlip, gFlip, bFlip, aFlip := doublePackFlipVals(bytesToPack)

		// Replace the last 4 bits of the NRGBA color
		// value with the appropriate "flip" value
		r := ((nrgba.R >> 4) << 4) + rFlip
		g := ((nrgba.G >> 4) << 4) + gFlip
		// TODO: This is wrong for case where bytesToPack is len 1,
		// I think. Should modify only if second byte is present.
		b := ((nrgba.B >> 4) << 4) + bFlip
		a := ((nrgba.A >> 4) << 4) + aFlip

		color := color.NRGBA{r, g, b, a}

		u.SetNRGBA(x, y, color)
	}
}

func doublePackFlipVals(bytesToPack []byte) (rFlip uint8, gFlip uint8, bFlip uint8, aFlip uint8) {
	rFlip, gFlip = splitByteInTwo(bytesToPack[0])

	if len(bytesToPack) == 2 {
		bFlip, aFlip = splitByteInTwo(bytesToPack[1])
	}

	return
}

func bytesInRange(message []byte, index int) []byte {
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
