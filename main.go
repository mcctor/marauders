package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/mcctor/marauders/db"
	"github.com/mcctor/marauders/http"
	_ "github.com/mcctor/marauders/http/users/handlers"
)

func main() {
	log.Println("Started Marauders server at port 8080 ...")
	if err := http.Server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
