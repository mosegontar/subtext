package underbyte

import (
	"image/color"
	"testing"

	"github.com/mosegontar/underbyte/underbytetest"
)

func TestEncode(t *testing.T) {
	t.Run("it correctly modifies the UnderbyteImage pixel values", func(t *testing.T) {
		u := UnderbyteImage{underbytetest.NewImage(10, 2)}
		u.Encode([]byte("hello"))

		checkColors := underbytetest.PixelColorChecker(u.NRGBA, t)

		// header, encodes message size (5)
		checkColors([4]int{0, 0, 0, 0}, 0, 0)
		checkColors([4]int{0, 0, 0, 0}, 1, 0)
		checkColors([4]int{0, 0, 0, 0}, 2, 0)
		checkColors([4]int{0, 0, 1, 1}, 3, 0)

		// h
		checkColors([4]int{1, 2, 2, 0}, 4, 0)
		// e
		checkColors([4]int{1, 2, 1, 1}, 5, 0)
		// l
		checkColors([4]int{1, 2, 3, 0}, 6, 0)
		// 1
		checkColors([4]int{1, 2, 3, 0}, 7, 0)
		// o
		checkColors([4]int{1, 2, 3, 3}, 8, 0)
	})

	t.Run("it correctly modifies existing pixel values", func(t *testing.T) {
		u := UnderbyteImage{underbytetest.NewImage(10, 2)}
		u.SetNRGBA(1, 0, color.NRGBA{100, 20, 19, 255})
		u.SetNRGBA(5, 0, color.NRGBA{5, 4, 244, 0})
		u.Encode([]byte("hello"))

		checkColors := underbytetest.PixelColorChecker(u.NRGBA, t)

		checkColors([4]int{100, 20, 16, 252}, 1, 0)
		checkColors([4]int{5, 6, 245, 1}, 5, 0)
	})

	t.Run("double packs message bytes when message length is greater than available pixels", func(t *testing.T) {
		u := UnderbyteImage{underbytetest.NewImage(4, 2)}
		err := u.Encode([]byte("hello"))

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			t.FailNow()
		}

		checkColors := underbytetest.PixelColorChecker(u.NRGBA, t)

		// header, encodes message size (5)
		checkColors([4]int{0, 0, 0, 0}, 0, 0)
		checkColors([4]int{0, 0, 0, 0}, 1, 0)
		checkColors([4]int{0, 0, 0, 0}, 2, 0)
		checkColors([4]int{0, 0, 1, 1}, 3, 0)

		//                 h     e
		checkColors([4]int{6, 8, 6, 5}, 0, 1)
		//                 l      l
		checkColors([4]int{6, 12, 6, 12}, 1, 1)
		//                 o
		checkColors([4]int{6, 15, 0, 0}, 2, 1)
	})
}
