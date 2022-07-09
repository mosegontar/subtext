package underbyte

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/mosegontar/underbyte/underbytetest"
)

func TestDecodeMessage(t *testing.T) {
	t.Run("correctly decodes an embedded message", func(t *testing.T) {
		message := []byte("hi how are you")

		newImage := underbytetest.BlankImage(300, 300)

		underbyteImage := UnderbyteImage{image: newImage, dimensions: newImage.Bounds().Size()}
		underbytetest.FillPixels(underbyteImage.image, underbyteImage.dimensions.X, underbyteImage.dimensions.Y)

		underbyteImage.EncodeMessage(message)

		buff := new(bytes.Buffer)
		underbyteImage.DecodeMessage(buff)

		expected := message
		actual := buff.Bytes()

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected '%v', actual '%v'", expected, actual)
		}
	})

	t.Run("correctly decodes an embedded message whose size must be represented in more than one byte", func(t *testing.T) {
		message := []byte(strings.Repeat("A", 256) + strings.Repeat("B", 256))
		newImage := underbytetest.BlankImage(300, 300)
		underbyteImage := UnderbyteImage{image: newImage, dimensions: newImage.Bounds().Size()}
		underbytetest.FillPixels(underbyteImage.image, underbyteImage.dimensions.X, underbyteImage.dimensions.Y)

		underbyteImage.EncodeMessage(message)

		buff := new(bytes.Buffer)
		underbyteImage.DecodeMessage(buff)

		expected := message
		actual := buff.Bytes()

		if len(actual) != 512 || !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected '%v', actual '%v'", expected, actual)
		}

	})
}

func BenchmarkDecodeMessage(b *testing.B) {
	newImage := underbytetest.BlankImage(5000, 5000)
	message := []byte(strings.Repeat("Z", 4750*4750))

	underbyteImage := UnderbyteImage{
		image:      newImage,
		dimensions: newImage.Bounds().Size(),
	}

	underbytetest.FillPixels(
		underbyteImage.image,
		underbyteImage.dimensions.X,
		underbyteImage.dimensions.Y,
	)

	err := underbyteImage.EncodeMessage(message)
	if err != nil {
		b.Fatal(err.Error())
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		buff := new(bytes.Buffer)
		b.StartTimer()
		underbyteImage.DecodeMessage(buff)

	}
}

