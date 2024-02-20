package www

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/sfomuseum/go-activitypub"
	"golang.org/x/image/font/gofont/goregular"
)

type IconHandlerOptions struct {
	AccountsDatabase activitypub.AccountsDatabase
	URIs             *activitypub.URIs
}

func IconHandler(opts *IconHandlerOptions) (http.Handler, error) {

	font_size := 48.0
	im_w := 48
	im_h := 48

	f, err := truetype.Parse(goregular.TTF)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse font, %w", err)
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size: font_size,
	})

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := LoggerWithRequest(req, nil)

		account_name, _, err := activitypub.ParseAddressFromRequest(req)

		if err != nil {
			logger.Error("Failed to parse address from request", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("account name", account_name)

		acct, err := opts.AccountsDatabase.GetAccountWithName(ctx, account_name)

		if err != nil {

			logger.Error("Failed to retrieve account", "error", err)

			if err == activitypub.ErrNotFound {
				http.Error(rsp, "Not found", http.StatusNotFound)
				return
			}

			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		logger = logger.With("account id", acct.Id)

		data := []byte(account_name)
		hash := fmt.Sprintf("%x", md5.Sum(data))
		hex := hash[0:6]

		logger = logger.With("hex", hex)
		values, err := strconv.ParseUint(string(hex), 16, 32)

		if err != nil {
			logger.Error("Failed to parse hex value", "err", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		r := uint8(values >> 16)
		g := uint8((values >> 8) & 0xFF)
		b := uint8(values & 0xFF)

		im := image.NewRGBA(image.Rect(0, 0, im_w, im_h)) // x1,y1,  x2,y2 of background rectangle
		im_c := color.RGBA{r, g, b, 255}                  //  R, G, B, Alpha

		draw.Draw(im, im.Bounds(), &image.Uniform{im_c}, image.ZP, draw.Src)

		// https://pkg.go.dev/github.com/fogleman/gg

		dc := gg.NewContext(im_w, im_h)
		dc.DrawImage(im, 0, 0)

		dc.SetFontFace(face)

		if err != nil {
			logger.Error("Failed to load font", "error", err)

			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		x := float64(im_w / 2)
		y := float64((im_w / 2) - 7)

		max_w := float64(im_w)
		dc.SetColor(color.White)

		text := strings.ToUpper(account_name[0:1])

		dc.DrawStringWrapped(text, x, y, 0.5, 0.5, max_w, 1.5, gg.AlignCenter)

		final_im := dc.Image()

		rsp.Header().Set("Content-type", "image/png")

		err = png.Encode(rsp, final_im)

		if err != nil {
			logger.Error("Failed to encode PNG icon", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		return
	}

	return http.HandlerFunc(fn), nil
}
