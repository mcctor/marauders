package handlers

import (
	"github.com/mcctor/marauders/http/users/serializers"
	"net/http"
)

const firstPage = 1

func users(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		usersGetHandler(writer, request)
	case http.MethodPost:
		usersPostHandler(writer, request)
	}
}

func usersGetHandler(writer http.ResponseWriter, _ *http.Request) {
	serializedUsers, err := serializers.PaginatedUserItemsSerializer(firstPage)
	if err != nil {
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}
	writer.Write(serializedUsers)
}

func usersPostHandler(writer http.ResponseWriter, request *http.Request) {
	//todo: Start by creating a serializer for a single user item
}
