package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mcctor/marauders/http/users/serializers"
)

func allUsersPaginated(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	curPageStr := vars["page_number"]
	curPage, err := strconv.ParseInt(curPageStr, 10, 64)
	if err != nil {
		http.Error(writer, "", http.StatusNotFound)
		return
	}
	serializedItems, err := serializers.PaginatedUserItemsSerializer(int(curPage))
	if err != nil {
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}
	writer.Write(serializedItems)
}
