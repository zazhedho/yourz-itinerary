package serviceitineraryday

import (
	"context"
	"testing"
	"time"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	domaintrip "yourz-itinerary/internal/domain/trip"
	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	"yourz-itinerary/internal/dto"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/pkg/filter"
)

func TestNewItineraryDayService(t *testing.T) {
	svc := NewItineraryDayService(nil, nil, nil)
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

func TestCreateDayExpandsTripDateRange(t *testing.T) {
	ctx := context.Background()
	start := mustParseDate(t, "2026-06-26")
	end := mustParseDate(t, "2026-06-27")
	memberRepo := &stubTripMemberRepo{member: domaintripmember.TripMember{Id: "member-1", TripId: "trip-1", UserId: "user-1", Role: serviceshared.TripRoleEditor}}
	dayRepo := &stubItineraryDayRepo{days: []domainitineraryday.ItineraryDay{
		{Id: "day-1", TripId: "trip-1", Date: &start, DayNumber: 1},
		{Id: "day-2", TripId: "trip-1", Date: &end, DayNumber: 2},
	}}
	tripRepo := &stubTripRepo{trip: domaintrip.Trip{Id: "trip-1", StartDate: &start, EndDate: &end}}
	svc := NewItineraryDayService(memberRepo, dayRepo, tripRepo)

	_, err := svc.CreateDay(ctx, "user-1", "trip-1", dto.CreateItineraryDayRequest{
		DayNumber: 3,
		Date:      "2026-06-28",
	})

	if err != nil {
		t.Fatalf("CreateDay returned error: %v", err)
	}
	if tripRepo.updated.StartDate == nil || formatDate(*tripRepo.updated.StartDate) != "2026-06-26" {
		t.Fatalf("start date = %v, want 2026-06-26", tripRepo.updated.StartDate)
	}
	if tripRepo.updated.EndDate == nil || formatDate(*tripRepo.updated.EndDate) != "2026-06-28" {
		t.Fatalf("end date = %v, want 2026-06-28", tripRepo.updated.EndDate)
	}
}

func TestDeleteDayShrinksTripDateRange(t *testing.T) {
	ctx := context.Background()
	start := mustParseDate(t, "2026-06-26")
	end := mustParseDate(t, "2026-06-28")
	remainingEnd := mustParseDate(t, "2026-06-27")
	memberRepo := &stubTripMemberRepo{member: domaintripmember.TripMember{Id: "member-1", TripId: "trip-1", UserId: "user-1", Role: serviceshared.TripRoleEditor}}
	dayRepo := &stubItineraryDayRepo{
		day: domainitineraryday.ItineraryDay{Id: "day-3", TripId: "trip-1"},
		days: []domainitineraryday.ItineraryDay{
			{Id: "day-1", TripId: "trip-1", Date: &start, DayNumber: 1},
			{Id: "day-2", TripId: "trip-1", Date: &remainingEnd, DayNumber: 2},
		},
	}
	tripRepo := &stubTripRepo{trip: domaintrip.Trip{Id: "trip-1", StartDate: &start, EndDate: &end}}
	svc := NewItineraryDayService(memberRepo, dayRepo, tripRepo)

	err := svc.DeleteDay(ctx, "user-1", "day-3")

	if err != nil {
		t.Fatalf("DeleteDay returned error: %v", err)
	}
	if tripRepo.updated.EndDate == nil || formatDate(*tripRepo.updated.EndDate) != "2026-06-27" {
		t.Fatalf("end date = %v, want 2026-06-27", tripRepo.updated.EndDate)
	}
}

func mustParseDate(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := serviceshared.ParseDate(value)
	if err != nil {
		t.Fatalf("ParseDate(%q): %v", value, err)
	}
	return parsed
}

func formatDate(value time.Time) string {
	return value.Format("2006-01-02")
}

type stubItineraryDayRepo struct {
	day  domainitineraryday.ItineraryDay
	days []domainitineraryday.ItineraryDay
}

func (r *stubItineraryDayRepo) Store(_ context.Context, day domainitineraryday.ItineraryDay) error {
	r.day = day
	r.days = append(r.days, day)
	return nil
}

func (r *stubItineraryDayRepo) GetByID(_ context.Context, _ string) (domainitineraryday.ItineraryDay, error) {
	return r.day, nil
}

func (r *stubItineraryDayRepo) GetAll(_ context.Context, _ filter.BaseParams) ([]domainitineraryday.ItineraryDay, int64, error) {
	return nil, 0, nil
}

func (r *stubItineraryDayRepo) Update(_ context.Context, day domainitineraryday.ItineraryDay) error {
	r.day = day
	return nil
}

func (r *stubItineraryDayRepo) Delete(_ context.Context, _ string) error { return nil }

func (r *stubItineraryDayRepo) SoftDelete(_ context.Context, _ string, _ string) error { return nil }

func (r *stubItineraryDayRepo) ListByTrip(_ context.Context, _ string) ([]domainitineraryday.ItineraryDay, error) {
	return r.days, nil
}

type stubTripRepo struct {
	trip    domaintrip.Trip
	updated domaintrip.Trip
}

func (r *stubTripRepo) Store(_ context.Context, trip domaintrip.Trip) error {
	r.trip = trip
	return nil
}

func (r *stubTripRepo) GetByID(_ context.Context, _ string) (domaintrip.Trip, error) {
	return r.trip, nil
}

func (r *stubTripRepo) GetAll(_ context.Context, _ filter.BaseParams) ([]domaintrip.Trip, int64, error) {
	return nil, 0, nil
}

func (r *stubTripRepo) Update(_ context.Context, trip domaintrip.Trip) error {
	r.updated = trip
	r.trip = trip
	return nil
}

func (r *stubTripRepo) Delete(_ context.Context, _ string) error { return nil }

func (r *stubTripRepo) SoftDelete(_ context.Context, _ string, _ string) error { return nil }

func (r *stubTripRepo) CreateTrip(_ context.Context, trip domaintrip.Trip, _ domaintripmember.TripMember, _ ...domainitineraryday.ItineraryDay) (domaintrip.Trip, error) {
	return trip, nil
}

func (r *stubTripRepo) ListByMember(_ context.Context, _ string) ([]domaintrip.Trip, int64, error) {
	return nil, 0, nil
}

type stubTripMemberRepo struct {
	member domaintripmember.TripMember
}

func (r *stubTripMemberRepo) Store(_ context.Context, member domaintripmember.TripMember) error {
	r.member = member
	return nil
}

func (r *stubTripMemberRepo) GetByID(_ context.Context, _ string) (domaintripmember.TripMember, error) {
	return r.member, nil
}

func (r *stubTripMemberRepo) GetAll(_ context.Context, _ filter.BaseParams) ([]domaintripmember.TripMember, int64, error) {
	return nil, 0, nil
}

func (r *stubTripMemberRepo) Update(_ context.Context, member domaintripmember.TripMember) error {
	r.member = member
	return nil
}

func (r *stubTripMemberRepo) Delete(_ context.Context, _ string) error { return nil }

func (r *stubTripMemberRepo) SoftDelete(_ context.Context, _ string, _ string) error { return nil }

func (r *stubTripMemberRepo) GetByTripAndUser(_ context.Context, _ string, _ string) (domaintripmember.TripMember, error) {
	return r.member, nil
}

func (r *stubTripMemberRepo) GetActiveByTripAndUser(_ context.Context, _ string, _ string) (domaintripmember.TripMember, error) {
	return r.member, nil
}

func (r *stubTripMemberRepo) ListByTrip(_ context.Context, _ string) ([]domaintripmember.TripMember, error) {
	return []domaintripmember.TripMember{r.member}, nil
}
