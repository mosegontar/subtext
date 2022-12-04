package underbyte

import (
	"bytes"
	"errors"
	"io"
	"math"
)

type unpacker func(*bytes.Buffer, *PixelCursor)

func (u *UnderbyteImage) Decode(w io.Writer) error {
	cursor := NewPixelCursor(headerSize, 0)

	messageLength := u.getMessageLengthFromHeader(cursor)

	totalLength := headerSize + messageLength

	var unpack unpacker
	var truncate bool
	if totalLength > 2*u.pixelCount() {
		return errors.New("calculated encoded data size is greater than pixel count")
	} else if totalLength > u.pixelCount() {
		halfLength := math.Round(float64(messageLength) / 2)
		cursor = NewPixelCursor(headerSize+int(halfLength), cursor.position())
		unpack = u.doubleUnpack
		truncate = messageLength%2 == 1
	} else {
		unpack = u.singleUnpack
		cursor = NewPixelCursor(headerSize+messageLength, cursor.position())

	}

	buff := new(bytes.Buffer)
	unpack(buff, cursor)

	messageBytes := buff.Bytes()
	if truncate {
		w.Write(messageBytes[:len(messageBytes)-1])
	} else {
		w.Write(messageBytes)
	}

	return nil
}

func (u *UnderbyteImage) getMessageLengthFromHeader(cursor *PixelCursor) int {
	buff := new(bytes.Buffer)
	u.singleUnpack(buff, cursor)

	return headerBytesToInt(buff.Bytes())
}

func (u *UnderbyteImage) singleUnpack(buff *bytes.Buffer, cursor *PixelCursor) {
	nthPixel, ok := cursor.next()

	for ok {
		x, y := u.nthPixelCoordinates(nthPixel)
		c := u.colorAtPixel(x, y)

		r := (c.R & 3 << 6)
		g := (c.G & 3 << 4)
		b := (c.B & 3 << 2)
		a := (c.A & 3)

		rgba := r | g | b | a

		buff.WriteByte(rgba)

		nthPixel, ok = cursor.next()
	}

}

func (u *UnderbyteImage) doubleUnpack(buff *bytes.Buffer, cursor *PixelCursor) {
	nthPixel, ok := cursor.next()

	for ok {
		x, y := u.nthPixelCoordinates(nthPixel)

		c := u.colorAtPixel(x, y)

		firstByte := (c.R&15)<<4 + (c.G & 15)
		secondByte := (c.B&15)<<4 + (c.A & 15)

		buff.WriteByte(firstByte)
		buff.WriteByte(secondByte)

		nthPixel, ok = cursor.next()
	}
}
