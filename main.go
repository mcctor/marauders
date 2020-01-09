package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/mcctor/marauders/db"
	"github.com/mcctor/marauders/http"
)

func startServer() {
	log.Println("Started server at port: 8080 ...")
	if err := http.Server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	startServer()
}
