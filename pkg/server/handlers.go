package server

import (
	"io/fs"
	"net/http"
)

// Routes instatiates http Handlers and associated patterns on the server.
func (s *Server) Routes() error {
	scfs, err := fs.Sub(s.FileSystem, StaticDirStr)
	if err != nil {
		return err
	}

	s.MuxRouter.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(scfs))))
	s.MuxRouter.HandleFunc("/", s.HandleIndex())

	return nil
}

func (s *Server) HandleIndex() http.HandlerFunc {
	indexTemplateFragments := []string{
		s.TemplateFragments.Base["root.html"],
		s.TemplateFragments.Base["layout.html"],
		s.TemplateFragments.Base["head.html"],
		s.TemplateFragments.Components["footer.html"],
		s.TemplateFragments.Components["nav.html"],
		s.TemplateFragments.Components["header.html"],
		s.TemplateFragments.Views["index.html"],
	}

	indexTemplate := s.BuildTemplates("index", nil, indexTemplateFragments...)

	return func(w http.ResponseWriter, r *http.Request) {
		htmlBytes, err := SafeTmplExec(indexTemplate, "root", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		SendHTML(w, htmlBytes)
	}
}
