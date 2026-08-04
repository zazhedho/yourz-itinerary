package servicetrip

import (
	"context"
	"errors"
	"testing"
	"time"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	domainitineraryitem "yourz-itinerary/internal/domain/itineraryitem"
	domaintrip "yourz-itinerary/internal/domain/trip"
	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	domainuser "yourz-itinerary/internal/domain/user"
	"yourz-itinerary/internal/dto"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/pkg/filter"
)

type tripServiceRepoStub struct {
	trip    domaintrip.Trip
	trips   []domaintrip.Trip
	err     error
	created bool
	updated bool
	deleted bool
}

func (s *tripServiceRepoStub) Store(context.Context, domaintrip.Trip) error { return nil }
func (s *tripServiceRepoStub) GetByID(context.Context, string) (domaintrip.Trip, error) {
	return s.trip, s.err
}
func (s *tripServiceRepoStub) GetAll(context.Context, filter.BaseParams) ([]domaintrip.Trip, int64, error) {
	return s.trips, int64(len(s.trips)), s.err
}
func (s *tripServiceRepoStub) Update(context.Context, domaintrip.Trip) error {
	s.updated = true
	return s.err
}
func (s *tripServiceRepoStub) Delete(context.Context, string) error { return nil }
func (s *tripServiceRepoStub) SoftDelete(context.Context, string, string) error {
	s.deleted = true
	return s.err
}
func (s *tripServiceRepoStub) CreateTrip(_ context.Context, trip domaintrip.Trip, _ domaintripmember.TripMember, _ ...domainitineraryday.ItineraryDay) (domaintrip.Trip, error) {
	s.created = true
	return trip, s.err
}
func (s *tripServiceRepoStub) ListByMember(context.Context, string) ([]domaintrip.Trip, int64, error) {
	return s.trips, int64(len(s.trips)), s.err
}

type tripServiceMemberRepoStub struct {
	member  domaintripmember.TripMember
	members []domaintripmember.TripMember
	err     error
}

func (s *tripServiceMemberRepoStub) Store(context.Context, domaintripmember.TripMember) error {
	return nil
}
func (s *tripServiceMemberRepoStub) GetByID(context.Context, string) (domaintripmember.TripMember, error) {
	return s.member, s.err
}
func (s *tripServiceMemberRepoStub) GetAll(context.Context, filter.BaseParams) ([]domaintripmember.TripMember, int64, error) {
	return s.members, int64(len(s.members)), s.err
}
func (s *tripServiceMemberRepoStub) Update(context.Context, domaintripmember.TripMember) error {
	return nil
}
func (s *tripServiceMemberRepoStub) Delete(context.Context, string) error             { return nil }
func (s *tripServiceMemberRepoStub) SoftDelete(context.Context, string, string) error { return nil }
func (s *tripServiceMemberRepoStub) GetByTripAndUser(context.Context, string, string) (domaintripmember.TripMember, error) {
	return s.member, s.err
}
func (s *tripServiceMemberRepoStub) GetActiveByTripAndUser(context.Context, string, string) (domaintripmember.TripMember, error) {
	return s.member, s.err
}
func (s *tripServiceMemberRepoStub) ListByTrip(context.Context, string) ([]domaintripmember.TripMember, error) {
	return s.members, s.err
}

type tripServiceDayRepoStub struct {
	days    []domainitineraryday.ItineraryDay
	err     error
	updated int
	created int
	deleted int
}

func (s *tripServiceDayRepoStub) Store(context.Context, domainitineraryday.ItineraryDay) error {
	s.created++
	return nil
}
func (s *tripServiceDayRepoStub) GetByID(context.Context, string) (domainitineraryday.ItineraryDay, error) {
	return domainitineraryday.ItineraryDay{}, s.err
}
func (s *tripServiceDayRepoStub) GetAll(context.Context, filter.BaseParams) ([]domainitineraryday.ItineraryDay, int64, error) {
	return s.days, int64(len(s.days)), s.err
}
func (s *tripServiceDayRepoStub) Update(context.Context, domainitineraryday.ItineraryDay) error {
	s.updated++
	return nil
}
func (s *tripServiceDayRepoStub) Delete(context.Context, string) error { return nil }
func (s *tripServiceDayRepoStub) SoftDelete(context.Context, string, string) error {
	s.deleted++
	return nil
}
func (s *tripServiceDayRepoStub) ListByTrip(context.Context, string) ([]domainitineraryday.ItineraryDay, error) {
	return s.days, s.err
}

type tripServiceItemRepoStub struct {
	items map[string][]domainitineraryitem.ItineraryItem
	err   error
}

