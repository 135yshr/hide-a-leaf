# hide-a-leaf

## Requirement

* Go 1.14+

## Installation

you need to have Go's runtime installed. Please install it by referring to the [Official Site] (https://golang.org/) in advance.

```bash
go get github.com/135yshr/hide-a-leaf
go install github.com/135yshr/hide-a-leaf
```

## Usage

The program provides two functions, an encode function that hides the data you want to hide in the specified PNG image and a decode function that retrieves the hidden data from the image.

### encode

If you specify the original PNG file and the file you want to hide as parameters, encode.png will be created.

```bash
hide-a-leaf encode <cover file> <hide data>
```

If you want to hide the character string, specify the `-text` parameter and handle the second parameter as a character string.

```bash
hide-a-leaf -text encode {元になるPNGファイル} {隠したい文字列}
```

### decode

If you specify a PNG file that hides the data you do not want to show and perform decoding, decode.png will be created.

```bash
hide-a-leaf decode <stego file>
```

If you want to get it as a string, specify the -text parameter.

```bash
hide-a-leaf -text decode <stego file>
```

## Author

- 135yshr <isago@oreha.dev>

## License

hide-a-leaf is under [Apache 2.0](LICENSE)
