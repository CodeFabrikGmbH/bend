package httpHandler

import (
	"code-fabrik.com/bend/infrastructure/markdown"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ReadMePage struct {
	MarkdownFile string
}

func (rs ReadMePage) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
		}
	}()

	data, err := ioutil.ReadFile(rs.MarkdownFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	markdown.PresentMarkdown(w, data)
}