func (s *tripServiceItemRepoStub) Store(context.Context, domainitineraryitem.ItineraryItem) error {
	return nil
}
func (s *tripServiceItemRepoStub) GetByID(context.Context, string) (domainitineraryitem.ItineraryItem, error) {
	return domainitineraryitem.ItineraryItem{}, s.err
}
func (s *tripServiceItemRepoStub) GetAll(context.Context, filter.BaseParams) ([]domainitineraryitem.ItineraryItem, int64, error) {
	return nil, 0, nil
}
func (s *tripServiceItemRepoStub) Update(context.Context, domainitineraryitem.ItineraryItem) error {
	return nil
}
func (s *tripServiceItemRepoStub) Delete(context.Context, string) error             { return nil }
func (s *tripServiceItemRepoStub) SoftDelete(context.Context, string, string) error { return nil }
func (s *tripServiceItemRepoStub) GetByDay(_ context.Context, dayId string) ([]domainitineraryitem.ItineraryItem, error) {
	return s.items[dayId], s.err
}
func (s *tripServiceItemRepoStub) GetByIDs(context.Context, []string) ([]domainitineraryitem.ItineraryItem, error) {
	return nil, nil
}
func (s *tripServiceItemRepoStub) Reorder(context.Context, string, []domainitineraryitem.ItineraryItem) error {
	return nil
}

type tripServiceUserRepoStub struct {
	users map[string]domainuser.Users
	err   error
}

func (s *tripServiceUserRepoStub) Store(context.Context, domainuser.Users) error { return nil }
func (s *tripServiceUserRepoStub) GetByID(_ context.Context, id string) (domainuser.Users, error) {
	return s.users[id], s.err
}
func (s *tripServiceUserRepoStub) GetAll(context.Context, filter.BaseParams) ([]domainuser.Users, int64, error) {
	return nil, 0, nil
}
func (s *tripServiceUserRepoStub) Update(context.Context, domainuser.Users) error   { return nil }
func (s *tripServiceUserRepoStub) Delete(context.Context, string) error             { return nil }
func (s *tripServiceUserRepoStub) SoftDelete(context.Context, string, string) error { return nil }
func (s *tripServiceUserRepoStub) GetByEmail(context.Context, string) (domainuser.Users, error) {
	return domainuser.Users{}, nil
}
func (s *tripServiceUserRepoStub) GetByPhone(context.Context, string) (domainuser.Users, error) {
	return domainuser.Users{}, nil
}

func tripServiceForTest() (*TripService, *tripServiceRepoStub, *tripServiceDayRepoStub) {
	tripRepo := &tripServiceRepoStub{}
	dayRepo := &tripServiceDayRepoStub{}
	return NewTripService(tripRepo, &tripServiceMemberRepoStub{member: domaintripmember.TripMember{Id: "member-1", TripId: "trip-1", UserId: "user-1", Role: serviceshared.TripRoleOwner}}, dayRepo, &tripServiceItemRepoStub{items: map[string][]domainitineraryitem.ItineraryItem{}}, &tripServiceUserRepoStub{users: map[string]domainuser.Users{"user-1": {Id: "user-1", Name: "Owner"}}}), tripRepo, dayRepo
}

func TestTripServiceCreateGetListUpdateDelete(t *testing.T) {
	svc, tripRepo, dayRepo := tripServiceForTest()
	created, err := svc.CreateTrip(context.Background(), "user-1", dto.CreateTripRequest{Title: " Bali ", Destination: "Bali", StartDate: "2026-08-01", EndDate: "2026-08-02", Timezone: "Asia/Jakarta", CurrencyCode: "idr"})
	if err != nil || !tripRepo.created || created.Title != "Bali" || created.CurrencyCode != "IDR" {
		t.Fatalf("create: %+v err=%v", created, err)
	}

	tripRepo.trip = domaintrip.Trip{Id: "trip-1", OwnerId: "user-1", Title: "Old", CurrencyCode: "IDR"}
	memberRepo := svc.memberRepo.(*tripServiceMemberRepoStub)
	memberRepo.members = []domaintripmember.TripMember{{Id: "member-1", TripId: "trip-1", UserId: "user-1", Role: serviceshared.TripRoleOwner}}
	day := domainitineraryday.ItineraryDay{Id: "day-1", TripId: "trip-1", DayNumber: 1}
	dayRepo.days = []domainitineraryday.ItineraryDay{day}
	detail, err := svc.GetTripDetail(context.Background(), "user-1", "trip-1")
	if err != nil || detail.Id != "trip-1" || len(detail.Members) != 1 || len(detail.Days) != 1 {
		t.Fatalf("detail: %+v err=%v", detail, err)
	}

	tripRepo.trips = []domaintrip.Trip{tripRepo.trip}
	list, err := svc.ListTrips(context.Background(), "user-1")
	if err != nil || len(list) != 1 || list[0].MemberCount != 1 {
		t.Fatalf("list: %+v err=%v", list, err)
	}

	updated, err := svc.UpdateTrip(context.Background(), "user-1", "trip-1", dto.UpdateTripRequest{Title: " New ", Destination: " New Place ", Status: "planned", Timezone: "UTC", CurrencyCode: "usd"})
	if err != nil || !tripRepo.updated || updated.Title != "New" || updated.CurrencyCode != "USD" {
		t.Fatalf("update: %+v err=%v", updated, err)
	}
	if err := svc.DeleteTrip(context.Background(), "user-1", "trip-1"); err != nil || !tripRepo.deleted {
		t.Fatalf("delete: err=%v deleted=%v", err, tripRepo.deleted)
	}
	if err := svc.DeleteTrip(context.Background(), "other-user", "trip-1"); !errors.Is(err, serviceshared.ErrAccessDenied) {
		t.Fatalf("delete access err=%v", err)
	}
	if dayRepo.created < 0 {
		t.Fatal("unreachable")
	}
}

