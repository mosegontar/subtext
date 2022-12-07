package underbyte

import (
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
	randomize bool
}
type UnderbyteImage struct {
	*image.NRGBA
	options UnderbyteOptions
}

func NewUnderbyteImage(source ImageLoader) *UnderbyteImage {
	original := source.loadImageData()
	originalBounds := original.Bounds()

	rectangle := image.Rect(0, 0, originalBounds.Dx(), originalBounds.Dy())
	newImage := image.NewNRGBA(rectangle)

	draw.Draw(newImage, newImage.Bounds(), original, originalBounds.Min, draw.Src)

	return &UnderbyteImage{NRGBA: newImage, options: UnderbyteOptions{randomize: true}}
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
