package server

import (
	"net/http"
)

// Routes instatiates http Handlers and associated patterns on the server.
func (s *Server) Routes() {
	s.Router.Handle("/static/", http.FileServer(http.FS(s.FileSystem)))
	s.Router.HandleFunc("/", s.HandleIndex())
}

func (s *Server) HandleIndex() http.HandlerFunc {
	tmpl := IndexTemplate(s)

	data := map[string]interface{}{
		"AppData":   s.AppData,
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
