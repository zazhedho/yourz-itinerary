package serviceitineraryitem

import (
	"context"
	"strings"
	"time"

	domainitineraryitem "yourz-itinerary/internal/domain/itineraryitem"
	"yourz-itinerary/internal/dto"
	interfaceitineraryday "yourz-itinerary/internal/interfaces/itineraryday"
	interfaceitineraryitem "yourz-itinerary/internal/interfaces/itineraryitem"
	interfacetripmember "yourz-itinerary/internal/interfaces/tripmember"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/utils"
)

type ItineraryItemService struct {
	memberRepo interfacetripmember.RepoTripMemberInterface
	dayRepo    interfaceitineraryday.RepoItineraryDayInterface
	itemRepo   interfaceitineraryitem.RepoItineraryItemInterface
}

func NewItineraryItemService(
	memberRepo interfacetripmember.RepoTripMemberInterface,
	dayRepo interfaceitineraryday.RepoItineraryDayInterface,
	itemRepo interfaceitineraryitem.RepoItineraryItemInterface,
) *ItineraryItemService {
	return &ItineraryItemService{
		memberRepo: memberRepo,
		dayRepo:    dayRepo,
		itemRepo:   itemRepo,
	}
}

func (s *ItineraryItemService) CreateItem(ctx context.Context, userId, dayId string, req dto.CreateItineraryItemRequest) (dto.ItineraryItemResponse, error) {
	day, err := s.dayRepo.GetByID(ctx, dayId)
	if err != nil {
		return dto.ItineraryItemResponse{}, serviceshared.ErrDayNotFound
	}

	if err := s.checkMemberEditAccess(ctx, day.TripId, userId); err != nil {
		return dto.ItineraryItemResponse{}, err
	}

	if err := validateCoordinates(req.Latitude, req.Longitude); err != nil {
		return dto.ItineraryItemResponse{}, err
	}

	sortOrder := req.SortOrder
	if sortOrder <= 0 {
		existingItems, _ := s.itemRepo.GetByDay(ctx, dayId)
		sortOrder = len(existingItems) + 1
	}

	now := time.Now()
	item := domainitineraryitem.ItineraryItem{
		Id:           utils.CreateUUID(),
		DayId:        dayId,
		Title:        strings.TrimSpace(req.Title),
		CostEstimate: req.CostEstimate,
		SortOrder:    sortOrder,
		CreatedBy:    userId,
		UpdatedBy:    userId,
		CreatedAt:    now,
	}

	if req.Description != "" {
		item.Description = new(strings.TrimSpace(req.Description))
	}
	if req.LocationName != "" {
		item.LocationName = new(strings.TrimSpace(req.LocationName))
	}
	if req.Latitude != nil {
		item.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		item.Longitude = req.Longitude
	}
	if req.StartTime != "" {
		st, err := parseClockTime(req.StartTime)
		if err != nil {
			return dto.ItineraryItemResponse{}, err
		}
		item.StartTime = &st
	}
	if req.EndTime != "" {
		et, err := parseClockTime(req.EndTime)
		if err != nil {
			return dto.ItineraryItemResponse{}, err
		}
		item.EndTime = &et
	}

	if err := s.itemRepo.Store(ctx, item); err != nil {
		return dto.ItineraryItemResponse{}, err
	}

	return serviceshared.ItineraryItemToResponse(item), nil
}

func (s *ItineraryItemService) UpdateItem(ctx context.Context, userId, itemId string, req dto.UpdateItineraryItemRequest) (dto.ItineraryItemResponse, error) {
	item, err := s.itemRepo.GetByID(ctx, itemId)
	if err != nil {
		return dto.ItineraryItemResponse{}, ErrItemNotFound
	}

	day, err := s.dayRepo.GetByID(ctx, item.DayId)
	if err != nil {
		return dto.ItineraryItemResponse{}, serviceshared.ErrDayNotFound
	}

	if err := s.checkMemberEditAccess(ctx, day.TripId, userId); err != nil {
		return dto.ItineraryItemResponse{}, err
	}

	if err := validateCoordinates(req.Latitude, req.Longitude); err != nil {
		return dto.ItineraryItemResponse{}, err
	}

	if req.Title != "" {
		item.Title = strings.TrimSpace(req.Title)
	}
	if req.Description != "" {
		item.Description = new(strings.TrimSpace(req.Description))
	}
	if req.LocationName != "" {
		item.LocationName = new(strings.TrimSpace(req.LocationName))
	}
	if req.Latitude != nil {
		item.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		item.Longitude = req.Longitude
	}
	if req.StartTime != "" {
		st, err := parseClockTime(req.StartTime)
		if err != nil {
			return dto.ItineraryItemResponse{}, err
		}
		item.StartTime = &st
	}
	if req.EndTime != "" {
		et, err := parseClockTime(req.EndTime)
		if err != nil {
			return dto.ItineraryItemResponse{}, err
		}
		item.EndTime = &et
	}
	if req.CostEstimate >= 0 {
		item.CostEstimate = req.CostEstimate
	}
	if req.SortOrder > 0 {
		item.SortOrder = req.SortOrder
	}

	item.UpdatedBy = userId
	item.UpdatedAt = new(time.Now())

	if err := s.itemRepo.Update(ctx, item); err != nil {
		return dto.ItineraryItemResponse{}, err
	}

	return serviceshared.ItineraryItemToResponse(item), nil
}

func (s *ItineraryItemService) DeleteItem(ctx context.Context, userId, itemId string) error {
	item, err := s.itemRepo.GetByID(ctx, itemId)
	if err != nil {
		return ErrItemNotFound
	}

	day, err := s.dayRepo.GetByID(ctx, item.DayId)
	if err != nil {
		return serviceshared.ErrDayNotFound
	}

	if err := s.checkMemberEditAccess(ctx, day.TripId, userId); err != nil {
		return err
	}

	return s.itemRepo.SoftDelete(ctx, item.Id, userId)
}

func (s *ItineraryItemService) ReorderItems(ctx context.Context, userId, dayId string, req dto.ReorderItineraryItemsRequest) error {
	if len(req.ItemIds) == 0 {
		return ErrReorderEmpty
	}

	day, err := s.dayRepo.GetByID(ctx, dayId)
	if err != nil {
		return serviceshared.ErrDayNotFound
	}

	if err := s.checkMemberEditAccess(ctx, day.TripId, userId); err != nil {
		return err
	}

	items, err := s.itemRepo.GetByIDs(ctx, req.ItemIds)
	if err != nil {
		return err
	}

	foundCount := 0
	for _, item := range items {
		if item.DayId != dayId {
			return ErrReorderDifferentDay
		}
		foundCount++
	}

	if foundCount != len(req.ItemIds) {
		return ErrReorderItemsNotFound
	}

	for i, itemId := range req.ItemIds {
		for j := range items {
			if items[j].Id == itemId {
				items[j].SortOrder = i + 1
				items[j].UpdatedBy = userId
				break
			}
		}
	}

	return s.itemRepo.Reorder(ctx, dayId, items)
}

func (s *ItineraryItemService) checkMemberEditAccess(ctx context.Context, tripId, userId string) error {
	member, err := s.memberRepo.GetActiveByTripAndUser(ctx, tripId, userId)
	if err != nil || member.Id == "" {
		return serviceshared.ErrNotMember
	}
	if !serviceshared.CanEditTrip(member.Role) {
		return serviceshared.ErrAccessDenied
	}
	return nil
}
