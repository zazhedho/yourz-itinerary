package repositoryitineraryitem

import (
	"context"
	"time"

	domainitineraryitem "yourz-itinerary/internal/domain/itineraryitem"
	interfaceitineraryitem "yourz-itinerary/internal/interfaces/itineraryitem"
	repositorygeneric "yourz-itinerary/internal/repositories/generic"
	"yourz-itinerary/pkg/filter"

	"gorm.io/gorm"
)

type repo struct {
	*repositorygeneric.GenericRepository[domainitineraryitem.ItineraryItem]
}

func NewItineraryItemRepo(db *gorm.DB) interfaceitineraryitem.RepoItineraryItemInterface {
	return &repo{GenericRepository: repositorygeneric.New[domainitineraryitem.ItineraryItem](db)}
}

func (r *repo) GetByDay(ctx context.Context, dayId string) ([]domainitineraryitem.ItineraryItem, error) {
	var items []domainitineraryitem.ItineraryItem
	err := r.DB.WithContext(ctx).Where("day_id = ?", dayId).Order("sort_order ASC").Find(&items).Error
	return items, err
}

func (r *repo) GetByIDs(ctx context.Context, ids []string) ([]domainitineraryitem.ItineraryItem, error) {
	var items []domainitineraryitem.ItineraryItem
	err := r.DB.WithContext(ctx).Where("id IN ?", ids).Find(&items).Error
	return items, err
}

func (r *repo) Reorder(ctx context.Context, dayId string, items []domainitineraryitem.ItineraryItem) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&domainitineraryitem.ItineraryItem{}).Where("id = ?", item.Id).Updates(map[string]interface{}{
				"sort_order": item.SortOrder,
				"updated_by": item.UpdatedBy,
				"updated_at": time.Now(),
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repo) GetAll(ctx context.Context, params filter.BaseParams) ([]domainitineraryitem.ItineraryItem, int64, error) {
	return r.GenericRepository.GetAll(ctx, params, repositorygeneric.QueryOptions{})
}
