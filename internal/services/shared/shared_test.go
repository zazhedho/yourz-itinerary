package serviceshared

import (
	"errors"
	"testing"
	"time"

	domainday "yourz-itinerary/internal/domain/itineraryday"
	domainitem "yourz-itinerary/internal/domain/itineraryitem"
	domainmember "yourz-itinerary/internal/domain/tripmember"
	domainuser "yourz-itinerary/internal/domain/user"
)

func TestRoleAccessHelpers(t *testing.T) {
	for _, role := range []string{TripRoleOwner, TripRoleEditor, TripRoleViewer} {
		if !CanViewTrip(role) {
			t.Errorf("%s should view", role)
		}
	}
	for _, role := range []string{TripRoleOwner, TripRoleEditor} {
		if !CanEditTrip(role) {
			t.Errorf("%s should edit", role)
		}
	}
	if !CanManageTripMembers(TripRoleOwner) || CanManageTripMembers(TripRoleEditor) || CanEditTrip(TripRoleViewer) || CanViewTrip("unknown") {
		t.Fatal("role access mismatch")
	}
}

func TestParseDateAndResponseMappers(t *testing.T) {
	if got, err := ParseDate(" 2026-08-04 "); err != nil || got.Format("2006-01-02") != "2026-08-04" {
		t.Fatalf("parse date: %v %v", got, err)
	}
	if _, err := ParseDate("bad"); !errors.Is(err, ErrInvalidDate) {
		t.Fatalf("invalid date err=%v", err)
	}
	now := time.Date(2026, time.August, 4, 10, 0, 0, 0, time.UTC)
	updated := now.Add(time.Hour)
	deletedBy := "admin"
	member := domainmember.TripMember{Id: "m-1", TripId: "t-1", UserId: "u-1", Role: TripRoleEditor, CreatedAt: now, UpdatedAt: &updated, DeletedBy: &deletedBy}
	member.DeletedAt.Valid = true
	member.DeletedAt.Time = updated
	memberResponse := TripMemberToResponseWithUser(member, domainuser.Users{Name: "Alice", Email: "alice@example.com", AvatarURL: "avatar"})
	if memberResponse.UserName != "Alice" || memberResponse.UserEmail != "alice@example.com" || memberResponse.DeletedBy == nil || memberResponse.DeletedAt == nil {
		t.Fatalf("member mapping: %+v", memberResponse)
	}

	title := "Day title"
	start := "09:00"
	item := domainitem.ItineraryItem{Id: "i-1", DayId: "d-1", Title: "Place", Description: &title, StartTime: &start, CreatedAt: now, UpdatedAt: &updated, DeletedBy: &deletedBy}
	item.DeletedAt.Valid = true
	item.DeletedAt.Time = updated
	dayDate := now
	day := domainday.ItineraryDay{Id: "d-1", TripId: "t-1", DayNumber: 1, Date: &dayDate, Title: &title, CreatedAt: now, UpdatedAt: &updated, DeletedBy: &deletedBy}
	day.DeletedAt.Valid = true
	day.DeletedAt.Time = updated
	dayResponse := ItineraryDayToResponse(day, []domainitem.ItineraryItem{item})
	if dayResponse.Date == nil || len(dayResponse.Items) != 1 || dayResponse.Items[0].UpdatedAt == nil || dayResponse.Items[0].DeletedAt == nil {
		t.Fatalf("day mapping: %+v", dayResponse)
	}
	if ItineraryItemToResponse(item).Id != "i-1" {
		t.Fatal("item mapping failed")
	}
}
