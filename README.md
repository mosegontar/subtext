# underbyte

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
	underbyte -file image.jpg -message "hello there" -out encoded_image.jpg
	underbyte -file image.jpg -message "hello there" > encoded_image.jpg
	cat somefile.txt | underbyte -file image.jpg > encoded_image.jpg
  Decoding:
  	underbyte -decode -file encoded_image.jpg
	underbyte -decode -file encoded_image.jpg -out decoded_image.jpg
```
