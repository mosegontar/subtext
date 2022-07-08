package underbyte

import (
	"fmt"
	"image/color"
	"reflect"
	"strings"
	"testing"

	"github.com/mosegontar/underbyte/underbytetest"
)

func TestEncodeMessage(t *testing.T) {
	t.Run("sets the image pixels correctly", func(t *testing.T) {
		message := []byte("hello")

		newImage := underbytetest.BlankImage(10, 1)

		underbyteImage := UnderbyteImage{image: newImage, dimensions: newImage.Bounds().Size()}

		checkColors := underbytetest.PixelColorChecker(underbyteImage.image, t)

		// Confirm all pixel RGBA values are set to 0 before encoding
		for i := 0; i < 10; i++ {
			checkColors([4]int{0, 0, 0, 0}, i, 0)
		}

		underbyteImage.EncodeMessage(message)

		// header prefix
		checkColors([4]int{0, 0, 0, 1}, 0, 0)

		// header suffix
		checkColors([4]int{0, 0, 1, 1}, 1, 0)

		// h
		checkColors([4]int{1, 2, 2, 0}, 2, 0)

		// e
		checkColors([4]int{1, 2, 1, 1}, 3, 0)

		// l
		checkColors([4]int{1, 2, 3, 0}, 4, 0)

		// l
		checkColors([4]int{1, 2, 3, 0}, 5, 0)

		// o
		checkColors([4]int{1, 2, 3, 3}, 6, 0)

		for i := 7; i < 10; i++ {
			checkColors([4]int{0, 0, 0, 0}, i, 0)
		}

	})

	t.Run("sets the image pixels correctly when there are RGBA values greater than 0", func(t *testing.T) {
		message := []byte("hello")

		newImage := underbytetest.BlankImage(10, 1)

		underbyteImage := UnderbyteImage{image: newImage, dimensions: newImage.Bounds().Size()}
		underbyteImage.image.SetNRGBA(4, 0, color.NRGBA{121, 255, 28, 4})

		underbyteImage.EncodeMessage(message)

		checkColors := underbytetest.PixelColorChecker(underbyteImage.image, t)
		checkColors([4]int{121, 254, 31, 4}, 4, 0)
	})

	t.Run("does not modify pixels that are outside the image dimensions and returns an error", func(t *testing.T) {
		message := []byte("hello")

		newImage := underbytetest.BlankImage(1, 1)

		underbyteImage := UnderbyteImage{image: newImage, dimensions: newImage.Bounds().Size()}

		err := underbyteImage.EncodeMessage(message)

		checkColors := underbytetest.PixelColorChecker(underbyteImage.image, t)
		checkColors([4]int{0, 0, 0, 0}, 2, 0)

		if err == nil {
			t.Errorf("expected an error but did not get one")
		}
	})
}

func TestBuildHeader(t *testing.T) {
	t.Run("message with 0 bytes returns single byte with value of 0", func(t *testing.T) {
		message := []byte{}
		expected := []byte{byte(0)}
		actual := buildHeader(message)

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	})

	t.Run("message with more than one byte", func(t *testing.T) {
		var testcases = []struct {
			message      []byte
			headerPrefix byte
			headerSuffix []byte
		}{
			{
				message:      []byte("a"),
				headerPrefix: uint8(1),
				headerSuffix: []byte{1},
			},
			{
				message:      []byte(strings.Repeat("a", 255)),
				headerPrefix: uint8(1),
				headerSuffix: []byte{255},
			},
			{
				message:      []byte(strings.Repeat("a", 256)),
				headerPrefix: uint8(2),
				headerSuffix: []byte{0, 1},
			},
			{
				message:      []byte(strings.Repeat("😇", 65432)),
				headerPrefix: uint8(3),
				headerSuffix: []byte{96, 254, 3},
			},
		}

		for _, testcase := range testcases {
			t.Run(fmt.Sprintf("with message %d bytes long, headerPrefix is correct", len(testcase.message)), func(t *testing.T) {
				header := buildHeader(testcase.message)
				if header[0] != testcase.headerPrefix {
					t.Errorf("expected headerPrefix %v, got %v", testcase.headerPrefix, header[0])
				}

				if header[0] != uint8(len(header[1:])) {
					t.Errorf("headerPrefix %v does not match number of remaining header bytes %v", testcase.headerPrefix, uint8(len(header[1:])))
				}
			})

			t.Run(fmt.Sprintf("with message %d bytes long, headersuffix is correct", len(testcase.message)), func(t *testing.T) {
				header := buildHeader(testcase.message)

				if !reflect.DeepEqual(header[1:], testcase.headerSuffix) {
					t.Errorf("expected headerSuffix %v, got %v", testcase.headerSuffix, header[1:])
				}
			})

		}
	})
}
