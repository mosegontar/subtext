package underbyte

import (
	"image/color"
	"testing"

	"github.com/mosegontar/underbyte/underbytetest"
)

func TestEncodeMessage(t *testing.T) {
	t.Run("sets the image pixels correctly", func(t *testing.T) {
		message := []byte("hello")

		newImage := underbytetest.BlankImage(5, 1)

		underbyteImage := UnderbyteImage{
			image:      newImage,
			dimensions: newImage.Bounds().Size(),
		}

		checkColors := underbytetest.PixelColorChecker(underbyteImage.image, t)

		// Confirm all pixel RGBA values are set to 0 before encoding
		for i := 0; i < 10; i++ {
			checkColors([4]int{0, 0, 0, 0}, i, 0)
		}

		underbyteImage.EncodeMessage(message)

		// header with size in bytes of message.
		// header is 4 bytes long, so takes up
		// two pixels.
		checkColors([4]int{0, 0, 0, 0}, 0, 0)
		checkColors([4]int{0, 0, 0, 5}, 1, 0)

		//                 h     e
		checkColors([4]int{6, 8, 6, 5}, 2, 0)
		//                 l      l
		checkColors([4]int{6, 12, 6, 12}, 3, 0)

		//                 o
		checkColors([4]int{6, 15, 0, 0}, 4, 0)

		for i := 5; i < 10; i++ {
			checkColors([4]int{0, 0, 0, 0}, i, 0)
		}

	})

	t.Run("sets the image pixels correctly using SinglePackStrategy", func(t *testing.T) {
		message := []byte("hello")

		newImage := underbytetest.BlankImage(10, 1)

		underbyteImage := UnderbyteImage{
			image:      newImage,
			dimensions: newImage.Bounds().Size(),
		}

		checkColors := underbytetest.PixelColorChecker(underbyteImage.image, t)

		// Confirm all pixel RGBA values are set to 0 before encoding
		for i := 0; i < 10; i++ {
			checkColors([4]int{0, 0, 0, 0}, i, 0)
		}

		underbyteImage.EncodeMessage(message)

		// header with size in bytes of message
		checkColors([4]int{0, 0, 0, 0}, 0, 0)
		checkColors([4]int{0, 0, 0, 5}, 1, 0)

		//                 h
		checkColors([4]int{1, 2, 2, 0}, 2, 0)
		//                 e
		checkColors([4]int{1, 2, 1, 1}, 3, 0)
		//                 l
		checkColors([4]int{1, 2, 3, 0}, 4, 0)
		//		   l
		checkColors([4]int{1, 2, 3, 0}, 5, 0)
		//                 o
		checkColors([4]int{1, 2, 3, 3}, 6, 0)

		for i := 7; i < 10; i++ {
			checkColors([4]int{0, 0, 0, 0}, i, 0)
		}
	})

	t.Run("sets the image pixels correctly when there are RGBA values greater than 0", func(t *testing.T) {
		message := []byte("hello")

		newImage := underbytetest.BlankImage(5, 1)

		underbyteImage := UnderbyteImage{image: newImage, dimensions: newImage.Bounds().Size()}
		underbyteImage.image.SetNRGBA(4, 0, color.NRGBA{121, 255, 28, 4})

		underbyteImage.EncodeMessage(message)

		checkColors := underbytetest.PixelColorChecker(underbyteImage.image, t)
		checkColors([4]int{118, 255, 16, 0}, 4, 0)
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
