package db

import (
	"crypto/sha256"
	"fmt"

	"github.com/mcctor/marauders/utils"
)

const saltLength = 50

type Password struct {
	User     string
	Salt     string
	Hash     string
	Created  string
	Modified string
}

// save commits the fields of a newly constructed password struct to the database.
func (p *Password) save() error {
	_, err := db.Exec(
		"INSERT INTO passwords (user, salt, hash) VALUES (?, ?, ?)",
		p.User, p.Salt, p.Hash)
	if err != nil {
		return fmt.Errorf("failed to save password for user<%s>: %v", p.User, err)
	}
	return nil
}

// update commits the changes of pre-existing password struct to the database
func (p *Password) update() error {
	_, err := db.Exec("UPDATE passwords SET user = ?, salt = ?, hash = ? WHERE user = ?",
		p.User, p.Salt, p.Hash, p.User)
	if err != nil {
		return fmt.Errorf("failed to update password for user<%s>: %v", p.User, err)
	}
	return nil
}

// CheckIf examines whether the passed password matches the currently existing
// password which is in the form of a salted hash
func (p *Password) CheckIf(password string) (matches bool) {
	matches = p.Hash == generateHash(p.Salt, password)
	return
}

// ChangeTo replaces the preexisting Hash with a new sha512 hexadecimal
// dump based on the new password
func (p *Password) ChangeTo(newPassword string) error {
	p.Hash = generateHash(p.Salt, newPassword)
	return p.update()
}

// newPasswordFor creates a new row in the passwords table. Typically, it is
// meant to be called when registering a new user. If changing an existing
// password is what is required, use the func (p *Password) ChangeTo method.
func newPasswordFor(username string, password string) (*Password, error) {
	salt := utils.GenerateKey(saltLength)
	hash := generateHash(salt, password)
	passwordStruct := &Password{
		User: username,
		Salt: salt,
		Hash: hash,
	}
	err := passwordStruct.save()
	if err != nil {
		return &Password{}, err
	}
	return passwordStruct, nil
}

// getPasswordFor fetches the corresponding password row as a struct for
// the passed username
func getPasswordFor(username string) (*Password, error) {
	password := &Password{}
	err := db.Get(password, "SELECT * FROM passwords WHERE user = ?", username)
	if err != nil {
		return &Password{}, fmt.Errorf("failed to get password for user<%s>: %v", username, err)
	}
	return password, nil
}

// generateHash creates a new sha512 hexadecimal dump based on the passed
// salt and password and returns it as a hash
func generateHash(salt string, password string) (hash string) {
	saltedPass := salt + password
	return fmt.Sprintf("%x", sha256.Sum256([]byte(saltedPass)))
}

