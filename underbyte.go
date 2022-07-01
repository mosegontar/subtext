package main

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"os"
)

type UnderbyteImage struct {
	image      *image.NRGBA
	dimensions image.Point
}

func NewUnderbyteImage(sourceFilename string) *UnderbyteImage {
	original := imageData(sourceFilename)
	originalBounds := original.Bounds()

	rectangle := image.Rect(0, 0, originalBounds.Dx(), originalBounds.Dy())
	newImage := image.NewNRGBA(rectangle)

	draw.Draw(newImage, newImage.Bounds(), original, originalBounds.Min, draw.Src)

	return &UnderbyteImage{image: newImage, dimensions: originalBounds.Size()}
}

func (u *UnderbyteImage) WriteImage(w io.Writer) {
	png.Encode(w, u.image)
}

func (u *UnderbyteImage) colorAtPixel(x, y int) color.NRGBA {
	return u.image.NRGBAAt(x, y)
}

func (u *UnderbyteImage) nthPixelCoordinates(n int) (x, y int) {
	x = n / u.dimensions.Y
	y = n % u.dimensions.Y
	return
}

func openImage(filepath string) *os.File {
	imgfile, err := os.Open(filepath)
	if err != nil {
		panic(err.Error())
	}
	return imgfile
}

func decodeImage(f *os.File) image.Image {
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err.Error())
	}

	return img
}

func imageData(filepath string) image.Image {
	imgfile := openImage(filepath)
	defer imgfile.Close()
	img := decodeImage(imgfile)

	return img
}
