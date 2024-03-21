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
	s := server.NewServer(embedFS, "8080")
	s.AppData = server.AppData{
		Title:   "Connors Go HTMX Starter",
		DevMode: false, // load from env
	}
	s.Routes()

	log.Println("[ 💿 Spinning up server on http://localhost:" + s.Port + " ]")
	if err := http.ListenAndServe(":"+s.Port, s.Router); err != nil {
		log.Fatal("Error starting server", err)
	}
}
