package html

import (
	"context"
	"embed"

	sfomuseum_html "github.com/sfomuseum/go-template/html"
	"html/template"
)

//go:embed *.html
var FS embed.FS

func LoadTemplates(ctx context.Context) (*template.Template, error) {

	return sfomuseum_html.LoadTemplates(ctx, FS)
}
