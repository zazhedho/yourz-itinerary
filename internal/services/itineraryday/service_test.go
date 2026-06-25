package serviceitineraryday

import (
	"testing"

	serviceshared "yourz-itinerary/internal/services/shared"
)

func TestNewItineraryDayService(t *testing.T) {
	svc := NewItineraryDayService(nil, nil)
	if svc == nil {
		t.Fatal("NewItineraryDayService returned nil")
	}
}

func TestDayErrorsDistinct(t *testing.T) {
	errors := [...]error{serviceshared.ErrDayNotFound, serviceshared.ErrTripNotFound, serviceshared.ErrInvalidDate}
	for i, e := range errors {
		for j, o := range errors {
			if i != j && e == o { //nolint:errorlint
				t.Errorf("errors[%d] and errors[%d] are same pointer", i, j)
			}
		}
	}
}

func TestParseDate(t *testing.T) {
	_, err := serviceshared.ParseDate("2026-06-25")
	if err != nil {
		t.Errorf("valid date should parse: %v", err)
	}
}
