package interfaceitineraryday

import (
	"context"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	interfacegeneric "yourz-itinerary/internal/interfaces/generic"
)

type RepoItineraryDayInterface interface {
	interfacegeneric.GenericRepository[domainitineraryday.ItineraryDay]

	ListByTrip(ctx context.Context, tripId string) ([]domainitineraryday.ItineraryDay, error)
}
