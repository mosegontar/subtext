package underbyte

import (
	"bytes"
	"testing"

	"github.com/mosegontar/underbyte/underbytetest"
)

func TestDecode(t *testing.T) {
	t.Run("it correctly decodes the encoded message", func(t *testing.T) {
		u := UnderbyteImage{NRGBA: underbytetest.NewImage(10, 2)}
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
		u := UnderbyteImage{NRGBA: underbytetest.NewImage(10, 2)}
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

	t.Run("decoding when randomization option is true", func(t *testing.T) {
		u := UnderbyteImage{
			NRGBA: underbytetest.NewImage(10, 2),
			options: UnderbyteOptions{
				randomize: true,
			},
		}

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
