package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/connorkuljis/connors-go-htmx-starter/pkg/server"
)

//go:embed web/templates/* web/static/*
var embedFS embed.FS

func main() {
	server := server.NewServer(embedFS, "8080")
	server.Routes()

	log.Println("[ 💿 Spinning up server on http://localhost:" + server.Port + " ]")
	if err := http.ListenAndServe(":"+server.Port, server.Router); err != nil {
		log.Fatal("Error starting server", err)
	}
}
