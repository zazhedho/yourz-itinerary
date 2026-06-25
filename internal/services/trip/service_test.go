package servicetrip

import (
	"testing"

	serviceshared "yourz-itinerary/internal/services/shared"
)

func TestNewTripService(t *testing.T) {
	svc := NewTripService(nil, nil, nil, nil)
	if svc == nil {
		t.Fatal("NewTripService returned nil")
	}
}

func TestTripErrorsDistinct(t *testing.T) {
	errors := [...]error{serviceshared.ErrTripNotFound, ErrInvalidTimezone, ErrInvalidCurrency, serviceshared.ErrInvalidDate, ErrInvalidDateRange}
	for i, e := range errors {
		for j, o := range errors {
			if i != j && e == o { //nolint:errorlint
				t.Errorf("errors[%d] and errors[%d] are same pointer", i, j)
			}
		}
	}
}

func TestIsValidCurrencyCode(t *testing.T) {
	if !isValidCurrencyCode("IDR") {
		t.Error("IDR should be valid")
	}
	if !isValidCurrencyCode("USD") {
		t.Error("USD should be valid")
	}
	if isValidCurrencyCode("idr") {
		t.Error("lowercase should be invalid")
	}
	if isValidCurrencyCode("ID") {
		t.Error("2-char should be invalid")
	}
	if isValidCurrencyCode("IDR!") {
		t.Error("special chars should be invalid")
	}
}

func TestParseDate(t *testing.T) {
	_, err := serviceshared.ParseDate("2026-06-25")
	if err != nil {
		t.Errorf("valid date should parse: %v", err)
	}
	_, err = serviceshared.ParseDate("25-06-2026")
	if err == nil {
		t.Error("invalid date format should error")
	}
}
