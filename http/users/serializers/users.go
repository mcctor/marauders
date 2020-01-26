package serializers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mcctor/marauders/db"
	"github.com/mcctor/marauders/http"
	"github.com/mcctor/marauders/utils"
)

const (
	version        = "1.0"
	resultsPerPage = 10
	firstPage      = 1
)

var href = http.ServerAddr + "/v1/users/"

func PaginatedUserItemsSerializer(curPage int) ([]byte, error) {
	actualPage := curPage - 1 // has to do with SQL's LIMIT's offset starting at 0
	users, err := db.GetUsersByPage(actualPage, resultsPerPage)
	if err != nil {
		return nil, fmt.Errorf("failed to paginate user items: %v", err)
	}
	return userCollectionSerializer(newUserCollection(users, curPage))
}

func userCollectionSerializer(collection utils.Collection) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(collection)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize user collection: %v", err)
	}
	return buffer.Bytes(), nil
}

func newUserCollection(userItems []*db.User, curPage int) (collection utils.Collection) {
	collection = utils.Collection{
		Collection: utils.ItemsCollection{
			Version:  version,
			Href:     href,
			Items:    serializeUserItems(userItems),
			Links:    userCollectionPaginationLinks(curPage),
			Queries:  userCollectionQueries(),
			Template: userCollectionTemplate(),
		},
	}
	return
}

func serializeUserItems(userItems []*db.User) (serializedItems []utils.CollectionItem) {
	for _, item := range userItems {
		itemSlug := fmt.Sprintf(href+"%s/", item.Username)
		serializedItems = append(serializedItems, utils.CollectionItem{
			Href: itemSlug,
			Data: []utils.DataField{
				{"username", "username", item.Username},
				{"first name", "fname", item.Fname.String},
				{"last name", "lname", item.Lname.String},
				{"email address", "email", item.Email},
				{"phone number", "phone", item.Phone.String},
			},
			Links: []utils.CollectionLink{
				{itemSlug + "billings/", "billings", "link"},
				{itemSlug + "auth-token/", "authentication token", "link"},
				{itemSlug + "invitation-links/", "invitation links", "link"},
				{itemSlug + "devices/", "owned devices", "link"},
				{itemSlug + "cloaks/", "owned cloaks", "link"},
			},
		})
	}
	return
}

func userCollectionPaginationLinks(curPage int) (links []utils.CollectionLink) {
	nextPageNum := curPage + 1
	prevPageNum := curPage - 1
	resultOffset := curPage * resultsPerPage

	if curPage == firstPage && hasReachedLastPage(resultOffset) {
		return
	} else if curPage == firstPage {
		firstPageLink := utils.CollectionLink{Href: fmt.Sprintf(href+"page/%d/", nextPageNum), Rel: "next", Render: "link"}
		links = append(links, firstPageLink)
		return
	}

	if hasReachedLastPage(resultOffset) {
		lastPageLink := utils.CollectionLink{Href: fmt.Sprintf(href+"page/%d/", prevPageNum), Rel: "prev", Render: "link"}
		links = append(links, lastPageLink)
		return
	}
	nextPageLink := utils.CollectionLink{Href: fmt.Sprintf(href+"page/%d/", nextPageNum), Rel: "next", Render: "link"}
	prevPageLink := utils.CollectionLink{Href: fmt.Sprintf(href+"page/%d/", prevPageNum), Rel: "prev", Render: "link"}
	links = append(links, nextPageLink, prevPageLink)
	return
}

func userCollectionQueries() (queries []utils.CollectionQuery) {
	queries = []utils.CollectionQuery{
		{href + "search/", "search", "search for user by username",
			[]utils.DataField{{"username", "username", ""}},
		}}
	return
}

func userCollectionTemplate() (userTemplate utils.ItemTemplate) {
	userTemplate.Data = []utils.DataField{
		{"username", "username", ""},
		{"first name", "fname", ""},
		{"last name", "lname", ""},
		{"email address", "email", ""},
		{"phone number", "phone", ""},
	}
	return
}

func hasReachedLastPage(resultOffset int) bool {
	if db.UserCount() <= resultOffset {
		return true
	} else {
		return false
	}
}
