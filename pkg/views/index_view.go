package views

import (
	"html/template"
	"log"

	"github.com/connorkuljis/connors-go-htmx-starter/pkg/webserver"
)

func IndexView(s *webserver.WebServer) *template.Template {
	name := "index.html"
	index := []string{
		s.TemplateStore.Base["root.html"],
		s.TemplateStore.Base["head.html"],
		s.TemplateStore.Base["layout.html"],
		s.TemplateStore.Components["nav.html"],
		s.TemplateStore.Components["header.html"],
		s.TemplateStore.Components["footer.html"],
		s.TemplateStore.Views["index.html"],
	}

	t, err := webserver.ParseTemplates(name, nil, s.FileSystem, index...)
	if err != nil {
		log.Fatal(err)
	}

	return t
}
