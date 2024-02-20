package server

import (
	"io/fs"
	"log"
	"net/http"
	"text/template"
)

func NewServer(port string, router *http.ServeMux, templatesDir string, staticDir string, fileSystem fs.FS) server {
	return server{
		Router:       router,
		Port:         port,
		TemplatesDir: templatesDir,
		StaticDir:    staticDir,
		FileSystem:   fileSystem,
	}
}

// Common html files used to compile a base view to render a view
func baseView() View {
	baseView := View{rootHTML, headHTML, layoutHTML, navHTML, footerHTML}
	return baseView
}

// Returns the template for the index view
func (s *server) indexView() *template.Template {
	return compile(s, "index", append(baseView(), indexHTML), nil)
}

// compiles a template from a view.
func compile(s *server, name string, view View, funcs template.FuncMap) *template.Template {
	// give the template a name
	tmpl := template.New(name)

	// custom template functions if exists
	if funcs != nil {
		tmpl.Funcs(funcs)
	}

	// create a collection of files from the view
	var files []string
	for _, htmlFile := range view {
		files = append(files, string(htmlFile))
	}

	// generate a template from the files in the server fs (usually embedded)
	tmpl, err := tmpl.ParseFS(s.FileSystem, files...)
	if err != nil {
		log.Fatal(err)
	}

	return tmpl
}
