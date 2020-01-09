package db

import "fmt"

type LocationSnapshot struct {
	DeviceID  int    `db:"device_id"`
	TimeStamp string `db:"time_stamp"`
	Latitude  float64
	Longitude float64
}

// saves commits the struct's fields to the location_snapshots table in the database.
func (snapshot LocationSnapshot) save() error {
	_, err := db.Exec(
		"INSERT INTO location_snapshots (device_id, time_stamp, latitude, longitude) VALUES (?, ?, ?, ?)",
		snapshot.DeviceID, snapshot.TimeStamp, snapshot.Latitude, snapshot.Longitude,
	)
	if err != nil {
		return fmt.Errorf("failed to save locationsnapshot for device<%d>: %v", snapshot.DeviceID, err)
	}
	return nil
}
