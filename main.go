package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func parseMessage(message string) []byte {
	if message != "" {
		return []byte(message)
	}

	byts, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err.Error())
	}

	return byts
}

func decodeMessage(filepath string, outFile *os.File) {
	sourceImage := NewUnderbyteImage(filepath)
	sourceImage.DecodeMessage(outFile)
}

func encodeMessage(message string, inputPath string, outFile *os.File) {
	messageBytes := parseMessage(message)
	underbyteImage := NewUnderbyteImage(inputPath)

	err := underbyteImage.EncodeMessage(messageBytes)
	if err != nil {
		log.Fatal(err.Error())
	}

	underbyteImage.WriteImage(outFile)
}

func outputFile(outputPath string) (*os.File, error) {
	var f *os.File
	var err error

	if outputPath != "" {
		f, err = os.Create(outputPath)
		if err != nil {
			return nil, err
		}
		return f, nil
	}

	return os.Stdout, nil
}

func main() {
	filepath := flag.String("file", "", "Input image filepath")
	outpath := flag.String("out", "", "Output filepath for encoded image (STDOUT used if not specified)")
	decode := flag.Bool("decode", false, "Decode message from image instead of encoding (default false)")
	message := flag.String("message", "", "message to encode (STDIN used if not specified)")

	var defaultUsage func()
	defaultUsage = flag.Usage
	flag.Usage = func() {
		defaultUsage()
		exampleBlock := `
Examples:
  Encoding:
	underbyte -file image.png -message "hello there" -out encoded_image.png
	underbyte -file image.png -message "hello there" > encoded_image.png
	cat somefile.txt | underbyte -file image.png > encoded_image.png
  Decoding:
  	underbyte -decode -file encoded_image.png
	underbyte -decode -file encoded_image.png -out decoded_image.png
`
		fmt.Fprintln(flag.CommandLine.Output(), exampleBlock)
	}

	flag.Parse()

	f, err := outputFile(*outpath)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	if *decode {
		decodeMessage(*filepath, f)
	} else {
		encodeMessage(*message, *filepath, f)
	}
}
