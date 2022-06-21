package lib

import (
	"fmt"
	"image/color"
	"math"
	"strings"
)

func DecodePixelsToBytes(filename string) {
	newImage, dimensions := drawCopy(filename)
	var message string

	for x := 0; x < dimensions.X; x++ {
		for y := 0; y < dimensions.Y; y++ {
			nrgba := newImage.NRGBAAt(x, y)
			r, g, b, a := nrgba.R, nrgba.G, nrgba.B, nrgba.A
			var val uint8
			val += (r & 3) << 6
			val += (g & 3) << 4
			val += (b & 3) << 2
			val += (a & 3)

			if string(val) == "\000" {
				fmt.Println(message)
				return
			} else {
				message += string(val)
			}

		}
	}
}

func DecodePixelsToNibbles(filename string) {
	newImage, dimensions := drawCopy(filename)

	var collection []string
	var nibbles = make([]uint8, 0)
	for x := 0; x < dimensions.X; x++ {

		for y := 0; y < dimensions.Y; y++ {
			nrgba := newImage.NRGBAAt(x, y)
			nibble := pixelToNibble(nrgba)

			if len(nibbles) == 0 {
				nibbles = append(nibbles, nibble)
			} else if len(nibbles) == 1 {
				nibbles[0] = (nibbles[0] << 4) + nibble
				if string(nibbles[0]) == "\000" {
					fmt.Println(strings.Join(collection, ""))
					return
				}
				collection = append(collection, string(nibbles[0]))
				nibbles = []uint8{}
			}

		}
	}
	fmt.Println(strings.Join(collection, ""))
}

func pixelToNibble(nrgba color.NRGBA) uint8 {
	r, g, b, a := nrgba.R, nrgba.G, nrgba.B, nrgba.A

	var nibble uint8

	nibble += (a & 1) * uint8(math.Pow(2, 0))
	nibble += (b & 1) * uint8(math.Pow(2, 1))
	nibble += (g & 1) * uint8(math.Pow(2, 2))
	nibble += (r & 1) * uint8(math.Pow(2, 3))

	return nibble
}
