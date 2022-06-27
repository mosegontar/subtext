package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"
)

func drawImage(img *image.NRGBA) {
	draw.Draw(img, img.Bounds(), img, img.Bounds().Min, draw.Src)
}

func pixelColorChecker(img *image.NRGBA, t *testing.T) func([4]int, int, int) {
	return func(expectedRGBAVals [4]int, x, y int) {
		t.Helper()

		expectedColors := color.NRGBA{
			R: uint8(expectedRGBAVals[0]),
			G: uint8(expectedRGBAVals[1]),
			B: uint8(expectedRGBAVals[2]),
			A: uint8(expectedRGBAVals[3]),
		}

		actualColors := img.NRGBAAt(x, y)

		if !reflect.DeepEqual(expectedColors, actualColors) {
			t.Errorf("expected %v, actual %v", expectedColors, actualColors)
		}
	}
}

func TestEncodeMessage(t *testing.T) {
	t.Run("sets the image pixels correctly", func(t *testing.T) {
		message := []byte("hello")

		rectangle := image.Rect(0, 0, 10, 1)
		newImage := image.NewNRGBA(rectangle)
		drawImage(newImage)

		underbyteImage := UnderbyteImage{image: newImage, dimensions: newImage.Bounds().Size()}

		checkColors := pixelColorChecker(underbyteImage.image, t)

		// Confirm all pixel RGBA values are set to 0 before encoding
		for i := 0; i < 10; i++ {
			checkColors([4]int{0, 0, 0, 0}, i, 0)
		}

		underbyteImage.EncodeMessage(message)

		// h
		checkColors([4]int{1, 2, 2, 0}, 0, 0)

		// e
		checkColors([4]int{1, 2, 1, 1}, 1, 0)

		// l
		checkColors([4]int{1, 2, 3, 0}, 2, 0)

		// l
		checkColors([4]int{1, 2, 3, 0}, 3, 0)

		// o
		checkColors([4]int{1, 2, 3, 3}, 4, 0)

		for i := 5; i < 10; i++ {
			checkColors([4]int{0, 0, 0, 0}, i, 0)
		}

	})

	t.Run("sets the image pixels correctly when there are RGBA values greater than 0", func(t *testing.T) {
		message := []byte("hello")

		rectangle := image.Rect(0, 0, 10, 1)
		newImage := image.NewNRGBA(rectangle)
		drawImage(newImage)

		underbyteImage := UnderbyteImage{image: newImage, dimensions: newImage.Bounds().Size()}
		underbyteImage.image.SetNRGBA(2, 0, color.NRGBA{121, 255, 28, 4})

		underbyteImage.EncodeMessage(message)

		checkColors := pixelColorChecker(underbyteImage.image, t)
		checkColors([4]int{121, 254, 31, 4}, 2, 0)
	})

	t.Run("does not modify pixels that are outside the image dimensions", func(t *testing.T) {
		message := []byte("hello")

		rectangle := image.Rect(0, 0, 1, 1)
		newImage := image.NewNRGBA(rectangle)
		drawImage(newImage)

		underbyteImage := UnderbyteImage{image: newImage, dimensions: newImage.Bounds().Size()}

		underbyteImage.EncodeMessage(message)

		checkColors := pixelColorChecker(underbyteImage.image, t)
		checkColors([4]int{0, 0, 0, 0}, 2, 0)
	})

}

func TestDecodeMessage(t *testing.T) {
	t.Run("correctly decodes an embedded message", func(t *testing.T) {
		message := []byte("hi how are you")

		rectangle := image.Rect(0, 0, 300, 300)
		newImage := image.NewNRGBA(rectangle)
		drawImage(newImage)

		underbyteImage := UnderbyteImage{image: newImage, dimensions: newImage.Bounds().Size()}
		underbyteImage.EncodeMessage(message)

		buff := new(bytes.Buffer)
		underbyteImage.DecodeMessage(buff)

		expected := "hi how are you"
		actual := string(buff.Bytes())

		if expected != actual {
			t.Errorf("expected '%v', actual '%v'", []byte(expected), []byte(actual))
		}
	})
}
