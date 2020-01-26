package handlers

import marauderhttp "github.com/mcctor/marauders/http"

func init() {
	marauderhttp.Router.HandleFunc("/v1/users/", users).
		Methods("GET", "POST")
	marauderhttp.Router.HandleFunc("/v1/users/page/{page_number}/", allUsersPaginated).
		Methods("GET")

	// associate endpoints to new sub-router with its own permissions
	// middleware
	usersRouter := marauderhttp.Router.PathPrefix("/v1/users").Subrouter()
	usersRouter.HandleFunc("/{username}/", users).
		Methods("GET", "PUT", "DELETE")

	usersRouter.HandleFunc("/{username}/billings/", userBillings).
		Methods("GET")

	usersRouter.HandleFunc("/{username}/billings/{billing_id}/", userBilling).
		Methods("GET")

	usersRouter.HandleFunc("/{username}/cloaks/", userCloaks).
		Methods("GET", "POST")

	usersRouter.HandleFunc("/{username}/cloaks/{cloak_id}/", userCloak).
		Methods("GET", "PUT", "DELETE")

	usersRouter.HandleFunc("/{username}/invitation-links/", userInvitationLinks).
		Methods("GET", "POST")

	usersRouter.HandleFunc("/{username}/invitation-links/{invitation_link_id}/", userInvitationLink).
		Methods("GET", "DELETE")

	usersRouter.HandleFunc("/{username}/auth-token/", userAuthToken).
		Methods("POST")

	usersRouter.HandleFunc("/{username}/devices/", userDevices).
		Methods("GET", "POST")

	usersRouter.HandleFunc("/{username}/devices/{device_id}/", userDevice).
		Methods("GET", "DELETE")

	usersRouter.HandleFunc("/{username}/devices/{device_id}/location-history/", userDeviceLocData).
		Methods("GET", "POST")

	// register middleware that ensures only owning users can access the private endpoints
	usersRouter.Use(marauderhttp.ApplyOwnerPermission)
}
