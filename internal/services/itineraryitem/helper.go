package serviceitineraryitem

import (
	"errors"
	"strings"
	"time"
	domainitineraryitem "yourz-itinerary/internal/domain/itineraryitem"
	"yourz-itinerary/internal/dto"
)

var (
	ErrItemNotFound         = errors.New("itinerary item not found")
	ErrDayNotFound          = errors.New("itinerary day not found")
	ErrTripNotFound         = errors.New("trip not found")
	ErrInvalidTime          = errors.New("invalid time; must use HH:MM")
	ErrInvalidCoordinates   = errors.New("latitude and longitude must both be provided when specifying coordinates")
	ErrInvalidLatitude      = errors.New("latitude must be between -90 and 90")
	ErrInvalidLongitude     = errors.New("longitude must be between -180 and 180")
	ErrReorderDifferentDay  = errors.New("all items must belong to the same day")
	ErrReorderEmpty         = errors.New("item_ids must not be empty")
	ErrReorderItemsNotFound = errors.New("one or more items not found")
)

func validateCoordinates(lat, lng *float64) error {
	if lat != nil && lng == nil {
		return ErrInvalidCoordinates
	}
	if lat == nil && lng != nil {
		return ErrInvalidCoordinates
	}
	if lat != nil {
		if *lat < -90 || *lat > 90 {
			return ErrInvalidLatitude
		}
	}
	if lng != nil {
		if *lng < -180 || *lng > 180 {
			return ErrInvalidLongitude
		}
	}
	return nil
}

func parseClockTime(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if _, err := time.Parse("15:04", trimmed); err != nil {
		return "", ErrInvalidTime
	}
	return trimmed, nil
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
