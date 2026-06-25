package interfaceitineraryitem

import (
	"context"

	domainitineraryitem "yourz-itinerary/internal/domain/itineraryitem"
	interfacegeneric "yourz-itinerary/internal/interfaces/generic"
)

type RepoItineraryItemInterface interface {
	interfacegeneric.GenericRepository[domainitineraryitem.ItineraryItem]

	GetByDay(ctx context.Context, dayId string) ([]domainitineraryitem.ItineraryItem, error)
	GetByIDs(ctx context.Context, ids []string) ([]domainitineraryitem.ItineraryItem, error)
	Reorder(ctx context.Context, dayId string, items []domainitineraryitem.ItineraryItem) error
}
