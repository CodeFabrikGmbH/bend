package application

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"io/ioutil"
	"net/http"
)

type ReadmeService struct {
}

func (rs ReadmeService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
		}
	}()

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	data, err := ioutil.ReadFile("README.md")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	html := markdown.ToHTML(data, nil, renderer)
	w.Write(html)
}
