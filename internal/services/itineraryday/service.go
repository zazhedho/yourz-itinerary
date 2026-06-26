package serviceitineraryday

import (
	"context"
	"strings"
	"time"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	"yourz-itinerary/internal/dto"
	interfaceitineraryday "yourz-itinerary/internal/interfaces/itineraryday"
	interfacetrip "yourz-itinerary/internal/interfaces/trip"
	interfacetripmember "yourz-itinerary/internal/interfaces/tripmember"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/utils"
)

type ItineraryDayService struct {
	memberRepo interfacetripmember.RepoTripMemberInterface
	dayRepo    interfaceitineraryday.RepoItineraryDayInterface
	tripRepo   interfacetrip.RepoTripInterface
}

func NewItineraryDayService(
	memberRepo interfacetripmember.RepoTripMemberInterface,
	dayRepo interfaceitineraryday.RepoItineraryDayInterface,
	tripRepo interfacetrip.RepoTripInterface,
) *ItineraryDayService {
	return &ItineraryDayService{memberRepo: memberRepo, dayRepo: dayRepo, tripRepo: tripRepo}
}

func (s *ItineraryDayService) CreateDay(ctx context.Context, userId, tripId string, req dto.CreateItineraryDayRequest) (dto.ItineraryDayResponse, error) {
	member, err := s.memberRepo.GetActiveByTripAndUser(ctx, tripId, userId)
	if err != nil || member.Id == "" {
		return dto.ItineraryDayResponse{}, serviceshared.ErrNotMember
	}
	if !serviceshared.CanEditTrip(member.Role) {
		return dto.ItineraryDayResponse{}, serviceshared.ErrAccessDenied
	}

	now := time.Now()
	day := domainitineraryday.ItineraryDay{
		Id:        utils.CreateUUID(),
		TripId:    tripId,
		DayNumber: req.DayNumber,
		CreatedBy: userId,
		UpdatedBy: userId,
		CreatedAt: now,
	}

	if req.Title != "" {
		day.Title = new(strings.TrimSpace(req.Title))
	}
	if req.Date != "" {
		parsed, err := serviceshared.ParseDate(req.Date)
		if err != nil {
			return dto.ItineraryDayResponse{}, err
		}
		day.Date = &parsed
	}

	if err := s.dayRepo.Store(ctx, day); err != nil {
		return dto.ItineraryDayResponse{}, err
	}
	if err := s.syncTripDateRange(ctx, tripId, userId); err != nil {
		return dto.ItineraryDayResponse{}, err
	}

	return serviceshared.ItineraryDayToResponse(day, nil), nil
}

func (s *ItineraryDayService) UpdateDay(ctx context.Context, userId, dayId string, req dto.UpdateItineraryDayRequest) (dto.ItineraryDayResponse, error) {
	day, err := s.dayRepo.GetByID(ctx, dayId)
	if err != nil {
		return dto.ItineraryDayResponse{}, serviceshared.ErrDayNotFound
	}

	member, err := s.memberRepo.GetActiveByTripAndUser(ctx, day.TripId, userId)
	if err != nil || member.Id == "" {
		return dto.ItineraryDayResponse{}, serviceshared.ErrNotMember
	}
	if !serviceshared.CanEditTrip(member.Role) {
		return dto.ItineraryDayResponse{}, serviceshared.ErrAccessDenied
	}

	if req.Title != "" {
		day.Title = new(strings.TrimSpace(req.Title))
	}
	if req.DayNumber > 0 {
		day.DayNumber = req.DayNumber
	}
	if req.Date != "" {
		parsed, err := serviceshared.ParseDate(req.Date)
		if err != nil {
			return dto.ItineraryDayResponse{}, err
		}
		day.Date = &parsed
	}

	day.UpdatedBy = userId
	day.UpdatedAt = new(time.Now())

	if err := s.dayRepo.Update(ctx, day); err != nil {
		return dto.ItineraryDayResponse{}, err
	}
	if err := s.syncTripDateRange(ctx, day.TripId, userId); err != nil {
		return dto.ItineraryDayResponse{}, err
	}

	return serviceshared.ItineraryDayToResponse(day, nil), nil
}

func (s *ItineraryDayService) DeleteDay(ctx context.Context, userId, dayId string) error {
	day, err := s.dayRepo.GetByID(ctx, dayId)
	if err != nil {
		return serviceshared.ErrDayNotFound
	}

	member, err := s.memberRepo.GetActiveByTripAndUser(ctx, day.TripId, userId)
	if err != nil || member.Id == "" {
		return serviceshared.ErrNotMember
	}
	if !serviceshared.CanEditTrip(member.Role) {
		return serviceshared.ErrAccessDenied
	}

	if err := s.dayRepo.SoftDelete(ctx, day.Id, userId); err != nil {
		return err
	}

	return s.syncTripDateRange(ctx, day.TripId, userId)
}

func (s *ItineraryDayService) syncTripDateRange(ctx context.Context, tripId, userId string) error {
	if s.tripRepo == nil {
		return nil
	}

	trip, err := s.tripRepo.GetByID(ctx, tripId)
	if err != nil {
		return serviceshared.ErrTripNotFound
	}

	days, err := s.dayRepo.ListByTrip(ctx, tripId)
	if err != nil {
		return err
	}

	var startDate *time.Time
	var endDate *time.Time
	for _, day := range days {
		if day.Date == nil {
			continue
		}
		date := *day.Date
		if startDate == nil || date.Before(*startDate) {
			startDate = &date
		}
		if endDate == nil || date.After(*endDate) {
			endDate = &date
		}
	}

	if sameDate(trip.StartDate, startDate) && sameDate(trip.EndDate, endDate) {
		return nil
	}

	now := time.Now()
	trip.StartDate = startDate
	trip.EndDate = endDate
	trip.UpdatedBy = userId
	trip.UpdatedAt = &now

	return s.tripRepo.Update(ctx, trip)
}

func sameDate(left, right *time.Time) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.Format("2006-01-02") == right.Format("2006-01-02")
}
