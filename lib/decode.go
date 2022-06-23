package lib

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"
	"sync"
)

func processRow(wg *sync.WaitGroup, column int, img *image.NRGBA, maxRow int) (string, bool) {
	defer wg.Done()

	var message string

	for y := 0; y < maxRow; y++ {
		nrgba := img.NRGBAAt(column, y)
		r, g, b, a := nrgba.R, nrgba.G, nrgba.B, nrgba.A

		var val uint8
		// Reconstruct the byte using the
		// first two bits of each color value.
		// Shift bits to the left such that
		// the values parsed from r become the
		// the most significant bits of the byte
		// and the values parsed from a become the least.
		val += (r & 3) << 6
		val += (g & 3) << 4
		val += (b & 3) << 2
		val += (a & 3)

		if string(val) == "\000" {
			message += string(val)
			// Return true to indicate that this substring is the final one in the encoded message.
			return message, true
		} else {
			message += string(val)
		}

	}
	return message, false
}

func DecodePixelsToBytes(filename string) {
	newImage, dimensions := drawCopy(filename)

	var wg sync.WaitGroup

	columns := make([]string, dimensions.X)

	var lastMessageRow int

	for x := 0; x < dimensions.X; x++ {
		wg.Add(1)

		go func(column int) {
			message, terminalColumn := processRow(&wg, column, newImage, dimensions.Y)
			columns[column] = message

			if terminalColumn {
				lastMessageRow = column
			}
		}(x)
	}

	wg.Wait()

	message := strings.Join(columns[:lastMessageRow], "")
	fmt.Printf("%s", message)
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
