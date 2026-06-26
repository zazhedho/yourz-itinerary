package interfacetrip

import (
	"context"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	domaintrip "yourz-itinerary/internal/domain/trip"
	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	interfacegeneric "yourz-itinerary/internal/interfaces/generic"
)

type RepoTripInterface interface {
	interfacegeneric.GenericRepository[domaintrip.Trip]

	CreateTrip(ctx context.Context, trip domaintrip.Trip, member domaintripmember.TripMember, days ...domainitineraryday.ItineraryDay) (domaintrip.Trip, error)
	ListByMember(ctx context.Context, userId string) ([]domaintrip.Trip, int64, error)
}
