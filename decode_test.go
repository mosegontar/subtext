package underbyte

import (
	"bytes"
	"testing"

	"github.com/mosegontar/underbyte/underbytetest"
)

func TestDecode(t *testing.T) {
	t.Run("it correctly decodes the encoded message", func(t *testing.T) {
		u := UnderbyteImage{underbytetest.NewImage(10, 2)}
		u.Encode([]byte("hello world"))

		buff := new(bytes.Buffer)
		u.Decode(buff)

		expected := "hello world"
		actual := buff.String()

		if actual != expected {
			t.Errorf("expected %s got %s", expected, actual)
		}
	})

	t.Run("it correctly a decodes a message that is longer than the total pixel count", func(t *testing.T) {
		u := UnderbyteImage{underbytetest.NewImage(10, 2)}
		message := []byte("hello world!!!!!!")
		u.Encode(message)

		buff := new(bytes.Buffer)
		u.Decode(buff)

		expected := string(message)
		actual := buff.String()

		if actual != expected {
			t.Errorf("expected %s got %s", expected, actual)
		}

	})
}

func Test_singleUnpack(t *testing.T) {
	t.Run("it writes the expected bytes to the buffer", func(t *testing.T) {
		u := UnderbyteImage{underbytetest.NewImage(10, 2)}
		u.Encode([]byte("hello"))

		buff := new(bytes.Buffer)
		cursor := NewPixelCursor(headerSize+len("hello"), headerSize)
		u.singleUnpack(buff, cursor)

		expected := "hello"
		actual := buff.String()

		if actual != expected {
			t.Errorf("expected %v actual %v", expected, actual)
		}
	})

}

func Test_getMessageLengthFromHeader(t *testing.T) {
	t.Run("it returns the correct integer value", func(t *testing.T) {
		u := UnderbyteImage{underbytetest.NewImage(10, 2)}
		u.Encode([]byte("hello world"))

		cursor := NewPixelCursor(headerSize, 0)
		actual := u.getMessageLengthFromHeader(cursor)
		expected := 11

		if actual != expected {
			t.Errorf("expected %d actual %d", expected, actual)
		}
	})
}
