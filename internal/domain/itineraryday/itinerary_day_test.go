package domainitineraryday

import "testing"

func TestItineraryDayTableName(t *testing.T) {
	if (ItineraryDay{}).TableName() != "itinerary_days" {
		t.Fatal("unexpected table name")
	}
}
