package handlers

import (
	"github.com/gorilla/mux"
	"github.com/mcctor/marauders/http/users/serializers"
	"log"
	"net/http"
	"strconv"
)

func allUsersPaginated(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	pageNumber, err := strconv.ParseInt(vars["page_number"], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	paginatedUserItems := serializers.PaginatedUserItemCollectionSerializer(int(pageNumber))
	writer.Write([]byte(paginatedUserItems))
}
