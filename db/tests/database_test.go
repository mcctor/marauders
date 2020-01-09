package tests

import (
	"database/sql"
	"testing"
	"time"

	"github.com/mcctor/marauders/db"
	"github.com/mcctor/marauders/utils"
)

func TestCreateNewUser(t *testing.T) {
	t.Log("Given the need to test for the successful creation of a new user.")
	{
		john, err := db.NewUser("john", "john@somewhere.com")
		if err != nil {
			t.Fatal("\t\tShould be able to successfully create a new user:", failMark, err)
		}
		t.Log("\t\tShould be able to successfully create a new user:", passMark, john)
	}
}

func TestGetExistingUser(t *testing.T) {
	t.Log("Given the need to test for the successful fetching of an existing user.")
	{
		existingUsername := "john"
		john, err := db.GetUser(existingUsername)
		if err != nil {
			t.Fatal("\t\tShould be able to successfully fetch an existing user:", failMark, err)
		}
		t.Log("\t\tShould be able to successfully fetch an existing user:", passMark, john)
	}
}

func TestGetAllUsers(t *testing.T) {
	t.Log("Given the need to test for the successful fetching of all the existing users.")
	{
		allUsers, err := db.GetUsers(5)
		if err != nil {
			t.Fatal("\t\tShould be able to successfully fetch all existing users:", failMark, err)
		}
		t.Log("\t\tShould be able to successfully fetch all existing users:", passMark, allUsers)
	}
}

func TestUpdateExistingUser(t *testing.T) {
	t.Log("Given the need to test the successful updating of the columns for an existing user.")
	{
		existingUsername := "john"
		john, _ := db.GetUser(existingUsername)

		john.Fname = sql.NullString{
			String: "Jonathan",
			Valid:  true,
		}
		john.Lname = sql.NullString{
			String: "Baker",
			Valid:  true,
		}
		john.Email = "jonathanbaker@something.com"

		err := john.Update()
		if err != nil {
			t.Fatal("\t\tShould successfully update the changed columns:", failMark, err)
		}
		t.Log("\t\tShould successfully update the changed columns:", passMark, john)
	}
}

func TestCreateNewDeviceForUser(t *testing.T) {
	t.Log("Given the need to test the successful creation of a new device for an existing user.")
	{
		existingUsername := "john"
		john, _ := db.GetUser(existingUsername)

		johnDevice, err := john.NewDevice(444)
		if err != nil {
			t.Fatal("\t\tShould successfully create a new device for an existing user:", failMark, err)
		}
		t.Log("\t\tShould successfully create a new device for an existing user:", passMark, johnDevice)
	}
}

func TestGetDevicesForUser(t *testing.T) {
	t.Log("Given the need to test the successful fetching of devices belonging to an existing user.")
	{
		existingUsername := "john"
		john, _ := db.GetUser(existingUsername)

		johnDevices, err := john.Devices()
		if err != nil {
			t.Fatal("\t\tShould successfully fetch an existing user's devices:", failMark, err)
		}
		t.Log("\t\tShould successfully fetch an existing user's devices:", passMark, johnDevices)
	}
}

func TestGetUserBillings(t *testing.T) {
	t.Log("Given the need to test the successful fetching of an existing user's bills.")
	{
		existingUsername := "john"
		john, _ := db.GetUser(existingUsername)

		johnBills, err := john.Billings(5)
		if err != nil {
			t.Fatal("\t\tShould successfully fetch an existing user's bills:", failMark, err)
		}
		t.Log("\t\tShould successfully fetch an existing user's bills:", passMark, johnBills)
	}
}

func TestCreditUser(t *testing.T) {
	t.Log("Given the need to test the successful crediting of an existing user's billing account.")
	{
		existingUsername := "john"
		john, _ := db.GetUser(existingUsername)

		newBill, err := john.CreditUserAmount(100)
		if err != nil {
			t.Fatal("\t\tShould successfully credit the given amount to a user's account:", failMark, err)
		}
		t.Log("\t\tShould successfully credit the given amount to a user's account:", passMark, newBill)
	}
}

