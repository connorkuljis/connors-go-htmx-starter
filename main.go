package main

import (
	"embed"
	"log"

	"github.com/connorkuljis/connors-go-htmx-starter/pkg/server"
)

const (
	TemplatesDir = "www/templates"
	StaticDir    = "www/static"
	Port         = "8080"
)

//go:embed www
var embedFS embed.FS

func main() {
	s, err := server.NewServer(embedFS, Port, TemplatesDir, StaticDir)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Routes(); err != nil {
		log.Fatal(err)
	}

	if err = s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
