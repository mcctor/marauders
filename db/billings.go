package db

import (
	"fmt"
	"time"

	"github.com/mcctor/marauders/utils"
)

const LATEST = 0

type Billing struct {
	TimeStamp string `db:"time_stamp"`
	User      string
	Debit     float64
	Credit    float64
}

// update inserts a new row with the struct's current fields
func (b Billing) save() error {
	_, err := db.Exec("INSERT INTO billings (time_stamp, user, debit, credit) VALUES (?, ?, ?, ?)",
		b.TimeStamp, b.User, b.Debit, b.Credit)
	if err != nil {
		return fmt.Errorf("could not save Billing info for user<%s>: %v", b.User, err)
	}
	return nil
}

// getBillsFor returns a slice of Billing structs that represents the entire
// billing history of the passed user. They are arranged from the latest to the oldest,
// with the first index representing the latest billing information.
func getBillsFor(username string, lim int) (billings []Billing, err error) {
	err = db.Select(&billings, "SELECT * FROM billings WHERE user = ? ORDER BY time_stamp DESC LIMIT ?",
		username, lim)
	if err != nil {
		return []Billing{}, fmt.Errorf("could not get bills for user<%s>: %v", username, err)
	}
	return billings, nil
}

// billUserFor creates a new Billing struct in the billings table for the passed username.
// An amount is also specified, and the credit flag stipulated. If it is true, the amount is
// credited, if false, the amount debited.
func billUserFor(username string, amount float64, credit bool) (current Billing, err error) {
	billings, err := getBillsFor(username, 1)
	if err != nil {
		return Billing{}, fmt.Errorf("could not add bill for user<%s>: %v", username, err)
	}
	current = billings[LATEST]
	current.TimeStamp = time.Now().Format(utils.TimeFormat)
	if credit {
		current.Credit += amount
	} else {
		current.Debit += amount
	}
	err = current.save()
	if err != nil {
		return Billing{}, fmt.Errorf("could not add bill for user<%s>: %v", username, err)
	}
	return current, nil
}

// newBillings creates the very first billing row for a newly registered user.
// Its usage should be limited to new user having just signed up, if new bills are
// to be added, use the func billUserFor()
func newBillings(username string) (Billing, error) {
	newBillingRow := Billing{
		TimeStamp: time.Now().UTC().Format(utils.TimeFormat),
		User:      username,
		Debit:     0.0,
		Credit:    0.0,
	}
	err := newBillingRow.save()
	if err != nil {
		return Billing{}, fmt.Errorf("could not create new billings row for user<%s>: %v", username, err)
	}
	return newBillingRow, nil
}
