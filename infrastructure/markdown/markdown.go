package markdown

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

// RenderMarkdown converts markdown source to an HTML fragment for embedding into
// a page template.
func RenderMarkdown(markdownData []byte) []byte {
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.ToHTML(markdownData, nil, renderer)
}
