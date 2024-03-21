package server

import (
	"net/http"
)

func (s *Server) HandleIndex() http.HandlerFunc {
	tmpl := IndexTemplate(s)

	data := map[string]interface{}{
		"AppData":   s.AppData,
		"PageTitle": "Index",
		"Username":  "connorkuljis",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := SafeTmplParse(tmpl, "root", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		SendHTML(w, buf)
	}
}
