// file: html.go
// purpose: type safe html files

package main

import (
	"log"
	"text/template"
)

type HTMLFile string

type View []HTMLFile

// base html files
const (
	HeadHTML   HTMLFile = "templates/head.html"
	LayoutHTML HTMLFile = "templates/layout.html"
	RootHTML   HTMLFile = "templates/root.html"
)

// component html files
const (
	HeroHTML   HTMLFile = "templates/components/hero.html"
	FooterHTML HTMLFile = "templates/components/footer.html"
	NavHTML    HTMLFile = "templates/components/nav.html"
)

// view html files
const (
	indexHTML HTMLFile = "templates/views/index.html"
)

// views
var (
	baseView  View = View{RootHTML, HeadHTML, LayoutHTML, NavHTML, FooterHTML}
	indexView View = append(baseView, indexHTML)
)

// compiles a template from a view.
func (s *Server) compile(name string, view View, funcs template.FuncMap) *template.Template {
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
