package repositoryitineraryday

import (
	"context"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	interfaceitineraryday "yourz-itinerary/internal/interfaces/itineraryday"
	repositorygeneric "yourz-itinerary/internal/repositories/generic"
	"yourz-itinerary/pkg/filter"

	"gorm.io/gorm"
)

type repo struct {
	*repositorygeneric.GenericRepository[domainitineraryday.ItineraryDay]
}

func NewItineraryDayRepo(db *gorm.DB) interfaceitineraryday.RepoItineraryDayInterface {
	return &repo{GenericRepository: repositorygeneric.New[domainitineraryday.ItineraryDay](db)}
}

func (r *repo) ListByTrip(ctx context.Context, tripId string) ([]domainitineraryday.ItineraryDay, error) {
	var days []domainitineraryday.ItineraryDay
	err := r.DB.WithContext(ctx).Where("trip_id = ?", tripId).Order("day_number ASC").Find(&days).Error
	return days, err
}

func (r *repo) GetAll(ctx context.Context, params filter.BaseParams) ([]domainitineraryday.ItineraryDay, int64, error) {
	return r.GenericRepository.GetAll(ctx, params, repositorygeneric.QueryOptions{})
}
