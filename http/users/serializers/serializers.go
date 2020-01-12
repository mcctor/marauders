package serializers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mcctor/marauders/db"
	"github.com/mcctor/marauders/http"
)

const (
	firstPage     = 1
	resultPerPage = 10
)

func userItemSerializer(user *db.User) string {
	slug := fmt.Sprintf("/v1/users/%s/", user.Username)
	href := http.ServerAddr + slug
	return fmt.Sprintf(userItemPattern, href, user.Username, user.Fname.String, user.Lname.String,
		user.Email, user.Phone.String, user.Created)
}

// UserItemCollectionSerializer serializer user structs using the collection pattern.
func UserItemCollectionSerializer(users []*db.User, paginationLinks string) string {
	var items bytes.Buffer
	for index, user := range users {
		if index == 0 {
			items.WriteString(userItemSerializer(user))
		} else {
			items.WriteString(",")
			items.WriteString(userItemSerializer(user))
		}
	}
	return fmt.Sprintf(userItemCollectionsPattern, items.String(), paginationLinks)
}

func PaginatedUserItemCollectionSerializer(curPage int) string {
	actualPage := curPage - 1
	pageUserItems, err := db.GetUsersByPage(actualPage, resultPerPage)
	if err != nil {
		log.Fatalf("could not get users by paginator: %v", err)
	}
	paginationLinks := getCurrentPagePaginationLinks(curPage)
	linksJsonString := convertPaginationLinksToJsonString(paginationLinks)

	return UserItemCollectionSerializer(pageUserItems, linksJsonString)
}

func convertPaginationLinksToJsonString(paginationLinks []UserRelationLinks) string {
	if paginationLinks == nil {
		return ""
	}

	linksJson, err := json.Marshal(paginationLinks)
	if err != nil {
		log.Fatalf("could not marshal all users pagination links: %v", err)
	}
	var linksJsonString bytes.Buffer
	linksJsonString.Write(linksJson)

	return linksJsonString.String()
}

func getCurrentPagePaginationLinks(curPage int) (collectionLinks []UserRelationLinks) {
	nextPageNum := curPage + 1
	prevPageNum := curPage - 1
	resultOffset := curPage * resultPerPage

	if curPage == firstPage && hasReachedLastPage(resultOffset) {
		return
	} else if curPage == firstPage {
		firstPageLink := UserRelationLinks{Href: fmt.Sprintf("/page/%d/", nextPageNum), Rel: "next", Render: "link"}
		collectionLinks = append(collectionLinks, firstPageLink)
		return
	}

	if hasReachedLastPage(resultOffset) {
		lastPageLink := UserRelationLinks{Href: fmt.Sprintf("/page/%d/", prevPageNum), Rel: "prev", Render: "link"}
		collectionLinks = append(collectionLinks, lastPageLink)
		return
	}
	nextPageLink := UserRelationLinks{Href: fmt.Sprintf("/page/%d/", nextPageNum), Rel: "next", Render: "link"}
	prevPageLink := UserRelationLinks{Href: fmt.Sprintf("/page/%d/", prevPageNum), Rel: "prev", Render: "link"}
	collectionLinks = append(collectionLinks, nextPageLink, prevPageLink)
	return
}

func hasReachedLastPage(resultOffset int) bool {
	userCount, err := db.AllUserCount()
	if err != nil {
		log.Printf("failed to count all the users in the database: %v\n", err)
	}
	if userCount <= resultOffset {
		return true
	} else {
		return false
	}
}
