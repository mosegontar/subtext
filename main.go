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
	sourceImage := NewSubtextImage(filepath)
	sourceImage.DecodeMessage(outFile)
}

func encodeMessage(message string, inputPath string, outFile *os.File) {
	messageBytes := parseMessage(message)
	subtextImage := NewSubtextImage(inputPath)

	subtextImage.EncodeMessage(messageBytes)

	subtextImage.WriteImage(outFile)
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
	filepath := flag.String("f", "", "input image filepath")
	outpath := flag.String("o", "", "output filepath for encoded image")
	decode := flag.Bool("d", false, "decode message from image")
	message := flag.String("m", "", "message to encode")

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
