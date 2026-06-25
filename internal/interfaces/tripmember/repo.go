package interfacetripmember

import (
	"context"

	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	interfacegeneric "yourz-itinerary/internal/interfaces/generic"
)

type RepoTripMemberInterface interface {
	interfacegeneric.GenericRepository[domaintripmember.TripMember]

	GetByTripAndUser(ctx context.Context, tripId, userId string) (domaintripmember.TripMember, error)
	GetActiveByTripAndUser(ctx context.Context, tripId, userId string) (domaintripmember.TripMember, error)
	ListByTrip(ctx context.Context, tripId string) ([]domaintripmember.TripMember, error)
}
