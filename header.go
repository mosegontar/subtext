package underbyte

import (
	"encoding/binary"
	"math"
)

type MessageHeader struct {
	size     int
	strategy PackingStrategy
}

func (m *MessageHeader) Bytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(m.size))
	return buf
}

func (m *MessageHeader) messageOffset() int {
	n := len(m.Bytes())
	pixelCount := math.Round(float64(n) / 2)
	return int(pixelCount)
}

func newHeader(message []byte) MessageHeader {
	size := len(message)

	header := MessageHeader{
		size:     size,
		strategy: DoublePackStrategy{},
	}

	return header
}
