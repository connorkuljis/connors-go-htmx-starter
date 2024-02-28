package server

import (
	"net/http"
	"text/template"
)

// Routes instatiates http Handlers and associated patterns on the server.
func (s *server) Routes() {
	s.Router.Handle("/static/", http.FileServer(http.FS(s.FileSystem)))

	s.Router.HandleFunc("/", s.handleIndex(s.getIndexTemplate()))
}

func (s *server) GetSiteData() map[string]interface{} {
	return map[string]interface{}{
		"Title":          s.Title,
		"DevModeEnabled": s.DevModeEnabled,
	}
}

func (s *server) handleIndex(tmpl *template.Template) http.HandlerFunc {
	data := map[string]interface{}{
		"PageTitle": "Index",
		"SiteData":  s.GetSiteData(),
		"Username":  "connorkuljis",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := safeTmplParse(tmpl, "root", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendHTML(w, buf)
	}
}
