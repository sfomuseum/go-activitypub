package html

import (
	"strings"

	"golang.org/x/net/html"
)

// HtmlToText converts HTML to plain text by stripping all HTML tags.
func HtmlToText(htmlStr string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", err
	}
	var textContent strings.Builder
	extractText(doc, &textContent)
	return textContent.String(), nil
}

// extractText recursively traverses HTML nodes and appends text to the result.
func extractText(n *html.Node, textContent *strings.Builder) {
	// If this is a text node, append it to the text content.
	if n.Type == html.TextNode {
		textContent.WriteString(n.Data)
	}

	// Recursively apply to all child nodes.
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, textContent)
	}
}
