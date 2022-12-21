package underbyte

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
)

type UnderbyteOptions struct {
	Randomize bool
	Seed      string
}
type UnderbyteImage struct {
	*image.NRGBA
	options UnderbyteOptions
}

func NewUnderbyteImage(source ImageLoader, options *UnderbyteOptions) *UnderbyteImage {
	original := source.loadImageData()
	originalBounds := original.Bounds()

	rectangle := image.Rect(0, 0, originalBounds.Dx(), originalBounds.Dy())
	newImage := image.NewNRGBA(rectangle)

	draw.Draw(newImage, newImage.Bounds(), original, originalBounds.Min, draw.Src)

	if options == nil {
		options = &UnderbyteOptions{Randomize: true}
	}

	return &UnderbyteImage{NRGBA: newImage, options: *options}
}

func (u *UnderbyteImage) WriteImage(w io.Writer) {
	png.Encode(w, u)
}

func (u *UnderbyteImage) colorAtPixel(x, y int) color.NRGBA {
	return u.NRGBAAt(x, y)
}

func (u *UnderbyteImage) maxXCoordinate() int {
	return u.Bounds().Size().X
}

func (u *UnderbyteImage) maxYCoordinate() int {
	return u.Bounds().Size().Y
}

func (u *UnderbyteImage) nthPixelCoordinates(n int) (x, y int) {
	x = n % u.maxXCoordinate()
	y = n / u.maxXCoordinate()
	return
}

func (u *UnderbyteImage) pixelCount() int {
	return u.maxXCoordinate() * u.maxYCoordinate()
}

func (u *UnderbyteImage) seedFromHeaderPixels() int64 {
	values := []byte{}

	for i := 0; i < headerSize; i++ {
		x, y := u.nthPixelCoordinates(i)
		color := u.colorAtPixel(x, y)
		values = append(values, color.G)
		values = append(values, color.B)
	}

	n := binary.BigEndian.Uint64(values)

	return int64(n)
}

func (u *UnderbyteImage) randomizationSeed() int64 {
	if u.options.Seed != "" {
		return toInt64(u.options.Seed)
	} else {
		return u.seedFromHeaderPixels()
	}
}

func toInt64(s string) int64 {
	var n int64

	h := sha256.New()
	h.Write([]byte(s))
	hsum := h.Sum(nil)

	buf := bytes.NewReader(hsum)
	err := binary.Read(buf, binary.BigEndian, &n)
	if err != nil {
		panic(err)
	}

	return n
}