func TestDebitUser(t *testing.T) {
	t.Log("Given the need to test the successful debiting of an existing user's billing account.")
	{
		// sleep for a bit to avoid conflict in the time_stamp primary key
		// for the billings table
		time.Sleep(1 * time.Second)

		existingUsername := "john"
		john, _ := db.GetUser(existingUsername)
		newBill, err := john.DebitUserAmount(100)
		if err != nil {
			t.Fatal("\t\tShould successfully debit the given amount to a user's account:", failMark, err)
		}
		t.Log("\t\tShould successfully debit the given amount to a user's account:", passMark, newBill)
	}
}

func TestNewAuthTokenForUser(t *testing.T) {
	t.Log("Given the need to test the successful creation of an auth token for an existing user.")
	{
		existingUsername := "john"
		john, _ := db.GetUser(existingUsername)
		johnAuthToken, err := john.NewAuthToken()
		if err != nil {
			t.Fatal("\t\tShould successfully create a new authentication token for a user:", failMark, err)
		}
		t.Log("\t\tShould successfully create a new authentication token for a user:", passMark, johnAuthToken)
	}
}

func TestFetchingAuthTokenForUser(t *testing.T) {
	t.Log("Given the need to test the successful fetching of an auth token for an existing user.")
	{
		existingUsername := "john"
		john, _ := db.GetUser(existingUsername)
		johnAuthToken, err := john.AuthToken()
		if err != nil {
			t.Fatal("\t\tShould successfully fetch the auth token of an existing user:", failMark, err)
		}
		t.Log("\t\tShould successfully fetch the auth token of an existing user:", passMark, johnAuthToken)
	}
}

func TestCreateNewPasswordForUser(t *testing.T) {
	t.Log("Given the need to test the successful creation of a password for an existing user.")
	{
		existingUsername := "john"
		john, _ := db.GetUser(existingUsername)
		johnPassword, err := john.NewPassword("somepassword1234")
		if err != nil {
			t.Fatal("\t\tShould successfully create a password for an existing user:", failMark, err)
		}
		t.Log("\t\tShould successfully create a password for an existing user:", passMark, johnPassword)
	}
}

func TestFetchPasswordForUser(t *testing.T) {
	t.Log("Given the need to test the successful fetching of the password for an existing user.")
	{
		existingUsername := "john"
		john, _ := db.GetUser(existingUsername)
		johnPassword, err := john.Password()
		if err != nil {
			t.Fatal("\t\tShould successfully fetch the password for an existing user:", failMark, err)
		}
		t.Log("\t\tShould successfully fetch the password for an existing user:", passMark, johnPassword)
	}
}

func TestCreateNewCloakForUser(t *testing.T) {
	t.Log("Given the need to test the successful creation of a new cloak for an existing user.")
	{
		existingUsername := "john"
		john, _ := db.GetUser(existingUsername)

		wakeTime, _ := time.Parse(utils.TimeFormat, "2001-01-01 06:00:00")
		sleepTime, _ := time.Parse(utils.TimeFormat, "2001-01-01 18:00:00")
		duration, _ := time.Parse(utils.TimeFormat, "2019-11-12 07:00:00")

		cloak, err := john.NewCloak("random cloak", "awesome", wakeTime, sleepTime, duration,
			"pinpoint", 13, true, true, false,
			true)
		if err != nil {
			t.Fatal("\t\tShould successfully create a new cloak for the existing user:", failMark, err)
		}
		t.Log("\t\tShould successfully create a new cloak for the existing user:", passMark, cloak)
	}
}

func TestFetchCloaksForUser(t *testing.T) {
	t.Log("Given the need to test the successful fetching of the cloaks owned by an existing user.")
	{
		existingUsername := "john"
		john, _ := db.GetUser(existingUsername)

		cloaks, err := john.Cloaks()
		if err != nil {
			t.Fatal("\t\tShould successfully fetch the cloaks owned by the existing user:", failMark, err)
		}
		t.Log("\t\tShould successfully fetch the cloaks owned by the existing user:", passMark, cloaks[0])
	}
}

