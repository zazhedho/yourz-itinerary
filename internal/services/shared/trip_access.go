package serviceshared

import "errors"

const (
	TripRoleOwner  = "owner"
	TripRoleEditor = "editor"
	TripRoleViewer = "viewer"
)

var (
	ErrNotMember    = errors.New("not a trip member")
	ErrAccessDenied = errors.New("access denied")
)

func CanViewTrip(role string) bool {
	return role == TripRoleOwner || role == TripRoleEditor || role == TripRoleViewer
}

func CanEditTrip(role string) bool {
	return role == TripRoleOwner || role == TripRoleEditor
}

func CanManageTripMembers(role string) bool {
	return role == TripRoleOwner
}
