package server

import (
	"net/http"
	"text/template"
)

func (s *server) handlePartialProjects(tmpl *template.Template) http.HandlerFunc {
	type ViewData struct {
		Projects []string
	}

	data := ViewData{
		Projects: []string{"block-cli", "green-tiles"},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "partial", data)
	}
}
