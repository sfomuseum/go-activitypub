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

	"github.com/sfomuseum/go-activitypub"
)

type IconHandlerOptions struct {
	AccountsDatabase activitypub.AccountsDatabase
	URIs             *activitypub.URIs
}

func IconHandler(opts *IconHandlerOptions) (http.Handler, error) {

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

		im := image.NewRGBA(image.Rect(0, 0, 220, 220)) // x1,y1,  x2,y2 of background rectangle
		im_c := color.RGBA{r, g, b, 255}                //  R, G, B, Alpha

		draw.Draw(im, im.Bounds(), &image.Uniform{im_c}, image.ZP, draw.Src)

		// Add text...
		// https://josemyduarte.github.io/2021-02-28-quotes-on-images-with-go/

		rsp.Header().Set("Content-type", "image/png")

		err = png.Encode(rsp, im)

		if err != nil {
			logger.Error("Failed to encode PNG icon", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		return
	}

	return http.HandlerFunc(fn), nil
}
