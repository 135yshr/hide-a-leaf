package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"os"

	"image/png"

	leaf "github.com/135yshr/hide-a-leaf"
)

var (
	textMode bool
)

func main() {
	flag.Parse()

	mode := flag.Args()[0]
	coverPath := flag.Args()[1]
	cover, err := openImage(coverPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch mode {
	case "encode":
		if flag.NArg() < 3 {
			printHelp("Error in the parameter specification")
			return
		}
		text := flag.Args()[2]
		var data []byte
		if textMode {
			data = []byte(text)
		} else {
			bs, err := ioutil.ReadFile(text)
			if err != nil {
				fmt.Println(err)
				return
			}
			enc := base64.StdEncoding.EncodeToString(bs)
			data = []byte(enc)
		}
		img, err := leaf.Encode(cover, data)
		if err != nil {
			fmt.Println(err)
			return
		}
		f, err := os.Create("encode.png")
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := png.Encode(f, img); err != nil {
			fmt.Println(err)
			return
		}

	case "decode":
		data := leaf.Decode(cover)
		if textMode {
			fmt.Println(string((data)))
		} else {
			imgDat, err := base64.StdEncoding.DecodeString(string(data))
			if err != nil {
				fmt.Println(err)
				return
			}

			f, err := os.Create("decode.png")
			if err != nil {
				fmt.Println(err)
				return
			}

			if _, err := f.Write(imgDat); err != nil {
				fmt.Println(err)
				return
			}
		}
	default:
		printHelp("Error in the parameter specification")
		return
	}
}

func init() {
	flag.BoolVar(&textMode, "text", false, "Treats the specified parameter as a character string and hides it in the image")
}

func printHelp(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Println("hide-a-leaf [encode|decode] forest.png leaf.png")
	flag.PrintDefaults()
}

func openImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}
