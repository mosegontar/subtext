# underbyte

Fun with steganography, images, and Go. Encode arbitrary bytes in image pixels.

Note that encoding only writes output images as PNGs (though JPEGs can be used as an input image).

This is a **toy** project I've worked on to learn Go while on parental leave, during late night baby-is-sleeping-but-might-wake-up shifts. For an example of a non-trivial implementation that actually handles JPEG encoding see https://github.com/owencm/secretbook (which also has an accompanying paper and talk).

## Installation
### Building from source
```
go build
```

## Usage
```
Usage of ./underbyte:
  -decode
    	Decode message from image instead of encoding (default false)
  -file string
    	Input image filepath
  -message string
    	message to encode (STDIN used if not specified)
  -out string
    	Output filepath for encoded image (STDOUT used if not specified)

Examples:
  Encoding:
	underbyte -file image.png -message "hello there" -out encoded_image.png
	underbyte -file image.png -message "hello there" > encoded_image.png
	cat somefile.txt | underbyte -file image.png > encoded_image.png
  Decoding:
  	underbyte -decode -file encoded_image.png
	underbyte -decode -file encoded_image.png -out decoded_image.png
```
