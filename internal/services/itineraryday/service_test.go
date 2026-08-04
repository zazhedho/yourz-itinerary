package serviceitineraryday

import (
	"context"
	"errors"
	"strings"
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

func TestItineraryDayServiceCreateValidationAndRepositoryErrors(t *testing.T) {
	ctx := context.Background()
	member := domaintripmember.TripMember{Id: "member-1", TripId: "trip-1", UserId: "user-1", Role: serviceshared.TripRoleEditor}
	date := mustParseDate(t, "2026-06-26")
	storeErr := errors.New("store failed")
	listErr := errors.New("list failed")
	updateErr := errors.New("trip update failed")

	tests := []struct {
		name      string
		member    domaintripmember.TripMember
		memberErr error
		day       stubItineraryDayRepo
		trip      *stubTripRepo
		req       dto.CreateItineraryDayRequest
		wantErr   error
	}{
		{name: "not member", wantErr: serviceshared.ErrNotMember},
		{name: "member lookup error", memberErr: errors.New("member lookup failed"), wantErr: serviceshared.ErrNotMember},
		{name: "access denied", member: domaintripmember.TripMember{Id: "member-1", Role: serviceshared.TripRoleViewer}, wantErr: serviceshared.ErrAccessDenied},
		{name: "invalid date", member: member, req: dto.CreateItineraryDayRequest{DayNumber: 1, Date: "bad"}, wantErr: serviceshared.ErrInvalidDate},
		{name: "store error", member: member, day: stubItineraryDayRepo{storeErr: storeErr}, req: dto.CreateItineraryDayRequest{DayNumber: 1, Title: "  Beach  "}, wantErr: storeErr},
		{name: "trip lookup error", member: member, day: stubItineraryDayRepo{days: []domainitineraryday.ItineraryDay{{Date: &date}}}, trip: &stubTripRepo{getErr: errors.New("trip missing")}, req: dto.CreateItineraryDayRequest{DayNumber: 1, Date: "2026-06-26"}, wantErr: serviceshared.ErrTripNotFound},
		{name: "day list error", member: member, day: stubItineraryDayRepo{listErr: listErr}, trip: &stubTripRepo{trip: domaintrip.Trip{Id: "trip-1"}}, req: dto.CreateItineraryDayRequest{DayNumber: 1, Date: "2026-06-26"}, wantErr: listErr},
		{name: "trip update error", member: member, day: stubItineraryDayRepo{days: []domainitineraryday.ItineraryDay{{Date: &date}}}, trip: &stubTripRepo{trip: domaintrip.Trip{Id: "trip-1"}, updateErr: updateErr}, req: dto.CreateItineraryDayRequest{DayNumber: 1, Date: "2026-06-26"}, wantErr: updateErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memberRepo := &stubTripMemberRepo{member: tt.member, activeErr: tt.memberErr}
			tripRepo := tt.trip
			if tripRepo == nil {
				tripRepo = &stubTripRepo{}
			}
			_, err := NewItineraryDayService(memberRepo, &tt.day, tripRepo).CreateDay(ctx, "user-1", "trip-1", tt.req)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestItineraryDayServiceUpdateAndDeleteErrors(t *testing.T) {
	ctx := context.Background()
	member := domaintripmember.TripMember{Id: "member-1", TripId: "trip-1", Role: serviceshared.TripRoleEditor}
	day := domainitineraryday.ItineraryDay{Id: "day-1", TripId: "trip-1"}
	updateErr := errors.New("day update failed")
	deleteErr := errors.New("day delete failed")
	listErr := errors.New("list failed")

	tests := []struct {
		name      string
		member    domaintripmember.TripMember
		memberErr error
		dayRepo   stubItineraryDayRepo
		tripRepo  *stubTripRepo
		wantErr   error
	}{
		{name: "update day not found", dayRepo: stubItineraryDayRepo{getErr: errors.New("missing")}, wantErr: serviceshared.ErrDayNotFound},
		{name: "update not member", dayRepo: stubItineraryDayRepo{day: day}, wantErr: serviceshared.ErrNotMember},
		{name: "update member lookup error", memberErr: errors.New("member lookup failed"), dayRepo: stubItineraryDayRepo{day: day}, wantErr: serviceshared.ErrNotMember},
		{name: "update access denied", member: domaintripmember.TripMember{Id: "member-1", Role: serviceshared.TripRoleViewer}, dayRepo: stubItineraryDayRepo{day: day}, wantErr: serviceshared.ErrAccessDenied},
		{name: "update invalid date", member: member, dayRepo: stubItineraryDayRepo{day: day}, wantErr: serviceshared.ErrInvalidDate},
		{name: "update repository error", member: member, dayRepo: stubItineraryDayRepo{day: day, updateErr: updateErr}, wantErr: updateErr},
		{name: "update sync trip error", member: member, dayRepo: stubItineraryDayRepo{day: day}, tripRepo: &stubTripRepo{getErr: errors.New("missing trip")}, wantErr: serviceshared.ErrTripNotFound},
		{name: "update sync list error", member: member, dayRepo: stubItineraryDayRepo{day: day, listErr: listErr}, tripRepo: &stubTripRepo{trip: domaintrip.Trip{Id: "trip-1"}}, wantErr: listErr},
		{name: "delete day not found", dayRepo: stubItineraryDayRepo{getErr: errors.New("missing")}, wantErr: serviceshared.ErrDayNotFound},
		{name: "delete not member", dayRepo: stubItineraryDayRepo{day: day}, wantErr: serviceshared.ErrNotMember},
		{name: "delete access denied", member: domaintripmember.TripMember{Id: "member-1", Role: serviceshared.TripRoleViewer}, dayRepo: stubItineraryDayRepo{day: day}, wantErr: serviceshared.ErrAccessDenied},
		{name: "delete repository error", member: member, dayRepo: stubItineraryDayRepo{day: day, softDeleteErr: deleteErr}, wantErr: deleteErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memberRepo := &stubTripMemberRepo{member: tt.member, activeErr: tt.memberErr}
			tripRepo := tt.tripRepo
			if tripRepo == nil {
				tripRepo = &stubTripRepo{trip: domaintrip.Trip{Id: "trip-1"}}
			}
			svc := NewItineraryDayService(memberRepo, &tt.dayRepo, tripRepo)
			var err error
			if strings.HasPrefix(tt.name, "update") {
				req := dto.UpdateItineraryDayRequest{}
				if tt.name == "update invalid date" {
					req.Date = "bad"
				}
				_, err = svc.UpdateDay(ctx, "user-1", "day-1", req)
			} else {
				err = svc.DeleteDay(ctx, "user-1", "day-1")
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestItineraryDayServiceNoTripRepositoryAndSameDate(t *testing.T) {
	ctx := context.Background()
	date := mustParseDate(t, "2026-06-26")
	memberRepo := &stubTripMemberRepo{member: domaintripmember.TripMember{Id: "member-1", Role: serviceshared.TripRoleEditor}}
	dayRepo := &stubItineraryDayRepo{}
	if _, err := NewItineraryDayService(memberRepo, dayRepo, nil).CreateDay(ctx, "user-1", "trip-1", dto.CreateItineraryDayRequest{DayNumber: 1, Title: "  Beach  "}); err != nil {
		t.Fatalf("CreateDay without trip repo: %v", err)
	}
	if dayRepo.day.Title == nil || *dayRepo.day.Title != "Beach" {
		t.Fatalf("trimmed title = %v", dayRepo.day.Title)
	}

	dayRepo = &stubItineraryDayRepo{days: []domainitineraryday.ItineraryDay{{Date: &date}}}
	tripRepo := &stubTripRepo{trip: domaintrip.Trip{Id: "trip-1", StartDate: &date, EndDate: &date}}
	if _, err := NewItineraryDayService(memberRepo, dayRepo, tripRepo).CreateDay(ctx, "user-1", "trip-1", dto.CreateItineraryDayRequest{DayNumber: 1, Date: "2026-06-26"}); err != nil {
		t.Fatalf("CreateDay with unchanged range: %v", err)
	}
	if tripRepo.updated.Id != "" {
		t.Fatal("trip should not be updated when dates are unchanged")
	}
}

func TestSameDate(t *testing.T) {
	date := mustParseDate(t, "2026-06-26")
	different := mustParseDate(t, "2026-06-27")
	for _, tt := range []struct {
		name        string
		left, right *time.Time
		want        bool
	}{
		{name: "both nil", want: true},
		{name: "left nil", right: &date, want: false},
		{name: "right nil", left: &date, want: false},
		{name: "same day different time", left: timePtr(date.Add(8 * time.Hour)), right: &date, want: true},
		{name: "different day", left: &date, right: &different, want: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := sameDate(tt.left, tt.right); got != tt.want {
				t.Fatalf("sameDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func timePtr(value time.Time) *time.Time { return &value }

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
	day           domainitineraryday.ItineraryDay
	days          []domainitineraryday.ItineraryDay
	getErr        error
	storeErr      error
	updateErr     error
	softDeleteErr error
	listErr       error
}

func (r *stubItineraryDayRepo) Store(_ context.Context, day domainitineraryday.ItineraryDay) error {
	if r.storeErr != nil {
		return r.storeErr
	}
	r.day = day
	r.days = append(r.days, day)
	return nil
}

func (r *stubItineraryDayRepo) GetByID(_ context.Context, _ string) (domainitineraryday.ItineraryDay, error) {
	if r.getErr != nil {
		return domainitineraryday.ItineraryDay{}, r.getErr
	}
	return r.day, nil
}

func (r *stubItineraryDayRepo) GetAll(_ context.Context, _ filter.BaseParams) ([]domainitineraryday.ItineraryDay, int64, error) {
	return nil, 0, nil
}

func (r *stubItineraryDayRepo) Update(_ context.Context, day domainitineraryday.ItineraryDay) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	r.day = day
	return nil
}

func (r *stubItineraryDayRepo) Delete(_ context.Context, _ string) error { return nil }

func (r *stubItineraryDayRepo) SoftDelete(_ context.Context, _ string, _ string) error {
	return r.softDeleteErr
}

func (r *stubItineraryDayRepo) ListByTrip(_ context.Context, _ string) ([]domainitineraryday.ItineraryDay, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	return r.days, nil
}

type stubTripRepo struct {
	trip      domaintrip.Trip
	updated   domaintrip.Trip
	getErr    error
	updateErr error
}

func (r *stubTripRepo) Store(_ context.Context, trip domaintrip.Trip) error {
	r.trip = trip
	return nil
}

func (r *stubTripRepo) GetByID(_ context.Context, _ string) (domaintrip.Trip, error) {
	if r.getErr != nil {
		return domaintrip.Trip{}, r.getErr
	}
	return r.trip, nil
}

func (r *stubTripRepo) GetAll(_ context.Context, _ filter.BaseParams) ([]domaintrip.Trip, int64, error) {
	return nil, 0, nil
}

func (r *stubTripRepo) Update(_ context.Context, trip domaintrip.Trip) error {
	if r.updateErr != nil {
		return r.updateErr
	}
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
	member    domaintripmember.TripMember
	activeErr error
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
	if r.activeErr != nil {
		return domaintripmember.TripMember{}, r.activeErr
	}
	return r.member, nil
}

func (r *stubTripMemberRepo) ListByTrip(_ context.Context, _ string) ([]domaintripmember.TripMember, error) {
	return []domaintripmember.TripMember{r.member}, nil
}
