package users

import (
	"fmt"
	"github.com/mcctor/marauders/db"
	"github.com/mcctor/marauders/http"
)

type UserSerializer struct {
	Href string  `json:"href"`
	Data db.User `json:"data"`
}

// NewUserSerializer returns a new serializer for the passed user struct.
func NewUserSerializer(user *db.User) UserSerializer {
	slug := fmt.Sprintf("https://%s/users/%s/", http.Server.Addr, user.Username)
	return UserSerializer{
		Href: slug,
		Data: *user,
	}
}
