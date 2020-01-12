package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/mcctor/marauders/db"
	marauderhttp "github.com/mcctor/marauders/http"
	"github.com/mcctor/marauders/http/users/serializers"
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
	serializedUsers := serializers.PaginatedUserItemCollectionSerializer(firstPage)
	writer.Write([]byte(serializedUsers))
}

func usersPostHandler(writer http.ResponseWriter, request *http.Request) {
	body, err := readRequestBody(request)
	if err != nil {
		log.Println(err)
		http.Error(writer, "{\"status\": \"failed to read request body\"}", http.StatusInternalServerError)
		return
	}
	userTemplate, err := unmarshalNewUserJsonBody(body)
	if err != nil {
		log.Println(err)
		http.Error(writer, "{\"status\": \"the json body sent is invalid\"}", http.StatusBadRequest)
		return
	}
	newUserFields := parseUserStructFields(userTemplate.Template.Data)
	newUser, err := createNewUserFrom(newUserFields)
	if err != nil {
		log.Println(err)
		http.Error(writer, "{\"status\": \"the given username is not available\"}", http.StatusBadRequest)
		return
	}
	newUserResourceSlug := fmt.Sprintf("/v1/users/%s/", newUser.Username)
	setContentCreatedHeader(newUserResourceSlug, writer)
	writer.Write([]byte(newUserResourceSerializer(newUser)))
}

func readRequestBody(request *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %v\n", err)
	}
	return body, nil
}

func unmarshalNewUserJsonBody(newUserData []byte) (serializers.UserTemplate, error) {
	var userTemplate serializers.UserTemplate
	err := json.Unmarshal(newUserData, &userTemplate)
	if err != nil {
		return serializers.UserTemplate{}, fmt.Errorf("failed to unmarshal new user json body: %v", err)
	}
	return userTemplate, nil
}

func parseUserStructFields(data []serializers.UserInputField) struct{ username, fname, lname, email, phone string } {
	var username, fname, lname, email, phone string
	for _, field := range data {
		switch field.Name {
		case "username":
			username = field.Value
		case "fname":
			fname = field.Value
		case "lname":
			lname = field.Value
		case "email":
			email = field.Value
		case "phone":
			phone = field.Value
		}
	}
	normalizedFields := normalizeFields(username, fname, lname, email, phone)
	return struct {
		username, fname, lname, email, phone string
	}{normalizedFields[0], normalizedFields[1], normalizedFields[2],
		normalizedFields[3], normalizedFields[4]}
}

func createNewUserFrom(userFields struct{ username, fname, lname, email, phone string }) (*db.User, error) {
	newUser, err := db.NewUser(userFields.username, userFields.email)
	if err != nil {
		return nil, fmt.Errorf("failed to create new user from request body fields: %v", err)
	}
	newUser.Fname = sql.NullString{String: userFields.fname, Valid: true}
	newUser.Lname = sql.NullString{String: userFields.lname, Valid: true}
	newUser.Phone = sql.NullString{String: userFields.phone, Valid: true}
	err = newUser.Update()
	if err != nil {
		return nil, fmt.Errorf("failed to save new user nullable fields: %v", err)
	}
	return newUser, nil
}

func newUserResourceSerializer(newUser *db.User) string {
	userSlice := []*db.User{newUser}
	return serializers.UserItemCollectionSerializer(userSlice, "")
}

func setContentCreatedHeader(newResourceSlug string, writer http.ResponseWriter) {
	contentLocationSlug := marauderhttp.ServerAddr + newResourceSlug
	writer.Header().Set("Content-Location", contentLocationSlug)
	writer.WriteHeader(http.StatusCreated)
}

func normalizeFields(fields ...string) []string {
	for index, field := range fields {
		fields[index] = strings.TrimSpace(strings.ToLower(field))
	}
	return fields
}
