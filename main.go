package main

import (
	"flag"
	"image/png"
	"io/ioutil"
	"os"
)

func parseMessage(message string) []byte {
	if message != "" {
		return []byte(message + "\000")
	}

	byts, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err.Error())
	}
	byts = append(byts, byte('\000'))
	return []byte(string(byts[:]))
}

func decodeMessage(filepath string) {
	DecodePixelsToBytes(filepath)
}

func encodeMessage(message string, inputPath string, outputPath string) {
	messageBytes := parseMessage(message)
	img := EncodeByteToPixel(inputPath, messageBytes)

	var outFile *os.File

	if outputPath == "" {
		outFile = os.Stdout
	} else {
		outFile, err := os.Create(outputPath)
		if err != nil {
			panic(err.Error())
		}
		defer outFile.Close()
	}

	png.Encode(outFile, img)
}

func main() {

	filepath := flag.String("f", "", "input image filepath")
	outpath := flag.String("o", "", "output image filepath")
	decode := flag.Bool("d", false, "decode message from image")
	message := flag.String("m", "", "message to encode")

	flag.Parse()

	if *decode {
		decodeMessage(*filepath)
	} else {
		encodeMessage(*message, *filepath, *outpath)
	}
}
