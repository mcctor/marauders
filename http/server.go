package http

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var Server *http.Server
var Router *mux.Router

func init() {
	Router = mux.NewRouter()
	Server = &http.Server{
		Addr:         "localhost:8080",
		Handler:      Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
}