func TestFetchInviteLinksForUser(t *testing.T) {
	t.Log("Given the need to test for the successful fetching of invite links created by an existing user.")
	{
		existingUsername := "john"
		john, _ := db.GetUser(existingUsername)

		cloaks, _ := john.Cloaks()
		someCloaks := cloaks[0]

		expiry, _ := time.Parse(utils.TimeFormat, "2020-01-10 06:00:00")
		_, err := someCloaks.NewInviteLink(existingUsername, 4, expiry)
		if err != nil {
			t.Fatal("\t\tShould be able to create new invite link for user:", failMark, err)
		}

		inviteLinks, err := john.InviteLinks()
		if err != nil {
			t.Fatal("\t\tShould be able to fetch the invite links for the existing user:", failMark, err)
		}
		t.Log("\t\tShould be able to fetch the invite links for the existing user:", passMark, inviteLinks[0])
	}
}

func TestUserDeletion(t *testing.T) {
	t.Log("Given the need to test for the successful deletion of an existing user.")
	{
		someone, _ := db.NewUser("someone", "")
		err := someone.Delete()
		if err != nil {
			t.Fatal("\t\tShould successfully delete the existing user:", failMark, err)
		}
		t.Log("\t\tShould successfully delete the existing user.", passMark)
	}
}

func TestUserAuthToken(t *testing.T) {
	existingUsername := "john"
	john, _ := db.GetUser(existingUsername)
	johnAuthToken, _ := john.AuthToken()


	t.Log("Given the need to test for the successful renewal of an auth token given its valid refresh token.")
	{
		refreshToken := johnAuthToken.RefreshToken
		newToken, err := johnAuthToken.Renew(refreshToken)
		if err != nil {
			t.Fatal("\t\tShould successfully renew old token given a valid refresh token:", failMark, err)
		}
		t.Log("\t\tShould successfully renew old token given a valid refresh token:", passMark, newToken)
	}

	t.Log("Given the need to test the failure of renewing an auth token with an invalid refresh token.")
	{
		invalidRefreshToken := "invalid refresh token"
		_, err := johnAuthToken.Renew(invalidRefreshToken)
		if err == nil {
			t.Fatal("\t\tShould fail in renewing an auth token given an invalid refresh token:", failMark)
		}
		t.Log("\t\tShould fail in renewing an auth token given an invalid refresh token:", passMark, err)
	}

	t.Log("Given the need to test for the validity of an auth token before its expiry.")
	{
		// should pass since this is a recently created auth token, and
		// its expiry is set some days after its creation
		if johnAuthToken.IsValid() {
			t.Log("\t\tShould be valid before its expiry", passMark)
		} else {
			t.Fatal("\t\tShould be valid before its expiry", failMark)
		}
	}

	t.Log("Given the need to test for the validity of an auth token after its expiry")
	{
		// set a time long in the past
		johnAuthToken.Expiry = "1000-01-01 12:00:00"

		if johnAuthToken.IsValid() {
			t.Fatal("\t\tShould not be valid after its expiry", failMark)
		} else {
			t.Log("\t\tShould not be valid after its expiry", passMark)
		}
	}
}

