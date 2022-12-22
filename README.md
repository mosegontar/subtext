# underbyte

Fun with steganography, images, and Go. Encode arbitrary bytes in image pixels.

Note that encoding only writes output images as PNGs (though JPEGs can be used as an input image).

## Example:
This image of Leo Tolstoy encodes an epub of the entirety of _War and Peace_.
![Leo Tolstoy](./assets/tolstoy_war_peace.png?raw=true "Leo Tolstoy")
It was created by running
```
cat war_and_peace.epub| ./underbyte -file tolstoy.jpg -out tolstoy_war_peace.png
```

You can decode and extract the epub file by running
```
./underbyte -decode -file tolstoy_war_peace.png -out decoded_war_peace.epub
```

### Secret keys

You can supply a secret key when coding and decoding. This secret is used to seed randomized pixel cursor
that determines which pixels contain which bytes of the encoded message.


## Installation
### Building from source
```
go build -o underbyte cmd/main.go
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
  -secret string
        secret key value used to encode and decode message

Examples:
  Encoding:
	underbyte -file image.png -message "hello there" -out encoded_image.png
	underbyte -file image.png -message "hello there" > encoded_image.png
	cat somefile.txt | underbyte -file image.png > encoded_image.png
  Decoding:
  	underbyte -decode -file encoded_image.png
	underbyte -decode -file encoded_image.png -out decoded_image.png
```
