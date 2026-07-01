package htmlTemplate

import (
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
)

// templates holds every page template, parsed once at startup instead of being
// read from disk and re-parsed on every request.
var templates *template.Template

// Load parses all templates matching pattern from the given filesystem. It must
// be called once during start-up before any template is presented.
func Load(fsys fs.FS, pattern string) error {
	parsed, err := template.ParseFS(fsys, pattern)
	if err != nil {
		return err
	}
	templates = parsed
	return nil
}

// PresentHtmlTemplate renders the named template (its base file name, e.g.
// "dashboard.html") to the response writer.
func PresentHtmlTemplate(w http.ResponseWriter, templateName string, templateData interface{}) {
	if templates == nil {
		slog.Error("templates not loaded", "template", templateName)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if err := templates.ExecuteTemplate(w, templateName, templateData); err != nil {
		slog.Error("template execution failed", "template", templateName, "err", err)
	}
}
