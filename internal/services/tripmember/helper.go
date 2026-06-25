package servicetripmember

import (
	"errors"
)

var (
	ErrMemberNotFound  = errors.New("trip member not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrDuplicateMember = errors.New("user is already a member of this trip")
	ErrOwnerRemove     = errors.New("cannot remove the trip owner")
	ErrOwnerLeave      = errors.New("owner cannot leave the trip; transfer ownership or delete the trip instead")
	ErrOwnerRoleChange = errors.New("cannot change the owner's role")
	ErrInvalidTripRole = errors.New("invalid trip role")
)
