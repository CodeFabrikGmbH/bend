package httpHandler

import (
	"code-fabrik.com/bend/infrastructure/htmlTemplate"
	"code-fabrik.com/bend/infrastructure/markdown"
	"html/template"
	"log/slog"
	"net/http"
)

type ReadMePage struct {
	Markdown []byte
}

type ReadMeViewData struct {
	Content template.HTML
}

func (rs ReadMePage) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			slog.Error("panic in ReadMePage", "recover", rec)
		}
	}()

	content := template.HTML(markdown.RenderMarkdown(rs.Markdown))
	htmlTemplate.PresentHtmlTemplate(w, "readme.html", ReadMeViewData{Content: content})
}
