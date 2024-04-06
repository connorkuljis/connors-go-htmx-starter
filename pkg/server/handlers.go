package server

import (
	"net/http"
)

func (s *Server) HandleIndex() http.HandlerFunc {
	fragments := []string{
		s.Fragments.Base["root.html"],
		s.Fragments.Base["layout.html"],
		s.Fragments.Base["head.html"],

		s.Fragments.Components["footer.html"],
		s.Fragments.Components["dev-tool.html"],
		s.Fragments.Components["nav.html"],
		s.Fragments.Components["header.html"],

		s.Fragments.Views["index.html"],
	}

	tmpl := s.BuildTemplates("index", nil, fragments...)

	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := SafeTmplExec(tmpl, "root", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		SendHTML(w, buf)
	}
}
