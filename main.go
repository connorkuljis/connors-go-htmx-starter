package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/connorkuljis/connors-go-htmx-starter/pkg/server"
)

//go:embed www/templates/* www/static/*
var embedFS embed.FS

const (
	StaticDirName    = "www/static"
	TemplatesDirName = "www/templates"
)

func main() {
	s := server.NewServer(embedFS, "8080", StaticDirName, TemplatesDirName)

	s.Routes()

	log.Println("[ 💿 Spinning up server on http://localhost:" + s.Port + " ]")
	if err := http.ListenAndServe(":"+s.Port, s.Router); err != nil {
		log.Fatal("Error starting server", err)
	}
}

// TODO: support environment variables for configuration
func envvars() {
	os.Getenv("PORT")    // PORT=8080
	os.Getenv("DEVMODE") // DEVMODE=1
}
