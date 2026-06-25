package interfaceitineraryitem

import (
	"context"

	"yourz-itinerary/internal/dto"
)

type ServiceItineraryItemInterface interface {
	CreateItem(ctx context.Context, userId, dayId string, req dto.CreateItineraryItemRequest) (dto.ItineraryItemResponse, error)
	UpdateItem(ctx context.Context, userId, itemId string, req dto.UpdateItineraryItemRequest) (dto.ItineraryItemResponse, error)
	DeleteItem(ctx context.Context, userId, itemId string) error
	ReorderItems(ctx context.Context, userId, dayId string, req dto.ReorderItineraryItemsRequest) error
}
