package underbyte

import (
	"bytes"
	"errors"
	"io"
	"math"
)

type UnderbyteReader interface {
	unpack(*UnderbyteImage, *PixelCursor) []byte
	maxPixels() int
	messagePixels() int
}

type DoublePackReader struct {
	messageLength int
	headerLength  int
}

func (dpr *DoublePackReader) truncate() bool {
	return dpr.messageLength%2 == 1
}

func (dpr *DoublePackReader) messagePixels() int {
	return int(math.Round(float64(dpr.messageLength) / 2))
}

func (dpr *DoublePackReader) maxPixels() int {
	return dpr.messagePixels() + dpr.headerLength
}

type SinglePackReader struct {
	messageLength int
	headerLength  int
}

func (spr *SinglePackReader) messagePixels() int {
	return spr.messageLength
}

func (spr *SinglePackReader) maxPixels() int {
	return spr.messagePixels() + spr.headerLength
}

func (u *UnderbyteImage) newMessageReader(headerLength, messageLength int) (UnderbyteReader, error) {
	totalLength := headerLength + messageLength

	if totalLength > 2*u.pixelCount() {
		return nil, errors.New("calculated encoded data size is greater than pixel count")
	}

	var reader UnderbyteReader
	if totalLength > u.pixelCount() {
		reader = &DoublePackReader{
			messageLength: messageLength,
			headerLength:  headerLength,
		}
	} else {
		reader = &SinglePackReader{
			messageLength: messageLength,
			headerLength:  headerLength,
		}
	}

	return reader, nil
}

func (u *UnderbyteImage) Decode(w io.Writer) error {
	/* Get message size from header bytes */
	cursor := NewPixelCursor(headerSize, 0)
	headerReader, err := u.newMessageReader(headerSize, 0)
	if err != nil {
		panic(err)
	}
	headerBytes := headerReader.unpack(u, cursor)
	messageLength := headerBytesToInt(headerBytes)

	/* Extract message */
	messageReader, err := u.newMessageReader(headerSize, messageLength)
	if err != nil {
		panic(err)
	}

	if u.options.Randomize {
		cursor = NewRandomizedPixelCursor(*u, cursor.position(), messageReader.messagePixels())
	} else {
		cursor = NewPixelCursor(messageReader.maxPixels(), cursor.position())
	}

	messageBytes := messageReader.unpack(u, cursor)

	w.Write(messageBytes)
	return nil
}

func (spr SinglePackReader) unpack(u *UnderbyteImage, cursor *PixelCursor) []byte {
	buff := new(bytes.Buffer)
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

	return buff.Bytes()

}

func (dpr *DoublePackReader) unpack(u *UnderbyteImage, cursor *PixelCursor) []byte {
	buff := new(bytes.Buffer)
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

	if dpr.truncate() {
		b := buff.Bytes()
		return b[:len(b)-1]
	} else {
		return buff.Bytes()
	}
}
