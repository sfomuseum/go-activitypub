// Package html defines methods and default templates for generating HTML pages.
package html

import (
	"context"
	"embed"

	sfomuseum_html "github.com/sfomuseum/go-template/html"
	"html/template"
)

//go:embed *.html
var FS embed.FS

// LoadTemplates will returns a `html/template.Template` instance containing the default templates provided by this package.
func LoadTemplates(ctx context.Context) (*template.Template, error) {

	return sfomuseum_html.LoadTemplates(ctx, FS)
}
