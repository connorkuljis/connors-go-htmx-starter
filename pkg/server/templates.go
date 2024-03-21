package server

import (
	"text/template"
)

// getIndexTemplate parses joined base and index view templates.
func IndexTemplate(s *Server) *template.Template {
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

	return BuildTemplates(s, "index.html", nil, view...)
}
