package servicetrip

import (
	"context"
	"strings"
	"time"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	domainitineraryitem "yourz-itinerary/internal/domain/itineraryitem"
	domaintrip "yourz-itinerary/internal/domain/trip"
	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	"yourz-itinerary/internal/dto"
	interfaceitineraryday "yourz-itinerary/internal/interfaces/itineraryday"
	interfaceitineraryitem "yourz-itinerary/internal/interfaces/itineraryitem"
	interfacetrip "yourz-itinerary/internal/interfaces/trip"
	interfacetripmember "yourz-itinerary/internal/interfaces/tripmember"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/utils"
)

type TripService struct {
	tripRepo   interfacetrip.RepoTripInterface
	memberRepo interfacetripmember.RepoTripMemberInterface
	dayRepo    interfaceitineraryday.RepoItineraryDayInterface
	itemRepo   interfaceitineraryitem.RepoItineraryItemInterface
}

func NewTripService(
	tripRepo interfacetrip.RepoTripInterface,
	memberRepo interfacetripmember.RepoTripMemberInterface,
	dayRepo interfaceitineraryday.RepoItineraryDayInterface,
	itemRepo interfaceitineraryitem.RepoItineraryItemInterface,
) *TripService {
	return &TripService{tripRepo: tripRepo, memberRepo: memberRepo, dayRepo: dayRepo, itemRepo: itemRepo}
}

func (s *TripService) CreateTrip(ctx context.Context, userId string, req dto.CreateTripRequest) (dto.TripDetailResponse, error) {
	now := time.Now()
	timezone := strings.TrimSpace(req.Timezone)
	if timezone == "" {
		timezone = "Asia/Jakarta"
	}
	if _, err := time.LoadLocation(timezone); err != nil {
		return dto.TripDetailResponse{}, ErrInvalidTimezone
	}

	currencyCode := strings.TrimSpace(strings.ToUpper(req.CurrencyCode))
	if currencyCode == "" {
		currencyCode = "IDR"
	}
	if !isValidCurrencyCode(currencyCode) {
		return dto.TripDetailResponse{}, ErrInvalidCurrency
	}

	tripId := utils.CreateUUID()
	memberId := utils.CreateUUID()

	trip := domaintrip.Trip{
		Id:           tripId,
		OwnerId:      userId,
		Title:        strings.TrimSpace(req.Title),
		Timezone:     timezone,
		CurrencyCode: currencyCode,
		Status:       "draft",
		CreatedBy:    userId,
		UpdatedBy:    userId,
		CreatedAt:    now,
	}

	if req.Destination != "" {
		trip.Destination = new(strings.TrimSpace(req.Destination))
	}
	if req.StartDate != "" {
		t, err := serviceshared.ParseDate(req.StartDate)
		if err != nil {
			return dto.TripDetailResponse{}, err
		}
		trip.StartDate = &t
	}
	if req.EndDate != "" {
		t, err := serviceshared.ParseDate(req.EndDate)
		if err != nil {
			return dto.TripDetailResponse{}, err
		}
		trip.EndDate = &t
	}

	if trip.StartDate != nil && trip.EndDate != nil && trip.EndDate.Before(*trip.StartDate) {
		return dto.TripDetailResponse{}, ErrInvalidDateRange
	}

	member := domaintripmember.TripMember{
		Id:        memberId,
		TripId:    tripId,
		UserId:    userId,
		Role:      serviceshared.TripRoleOwner,
		CreatedBy: userId,
		UpdatedBy: userId,
		CreatedAt: now,
	}

	createdTrip, err := s.tripRepo.CreateTrip(ctx, trip, member)
	if err != nil {
		return dto.TripDetailResponse{}, err
	}

	return tripToDetail(createdTrip, []domaintripmember.TripMember{member}, nil, nil), nil
}

func (s *TripService) GetTripDetail(ctx context.Context, userId, tripId string) (dto.TripDetailResponse, error) {
	member, err := s.memberRepo.GetActiveByTripAndUser(ctx, tripId, userId)
	if err != nil || member.Id == "" {
		return dto.TripDetailResponse{}, serviceshared.ErrNotMember
	}

	trip, err := s.tripRepo.GetByID(ctx, tripId)
	if err != nil {
		return dto.TripDetailResponse{}, serviceshared.ErrTripNotFound
	}

	members, days, itemsByDay, err := s.loadTripRelations(ctx, tripId)
	if err != nil {
		return dto.TripDetailResponse{}, err
	}

	return tripToDetail(trip, members, days, itemsByDay), nil
}

