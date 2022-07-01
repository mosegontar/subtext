package main

import (
	"bytes"
	"fmt"
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
			fmt.Printf("%s vs %s", string(expected), string(actual))
			t.Errorf("expected '%v', actual '%v'", expected, actual)
		}
	})
}

func BenchmarkDecodeMessage(b *testing.B) {
	newImage := underbytetest.BlankImage(5000, 5000)
	message := []byte(strings.Repeat("Z", 4500*4500))

	underbyteImage := UnderbyteImage{image: newImage, dimensions: newImage.Bounds().Size()}
	img := underbyteImage.image
	underbytetest.FillPixels(img, underbyteImage.dimensions.X, underbyteImage.dimensions.Y)

	underbyteImage.EncodeMessage(message)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		buff := new(bytes.Buffer)
		b.StartTimer()
		underbyteImage.DecodeMessage(buff)

	}
}
