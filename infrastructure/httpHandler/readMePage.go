package httpHandler

import (
	"code-fabrik.com/bend/infrastructure/markdown"
	"log/slog"
	"net/http"
)

type ReadMePage struct {
	Markdown []byte
}

func (rs ReadMePage) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			slog.Error("panic in ReadMePage", "recover", rec)
		}
	}()

	markdown.PresentMarkdown(w, rs.Markdown)
}
