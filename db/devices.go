package db

import "fmt"

type Device struct {
	ID      int `db:"id"`
	User    string
	Created string
}

// AssociatedCloaks returns a list of cloak to which this device is a
// a member of.
func (d Device) AssociatedCloaks() ([]*Cloak, error) {
	var joinedCloaks []*Cloak
	query := `
	SELECT id, user, name, description, active, wake, sleep, accuracy, duration, member_limit, member_visible,
	       creator_visible, everyone_visible, private, created, modified
	FROM cloaks INNER JOIN (
		SELECT cloak_id FROM associated_cloaks WHERE device_id = ?
	) AS cloak_assoc ON cloaks.id = cloak_assoc.cloak_id
`
	err := db.Select(&joinedCloaks, query, d.ID)
	if err != nil {
		return joinedCloaks, fmt.Errorf("could not fetch cloaks joined by device<%d>: %v", d.ID, err)
	}
	return joinedCloaks, nil
}

// LocationSnapshots returns the latest location snapshots for this particular device.
func (d Device) LocationSnapshots(lim int) ([]LocationSnapshot, error) {
	var allLocationSnapshots []LocationSnapshot
	err := db.Select(&allLocationSnapshots,
		"SELECT * FROM location_snapshots WHERE device_id = ? ORDER BY time_stamp DESC LIMIT ?",
		d.ID, lim)
	if err != nil {
		return allLocationSnapshots, fmt.Errorf("failed to fetch location history for device<%d>: %v",
			d.ID, err)
	}
	return allLocationSnapshots, nil
}

// NewLocationSnapshot adds the passed fields for a LocationSnapshot struct
// as a row in the LocationSnapshots table. If this fails, an error is
// returned.NewLocation
func (d Device) NewLocationSnapshot(timeStamp string, latitude, longitude float64) error {
	newSnapshot := LocationSnapshot{
		DeviceID:  d.ID,
		TimeStamp: timeStamp,
		Latitude:  latitude,
		Longitude: longitude,
	}
	return newSnapshot.save()
}

// AssociateToCloak associates this device to the cloak with the specified cloak_id.
func (d Device) AssociateToCloak(cloakID string) error {
	_, err := db.Exec("INSERT INTO associated_cloaks (cloak_id, device_id) VALUES (?, ?)", cloakID, d.ID)
	if err != nil {
		return fmt.Errorf("device<%d> for user<%s> could not join cloak<%s>: %v", d.ID, d.User, cloakID, err)
	}
	return nil
}

// DissociateFromCloak dissociates this device from the cloak with the specified cloak_id.
func (d Device) DissociateFromCloak(cloakID string) error {
	_, err := db.Exec("DELETE FROM associated_cloaks WHERE cloak_id = ? AND device_id = ?", cloakID, d.ID)
	if err != nil {
		return fmt.Errorf("device<%d> could not leave cloak<%s>: %v", d.ID, cloakID, err)
	}
	return nil
}

// DeleteDevice deletes the row from the table Device whose id column
// matches the passed deviceID. An error is returned if this
// process fails.
func (d Device) Delete() error {
	_, err := db.Exec("DELETE FROM devices WHERE id = ?", d.ID)
	if err != nil {
		return fmt.Errorf("failed to delete device<%d>: %v", d.ID, err)
	}
	return nil
}

// locationSnapshotsForCloak method returns the location history of device after having been
// filtered by the rules of the passed cloak.
func (d Device) locationSnapshotsForCloak(cloak *Cloak, lim int) (locationSnaps []LocationSnapshot, err error) {
	query := `
	SELECT * FROM location_snapshots
	WHERE device_id = ? AND (time_stamp > ? AND (TIME(time_stamp) BETWEEN ? AND ?))
	ORDER BY time_stamp DESC LIMIT ?
`
	err = db.Select(&locationSnaps, query, d.ID, cloak.Duration, cloak.Wake, cloak.Sleep, lim)
	if err != nil {
		return locationSnaps, fmt.Errorf("failed to get loc snapshots for device<%d> of cloak<%s>: %v",
			d.ID, cloak.ID, err)
	}
	return locationSnaps, nil
}

// newDeviceFor adds a new row to the table Device associating it to
// the passed username.
func newDeviceFor(username string, deviceID int) (Device, error) {
	createdDevice := Device{
		ID:   deviceID,
		User: username,
	}
	_, err := db.Exec("INSERT INTO devices (id, user) VALUES (?, ?)", deviceID, username)
	if err != nil {
		return Device{}, fmt.Errorf("failed to create new device for user<%s>: %v", username, err)
	}
	return createdDevice, nil
}

func getDevicesFor(username string) (devices []Device, err error) {
	err = db.Select(&devices, "SELECT * FROM devices WHERE user = ?", username)
	if err != nil {
		return devices, fmt.Errorf("failed to get devices for user<%s>: %v", username, err)
	}
	return devices, nil
}
