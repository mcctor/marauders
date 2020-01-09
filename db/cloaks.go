package db

import (
	"fmt"
	"time"

	"github.com/mcctor/marauders/utils"
)

const cloakIDLen = 20

type Cloak struct {
	ID              string `db:"id"`
	User            string
	Name            string
	Description     string
	Active          bool
	Wake            string
	Sleep           string
	Accuracy        string
	Duration        string
	MemberLimit     int  `db:"member_limit"`
	MemberVisible   bool `db:"member_visible"`
	CreatorVisible  bool `db:"creator_visible"`
	EveryoneVisible bool `db:"everyone_visible"`
	Private         bool
	Created         string
	Modified        string
}

func (c *Cloak) save() error {
	insertQuery := `
	INSERT INTO cloaks 
		(id, user, name, description, active, wake, sleep, accuracy, duration, member_limit,
		 member_visible, creator_visible, everyone_visible, private)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
	c.ID = utils.GenerateKey(cloakIDLen)
	_, err := db.Exec(insertQuery, c.ID, c.User, c.Name, c.Description, c.Active, c.Wake, c.Sleep,
		c.Accuracy, c.Duration, c.MemberLimit, c.MemberVisible, c.CreatorVisible, c.EveryoneVisible, c.Private)
	if err != nil {
		return fmt.Errorf("failed to save cloak<%s>: %v", c.ID, err)
	}
	return nil
}

// Update saves the current state of the struct fields into the associated
// row in the database.
func (c *Cloak) Update() error {
	updateQuery := `
	UPDATE cloaks SET name = ?, description = ?, wake = ?, sleep = ?, accuracy = ?, duration = ?, member_visible = ?,
	                  creator_visible = ?, everyone_visible = ?, member_limit = ?, private = ? WHERE id = ?
`

	_, err := db.Exec(updateQuery, c.Name, c.Description, c.Wake, c.Sleep, c.Accuracy, c.Duration, c.MemberVisible,
		c.CreatorVisible, c.EveryoneVisible, c.MemberLimit, c.Private, c.ID)
	if err != nil {
		return fmt.Errorf("failed to update cloak<%s>: %v", c.ID, err)
	}
	return nil
}

// Delete removes an cloak row with the corresponding CloakID from the cloaks table.
func (c *Cloak) Delete() error {
	_, err := db.Exec("DELETE FROM cloaks WHERE id = ?", c.ID)
	if err != nil {
		return fmt.Errorf("could not delete cloak<%s>: %v", c.ID, err)
	}
	return nil
}

// AddPermittedCloak allows the members of the cloak with the specified cloakID
// have access to the data this struct's cloak makes visible.
func (c *Cloak) AddPermittedCloak(cloakID string) error {
	_, err := db.Exec("INSERT INTO permitted_cloaks VALUES (null, ?, ?, CURRENT_TIMESTAMP)", c.ID, cloakID)
	if err != nil {
		return fmt.Errorf("failed to add cloak<%s> to permitted cloaks: %v", cloakID, err)
	}
	return nil
}

// RemovePermittedCloak disallows the members of the cloak with the specified cloakID
// from having access to the data this struct's cloaks makes visible.
func (c *Cloak) RemovePermittedCloak(cloakID string) error {
	_, err := db.Exec("DELETE FROM permitted_cloaks WHERE permitted_cloak_id = ?", cloakID)
	if err != nil {
		return fmt.Errorf("failed to remove cloak<%s> to permitted cloaks: %v", c.ID, err)
	}
	return nil
}

// NewInviteLink creates a new invite link for this specific cloak with
// the specified parameters
func (c *Cloak) NewInviteLink(creator string, countLimit int, expiry time.Time) (*CloakInviteLink, error) {
	return newCloakInviteLink(c.ID, creator, countLimit, expiry)
}

// AllMembers method fetches both associated members and actual members of this cloak
func (c *Cloak) AllMembers() (devices []Device, err error) {
	associatedMembers, err := c.AssociatedMembers()
	if err != nil {
		return devices, err
	}
	for _, device := range associatedMembers {
		devices = append(devices, device)
	}

	cloakMembers, err := c.Members()
	if err != nil {
		return devices, fmt.Errorf("could not fetch all members of cloak<%s>: %v", c.ID, err)
	}
	for _, device := range cloakMembers {
		devices = append(devices, device)
	}

	return devices, nil
}

// AssociatedMembers fetches members who belong to other cloaks but are
// permitted to see the location data given by this cloak
func (c *Cloak) AssociatedMembers() (devices []Device, err error) {
	var assocCloakIDs []string
	err = db.Select(&assocCloakIDs, `SELECT permitted_cloak_id FROM permitted_cloaks WHERE cloak_id = ?`, c.ID)
	if err != nil {
		return []Device{}, fmt.Errorf("could not get associated members for cloak<%s>: %v", c.ID, err)
	}

	for _, cloakID := range assocCloakIDs {
		cloak, _ := GetCloakByID(cloakID)
		cloakDevices, _ := cloak.Members()

		for _, device := range cloakDevices {
			devices = append(devices, device)
		}
	}

	return devices, nil
}

// Members returns a slice of devices associated to this particular cloak.
func (c *Cloak) Members() (devices []Device, err error) {
	query := `
	SELECT assoc_dev.device_id AS id, dev.user FROM
	(SELECT device_id FROM associated_cloaks WHERE cloak_id = ?) AS assoc_dev INNER JOIN devices AS dev
	ON assoc_dev.device_id = dev.id
`
	err = db.Select(&devices, query, c.ID)
	if err != nil {
		return devices, fmt.Errorf("failed to get members of cloak<%s>: %v", c.ID, err)
	}

	return devices, nil
}

// LocationSnapshotsForMember returns the location snapshots for the passed device
// that are allowed by the rules defined by this cloak.
func (c *Cloak) LocationSnapshotsForMember(device Device, lim int) ([]LocationSnapshot, error) {
	return device.locationSnapshotsForCloak(c, lim)
}

// newCloak creates a new cloak and commits it to the database
func newCloak(creator, name, description string, wake, sleep, duration time.Time,
	accuracy string, memberLimit int, memberVisible, creatorVisible, everyoneVisible, isPrivate bool) (*Cloak, error) {

	newCloak := &Cloak{
		User:            creator,
		Name:            name,
		Description:     description,
		Active:          true,
		Wake:            wake.UTC().Format(utils.TimeFormat),
		Sleep:           sleep.UTC().Format(utils.TimeFormat),
		Accuracy:        accuracy,
		Duration:        duration.UTC().Format(utils.TimeFormat),
		MemberLimit:     memberLimit,
		MemberVisible:   memberVisible,
		CreatorVisible:  creatorVisible,
		EveryoneVisible: everyoneVisible,
		Private:         isPrivate,
	}
	err := newCloak.save()
	if err != nil {
		return &Cloak{}, fmt.Errorf("failed to create new cloak<%s>: %v", newCloak.ID, err)
	}
	return newCloak, nil
}

// getCloaksFor returns all the cloak rows organized as a struct of *Cloak owned by the passed
// user.
func getCloaksFor(username string) (ownedCloaks []*Cloak, err error) {
	err = db.Select(&ownedCloaks, "SELECT * FROM cloaks WHERE user = ?", username)
	if err != nil {
		return ownedCloaks, fmt.Errorf("failed to get cloaks for user<%s>: <%v>", username, err)
	}
	return ownedCloaks, nil
}

// GetCloakByID returns the cloaks row which matches the cloakID passed in as a struct
func GetCloakByID(cloakID string) (fetchedCloak *Cloak, err error) {
	fetchedCloak = &Cloak{}
	err = db.Get(fetchedCloak, "SELECT * FROM cloaks WHERE id = ?", cloakID)
	if err != nil {
		return fetchedCloak, fmt.Errorf("failed to get cloak with id<%s>: %v", cloakID, err)
	}
	return fetchedCloak, nil
}
