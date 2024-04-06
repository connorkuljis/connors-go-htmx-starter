package main

import (
	"embed"
	"log"

	"github.com/connorkuljis/connors-go-htmx-starter/pkg/server"
)

//go:embed www
var embedFS embed.FS

func main() {
	s, err := server.NewServer(embedFS, "8080")
	if err != nil {
		log.Fatal(err)
	}

	s.Routes()
}
