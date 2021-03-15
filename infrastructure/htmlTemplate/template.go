package htmlTemplate

import (
	"html/template"
	"net/http"
)

func PresentHtmlTemplate(w http.ResponseWriter, templateName string, templateDate interface{}) {
	tmpl := template.Must(template.ParseFiles(templateName))
	_ = tmpl.Execute(w, templateDate)
}
