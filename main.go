package main

import (
	"flag"
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

func decodeMessage(filepath string, outFile *os.File) {
	sourceImage := NewUnderbyteImage(filepath)
	sourceImage.DecodeMessage(outFile)
}

func encodeMessage(message string, inputPath string, outFile *os.File) {
	messageBytes := parseMessage(message)
	underbyteImage := NewUnderbyteImage(inputPath)

	underbyteImage.EncodeMessage(messageBytes)

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
