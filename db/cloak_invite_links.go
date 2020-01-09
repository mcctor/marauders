package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mcctor/marauders/utils"
)

const linkSize = 12

type CloakInviteLink struct {
	Link       string
	CloakID    string `db:"cloak_id"`
	CreatedBy  string `db:"created_by"`
	Expiry     string
	Added      int
	CountLimit int `db:"count_limit"`
	Created    string
	Modified   string
}

// update commits the fields of a newly created struct to the database.
func (invite *CloakInviteLink) save() error {
	_, err := db.Exec(
		"INSERT INTO cloak_invite_links (link, cloak_id, created_by, expiry, count_limit) VALUES (?, ?, ?, ?, ?)",
		invite.Link, invite.CloakID, invite.CreatedBy, invite.Expiry, invite.CountLimit)
	if err != nil {
		return fmt.Errorf("failed to save invite link for cloak<%s>: %v", invite.CloakID, err)
	}
	return nil
}

// update commits the changes of a pre-existing struct to the database. The only field that
// can be changed is the 'Added' field.
func (invite *CloakInviteLink) update() error {
	_, err := db.Exec(
		"UPDATE cloak_invite_links SET added = ? WHERE cloak_id = ? AND link = ?",
		invite.Added, invite.CloakID, invite.Link)
	return err
}

// Invite associates an "invitable" entity to the cloak this invite link is for
// as long as the invite link is still valid.
func (invite *CloakInviteLink) Invite(entity utils.Invitee) error {
	if !(invite.isValid()) {
		return errors.New("the invite link is no longer valid")
	}
	err := entity.AssociateToCloak(invite.CloakID)
	if err != nil {
		return fmt.Errorf("could not invite entity through link <%s> for cloak<%s>; %v",
			invite.Link, invite.CloakID, err)
	}
	invite.Added++
	return invite.update()
}

// Delete deletes the cloak_invite_links row which has the passed cloakID and inviteLink
// as its column values.
func (invite *CloakInviteLink) Delete() error {
	_, err := db.Exec("DELETE FROM cloak_invite_links WHERE cloak_id = ? AND link = ?",
		invite.CloakID, invite.Link)
	if err != nil {
		return fmt.Errorf("could not delete cloak invite link for cloak<%s>: %v", invite.CloakID, err)
	}
	return nil
}

// isValid checks the validity of the InviteLink by ensuring its expiry time has not
// reached nor has its invite count limit met.
func (invite *CloakInviteLink) isValid() bool {
	expiry, err := time.Parse(utils.TimeFormat, invite.Expiry)
	if err != nil {
		log.Fatal(err)
	}
	return time.Now().UTC().Before(expiry) && invite.Added < invite.CountLimit
}

// newCloakInviteLink creates a new cloak invite based on the passed parameters. It then generates
// a link before committing the result to the database and returning the newly created struct.
func newCloakInviteLink(cloakID string, creator string, countLimit int, expiry time.Time) (*CloakInviteLink, error) {
	inviteLink := &CloakInviteLink{
		Link:       utils.GenerateKey(linkSize),
		CloakID:    cloakID,
		CreatedBy:  creator,
		Expiry:     expiry.UTC().Format(utils.TimeFormat),
		CountLimit: countLimit,
	}
	err := inviteLink.save()
	if err != nil {
		return &CloakInviteLink{}, fmt.Errorf("could not create new invite link for cloak<%s>: %v", cloakID, err)
	}

	return inviteLink, nil
}

// getCloakInviteLinksFor returns all the invite links created by the passed username.
func getCloakInviteLinksFor(username string) (inviteLinks []*CloakInviteLink, err error) {
	err = db.Select(&inviteLinks, "SELECT * FROM cloak_invite_links WHERE created_by = ?", username)
	if err != nil {
		return []*CloakInviteLink{}, fmt.Errorf("could not get cloak invite links for %s: %v", username, err)
	}
	return inviteLinks, nil
}

// GetCloakInviteBy returns a CloakInviteLink struct which matches the passed cloak_id and
// inviteLink
func GetCloakInviteBy(cloakID string, inviteLink string) (link *CloakInviteLink, err error) {
	link = &CloakInviteLink{}
	err = db.Get(link, "SELECT * FROM cloak_invite_links WHERE cloak_id = ? AND link = ?", cloakID, inviteLink)
	if err != nil {
		return &CloakInviteLink{}, fmt.Errorf("could not get cloak invite link for cloak<%s>: %v", cloakID, err)
	}
	return link, nil
}


