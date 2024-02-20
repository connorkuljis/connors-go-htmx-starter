// FILE: html.go
// purpose:

package server

import (
	"log"
	"text/template"
)

const (
	headHTML   = "templates/head.html"
	layoutHTML = "templates/layout.html"
	rootHTML   = "templates/root.html"
	heroHTML   = "templates/components/hero.html"
	footerHTML = "templates/components/footer.html"
	navHTML    = "templates/components/nav.html"
	indexHTML  = "templates/views/index.html"
)

type HTMLFile string

// A View is a collection of HTML Files
type View []HTMLFile

// Common html files used to compile a base view to render a view
func baseView() View {
	baseView := View{rootHTML, headHTML, layoutHTML, navHTML, footerHTML}
	return baseView
}

func (s *server) indexView() *template.Template {
	return s.compile("index", append(baseView(), indexHTML), nil)
}

// compiles a template from a view.
func (s *server) compile(name string, view View, funcs template.FuncMap) *template.Template {
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
