package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/mcctor/marauders/db"
	userConst "github.com/mcctor/marauders/http/users"
	"github.com/mcctor/marauders/http/users/serializers"
	"github.com/mcctor/marauders/utils"
	"io/ioutil"
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
	serializedUsers, err := serializers.PaginatedUserItemsSerializer(userConst.FirstPage)
	if err != nil {
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}
	writer.Write(serializedUsers)
}

func usersPostHandler(writer http.ResponseWriter, request *http.Request) {
	bodyBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}
	newUserFields, err := parseNewUserFields(bodyBytes)
	if err != nil {
		fmt.Println(err)
		http.Error(writer, "{\"status\": \"bad formatted json\"}", http.StatusBadRequest)
		return
	}
	newUser, err := createNewUserFromFields(newUserFields)
	if err != nil {
		http.Error(writer, "{\"status\": \"username already exists\"}", http.StatusBadRequest)
		return
	}
	newUserItem, err := serializers.UserItemSerializer(newUser)
	if err != nil {
		http.Error(writer, "", http.StatusInternalServerError)
		return
	}

	setContentCreatedHeader(fmt.Sprintf("%s/%s/", userConst.Href, newUser.Username), writer)
	writer.Write(newUserItem)
}

func parseNewUserFields(respBody []byte) (fields struct {
	Template utils.ItemTemplate `db:"template"`
}, err error) {
	err = json.Unmarshal(respBody, &fields)
	if err != nil {
		return fields, fmt.Errorf("failed to parse new user fields: %v", err)
	}
	return
}

func createNewUserFromFields(newUserFields struct {
	Template utils.ItemTemplate `db:"template"`
}) (newUser *db.User, err error) {
	var username, password, fname, lname, email, phone string
	for _, field := range newUserFields.Template.Data {
		switch field.Name {
		case "username":
			username = field.Value
		case "password":
			password = field.Value
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
	newUser, err = db.NewUser(username, email)
	if err != nil {
		return newUser, fmt.Errorf("failed to create new user from fields: %v", err)
	}
	newUser.Fname = sql.NullString{String: fname, Valid: true}
	newUser.Lname = sql.NullString{String: lname, Valid: true}
	newUser.Phone = sql.NullString{String: phone, Valid: true}
	err = newUser.Update()
	if err != nil {
		return newUser, fmt.Errorf("failed to save user fields: %v", err)
	}
	_, err = newUser.NewPassword(password)
	if err != nil {
		return newUser, fmt.Errorf("failed to save user password: %v", err)
	}
	return
}

func setContentCreatedHeader(newResourceSlug string, writer http.ResponseWriter) {
	writer.Header().Set("Content-Location", newResourceSlug)
	writer.WriteHeader(http.StatusCreated)
}
