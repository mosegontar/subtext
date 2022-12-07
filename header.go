package underbyte

import "encoding/binary"

const headerSize = 4

func messageHeader(message []byte) []byte {
	messageLength := len(message)
	buf := make([]byte, headerSize)

	// BigEndian to avoid ambiguity
	// with leading 0s in the buffer
	binary.BigEndian.PutUint32(buf, uint32(messageLength))

	return buf
}

func headerBytesToInt(b []byte) int {
	val := binary.BigEndian.Uint32(b)
	return int(val)
}

