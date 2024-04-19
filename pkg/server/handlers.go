package server

import (
	"net/http"
)

func (s *Server) HandleIndex() http.HandlerFunc {
	fragments := []string{
		s.TemplateFragments.Base["root.html"],
		s.TemplateFragments.Base["layout.html"],
		s.TemplateFragments.Base["head.html"],

		s.TemplateFragments.Components["footer.html"],
		s.TemplateFragments.Components["nav.html"],
		s.TemplateFragments.Components["header.html"],

		s.TemplateFragments.Views["index.html"],
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
