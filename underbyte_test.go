package underbyte

import (
	"image"
	"image/color"
	"reflect"
	"testing"

	"github.com/mosegontar/underbyte/underbytetest"
)

type SourceImageMock struct {
	image *image.NRGBA
}

func (s SourceImageMock) loadImageData() image.Image {
	return s.image
}

func TestColorAtPixel(t *testing.T) {
	t.Run("it returns the correct NRGBA values", func(t *testing.T) {
		newImage := underbytetest.NewImage(5, 1)
		newImage.SetNRGBA(1, 0, color.NRGBA{5, 3, 1, 0})
		newImage.SetNRGBA(2, 0, color.NRGBA{2, 3, 4, 234})

		u := UnderbyteImage{NRGBA: newImage}

		testcases := []struct {
			x     int
			y     int
			nrgba color.NRGBA
		}{
			{x: 0, y: 0, nrgba: color.NRGBA{0, 0, 0, 0}},
			{x: 1, y: 0, nrgba: color.NRGBA{5, 3, 1, 0}},
			{x: 2, y: 0, nrgba: color.NRGBA{2, 3, 4, 234}},
		}

		for _, testcase := range testcases {
			actual := u.colorAtPixel(testcase.x, testcase.y)
			expected := testcase.nrgba

			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("expected %v actual %v, %v", expected, actual, newImage.NRGBAAt(testcase.x, testcase.y))
			}

		}
	})

}

func TestMaxXCoordinate(t *testing.T) {
	t.Run("it retuns the correct X coordinate value", func(t *testing.T) {
		testcases := []struct {
			u UnderbyteImage
			x int
		}{
			{u: UnderbyteImage{NRGBA: underbytetest.NewImage(5, 1)}, x: 5},
			{u: UnderbyteImage{NRGBA: underbytetest.NewImage(1, 1)}, x: 1},
			{u: UnderbyteImage{NRGBA: underbytetest.NewImage(200, 1)}, x: 200},
		}

		for _, testcase := range testcases {
			actual := testcase.u.maxXCoordinate()
			expected := testcase.x

			if actual != expected {
				t.Errorf("expected %v actual %v", expected, actual)
			}
		}
	})
}

func TestMaxYCoordinate(t *testing.T) {
	t.Run("it returns the correct X coordinate value", func(t *testing.T) {
		testcases := []struct {
			u UnderbyteImage
			y int
		}{
			{u: UnderbyteImage{NRGBA: underbytetest.NewImage(5, 1)}, y: 1},
			{u: UnderbyteImage{NRGBA: underbytetest.NewImage(1, 3)}, y: 3},
			{u: UnderbyteImage{NRGBA: underbytetest.NewImage(200, 321)}, y: 321},
		}

		for _, testcase := range testcases {
			actual := testcase.u.maxYCoordinate()
			expected := testcase.y

			if actual != expected {
				t.Errorf("expected %v actual %v", expected, actual)
			}
		}
	})
}

func TestPixelCount(t *testing.T) {
	t.Run("it returns the total number of pixels in the image", func(t *testing.T) {
		testcases := []struct {
			u     UnderbyteImage
			count int
		}{
			{u: UnderbyteImage{NRGBA: underbytetest.NewImage(1, 1)}, count: 1},
			{u: UnderbyteImage{NRGBA: underbytetest.NewImage(7, 5)}, count: 35},
			{u: UnderbyteImage{NRGBA: underbytetest.NewImage(111, 2)}, count: 222},
		}

		for _, testcase := range testcases {
			actual := testcase.u.pixelCount()
			expected := testcase.count

			if actual != expected {
				t.Errorf("expected %v actual %v", expected, actual)
			}
		}

	})
}

func TestNthPixelCoordinates(t *testing.T) {
	t.Run("it returns the coordinates of the n-th pixel", func(t *testing.T) {
		u := UnderbyteImage{NRGBA: underbytetest.NewImage(5, 5)}

		testcases := []struct {
			nthPixel int
			x        int
			y        int
		}{
			{nthPixel: 0, x: 0, y: 0},
			{nthPixel: 1, x: 1, y: 0},
			{nthPixel: 4, x: 4, y: 0},
			{nthPixel: 5, x: 0, y: 1},
			{nthPixel: 9, x: 4, y: 1},
			{nthPixel: 24, x: 4, y: 4},
		}

		for _, testcase := range testcases {
			actualX, actualY := u.nthPixelCoordinates(testcase.nthPixel)
			expectedX := testcase.x
			expectedY := testcase.y

			if actualX != expectedX {
				t.Errorf("expected %v actual %v", expectedX, actualX)
			}

			if actualY != expectedY {
				t.Errorf("expected %v actual %v", expectedY, actualY)
			}
		}

	})
}

func TestNewUnderbyteImage(t *testing.T) {
	t.Run("it returns an UnderbyteImage with the correct dimensions", func(t *testing.T) {
		sourceImageMock := SourceImageMock{image: underbytetest.NewImage(15, 30)}
		u := NewUnderbyteImage(sourceImageMock)

		actualX := u.maxXCoordinate()
		actualY := u.maxYCoordinate()

		expectedX := 15
		expectedY := 30

		if actualX != expectedX || actualY != expectedY {
			t.Errorf("expected (%v, %v) actual (%v, %v)", expectedX, expectedY, actualX, actualY)
		}
	})

	t.Run("it returns an UnderbyteImage with the expected pixel values", func(t *testing.T) {
		img := underbytetest.NewImage(2, 2)

		color0_0 := color.NRGBA{0, 0, 0, 0}
		color1_0 := color.NRGBA{1, 2, 3, 255}
		color0_1 := color.NRGBA{2, 3, 4, 255}
		color1_1 := color.NRGBA{3, 4, 5, 255}

		img.SetNRGBA(0, 0, color0_0)
		img.SetNRGBA(1, 0, color1_0)
		img.SetNRGBA(0, 1, color0_1)
		img.SetNRGBA(1, 1, color1_1)

		sourceImageMock := SourceImageMock{image: img}
		u := NewUnderbyteImage(sourceImageMock)

		first := u.NRGBAAt(0, 0)
		second := u.NRGBAAt(1, 0)
		third := u.NRGBAAt(0, 1)
		fourth := u.NRGBAAt(1, 1)

		if !reflect.DeepEqual(first, color0_0) ||
			!reflect.DeepEqual(second, color1_0) ||
			!reflect.DeepEqual(third, color0_1) ||
			!reflect.DeepEqual(fourth, color1_1) {
			t.Errorf("expected (0, 0): %v, (1, 0): %v, (0, 1): %v, (1, 1): %v. actual: (0, 0): %v, (1, 0): %v, (0, 1): %v, (1, 1): %v", color0_0, color1_0, color0_1, color1_1, first, second, third, fourth)
		}
	})
}
