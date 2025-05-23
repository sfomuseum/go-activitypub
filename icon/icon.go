package icon

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log/slog"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/sfomuseum/go-activitypub/text"
	"golang.org/x/image/font/gofont/goregular"
)

type GenerateIconOptions struct {
	Label    string
	TrimTo   int
	FontSize float64
	Width    int
	Height   int
}

func GenerateIcon(ctx context.Context, opts *GenerateIconOptions) (image.Image, error) {

	logger := slog.Default()

	font_size := 82.0
	im_w := 200
	im_h := 200

	if opts.FontSize != 0 {
		font_size = opts.FontSize
	}

	if opts.Width != 0 {
		im_w = opts.Width
	}

	if opts.Height != 0 {
		im_h = opts.Height
	}

	f, err := truetype.Parse(goregular.TTF)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse font, %w", err)
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size: font_size,
	})

	im_c, err := text.TextToRGBAColor(opts.Label)

	if err != nil {
		logger.Error("Failed to derive colour from text", "err", err)
		return nil, err
	}

	im := image.NewRGBA(image.Rect(0, 0, im_w, im_h)) // x1,y1,  x2,y2 of background rectangle

	draw.Draw(im, im.Bounds(), &image.Uniform{im_c}, image.ZP, draw.Src)

	// https://pkg.go.dev/github.com/fogleman/gg

	dc := gg.NewContext(im_w, im_h)
	dc.DrawImage(im, 0, 0)

	dc.SetFontFace(face)

	if err != nil {
		logger.Error("Failed to load font", "error", err)
		return nil, err
	}

	x := float64(im_w / 2)
	y := float64((im_w / 2) - 7)

	max_w := float64(im_w)
	dc.SetColor(color.White)

	text := opts.Label

	if opts.TrimTo > 0 {
		text = strings.ToUpper(text[0:opts.TrimTo])
	}

	dc.DrawStringWrapped(text, x, y, 0.5, 0.5, max_w, 1.5, gg.AlignCenter)
	final_im := dc.Image()

	return final_im, nil
}