func (s *TripService) ListTrips(ctx context.Context, userId string) ([]dto.TripListResponse, error) {
	trips, _, err := s.tripRepo.ListByMember(ctx, userId)
	if err != nil {
		return nil, err
	}

	result := make([]dto.TripListResponse, 0, len(trips))
	for _, trip := range trips {
		members, days, _, err := s.loadTripRelations(ctx, trip.Id)
		if err != nil {
			return nil, err
		}
		result = append(result, tripToList(trip, len(members), len(days)))
	}

	return result, nil
}

func (s *TripService) UpdateTrip(ctx context.Context, userId, tripId string, req dto.UpdateTripRequest) (dto.TripDetailResponse, error) {
	trip, err := s.tripRepo.GetByID(ctx, tripId)
	if err != nil {
		return dto.TripDetailResponse{}, serviceshared.ErrTripNotFound
	}

	member, err := s.memberRepo.GetActiveByTripAndUser(ctx, tripId, userId)
	if err != nil || member.Id == "" {
		return dto.TripDetailResponse{}, serviceshared.ErrNotMember
	}
	if !serviceshared.CanEditTrip(member.Role) {
		return dto.TripDetailResponse{}, serviceshared.ErrAccessDenied
	}

	if req.Title != "" {
		trip.Title = strings.TrimSpace(req.Title)
	}
	if req.Destination != "" {
		trip.Destination = new(strings.TrimSpace(req.Destination))
	}
	if req.Status != "" {
		trip.Status = strings.TrimSpace(req.Status)
	}
	if req.Timezone != "" {
		tz := strings.TrimSpace(req.Timezone)
		if _, err := time.LoadLocation(tz); err != nil {
			return dto.TripDetailResponse{}, ErrInvalidTimezone
		}
		trip.Timezone = tz
	}
	if req.CurrencyCode != "" {
		cc := strings.TrimSpace(strings.ToUpper(req.CurrencyCode))
		if !isValidCurrencyCode(cc) {
			return dto.TripDetailResponse{}, ErrInvalidCurrency
		}
		trip.CurrencyCode = cc
	}
	if req.StartDate != "" {
		t, err := serviceshared.ParseDate(req.StartDate)
		if err != nil {
			return dto.TripDetailResponse{}, err
		}
		trip.StartDate = &t
	}
	if req.EndDate != "" {
		t, err := serviceshared.ParseDate(req.EndDate)
		if err != nil {
			return dto.TripDetailResponse{}, err
		}
		trip.EndDate = &t
	}

	if trip.StartDate != nil && trip.EndDate != nil && trip.EndDate.Before(*trip.StartDate) {
		return dto.TripDetailResponse{}, ErrInvalidDateRange
	}

	trip.UpdatedBy = userId
	trip.UpdatedAt = new(time.Now())

	if err := s.tripRepo.Update(ctx, trip); err != nil {
		return dto.TripDetailResponse{}, err
	}

	members, days, itemsByDay, err := s.loadTripRelations(ctx, tripId)
	if err != nil {
		return dto.TripDetailResponse{}, err
	}

	return tripToDetail(trip, members, days, itemsByDay), nil
}

func (s *TripService) DeleteTrip(ctx context.Context, userId, tripId string) error {
	trip, err := s.tripRepo.GetByID(ctx, tripId)
	if err != nil {
		return serviceshared.ErrTripNotFound
	}

	if trip.OwnerId != userId {
		return serviceshared.ErrAccessDenied
	}

	return s.tripRepo.SoftDelete(ctx, trip.Id, userId)
}

func (s *TripService) loadTripRelations(ctx context.Context, tripId string) ([]domaintripmember.TripMember, []domainitineraryday.ItineraryDay, map[string][]domainitineraryitem.ItineraryItem, error) {
	members, err := s.memberRepo.ListByTrip(ctx, tripId)
	if err != nil {
		return nil, nil, nil, err
	}

	days, err := s.dayRepo.ListByTrip(ctx, tripId)
	if err != nil {
		return nil, nil, nil, err
	}

	itemsByDay := make(map[string][]domainitineraryitem.ItineraryItem, len(days))
	for _, day := range days {
		items, err := s.itemRepo.GetByDay(ctx, day.Id)
		if err != nil {
			return nil, nil, nil, err
		}
		itemsByDay[day.Id] = items
	}

	return members, days, itemsByDay, nil
}
