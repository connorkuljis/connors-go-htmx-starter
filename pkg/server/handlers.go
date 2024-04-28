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
	indexTemplateFragments := []string{
		s.TemplateFragments[Base]["root.html"],
		s.TemplateFragments[Base]["layout.html"],
		s.TemplateFragments[Base]["head.html"],
		s.TemplateFragments["components"]["footer.html"],
		s.TemplateFragments["components"]["nav.html"],
		s.TemplateFragments["components"]["header.html"],
		s.TemplateFragments["views"]["index.html"],
	}

	indexTemplate := s.BuildTemplates("index", nil, indexTemplateFragments...)

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
