package htmlTemplate

import (
	"code-fabrik.com/bend/domain/loginpage"
	"html/template"
	"net/http"
)

type LoginPage struct {
}

func (dbp LoginPage) Present(w http.ResponseWriter, login loginpage.Login) {
	tmpl := template.Must(template.ParseFiles("resources/login.html"))
	_ = tmpl.Execute(w, login)
}