func TestTripServiceValidationAndAccessErrors(t *testing.T) {
	svc, _, _ := tripServiceForTest()
	for _, tc := range []struct {
		name string
		req  dto.CreateTripRequest
		want error
	}{
		{"timezone", dto.CreateTripRequest{Timezone: "No/Such"}, ErrInvalidTimezone},
		{"currency", dto.CreateTripRequest{CurrencyCode: "BAD!"}, ErrInvalidCurrency},
		{"date", dto.CreateTripRequest{StartDate: "bad"}, serviceshared.ErrInvalidDate},
		{"range", dto.CreateTripRequest{StartDate: "2026-08-02", EndDate: "2026-08-01"}, ErrInvalidDateRange},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.CreateTrip(context.Background(), "user-1", tc.req)
			if !errors.Is(err, tc.want) {
				t.Fatalf("err=%v want=%v", err, tc.want)
			}
		})
	}
	memberRepo := svc.memberRepo.(*tripServiceMemberRepoStub)
	memberRepo.member = domaintripmember.TripMember{}
	if _, err := svc.GetTripDetail(context.Background(), "user-1", "trip-1"); !errors.Is(err, serviceshared.ErrNotMember) {
		t.Fatalf("member err=%v", err)
	}
}

func TestNewTripService(t *testing.T) {
	svc := NewTripService(nil, nil, nil, nil, nil)
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

func TestBuildItineraryDaySyncPlanCreatesDateRange(t *testing.T) {
	start, _ := serviceshared.ParseDate("2026-07-10")
	end, _ := serviceshared.ParseDate("2026-07-12")
	now := time.Date(2026, 7, 1, 10, 0, 0, 0, time.UTC)

	plan := buildItineraryDaySyncPlan(domaintrip.Trip{
		Id:        "trip-1",
		StartDate: &start,
		EndDate:   &end,
	}, "user-1", now, nil)

	if len(plan.Create) != 3 || len(plan.Update) != 0 || len(plan.Delete) != 0 {
		t.Fatalf("unexpected sync plan: %+v", plan)
	}
	for i, day := range plan.Create {
		if day.TripId != "trip-1" || day.DayNumber != i+1 || day.CreatedBy != "user-1" || day.UpdatedBy != "user-1" {
			t.Fatalf("unexpected day[%d]: %+v", i, day)
		}
	}
	if got := plan.Create[2].Date.Format("2006-01-02"); got != "2026-07-12" {
		t.Fatalf("expected last day date 2026-07-12, got %s", got)
	}
}

func TestBuildItineraryDaySyncPlanUpdatesCreatesAndDeletesByDateRange(t *testing.T) {
	start, _ := serviceshared.ParseDate("2026-07-10")
	end, _ := serviceshared.ParseDate("2026-07-12")
	oldStart, _ := serviceshared.ParseDate("2026-07-01")
	oldSecond, _ := serviceshared.ParseDate("2026-07-02")
	oldFourth, _ := serviceshared.ParseDate("2026-07-04")
	now := time.Date(2026, 7, 1, 10, 0, 0, 0, time.UTC)

	plan := buildItineraryDaySyncPlan(domaintrip.Trip{
		Id:        "trip-1",
		StartDate: &start,
		EndDate:   &end,
	}, "user-1", now, []domainitineraryday.ItineraryDay{
		{Id: "day-1", TripId: "trip-1", DayNumber: 1, Date: &oldStart},
		{Id: "day-2", TripId: "trip-1", DayNumber: 2, Date: &oldSecond},
		{Id: "day-4", TripId: "trip-1", DayNumber: 4, Date: &oldFourth},
	})

	if len(plan.Update) != 2 || len(plan.Create) != 1 || len(plan.Delete) != 1 {
		t.Fatalf("unexpected sync plan: %+v", plan)
	}
	if plan.Update[0].DayNumber != 1 || plan.Update[0].Date.Format("2006-01-02") != "2026-07-10" {
		t.Fatalf("unexpected first updated day: %+v", plan.Update[0])
	}
	if plan.Update[1].DayNumber != 2 || plan.Update[1].Date.Format("2006-01-02") != "2026-07-11" {
		t.Fatalf("unexpected second updated day: %+v", plan.Update[1])
	}
	if plan.Create[0].DayNumber != 3 || plan.Create[0].Date.Format("2006-01-02") != "2026-07-12" {
		t.Fatalf("unexpected created day: %+v", plan.Create[0])
	}
	if plan.Delete[0].DayNumber != 4 {
		t.Fatalf("unexpected deleted day: %+v", plan.Delete[0])
	}
}

func TestTripToDetailIncludesTotalCostEstimate(t *testing.T) {
	day1 := domainitineraryday.ItineraryDay{Id: "day-1", TripId: "trip-1", DayNumber: 1}
	day2 := domainitineraryday.ItineraryDay{Id: "day-2", TripId: "trip-1", DayNumber: 2}

	result := tripToDetail(
		domaintrip.Trip{Id: "trip-1", Title: "Bali", CurrencyCode: "IDR"},
		nil,
		nil,
		[]domainitineraryday.ItineraryDay{day1, day2},
		map[string][]domainitineraryitem.ItineraryItem{
			"day-1": {
				{Id: "item-1", DayId: "day-1", CostEstimate: 150000},
				{Id: "item-2", DayId: "day-1", CostEstimate: 250000},
			},
			"day-2": {
				{Id: "item-3", DayId: "day-2", CostEstimate: 100000},
			},
		},
	)

	if result.TotalCostEstimate != 500000 {
		t.Fatalf("expected total cost estimate 500000, got %v", result.TotalCostEstimate)
	}
}

func TestApplyTripAccessMetadata(t *testing.T) {
	tests := []struct {
		name             string
		member           domaintripmember.TripMember
		wantRole         string
		wantEdit         bool
		wantManageMember bool
		wantDelete       bool
		wantLeave        bool
	}{
		{
			name:             "owner",
			member:           domaintripmember.TripMember{Id: "member-owner", Role: serviceshared.TripRoleOwner},
			wantRole:         serviceshared.TripRoleOwner,
			wantEdit:         true,
			wantManageMember: true,
			wantDelete:       true,
			wantLeave:        false,
		},
		{
			name:      "editor",
			member:    domaintripmember.TripMember{Id: "member-editor", Role: serviceshared.TripRoleEditor},
			wantRole:  serviceshared.TripRoleEditor,
			wantEdit:  true,
			wantLeave: true,
		},
		{
			name:      "viewer",
			member:    domaintripmember.TripMember{Id: "member-viewer", Role: serviceshared.TripRoleViewer},
			wantRole:  serviceshared.TripRoleViewer,
			wantLeave: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail := applyTripAccessMetadata(dtoTripDetailForTest(), tt.member)

			if detail.CurrentMemberId != tt.member.Id {
				t.Fatalf("expected current member id %s, got %s", tt.member.Id, detail.CurrentMemberId)
			}
			if detail.CurrentMemberRole != tt.wantRole {
				t.Fatalf("expected role %s, got %s", tt.wantRole, detail.CurrentMemberRole)
			}
			if detail.CanEdit != tt.wantEdit {
				t.Fatalf("expected can_edit %v, got %v", tt.wantEdit, detail.CanEdit)
			}
			if detail.CanManageMembers != tt.wantManageMember {
				t.Fatalf("expected can_manage_members %v, got %v", tt.wantManageMember, detail.CanManageMembers)
			}
			if detail.CanDelete != tt.wantDelete {
				t.Fatalf("expected can_delete %v, got %v", tt.wantDelete, detail.CanDelete)
			}
			if detail.CanLeave != tt.wantLeave {
				t.Fatalf("expected can_leave %v, got %v", tt.wantLeave, detail.CanLeave)
			}
		})
	}
}

func dtoTripDetailForTest() dto.TripDetailResponse {
	return dto.TripDetailResponse{Id: "trip-1"}
}
