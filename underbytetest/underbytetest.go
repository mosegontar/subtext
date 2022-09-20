package underbytetest

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"
)

func BlankImage(x, y int) *image.NRGBA {
	rectangle := image.Rect(0, 0, x, y)
	img := image.NewNRGBA(rectangle)
	draw.Draw(img, img.Bounds(), img, img.Bounds().Min, draw.Src)
	return img
}

func FillPixels(img *image.NRGBA, maxX, maxY int) {
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			img.SetNRGBA(x, y, color.NRGBA{1, 2, 3, 4})
		}
	}
}

func PixelColorChecker(img *image.NRGBA, t *testing.T) func([4]int, int, int) {
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
			t.Errorf("expected %v for x %d y %d, actual %v", expectedColors, x, y, actualColors)
		}
	}
}
