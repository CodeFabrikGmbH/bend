package htmlTemplate

import (
	"code-fabrik.com/bend/domain/dashboardpage"
	"html/template"
	"net/http"
)

type DashBoardPage struct {
}

func (dbp DashBoardPage) Present(w http.ResponseWriter, board dashboardpage.DashBoard) {
	tmpl := template.Must(template.ParseFiles("resources/dashboard.html"))
	_ = tmpl.Execute(w, board)
}
