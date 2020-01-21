package handlers

import (
	"net/http"
)

func users(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		usersGetHandler(writer, request)
	case http.MethodPost:
		usersPostHandler(writer, request)
	}
}

func usersGetHandler(writer http.ResponseWriter, _ *http.Request) {

}

func usersPostHandler(writer http.ResponseWriter, request *http.Request) {

}
