package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/connorkuljis/connors-go-htmx-starter/pkg/server"
)

//go:embed templates/* static/*
var embedFS embed.FS

const (
	port             = "8080"
	staticDirName    = "static"
	templatesDirName = "templates"
)

func main() {
	router := http.NewServeMux()

	s := server.NewServer(port, router, templatesDirName, staticDirName, embedFS)

	s.Routes()

	log.Println("[ 💿 Spinning up server on http://localhost:" + s.Port + " ]")

	if err := http.ListenAndServe(":"+s.Port, s.Router); err != nil {
		log.Fatal("Error starting server", err)
	}
}
