package server

import (
	"bytes"
	"log"
	"net/http"
	"text/template"
)

// getIndexTemplate parses joined base and index view templates.
func IndexTemplate(s *server) *template.Template {
	view := []string{
		s.Templates.BaseLayout.Head,
		s.Templates.BaseLayout.Root,
		s.Templates.BaseLayout.Layout,
		s.Templates.Components.DevTool,
		s.Templates.Components.Header,
		s.Templates.Components.Footer,
		s.Templates.Components.Nav,
		s.Templates.Views.Index,
	}

	return buildTemplates(s, "index.html", nil, view...)
}

// getProjectsTemplatePartial parses the projects template partial (incomplete html for htmx swapping)
// and the joined base and projects view templates including the projects template partial.
// func (s *server) getProjectsTemplatePartial() map[string]*template.Template {
// 	projects := "templates/partials/projects.html"
// 	view := "templates/views/projects.html"

// 	return map[string]*template.Template{
// 		"projects": buildTemplates(s, "projects", nil, projects),
// 		"full":     buildTemplates(s, "full", nil, s.getBaseTemplates(view, projects)...),
// 	}
// }

// buildTemplates is a fast way to parse a collection of templates in the server filesystem.
func buildTemplates(s *server, name string, funcs template.FuncMap, templates ...string) *template.Template {
	// give the template a name
	tmpl := template.New(name)

	// custom template functions if exists
	if funcs != nil {
		tmpl.Funcs(funcs)
	}

	// generate a template from the files in the server fs (usually embedded)
	tmpl, err := tmpl.ParseFS(s.FileSystem, templates...)
	if err != nil {
		log.Fatal(err)
	}

	return tmpl
}

// safeTmplParse executes a given template to a bytes buffer. It returns the resulting buffer or nil, err if any error occurred.
//
// Templates are checked for missing keys to prevent partial data being written to the writer.
func safeTmplParse(tmpl *template.Template, name string, data any) (bytes.Buffer, error) {
	var buf bytes.Buffer
	tmpl.Option("missingkey=error")
	err := tmpl.ExecuteTemplate(&buf, name, data)
	if err != nil {
		log.Print(err)
		return buf, err
	}
	return buf, nil
}

// sendHTML writes a buffer a response writer as html
func sendHTML(w http.ResponseWriter, buf bytes.Buffer) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	buf.WriteTo(w)
}
