package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mcctor/marauders/utils"
)

const (
	hoursToExpiry      = 120 * time.Hour
	tokenLength        = 40
	refreshTokenLength = 20
)

type AuthToken struct {
	User         string
	Token        string
	RefreshToken string `db:"refresh_token"`
	Expiry       string
	Modified     string
	Created      string
}

// Renew generates a new token if the passed refreshToken matches
// the one belonging to the struct. It then returns the new token if it
// succeeds, or an error if the refresh token does not match the current
// refresh token.
func (auth *AuthToken) Renew(refreshToken string) (newToken string, err error) {
	if refreshToken != auth.RefreshToken {
		return "", errors.New("refresh tokens do not match")
	}
	auth.Token = utils.GenerateKey(tokenLength)
	auth.Expiry = time.Now().UTC().Add(hoursToExpiry).Format(utils.TimeFormat)

	err = auth.save()
	if err != nil {
		return "", fmt.Errorf("could not renew token for user<%s>: %v", auth.User, err)
	}

	return auth.Token, nil
}

// isValid examines the expiry field by comparing it to the current
// time. If the current time is past the expiry stipulated in the struct
// false is returned, otherwise, true.
func (auth *AuthToken) IsValid() bool {
	expiry, err := time.Parse(utils.TimeFormat, auth.Expiry)
	if err != nil {
		log.Fatal(err)
	}
	return time.Now().Before(expiry)
}

// update persists the changes made to this struct to the auth_tokens
// table in the database. Various errors may result from this operation
// the most common being a UNIQUE KEY Integrity error.
func (auth *AuthToken) save() error {
	_, err := db.Exec(
		"UPDATE auth_tokens SET token = ?, refresh_token = ?, expiry = ? WHERE user = ?",
		auth.Token, auth.RefreshToken, auth.Expiry, auth.User,
	)
	if err != nil {
		return fmt.Errorf("failed to save authtoken for user<%s>: %v", auth.User, err)
	}
	return nil
}

// newAuthTokenFor creates a new row in the AuthToken tables belonging
// to the user associated with the passed username. This method will normally
// be called after the sign-up of a new user.
func newAuthTokenFor(username string) (*AuthToken, error) {
	authToken := &AuthToken{
		User:         username,
		Token:        utils.GenerateKey(tokenLength),
		RefreshToken: utils.GenerateKey(refreshTokenLength),
		Expiry:       time.Now().UTC().Add(hoursToExpiry).Format(utils.TimeFormat),
	}

	_, err := db.Exec(
		"INSERT INTO auth_tokens (user, token, refresh_token, expiry) VALUES (?, ?, ?, ?)",
		authToken.User, authToken.Token, authToken.RefreshToken, authToken.Expiry,
	)
	if err != nil {
		return &AuthToken{}, fmt.Errorf("could not create new authentication token for user<%s>: %v", username, err)
	}
	return authToken, nil
}

// getAuthTokenFor fetches the corresponding auth_token row whose user column
// matches the passed username.
func getAuthTokenFor(username string) (authToken *AuthToken, err error) {
	authToken = &AuthToken{}
	err = db.Get(authToken, "SELECT * FROM auth_tokens WHERE user = ?", username)
	if err != nil {
		return authToken, fmt.Errorf("failed to get auth token for User<%s>: %v", username, err)
	}
	return authToken, nil
}
