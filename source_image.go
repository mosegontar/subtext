package underbyte

import (
	"bytes"
	"image"
	"io"
	"os"
)

type SourceImagePath string
type SourceImageBytes []byte

type ImageLoader interface {
	loadImageData() image.Image
}

func (s SourceImagePath) loadImageData() image.Image {
	return imageData(string(s))
}

func (s SourceImageBytes) loadImageData() image.Image {
	reader := bytes.NewReader(s)
	return decodeImage(reader)
}

func imageData(filepath string) image.Image {
	imgfile := openImage(filepath)
	defer imgfile.Close()
	img := decodeImage(imgfile)

	return img
}

func openImage(filepath string) *os.File {
	imgfile, err := os.Open(filepath)
	if err != nil {
		panic(err.Error())
	}
	return imgfile
}

func decodeImage(r io.Reader) image.Image {
	img, _, err := image.Decode(r)
	if err != nil {
		panic(err.Error())
	}

	return img
}
