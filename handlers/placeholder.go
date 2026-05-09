package handlers

import (
	"html/template"
	"net/http"
	"strings"
)

var tmplPlaceholder *template.Template

func InitPlaceholder(tmpl *template.Template) {
	tmplPlaceholder = tmpl
}

func PlaceholderPage(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/")
	section := strings.ToUpper(name[:1]) + name[1:]
	tmplPlaceholder.ExecuteTemplate(w, "placeholder.html", map[string]interface{}{
		"Page":    name,
		"Section": section,
	})
}
