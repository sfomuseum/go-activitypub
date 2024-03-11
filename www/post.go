package www

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/uris"
)

type PostHandlerOptions struct {
	AccountsDatabase activitypub.AccountsDatabase
	PostsDatabase    activitypub.PostsDatabase
	URIs             *uris.URIs
	Templates        *template.Template
}

type PostHandlerVars struct {
	Post       *activitypub.Post
	PostBody   template.HTML
	Account    *activitypub.Account
	AccountURL string
	IconURL    string
}

func PostHandler(opts *PostHandlerOptions) (http.Handler, error) {

	post_t := opts.Templates.Lookup("post")

	if post_t == nil {
		return nil, fmt.Errorf("Failed to lookup 'post' template")
	}

	post_pat := opts.URIs.Post
	post_pat = strings.Replace(post_pat, "{resource}", "(?:[^\\/]+)", 1)
	post_pat = strings.Replace(post_pat, "{id}", "(\\d+)", 1)

	re_post, err := regexp.Compile(post_pat)

	if err != nil {
		return nil, fmt.Errorf("Failed to create post regular expression, %w", err)
	}

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

		// Basic sanity checking of post ID

		if !re_post.MatchString(req.URL.Path) {
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		m := re_post.FindStringSubmatch(req.URL.Path)

		str_id := m[1]
		post_id, err := strconv.ParseInt(str_id, 10, 64)

		if err != nil {
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("post id", post_id)

		// Get account

		account_name, host, err := activitypub.ParseAddressFromRequest(req)

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

		// Get post

		post, err := opts.PostsDatabase.GetPostWithId(ctx, post_id)

		if err != nil {
			logger.Error("Failed to retrieve post", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		if post.AccountId != acct.Id {
			logger.Error("Post is owned by different account", "post account id", post.AccountId)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		acct.PrivateKeyURI = "constant://?val="

		account_url := acct.AccountURL(ctx, opts.URIs)

		icon_path := uris.AssignResource(opts.URIs.Icon, acct.Name)
		icon_url := uris.NewURL(opts.URIs, icon_path)

		// Render template

		vars := PostHandlerVars{
			Account:    acct,
			Post:       post,
			PostBody:   template.HTML(post.Body),
			IconURL:    icon_url.String(),
			AccountURL: account_url.String(),
		}

		rsp.Header().Set("Content-Type", "text/html")

		err = post_t.Execute(rsp, vars)

		if err != nil {
			logger.Error("Failed to render template", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	return http.HandlerFunc(fn), nil
}
