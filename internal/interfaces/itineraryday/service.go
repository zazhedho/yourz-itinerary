package interfaceitineraryday

import (
	"context"

	"yourz-itinerary/internal/dto"
)

type ServiceItineraryDayInterface interface {
	CreateDay(ctx context.Context, userId, tripId string, req dto.CreateItineraryDayRequest) (dto.ItineraryDayResponse, error)
	UpdateDay(ctx context.Context, userId, dayId string, req dto.UpdateItineraryDayRequest) (dto.ItineraryDayResponse, error)
	DeleteDay(ctx context.Context, userId, dayId string) error
}
