package repositorytrip

import (
	"context"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	domaintrip "yourz-itinerary/internal/domain/trip"
	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	interfacetrip "yourz-itinerary/internal/interfaces/trip"
	repositorygeneric "yourz-itinerary/internal/repositories/generic"
	"yourz-itinerary/pkg/filter"

	"gorm.io/gorm"
)

type repo struct {
	*repositorygeneric.GenericRepository[domaintrip.Trip]
}

func NewTripRepo(db *gorm.DB) interfacetrip.RepoTripInterface {
	return &repo{GenericRepository: repositorygeneric.New[domaintrip.Trip](db)}
}

func (r *repo) CreateTrip(ctx context.Context, trip domaintrip.Trip, member domaintripmember.TripMember, days ...domainitineraryday.ItineraryDay) (domaintrip.Trip, error) {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&trip).Error; err != nil {
			return err
		}
		member.TripId = trip.Id
		if err := tx.Create(&member).Error; err != nil {
			return err
		}
		for i := range days {
			days[i].TripId = trip.Id
		}
		if len(days) > 0 {
			if err := tx.Create(&days).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return trip, err
}

func (r *repo) ListByMember(ctx context.Context, userId string) ([]domaintrip.Trip, int64, error) {
	var trips []domaintrip.Trip
	var total int64

	subQuery := r.DB.WithContext(ctx).
		Model(&domaintripmember.TripMember{}).
		Select("trip_id").
		Where("user_id = ?", userId)

	err := r.DB.WithContext(ctx).
		Model(&domaintrip.Trip{}).
		Where("id IN (?)", subQuery).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.DB.WithContext(ctx).
		Where("id IN (?)", subQuery).
		Order("created_at DESC").
		Find(&trips).Error
	if err != nil {
		return nil, 0, err
	}

	return trips, total, nil
}

func (r *repo) GetAll(ctx context.Context, params filter.BaseParams) ([]domaintrip.Trip, int64, error) {
	return r.GenericRepository.GetAll(ctx, params, repositorygeneric.QueryOptions{})
}
