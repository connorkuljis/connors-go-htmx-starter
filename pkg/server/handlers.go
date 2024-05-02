package server

import (
	"net/http"
	"time"
)

// Routes instatiates http Handlers and associated patterns on the server.
func (s *Server) Routes() error {
	s.MuxRouter.Handle("/static/", http.StripPrefix("/static/", s.StaticContentHandler))
	s.MuxRouter.HandleFunc("/", s.HandleIndex())

	return nil
}

func (s *Server) HandleIndex() http.HandlerFunc {
	indexPage := []string{
		"root.html",
		"head.html",
		"layout.html",
		"nav.html",
		"header.html",
		"footer.html",
		"index.html",
	}

	indexTemplate := s.ParseTemplates("index.html", nil, indexPage...)

	return func(w http.ResponseWriter, r *http.Request) {
		parcel := map[string]any{"Time": time.Now()}
		htmlBytes, err := SafeTmplExec(indexTemplate, "root", parcel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		SendHTML(w, htmlBytes)
	}
}
