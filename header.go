package underbyte

type MessageHeader struct {
	data []byte
}

func (m *MessageHeader) Bytes() []byte {
	return m.data
}

func (m *MessageHeader) messageOffset() uint {
	return uint(1 + m.data[0])
}

func (m *MessageHeader) messageEnd() uint {
	offset := uint(m.data[0])
	return offset + m.messageSize()
}

func (m *MessageHeader) messageSize() uint {
	return bytesToInt(m.data[1:])
}

func newHeader(message []byte) MessageHeader {
	size := len(message)

	headerSuffix := intToBytes(size)
	headerPrefix := uint8(len(headerSuffix))

	header := MessageHeader{data: []byte{headerPrefix}}
	header.data = append(header.data, headerSuffix...)

	return header
}

func intToBytes(val int) []byte {
	b := []byte{}

	for val > 0 {
		current := val & 255
		b = append(b, uint8(current))
		val = val >> 8
	}

	return b
}

// Convert byte slice to unsigned integer.
// Assumes bytes are in little endian order.
func bytesToInt(b []byte) (total uint) {
	for i := len(b) - 1; i >= 0; i-- {
		shift := i * 8
		total += uint(b[i]) << shift
	}
	return
}
