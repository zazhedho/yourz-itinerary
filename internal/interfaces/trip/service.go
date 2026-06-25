package interfacetrip

import (
	"context"

	"yourz-itinerary/internal/dto"
)

type ServiceTripInterface interface {
	CreateTrip(ctx context.Context, userId string, req dto.CreateTripRequest) (dto.TripDetailResponse, error)
	GetTripDetail(ctx context.Context, userId, tripId string) (dto.TripDetailResponse, error)
	ListTrips(ctx context.Context, userId string) ([]dto.TripListResponse, error)
	UpdateTrip(ctx context.Context, userId, tripId string, req dto.UpdateTripRequest) (dto.TripDetailResponse, error)
	DeleteTrip(ctx context.Context, userId, tripId string) error
}
