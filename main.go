package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	srv "github.com/connorkuljis/connors-go-htmx-starter/pkg/server"
)

//go:embed www/templates/* www/static/*
var embedFS embed.FS

func main() {
	s := srv.NewServer(embedFS, "8080")
	s.AppData = srv.AppData{
		Title:   "Connors Go HTMX Starter",
		DevMode: false, // load from env
	}
	s.Routes()

	log.Println("[ 💿 Spinning up server on http://localhost:" + s.Port + " ]")
	if err := http.ListenAndServe(":"+s.Port, s.Router); err != nil {
		log.Fatal("Error starting server", err)
	}
}

func port() {
	port := os.Getenv("PORT")
	devmode := os.Getenv("DEVMODE")
}
