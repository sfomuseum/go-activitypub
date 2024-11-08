package www

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aaronland/gocloud-blob/bucket"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/text"
	"github.com/sfomuseum/go-activitypub/uris"
	"golang.org/x/image/font/gofont/goregular"
)

// Reconcile with app/account/add/add.go

var re_data_url = regexp.MustCompile(`^data:image\/[^;]+;base64,(.*)`)
var re_http_url = regexp.MustCompile(`^https?\:\/\/(.*)`)

type IconHandlerOptions struct {
	AccountsDatabase database.AccountsDatabase
	URIs             *uris.URIs
	AllowRemote      bool
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

		t1 := time.Now()

		defer func() {
			logger.Info("Time to serve request", "ms", time.Since(t1).Milliseconds())
		}()

		if req.Method != http.MethodGet {
			logger.Error("Method not allowed")
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		account_name, host, err := ap.ParseAddressFromRequest(req)

		if err != nil {
			logger.Error("Failed to parse address from request", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("account name", account_name)

		if host != "" && host != opts.URIs.Hostname {
			logger.Error("Resouce has bunk hostname", "host", host)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

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

		// START OF check to see if there is a custom account icon image
		// START OF move this in to a separate package or something

		// Until then this is a twisty mess of if/else blocks because we always want
		// to fail over to the auto-generate icons if necessary. Womp womp...

		switch {
		case re_http_url.MatchString(acct.IconURI):

			logger.Debug("Account image is a pointer to a remote address", "address", acct.IconURI)

			if opts.AllowRemote {

				icon_u, err := url.Parse(acct.IconURI)

				if err != nil {
					logger.Error("Failed to parse remote address for account image", "error", err)
				} else {
					http.Redirect(rsp, req, icon_u.String(), http.StatusSeeOther)
					return
				}

			} else {
				logger.Error("Server configuration disallows remote account icon images")
			}

		case re_data_url.MatchString(acct.IconURI):

			logger.Debug("Account image matches base64-encoded data URL")

			m := re_data_url.FindStringSubmatch(acct.IconURI)
			b64 := m[1]

			data, err := base64.StdEncoding.DecodeString(b64)

			if err != nil {
				logger.Error("Failed to decode account image URL from base64", "error", err)
			} else {

				im_r := bytes.NewReader(data)

				im, _, err := image.Decode(im_r)

				if err != nil {
					logger.Error("Failed to decode account image URL as image", "error", err)
				} else {

					err := png.Encode(rsp, im)

					if err != nil {
						logger.Error("Failed to encode account image URL as PNG", "error", err)
					}
				}
			}

		case acct.IconURI != "":

			icon_u, err := url.Parse(acct.IconURI)

			if err != nil {
				logger.Error("Failed to parse icon url", "error", err)
			} else {

				// START OF put me in a function (probably in aaronland/gocloud-blob/bucket)
				root := filepath.Dir(icon_u.Path)
				fname := filepath.Base(icon_u.Path)

				root = strings.TrimLeft(root, "/")
				root = strings.TrimRight(root, "/")

				icon_q := icon_u.Query()
				icon_q.Set("prefix", fmt.Sprintf("%s/", root))

				icon_u.Path = ""
				icon_u.RawQuery = icon_q.Encode()

				bucket_uri := icon_u.String()
				// END OF put me in a function (probably in aaronland/gocloud-blob/bucket)

				logger = logger.With("bucket uri", bucket_uri)
				logger = logger.With("bucket key", fname)

				b, err := bucket.OpenBucket(ctx, bucket_uri)

				if err != nil {
					logger.Error("Failed to open bucket", "bucket uri", bucket_uri, "error", err)
				} else {

					defer b.Close()

					attrs, err := b.Attributes(ctx, fname)

					if err != nil {
						logger.Error("Failed to derive attributes for icon", "error", err)
					} else {

						r, err := b.NewReader(ctx, fname, nil)

						if err != nil {
							logger.Error("Failed to open icon file for reading", "filename", fname, "error", err)
						} else {

							defer r.Close()

							rsp.Header().Set("Content-Type", attrs.ContentType)

							_, err := io.Copy(rsp, r)

							if err != nil {
								logger.Error("Failed to copy icon file", "error", err)
							}

							return
						}
					}
				}
			}

		default:
			logger.Debug("No custom account image URL, generate on demand")
		}

		// END OF move this in to a separate package or something
		// END OF check to see if there is a custom account icon image

		// START OF replace with icon/GenerateIcon

		im_c, err := text.TextToRGBAColor(account_name)

		if err != nil {
			logger.Error("Failed to derive colour from text", "err", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		im := image.NewRGBA(image.Rect(0, 0, im_w, im_h)) // x1,y1,  x2,y2 of background rectangle

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

		// END OF replace with icon/GenerateIcon

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
