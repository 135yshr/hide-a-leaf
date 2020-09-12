package leaf

import (
	"image"
	"image/color"
)

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

// Encode - cover に data を隠す
func Encode(cover image.Image, data []byte) (image.Image, error) {
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
			newImg.SetNRGBA(x, y, hiding(baseColor, data, index))
			index++
		}
	}
	return newImg, nil
}

// Decode - 秘密データが埋め込まれた画像からデータを取得する
func Decode(stego image.Image) []byte {
	bs := []byte{}
	rect := stego.Bounds()
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			i := color.NRGBAModel.Convert(stego.At(x, y))
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
