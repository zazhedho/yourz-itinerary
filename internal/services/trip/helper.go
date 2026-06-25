package servicetrip

import (
	"errors"
	"strings"
	"time"
	"unicode"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	domainitineraryitem "yourz-itinerary/internal/domain/itineraryitem"
	domaintrip "yourz-itinerary/internal/domain/trip"
	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	"yourz-itinerary/internal/dto"
)

var (
	ErrTripNotFound     = errors.New("trip not found")
	ErrInvalidTimezone  = errors.New("invalid timezone")
	ErrInvalidCurrency  = errors.New("invalid currency code; must be a 3-letter uppercase ISO 4217 code")
	ErrInvalidDate      = errors.New("invalid date; must use YYYY-MM-DD")
	ErrInvalidDateRange = errors.New("end_date must be on or after start_date")
)

func isValidCurrencyCode(value string) bool {
	if len(value) != 3 {
		return false
	}
	for _, r := range value {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func parseDate(value string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, ErrInvalidDate
	}
	return parsed, nil
}

func tripToList(trip domaintrip.Trip, memberCount, dayCount int) dto.TripListResponse {
	tr := dto.TripListResponse{
		Id:           trip.Id,
		OwnerId:      trip.OwnerId,
		Title:        trip.Title,
		Destination:  trip.Destination,
		Timezone:     trip.Timezone,
		CurrencyCode: trip.CurrencyCode,
		Status:       trip.Status,
		MemberCount:  memberCount,
		DayCount:     dayCount,
		CreatedBy:    trip.CreatedBy,
		UpdatedBy:    trip.UpdatedBy,
		CreatedAt:    trip.CreatedAt.Format(time.RFC3339),
	}

	if trip.StartDate != nil {
		tr.StartDate = new(trip.StartDate.Format("2006-01-02"))
	}
	if trip.EndDate != nil {
		tr.EndDate = new(trip.EndDate.Format("2006-01-02"))
	}
	if trip.UpdatedAt != nil {
		tr.UpdatedAt = new(trip.UpdatedAt.Format(time.RFC3339))
	}
	if trip.DeletedBy != nil {
		tr.DeletedBy = trip.DeletedBy
	}
	if trip.DeletedAt.Valid {
		tr.DeletedAt = new(trip.DeletedAt.Time.Format(time.RFC3339))
	}

	return tr
}

func tripToDetail(trip domaintrip.Trip, members []domaintripmember.TripMember, days []domainitineraryday.ItineraryDay, itemsByDay map[string][]domainitineraryitem.ItineraryItem) dto.TripDetailResponse {
	tr := dto.TripDetailResponse{
		Id:           trip.Id,
		OwnerId:      trip.OwnerId,
		Title:        trip.Title,
		Destination:  trip.Destination,
		Timezone:     trip.Timezone,
		CurrencyCode: trip.CurrencyCode,
		Status:       trip.Status,
		CreatedBy:    trip.CreatedBy,
		UpdatedBy:    trip.UpdatedBy,
		CreatedAt:    trip.CreatedAt.Format(time.RFC3339),
	}

	if trip.StartDate != nil {
		tr.StartDate = new(trip.StartDate.Format("2006-01-02"))
	}
	if trip.EndDate != nil {
		tr.EndDate = new(trip.EndDate.Format("2006-01-02"))
	}
	if trip.UpdatedAt != nil {
		tr.UpdatedAt = new(trip.UpdatedAt.Format(time.RFC3339))
	}
	if trip.DeletedBy != nil {
		tr.DeletedBy = trip.DeletedBy
	}
	if trip.DeletedAt.Valid {
		tr.DeletedAt = new(trip.DeletedAt.Time.Format(time.RFC3339))
	}

	memberResponses := make([]dto.TripMemberResponse, 0, len(members))
	for _, m := range members {
		memberResponses = append(memberResponses, memberToResponse(m))
	}
	tr.Members = memberResponses

	dayResponses := make([]dto.ItineraryDayResponse, 0, len(days))
	for _, d := range days {
		dayResponses = append(dayResponses, dayToResponse(d, itemsByDay[d.Id]))
	}
	tr.Days = dayResponses

	return tr
}

func memberToResponse(m domaintripmember.TripMember) dto.TripMemberResponse {
	mr := dto.TripMemberResponse{
		Id:        m.Id,
		TripId:    m.TripId,
		UserId:    m.UserId,
		Role:      m.Role,
		CreatedBy: m.CreatedBy,
		UpdatedBy: m.UpdatedBy,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}

	if m.UpdatedAt != nil {
		mr.UpdatedAt = new(m.UpdatedAt.Format(time.RFC3339))
	}
	if m.DeletedBy != nil {
		mr.DeletedBy = m.DeletedBy
	}
	if m.DeletedAt.Valid {
		mr.DeletedAt = new(m.DeletedAt.Time.Format(time.RFC3339))
	}

	return mr
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
