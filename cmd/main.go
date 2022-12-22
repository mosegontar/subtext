package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mosegontar/underbyte"
)

type decodeOptions struct {
	filepath string
	outFile  *os.File
	secret   string
}

type encodeOptions struct {
	message  string
	filepath string
	outFile  *os.File
	secret   string
}

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

func decodeMessage(options decodeOptions) {
	opts := underbyte.UnderbyteOptions{
		Randomize: true,
		Secret:    options.secret,
	}

	s := underbyte.SourceImagePath(options.filepath)
	ub := underbyte.NewUnderbyteImage(s, &opts)

	ub.Decode(options.outFile)
}

func encodeMessage(options encodeOptions) {
	opts := underbyte.UnderbyteOptions{
		Randomize: true,
		Secret:    options.secret,
	}

	s := underbyte.SourceImagePath(options.filepath)
	underbyteImage := underbyte.NewUnderbyteImage(s, &opts)

	messageBytes := parseMessage(options.message)

	err := underbyteImage.Encode(messageBytes)
	if err != nil {
		log.Fatal(err.Error())
	}

	underbyteImage.WriteImage(options.outFile)
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
	secret := flag.String("secret", "", "secret key value used to encode and decode message")

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
		opts := decodeOptions{
			filepath: *filepath,
			outFile:  f,
			secret:   *secret,
		}
		decodeMessage(opts)
	} else {
		opts := encodeOptions{
			message:  *message,
			filepath: *filepath,
			outFile:  f,
			secret:   *secret,
		}
		encodeMessage(opts)
	}

}
