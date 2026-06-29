package servicetrip

import (
	"testing"
	"time"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	domainitineraryitem "yourz-itinerary/internal/domain/itineraryitem"
	domaintrip "yourz-itinerary/internal/domain/trip"
	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	"yourz-itinerary/internal/dto"
	serviceshared "yourz-itinerary/internal/services/shared"
)

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
