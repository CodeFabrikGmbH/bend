package htmlTemplate

import (
	"code-fabrik.com/bend/domain/loginpage"
	"html/template"
	"net/http"
)

type LoginPagePresenter struct {
}

func (dbp LoginPagePresenter) Present(w http.ResponseWriter, login loginpage.Login) {
	tmpl := template.Must(template.ParseFiles("resources/login.html"))
	_ = tmpl.Execute(w, login)
}
