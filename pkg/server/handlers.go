package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
)

// Routes instatiates http Handlers and associated patterns on the server.
func (s *Server) Routes() error {
	scfs, err := fs.Sub(s.FileSystem, StaticDirStr)
	if err != nil {
		return err
	}

	s.MuxRouter.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(scfs))))
	s.MuxRouter.HandleFunc("/", s.HandleIndex())
	s.MuxRouter.HandleFunc("/counter", s.HandleCounter())

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

	data := map[string]any{
		"Count": 0,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		htmlBytes, err := SafeTmplExec(indexTemplate, "root", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		SendHTML(w, htmlBytes)
	}
}

func (s *Server) HandleCounter() http.HandlerFunc {
	fragment := s.TemplateFragments.Views["index.html"]

	counterTemplate := s.BuildTemplates("counter", nil, fragment)

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		strCount := r.Form.Get("count")
		if strCount == "" {
			http.Error(w, "Error: no form value for 'count'", http.StatusInternalServerError)
			return
		}

		count, err := strconv.Atoi(strCount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		op := r.URL.Query().Get("op")
		switch op {
		case "add":
			count++
			fmt.Fprintf(w, "Add operation")
		case "subtract":
			count--
			fmt.Fprintf(w, "Subtract operation")
		default:
			fmt.Fprintf(w, "Updated count to %d", count)
		}

		data := map[string]any{
			"Count": count,
		}

		htmlBytes, err := SafeTmplExec(counterTemplate, "counter", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		SendHTML(w, htmlBytes)
	}
}
