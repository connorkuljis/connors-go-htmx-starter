package server

import (
	"net/http"
)

type AppData struct {
	Title   string
	DevMode bool
}

// Routes instatiates http Handlers and associated patterns on the server.
func (s *server) Routes() {
	appData := AppData{
		Title:   "Epic Title",
		DevMode: false,
	}

	s.Router.Handle("/static/", http.FileServer(http.FS(s.FileSystem)))
	s.Router.HandleFunc("/", s.handleIndex(appData))
}

func (s *server) handleIndex(appData AppData) http.HandlerFunc {
	tmpl := IndexTemplate(s)

	data := map[string]interface{}{
		"PageTitle": "Index",
		"AppData":   appData,
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
