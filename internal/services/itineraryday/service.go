package serviceitineraryday

import (
	"context"
	"strings"
	"time"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	"yourz-itinerary/internal/dto"
	interfaceitineraryday "yourz-itinerary/internal/interfaces/itineraryday"
	interfacetripmember "yourz-itinerary/internal/interfaces/tripmember"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/utils"
)

type ItineraryDayService struct {
	memberRepo interfacetripmember.RepoTripMemberInterface
	dayRepo    interfaceitineraryday.RepoItineraryDayInterface
}

func NewItineraryDayService(memberRepo interfacetripmember.RepoTripMemberInterface, dayRepo interfaceitineraryday.RepoItineraryDayInterface) *ItineraryDayService {
	return &ItineraryDayService{memberRepo: memberRepo, dayRepo: dayRepo}
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
		parsed, err := parseDate(req.Date)
		if err != nil {
			return dto.ItineraryDayResponse{}, err
		}
		day.Date = &parsed
	}

	if err := s.dayRepo.Store(ctx, day); err != nil {
		return dto.ItineraryDayResponse{}, err
	}

	return dayToResponse(day, nil), nil
}

func (s *ItineraryDayService) UpdateDay(ctx context.Context, userId, dayId string, req dto.UpdateItineraryDayRequest) (dto.ItineraryDayResponse, error) {
	day, err := s.dayRepo.GetByID(ctx, dayId)
	if err != nil {
		return dto.ItineraryDayResponse{}, ErrDayNotFound
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
		parsed, err := parseDate(req.Date)
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

	return dayToResponse(day, nil), nil
}

func (s *ItineraryDayService) DeleteDay(ctx context.Context, userId, dayId string) error {
	day, err := s.dayRepo.GetByID(ctx, dayId)
	if err != nil {
		return ErrDayNotFound
	}

	member, err := s.memberRepo.GetActiveByTripAndUser(ctx, day.TripId, userId)
	if err != nil || member.Id == "" {
		return serviceshared.ErrNotMember
	}
	if !serviceshared.CanEditTrip(member.Role) {
		return serviceshared.ErrAccessDenied
	}

	return s.dayRepo.SoftDelete(ctx, day.Id, userId)
}

func parseDate(value string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, ErrInvalidDate
	}
	return parsed, nil
}
