package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type User struct {
	Username string
	Fname    sql.NullString
	Lname    sql.NullString
	Email    string
	Phone    sql.NullString
	Created  string
	Modified string
}

// NewDevice creates a new device with the passed deviceID
// associated to the owning user of this struct and returns its
// device struct.
func (u *User) NewDevice(deviceID int) (Device, error) {
	return newDeviceFor(u.Username, deviceID)
}

// Devices returns a slice of devices that are owned by the user this
// struct represents.
func (u *User) Devices() (ownedDevices []Device, err error) {
	return getDevicesFor(u.Username)
}

// Billings returns a slice of bills that have been charged to the user
// this struct represents. The length of the slice is limited to the
// limit passed to the function
func (u *User) Billings(lim int) (billings []Billing, err error) {
	return getBillsFor(u.Username, lim)
}

// CreditUserAmount method credits the user represented by this struct
// the amount passed.
func (u *User) CreditUserAmount(amount float64) (Billing, error) {
	return billUserFor(u.Username, amount, true)
}

// DebitUserAmount method debits the user represented by this struct
// the amount passed.
func (u *User) DebitUserAmount(amount float64) (Billing, error) {
	return billUserFor(u.Username, amount, false)
}

// NewAuthToken creates a new authentication token for the user
// represented by this struct, and returns the auth token struct
func (u *User) NewAuthToken() (*AuthToken, error) {
	return newAuthTokenFor(u.Username)
}

// AuthToken returns the authentication struct for the user this
// struct represents.
func (u *User) AuthToken() (*AuthToken, error) {
	return getAuthTokenFor(u.Username)
}

// NewPassword creates a new password row using the passed
// password string and returns the resulting password struct
func (u *User) NewPassword(password string) (*Password, error) {
	return newPasswordFor(u.Username, password)
}

// Password returns a password struct that represents the user's
// hashed password
func (u *User) Password() (*Password, error) {
	return getPasswordFor(u.Username)
}

// NewCloak creates a new cloak that is associated to this struct's user
// using the passed in parameters.
func (u *User) NewCloak(name, description string, wake, sleep, duration time.Time, accuracy string,
	memberLimit int, memberVisible, creatorVisible, everyoneVisible, isPrivate bool) (*Cloak, error) {
	return newCloak(u.Username, name, description, wake, sleep, duration, accuracy,
		memberLimit, memberVisible, creatorVisible, everyoneVisible, isPrivate)
}

// Cloaks returns all the cloaks that are owned by the user
// represented by this struct
func (u *User) Cloaks() ([]*Cloak, error) {
	return getCloaksFor(u.Username)
}

// InviteLinks fetches all the invite links created by this user
func (u *User) InviteLinks() ([]*CloakInviteLink, error) {
	return getCloakInviteLinksFor(u.Username)
}

// Update function commits the struct fields for User into
// the users table. It can only be used in updating every other field
// except for the username column.
func (u *User) Update() error {

	_, err := db.Exec("UPDATE users SET fname = ?, lname = ?, email = ?, phone = ? WHERE username = ?",
		u.Fname.String, u.Lname.String, u.Email, u.Phone.String, u.Username)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to update information for user<%s>: %v", u.Username, err)
	}
	return nil
}

// Delete drops this struct's row from users table
func (u *User) Delete() error {
	_, err := db.Exec("DELETE FROM users WHERE username = ?", u.Username)
	if err != nil {
		return err
	}
	return nil
}

// String returns a shortened representation of this user struct
func (u *User) String() string {
	return fmt.Sprintf("User<%s>", u.Username)
}

// GetUsers fetches all the user rows from the database.
func GetUsers(lim int) (allUsers []*User, err error) {
	err = db.Select(&allUsers, "SELECT * FROM users ORDER BY created DESC LIMIT ?", lim)
	if err != nil {
		return allUsers, fmt.Errorf("failed to fetch all users: %v", err)
	}
	return allUsers, nil
}

// GetUsersByPage takes in an integer representing the current page,
// then returns a list of users for that particular page
func GetUsersByPage(curPage, resultPerPage int) (pageUsers []*User, err error) {
	from := curPage * resultPerPage

	err = db.Select(&pageUsers, "SELECT * FROM users LIMIT ?, ?", from, resultPerPage)
	if err != nil {
		return pageUsers, fmt.Errorf("failed to fetch paginated list of users: %v", err)
	}
	return pageUsers, nil
}

// UserCount returns the number of users currently saved in the database
func UserCount() int {
	var num int
	err := db.Get(&num, "SELECT COUNT(username) FROM users")
	if err != nil {
		log.Fatalf("failed to get user count: %v", err)
	}
	return num
}

// NewUser creates a new User struct and saves it to the user table.
func NewUser(username, email string) (*User, error) {
	newUser := &User{
		Username: username,
		Email:    email,
	}
	_, err := db.Exec("INSERT INTO users (username, email) VALUES (?, ?)", newUser.Username, newUser.Email)
	if err != nil {
		return &User{}, fmt.Errorf("failed to create new user<%s>: %v", username, err)
	}
	_, err = newBillings(username)
	if err != nil {
		return &User{}, fmt.Errorf("failed to create new user<%s>: %v", username, err)
	}

	return newUser, nil
}

// GetUser fetches the user row with the corresponding passed username and
// returns a user struct with the fetched information
func GetUser(username string) (existingUser *User, err error) {
	existingUser = &User{}
	err = db.Get(existingUser, "SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		return &User{}, err
	}
	return existingUser, nil
}
