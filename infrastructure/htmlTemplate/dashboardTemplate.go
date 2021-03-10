package htmlTemplate

import (
	"code-fabrik.com/bend/domain/dashboard"
	"html/template"
	"net/http"
)

type DashBoardPresenter struct {
}

func (dbp DashBoardPresenter) Present(w http.ResponseWriter, board dashboard.DashBoard) {
	tmpl := template.Must(template.ParseFiles("resources/dashboard.html"))
	_ = tmpl.Execute(w, board)
}
