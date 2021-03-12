package markdown

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"net/http"
)

func PresentMarkdown(w http.ResponseWriter, markdownData []byte) {
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	htmlData := markdown.ToHTML(markdownData, nil, renderer)
	_, _ = w.Write(htmlData)
}
