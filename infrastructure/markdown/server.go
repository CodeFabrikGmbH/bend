package markdown

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"io/ioutil"
	"net/http"
)

type Server struct {
	File string
}

func FileServer(file string) Server {
	return Server{File: file}
}

func (rs Server) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
		}
	}()

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	data, err := ioutil.ReadFile(rs.File)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	htmlData := markdown.ToHTML(data, nil, renderer)
	_, _ = w.Write(htmlData)
}
