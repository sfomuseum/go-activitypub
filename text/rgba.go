package text

import (
	"crypto/md5"
	"fmt"
	"image/color"
	"math"
	"strconv"
)

// Convert 'body' in to a RGB hex string and finally a `color.RGBA` instance.
func TextToRGBAColor(body string) (*color.RGBA, error) {

	data := []byte(body)
	hash := fmt.Sprintf("%x", md5.Sum(data))
	hex := hash[0:6]

	values, err := strconv.ParseUint(string(hex), 16, 32)

	if err != nil {
		return nil, err
	}

	r64 := values >> 16
	g64 := (values >> 8) & 0xFF
	b64 := values & 0xFF

	if r64 > math.MaxUint8 {
		return nil, err
	}

	if g64 > math.MaxUint8 {
		return nil, err
	}

	if b64 > math.MaxUint8 {
		return nil, err
	}

	r := uint8(r64)
	g := uint8(g64)
	b := uint8(b64)

	im_c := color.RGBA{r, g, b, 255} //  R, G, B, Alpha
	return &im_c, nil
}
