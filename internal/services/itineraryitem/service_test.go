package serviceitineraryitem

import "testing"

func TestNewItineraryItemService(t *testing.T) {
	svc := NewItineraryItemService(nil, nil, nil)
	if svc == nil {
		t.Fatal("NewItineraryItemService returned nil")
	}
}

func TestItemErrorsDistinct(t *testing.T) {
	errors := [...]error{
		ErrItemNotFound, ErrDayNotFound, ErrTripNotFound, ErrInvalidTime,
		ErrInvalidCoordinates, ErrInvalidLatitude, ErrInvalidLongitude,
		ErrReorderDifferentDay, ErrReorderEmpty, ErrReorderItemsNotFound,
	}
	for i, e := range errors {
		for j, o := range errors {
			if i != j && e == o { //nolint:errorlint
				t.Errorf("errors[%d] and errors[%d] are same pointer", i, j)
			}
		}
	}
}

func TestValidateCoordinates(t *testing.T) {
	if err := validateCoordinates(nil, nil); err != nil {
		t.Errorf("nil coords should be valid: %v", err)
	}

	lat := 45.0
	if err := validateCoordinates(&lat, nil); err == nil {
		t.Error("partial coords should error")
	}

	lng := 90.0
	if err := validateCoordinates(nil, &lng); err == nil {
		t.Error("partial coords should error")
	}

	if err := validateCoordinates(&lat, &lng); err != nil {
		t.Errorf("valid coords should pass: %v", err)
	}

	badLat := 91.0
	if err := validateCoordinates(&badLat, &lng); err == nil {
		t.Error("out-of-range lat should error")
	}

	badLng := -181.0
	if err := validateCoordinates(&lat, &badLng); err == nil {
		t.Error("out-of-range lng should error")
	}
}

func TestParseClockTime(t *testing.T) {
	_, err := parseClockTime("14:30")
	if err != nil {
		t.Errorf("valid time should parse: %v", err)
	}
	_, err = parseClockTime("25:90")
	if err == nil {
		t.Error("invalid time should error")
	}
}