func TestUserDevice(t *testing.T) {
	existingUsername := "john"
	john, _ := db.GetUser(existingUsername)
	johnDevices, _ := john.Devices()
	device := johnDevices[0]

	t.Log("Given the need to test for the successful creation of a location snapshot for an existing device.")
	{
		err := device.NewLocationSnapshot(time.Now().Format(utils.TimeFormat), 9.123, -12.244)
		if err != nil {
			t.Fatal("\t\tShould successfully create a new location snapshot:", failMark, err)
		}
		t.Log("\t\tShould successfully create a new location snapshot", passMark)
	}

	t.Log("Given the need to test the successful fetching of a device's location snapshot history.")
	{
		locSnaps, err := device.LocationSnapshots(1)
		if err != nil {
			t.Fatal("\t\tShould successfully fetch a device's location snapshot history:", failMark, err)
		}
		t.Log("\t\tShould successfully fetch a device's location snapshot history: ", passMark, locSnaps)
	}

	t.Log("Given the need to test for the successful association of a device to an existing cloak.")
	{
		existingCloaks, _ := john.Cloaks()
		existingCloakID := existingCloaks[0].ID

		err := device.AssociateToCloak(existingCloakID)
		if err != nil {
			t.Fatal("\t\tShould successfully associate device to the existing cloak:", failMark, err)
		}
		t.Log("\t\tShould successfully associate device to the existing cloak:", passMark)
	}

	t.Log("Given the need to test for the successful fetching of the cloaks an existing device is associated to.")
	{
		cloaks, err := device.AssociatedCloaks()
		if err != nil {
			t.Fatal("\t\tShould successfully fetch the cloaks an existing device is associated to:", failMark, err)
		}
		if len(cloaks) != 1 {
			t.Fatal("\t\tAssociated cloaks should be a total of 1: Found", len(cloaks), failMark)
		}
		t.Log("\t\tShould successfully fetch the cloaks an existing device is associated to", passMark, cloaks)
	}

	t.Log("Given the need to test the successful dissociation of a device from an associated cloak.")
	{
		existingCloaks, _ := john.Cloaks()
		existingCloakID := existingCloaks[0].ID

		newDevice, _ := john.NewDevice(555)
		_ = newDevice.AssociateToCloak(existingCloakID)

		err := newDevice.DissociateFromCloak(existingCloakID)
		if err != nil {
			t.Fatal("\t\tShould successfully dissociate device from an associated cloak:", failMark, err)
		}
		t.Log("\t\tShould successfully dissociate device from an associated cloak", passMark)
	}

	t.Log("Given the need to test the successful deletion of an existing device.")
	{
		existingDevices, _ := john.Devices()
		var existingDevice db.Device
		for _, device := range existingDevices {
			if device.ID == 555 {
				existingDevice = device
			}
		}

		err := existingDevice.Delete()
		if err != nil {
			t.Fatal("\t\tShould successfully delete an existing device:", failMark, err)
		}
		t.Log("\t\tShould successfully delete an existing device", passMark)
	}
}

func TestUserPassword(t *testing.T) {
	existingUsername := "john"
	existingPassword := "somepassword1234"
	john, _ := db.GetUser(existingUsername)
	johnPasswd, _ := john.Password()

	t.Log("Given the need to test the validity of a password for an existing user against a wrong password.")
	{
		if johnPasswd.CheckIf("wrongpasswd") {
			t.Fatal("\t\tShould not be the correct password", failMark)
		}
		t.Log("\t\tShould not be the correct password", passMark)
	}

	t.Log("Given the need to test the validity of a password for an existing user against the right password.")
	{
		if johnPasswd.CheckIf(existingPassword) {
			t.Log("\t\tShould be the correct password", passMark)
		} else {
			t.Fatal("\t\tShould be the correct password", failMark)
		}
	}

	t.Log("Given the need to test the successful changing of an existing user's password.")
	{
		err := johnPasswd.ChangeTo("someotherpassword1234")
		if err != nil {
			t.Fatal("\t\tShould successfully change the existing password into the new one:", failMark, err)
		}
		t.Log("\t\tShould successfully change the existing password into the new one:", passMark, johnPasswd.Hash)
	}
}

