package server

import (
	"log"
	"text/template"
)

const (
	headHTML            htmlFile = "templates/head.html"
	layoutHTML          htmlFile = "templates/layout.html"
	rootHTML            htmlFile = "templates/root.html"
	heroHTML            htmlFile = "templates/components/hero.html"
	footerHTML          htmlFile = "templates/components/footer.html"
	navHTML             htmlFile = "templates/components/nav.html"
	indexHTML           htmlFile = "templates/views/index.html"
	projectsHTML        htmlFile = "templates/views/projects.html"
	projectsHTMLPartial htmlFile = "templates/partials/projects.html"
)

type htmlFile string

type view []htmlFile

// Common html files used to compile a base view to render a view
func baseView() view {
	baseView := view{rootHTML, headHTML, layoutHTML, navHTML, footerHTML}
	return baseView
}

// Returns the template for the index view
func (s *server) getViewIndex() *template.Template {
	return compile(s, "index", append(baseView(), indexHTML), nil)
}

// Returns the template for the projects view
func (s *server) getViewProjects() *template.Template {
	return compile(s, "projects", append(baseView(), projectsHTML), nil)
}

func (s *server) getPartialProjects() *template.Template {
	return compile(s, "partial-projects", view{projectsHTMLPartial}, nil)
}

// compiles a template from a view.
// returns a template or fatals if template cannot be parsed
func compile(s *server, name string, htmlFiles view, funcs template.FuncMap) *template.Template {
	// give the template a name
	tmpl := template.New(name)

	// custom template functions if exists
	if funcs != nil {
		tmpl.Funcs(funcs)
	}

	// create a collection of filenames from the view
	var filenames []string
	for _, htmlFile := range htmlFiles {
		filenames = append(filenames, string(htmlFile))
	}

	// generate a template from the files in the server fs (usually embedded)
	tmpl, err := tmpl.ParseFS(s.FileSystem, filenames...)
	if err != nil {
		log.Fatal(err)
	}

	return tmpl
}
