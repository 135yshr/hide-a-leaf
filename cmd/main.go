package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"os"

	"image/color"
	"image/png"
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
		img, err := encode(cover, data)
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
		data := decode(cover)
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

func encode(cover image.Image, text []byte) (image.Image, error) {
	rect := cover.Bounds()
	newImg := image.NewNRGBA(image.Rectangle{rect.Min, rect.Max})
	index := 0
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			c1 := color.NRGBAModel.Convert(cover.At(x, y))
			baseColor, ok := c1.(color.NRGBA)
			if !ok {
				continue
			}
			newImg.SetNRGBA(x, y, hiding(baseColor, text, index))
			index++
		}
	}
	return newImg, nil
}

func hiding(c color.NRGBA, text []byte, n int) color.NRGBA {
	var r, g, b, a uint8
	if len(text) <= n || text[n] == 0 {
		r = c.R & 0xfc
		g = c.G & 0xfc
		b = c.B & 0xfc
		a = c.A & 0xfc
	} else {
		r = c.R&0xfc + text[n]&3
		g = c.G&0xfc + (text[n]>>2)&0x3
		b = c.B&0xfc + (text[n]>>4)&0x3
		a = c.A&0xfc + (text[n]>>6)&0x3
	}

	return color.NRGBA{r, g, b, a}
}

func decode(cover image.Image) []byte {
	bs := []byte{}
	rect := cover.Bounds()
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			i := color.NRGBAModel.Convert(cover.At(x, y))
			c, ok := i.(color.NRGBA)
			if !ok {
				continue
			}
			r := c.R & 0x3
			g := c.G & 0x3 << 2
			b := c.B & 0x3 << 4
			a := c.A & 0x3 << 6
			if (r + g + b + a) == 0 {
				continue
			}
			bs = append(bs, byte(r+g+b+a))
		}
	}
	return bs
}
