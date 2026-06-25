package serviceitineraryday

import (
	"context"
	"errors"
	"strings"
	"time"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	domainitineraryitem "yourz-itinerary/internal/domain/itineraryitem"
	"yourz-itinerary/internal/dto"
	interfaceitineraryday "yourz-itinerary/internal/interfaces/itineraryday"
	interfacetripmember "yourz-itinerary/internal/interfaces/tripmember"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/utils"
)

var (
	ErrDayNotFound  = errors.New("itinerary day not found")
	ErrTripNotFound = errors.New("trip not found")
	ErrInvalidDate  = errors.New("invalid date; must use YYYY-MM-DD")
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

func dayToResponse(d domainitineraryday.ItineraryDay, items []domainitineraryitem.ItineraryItem) dto.ItineraryDayResponse {
	dr := dto.ItineraryDayResponse{
		Id:        d.Id,
		TripId:    d.TripId,
		DayNumber: d.DayNumber,
		Title:     d.Title,
		CreatedBy: d.CreatedBy,
		UpdatedBy: d.UpdatedBy,
		CreatedAt: d.CreatedAt.Format(time.RFC3339),
	}

	if d.Date != nil {
		dr.Date = new(d.Date.Format("2006-01-02"))
	}
	if d.UpdatedAt != nil {
		dr.UpdatedAt = new(d.UpdatedAt.Format(time.RFC3339))
	}
	if d.DeletedBy != nil {
		dr.DeletedBy = d.DeletedBy
	}
	if d.DeletedAt.Valid {
		dr.DeletedAt = new(d.DeletedAt.Time.Format(time.RFC3339))
	}

	itemResponses := make([]dto.ItineraryItemResponse, 0, len(items))
	for _, item := range items {
		itemResponses = append(itemResponses, itemToResponse(item))
	}
	dr.Items = itemResponses

	return dr
}

func itemToResponse(item domainitineraryitem.ItineraryItem) dto.ItineraryItemResponse {
	ir := dto.ItineraryItemResponse{
		Id:           item.Id,
		DayId:        item.DayId,
		Title:        item.Title,
		Description:  item.Description,
		LocationName: item.LocationName,
		Latitude:     item.Latitude,
		Longitude:    item.Longitude,
		StartTime:    item.StartTime,
		EndTime:      item.EndTime,
		CostEstimate: item.CostEstimate,
		SortOrder:    item.SortOrder,
		CreatedBy:    item.CreatedBy,
		UpdatedBy:    item.UpdatedBy,
		CreatedAt:    item.CreatedAt.Format(time.RFC3339),
	}

	if item.UpdatedAt != nil {
		ir.UpdatedAt = new(item.UpdatedAt.Format(time.RFC3339))
	}
	if item.DeletedBy != nil {
		ir.DeletedBy = item.DeletedBy
	}
	if item.DeletedAt.Valid {
		ir.DeletedAt = new(item.DeletedAt.Time.Format(time.RFC3339))
	}

	return ir
}
