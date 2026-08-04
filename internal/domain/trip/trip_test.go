package domaintrip

import "testing"

func TestTripTableName(t *testing.T) {
	if (Trip{}).TableName() != "trips" {
		t.Fatal("unexpected table name")
	}
}
