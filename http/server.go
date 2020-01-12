package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var Server *http.Server
var Router *mux.Router
var ServerAddr string

func init() {
	// create and apply middleware to root router
	Router = mux.NewRouter()
	Router.Use(ApplyContentTypeCollectionsJson, ApplyGzipCompression)
	Server = &http.Server{
		Addr:         "localhost:8080",
		Handler:      Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	ServerAddr = fmt.Sprintf("http://%s", Server.Addr)
}
