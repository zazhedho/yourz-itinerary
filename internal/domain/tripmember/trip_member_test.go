package domaintripmember

import "testing"

func TestTripMemberTableName(t *testing.T) {
	if (TripMember{}).TableName() != "trip_members" {
		t.Fatal("unexpected table name")
	}
}
