package underbyte

import (
	"encoding/binary"
	"math"
)

type MessageHeader struct {
	size           int
	pixelByteRatio float64
}

func (m *MessageHeader) Bytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(m.size))
	return buf
}

func (m *MessageHeader) messageOffset() int {
	n := len(m.Bytes())
	pixelCount := math.Round(float64(n) * m.pixelByteRatio)
	return int(pixelCount)
}

func (m *MessageHeader) messageEnd() int {
	headerSize := binary.BigEndian.Uint32(m.Bytes())
	pixelCount := math.Round(float64(headerSize) * m.pixelByteRatio)
	return m.messageOffset() + int(pixelCount)
}

func newHeader(message []byte, ratio float64) MessageHeader {
	size := len(message)

	header := MessageHeader{
		size:           size,
		pixelByteRatio: math.Max(ratio, 0.5),
	}

	return header
}