func TestUserCloak(t *testing.T) {
	existingUsername := "john"
	john, _ := db.GetUser(existingUsername)
	johnCloaks, _ := john.Cloaks()
	existingCloak := johnCloaks[0]

	// create a new cloak
	wakeTime, _ := time.Parse(utils.TimeFormat, "2001-01-01 06:00:00")
	sleepTime, _ := time.Parse(utils.TimeFormat, "2001-01-01 18:00:00")
	duration, _ := time.Parse(utils.TimeFormat, "2019-09-09 07:00:00")
	newCloak, _ := john.NewCloak("second cloak", "very nice cloak", wakeTime, sleepTime, duration,
		"pinpoint", 15, true, true, false, true)

	t.Log("Given the need to test the successful fetching of an existing cloak.")
	{
		cloak, err := db.GetCloakByID(existingCloak.ID)
		if err != nil {
			t.Fatal("\t\tShould be able to successfully fetch the existing cloak:", failMark, err)
		}
		t.Log("\t\tShould be able to successfully fetch the existing cloak:", passMark, cloak)
	}

	t.Log("Given the need to test the successful fetching of the members of an existing cloaks.")
	{
		cloak, _ := db.GetCloakByID(existingCloak.ID)
		devices, err := cloak.Members()
		if err != nil {
			t.Fatal("\t\tShould be able to successfully fetch the members of an existing cloak:", failMark, err)
		}
		t.Log("\t\tShould be able to successfully fetch the members of an existing cloak:", passMark, devices)
	}

	t.Log("Given the need to test the successful fetching of the associated members of an existing cloak.")
	{
		cloak, _ := db.GetCloakByID(existingCloak.ID)
		devices, err := cloak.AssociatedMembers()
		if err != nil {
			t.Fatal("\t\tShould be able to successfully fetch the associated members of an existing cloak",
				failMark, err)
		}
		t.Log("\t\tShould be able to successfully fetch the associated members of an existing cloak",
			passMark, devices)
	}

	t.Log("Given the need to test the successful update of an existing cloak's column fields.")
	{
		cloak, _ := db.GetCloakByID(existingCloak.ID)

		// change a couple of fields
		cloak.Name = "The Fun Cloak"
		cloak.Active = false
		cloak.MemberLimit = 50
		cloak.Description = "A cloak to die for"
		cloak.Sleep = "2001-01-01 20:00:00"
		cloak.Accuracy = "city"
		cloak.Private = false
		err := cloak.Update()
		if err != nil {
			t.Fatal("\t\tShould be able to successfully update an existing cloak's column fields:", failMark, err)
		}
		t.Log("\t\tShould be able to successfully update an existing cloak's column fields", passMark)
	}

	t.Log("Given the need to test the successful association of a cloak to another cloak.")
	{
		cloak, _ := db.GetCloakByID(existingCloak.ID)

		err := cloak.AddPermittedCloak(newCloak.ID)
		if err != nil {
			t.Fatal("\t\tShould be able to successfully associate cloak to another cloak:", failMark, err)
		}
		t.Log("\t\tShould be able to successfully associate cloak to another cloak", passMark)
	}

	t.Log("Given the need to test the successful dissociation cloak from another cloak.")
	{
		cloak, _ := db.GetCloakByID(existingCloak.ID)

		err := cloak.RemovePermittedCloak(newCloak.ID)
		if err != nil {
			t.Fatal("\t\tShould be able to successfully dissociate cloak from another cloak:", failMark, err)
		}
		t.Log("\t\tShould be able to successfully dissociate cloak from another cloak", passMark)
	}

	t.Log("Given the need to successfully create a new invite link for an existing cloak.")
	{
		cloak, _ := db.GetCloakByID(existingCloak.ID)

		expiry, _ := time.Parse(utils.TimeFormat, "2020-10-01 08:00:00")
		inviteLink, err := cloak.NewInviteLink(john.Username, 10, expiry)
		if err != nil {
			t.Fatal("\t\tShould be able to create a new invite link for an existing cloak:", passMark, err)
		}
		t.Log("\t\tShould be able to create a new invite link for an existing cloak:", passMark, inviteLink)
	}

	t.Log("Given the need to successfully fetch the location snapshots for a user filtered by this cloak's rules")
	{
		cloak, _ := db.GetCloakByID(existingCloak.ID)
		johnDevices, _ := john.Devices()
		device := johnDevices[0]

		locSnaps, err := cloak.LocationSnapshotsForMember(device, 1)
		if err != nil {
			t.Fatal("\t\tShould be able to fetch the location snapshots for a user filtered by this cloak's rules:",
				failMark, err)
		}
		t.Log("\t\tShould be able to fetch the location snapshots for a user filtered by this cloak's rules:",
			passMark, locSnaps)
	}
}

func TestInviteDeviceByLink(t *testing.T) {
	existingUsername := "john"
	john, _ := db.GetUser(existingUsername)
	newDevice, _ := john.NewDevice(789)

	t.Log("Given the need to test the successful inviting of a device based on a valid invite link.")
	{
		links, _ := john.InviteLinks()
		link := links[0]

		err := link.Invite(newDevice)
		if err != nil {
			t.Fatal("\t\tShould be able to successfully invite a device through a valid link:", failMark, err)
		}
		t.Log("\t\tShould be able to successfully invite a device through a valid link", passMark)
	}
}