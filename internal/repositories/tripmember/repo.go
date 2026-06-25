package repositorytripmember

import (
	"context"

	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	interfacetripmember "yourz-itinerary/internal/interfaces/tripmember"
	repositorygeneric "yourz-itinerary/internal/repositories/generic"
	"yourz-itinerary/pkg/filter"

	"gorm.io/gorm"
)

type repo struct {
	*repositorygeneric.GenericRepository[domaintripmember.TripMember]
}

func NewTripMemberRepo(db *gorm.DB) interfacetripmember.RepoTripMemberInterface {
	return &repo{GenericRepository: repositorygeneric.New[domaintripmember.TripMember](db)}
}

func (r *repo) GetByTripAndUser(ctx context.Context, tripId, userId string) (domaintripmember.TripMember, error) {
	var member domaintripmember.TripMember
	err := r.DB.WithContext(ctx).
		Where("trip_id = ? AND user_id = ?", tripId, userId).
		First(&member).Error
	return member, err
}

func (r *repo) GetActiveByTripAndUser(ctx context.Context, tripId, userId string) (domaintripmember.TripMember, error) {
	var member domaintripmember.TripMember
	err := r.DB.WithContext(ctx).
		Where("trip_id = ? AND user_id = ?", tripId, userId).
		First(&member).Error
	return member, err
}

func (r *repo) ListByTrip(ctx context.Context, tripId string) ([]domaintripmember.TripMember, error) {
	var members []domaintripmember.TripMember
	err := r.DB.WithContext(ctx).Where("trip_id = ?", tripId).Find(&members).Error
	return members, err
}

func (r *repo) GetAll(ctx context.Context, params filter.BaseParams) ([]domaintripmember.TripMember, int64, error) {
	return r.GenericRepository.GetAll(ctx, params, repositorygeneric.QueryOptions{})
}
