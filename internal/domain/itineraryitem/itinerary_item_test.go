package domainitineraryitem

import "testing"

func TestItineraryItemTableName(t *testing.T) {
	if (ItineraryItem{}).TableName() != "itinerary_items" {
		t.Fatal("unexpected table name")
	}
}
