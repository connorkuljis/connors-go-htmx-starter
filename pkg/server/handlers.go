package server

import (
	"bytes"
	"log"
	"net/http"
	"text/template"
)

func (s *server) handleIndex(tmpl *template.Template) http.HandlerFunc {
	data := map[string]interface{}{
		"SiteData":  s.GetSiteData(),
		"PageTitle": "Index",
		"Username":  "connorkuljis",
		"Option":    "dev",
		"Offset":    "-2",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := executeTemplateToBuffer(tmpl, "root", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		renderHTML(w, buf)
	}
}

func (s *server) handleProjects(tmpl *template.Template) http.HandlerFunc {
	data := map[string]interface{}{
		"SiteData":  s.GetSiteData(),
		"Username":  "connorkuljis",
		"PageTitle": "Projects Page!",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := executeTemplateToBuffer(tmpl, "root", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		renderHTML(w, buf)
	}
}

// Checks if a template is broken / missing keys.
// IO writes to buffer, rather than the http.ResponseWriter.
// To prevent 200 OK sent if an error occurs.
func executeTemplateToBuffer(tmpl *template.Template, name string, data any) (bytes.Buffer, error) {
	var buf bytes.Buffer
	tmpl.Option("missingkey=error")
	err := tmpl.ExecuteTemplate(&buf, "root", data)
	if err != nil {
		log.Print(err)
		return buf, err
	}
	return buf, nil
}

// Sends buffer to http.ResponseWriter as HTML
func renderHTML(w http.ResponseWriter, buf bytes.Buffer) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	buf.WriteTo(w)
}
