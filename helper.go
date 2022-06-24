package main

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

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

func drawCopy(filename string) (*image.NRGBA, image.Point) {
	original := imageData(filename)
	originalBounds := original.Bounds()

	rectangle := image.Rect(0, 0, originalBounds.Dx(), originalBounds.Dy())
	newImage := image.NewNRGBA(rectangle)

	draw.Draw(newImage, newImage.Bounds(), original, originalBounds.Min, draw.Src)
	return newImage, originalBounds.Size()
}


