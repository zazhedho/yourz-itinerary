package serviceitineraryitem

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrItemNotFound         = errors.New("itinerary item not found")
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
