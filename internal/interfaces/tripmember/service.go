package interfacetripmember

import (
	"context"

	"yourz-itinerary/internal/dto"
)

type ServiceTripMemberInterface interface {
	AddMember(ctx context.Context, userId, tripId string, req dto.AddTripMemberRequest) (dto.TripMemberResponse, error)
	UpdateMemberRole(ctx context.Context, userId, tripId, memberId string, req dto.UpdateTripMemberRoleRequest) (dto.TripMemberResponse, error)
	RemoveMember(ctx context.Context, userId, tripId, memberId string) error
	LeaveTrip(ctx context.Context, userId, tripId string) error
}
