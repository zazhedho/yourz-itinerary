package servicetripmember

import (
	"errors"
	"time"
	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	"yourz-itinerary/internal/dto"
)

var (
	ErrMemberNotFound  = errors.New("trip member not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrDuplicateMember = errors.New("user is already a member of this trip")
	ErrOwnerRemove     = errors.New("cannot remove the trip owner")
	ErrOwnerLeave      = errors.New("owner cannot leave the trip; transfer ownership or delete the trip instead")
	ErrOwnerRoleChange = errors.New("cannot change the owner's role")
	ErrInvalidTripRole = errors.New("invalid trip role")
	ErrTripNotFound    = errors.New("trip not found")
)

func memberToResponse(m domaintripmember.TripMember) dto.TripMemberResponse {
	mr := dto.TripMemberResponse{
		Id:        m.Id,
		TripId:    m.TripId,
		UserId:    m.UserId,
		Role:      m.Role,
		CreatedBy: m.CreatedBy,
		UpdatedBy: m.UpdatedBy,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}

	if m.UpdatedAt != nil {
		mr.UpdatedAt = new(m.UpdatedAt.Format(time.RFC3339))
	}
	if m.DeletedBy != nil {
		mr.DeletedBy = m.DeletedBy
	}
	if m.DeletedAt.Valid {
		mr.DeletedAt = new(m.DeletedAt.Time.Format(time.RFC3339))
	}

	return mr
}
