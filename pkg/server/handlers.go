package server

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/connorkuljis/connors-go-htmx-starter/pkg/fragments"
)

func (s *Server) HandleIndex() http.HandlerFunc {
	tmpl := s.BuildTemplates("index", nil, fragments.IndexTemplate()...)

	data := map[string]interface{}{
		"PageTitle": "Index",
		"Username":  "connorkuljis",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := SafeTmplExec(tmpl, "root", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		SendHTML(w, buf)
	}
}

// safeTmplParse executes a given template to a bytes buffer. It returns the resulting buffer or nil, err if any error occurred.
//
// Templates are checked for missing keys to prevent partial data being written to the writer.
func SafeTmplExec(tmpl *template.Template, name string, data any) ([]byte, error) {
	var buf bytes.Buffer
	tmpl.Option("missingkey=error")
	err := tmpl.ExecuteTemplate(&buf, name, data)
	if err != nil {
		log.Print(err)
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}

// sendHTML writes a buffer a response writer as html
func SendHTML(w http.ResponseWriter, buf []byte) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	_, err := w.Write(buf)
	if err != nil {
		log.Println(err)
	}
}
