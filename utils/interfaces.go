package utils

// Invitee is an interface for any entity that can be able
// to associate itself to a cloak
type Invitee interface {
	AssociateToCloak(cloakID string) error
}
