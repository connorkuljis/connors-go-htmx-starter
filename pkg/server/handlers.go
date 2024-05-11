package server

import (
	"net/http"
	"time"

	"github.com/connorkuljis/connors-go-htmx-starter/pkg/middleware"
)

// Routes instatiates http Handlers and associated patterns on the server.
func (s *Server) Routes() {
	s.MuxRouter.Handle("/static/", http.StripPrefix("/static/", s.StaticContentHandler))
	s.MuxRouter.HandleFunc("/", middleware.Logging(s.Logger, s.HandleIndex()))
}

func (s *Server) HandleIndex() http.HandlerFunc {
	index := []string{
		s.TemplateStore.Base["root.html"],
		s.TemplateStore.Base["head.html"],
		s.TemplateStore.Base["layout.html"],
		s.TemplateStore.Components["nav.html"],
		s.TemplateStore.Components["header.html"],
		s.TemplateStore.Components["footer.html"],
		s.TemplateStore.Views["index.html"],
	}

	t := s.ParseTemplates("index.html", nil, index...)

	return func(w http.ResponseWriter, r *http.Request) {
		parcel := map[string]any{"Time": time.Now()}
		htmlBytes, err := SafeTmplExec(t, "root", parcel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		SendHTML(w, htmlBytes)
	}
}
