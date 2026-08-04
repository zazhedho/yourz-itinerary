package serviceitineraryitem

import (
	"context"
	"errors"
	"testing"

	domainday "yourz-itinerary/internal/domain/itineraryday"
	domainitem "yourz-itinerary/internal/domain/itineraryitem"
	domainmember "yourz-itinerary/internal/domain/tripmember"
	"yourz-itinerary/internal/dto"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/pkg/filter"
)

type itemDayRepoStub struct {
	day domainday.ItineraryDay
	err error
}

func (s *itemDayRepoStub) Store(context.Context, domainday.ItineraryDay) error { return nil }
func (s *itemDayRepoStub) GetByID(context.Context, string) (domainday.ItineraryDay, error) {
	return s.day, s.err
}
func (s *itemDayRepoStub) GetAll(context.Context, filter.BaseParams) ([]domainday.ItineraryDay, int64, error) {
	return nil, 0, nil
}
func (s *itemDayRepoStub) Update(context.Context, domainday.ItineraryDay) error { return nil }
func (s *itemDayRepoStub) Delete(context.Context, string) error                 { return nil }
func (s *itemDayRepoStub) SoftDelete(context.Context, string, string) error     { return nil }
func (s *itemDayRepoStub) ListByTrip(context.Context, string) ([]domainday.ItineraryDay, error) {
	return nil, nil
}

type itemRepoStub struct {
	items     []domainitem.ItineraryItem
	item      domainitem.ItineraryItem
	err       error
	getErr    error
	stored    bool
	updated   bool
	deleted   bool
	reordered bool
}

func (s *itemRepoStub) Store(_ context.Context, item domainitem.ItineraryItem) error {
	s.item = item
	s.stored = true
	return s.err
}
func (s *itemRepoStub) GetByID(context.Context, string) (domainitem.ItineraryItem, error) {
	return s.item, s.getErr
}
func (s *itemRepoStub) GetAll(context.Context, filter.BaseParams) ([]domainitem.ItineraryItem, int64, error) {
	return nil, 0, nil
}
func (s *itemRepoStub) Update(_ context.Context, item domainitem.ItineraryItem) error {
	s.item = item
	s.updated = true
	return s.err
}
func (s *itemRepoStub) Delete(context.Context, string) error { return nil }
func (s *itemRepoStub) SoftDelete(context.Context, string, string) error {
	s.deleted = true
	return s.err
}
func (s *itemRepoStub) GetByDay(context.Context, string) ([]domainitem.ItineraryItem, error) {
	return s.items, nil
}
func (s *itemRepoStub) GetByIDs(context.Context, []string) ([]domainitem.ItineraryItem, error) {
	return s.items, s.err
}
func (s *itemRepoStub) Reorder(context.Context, string, []domainitem.ItineraryItem) error {
	s.reordered = true
	return s.err
}

type itemMemberRepoStub struct {
	member domainmember.TripMember
	err    error
}

func (s *itemMemberRepoStub) Store(context.Context, domainmember.TripMember) error { return nil }
func (s *itemMemberRepoStub) GetByID(context.Context, string) (domainmember.TripMember, error) {
	return s.member, s.err
}
func (s *itemMemberRepoStub) GetAll(context.Context, filter.BaseParams) ([]domainmember.TripMember, int64, error) {
	return nil, 0, nil
}
func (s *itemMemberRepoStub) Update(context.Context, domainmember.TripMember) error { return nil }
func (s *itemMemberRepoStub) Delete(context.Context, string) error                  { return nil }
func (s *itemMemberRepoStub) SoftDelete(context.Context, string, string) error      { return nil }
func (s *itemMemberRepoStub) GetByTripAndUser(context.Context, string, string) (domainmember.TripMember, error) {
	return s.member, s.err
}
func (s *itemMemberRepoStub) GetActiveByTripAndUser(context.Context, string, string) (domainmember.TripMember, error) {
	return s.member, s.err
}
func (s *itemMemberRepoStub) ListByTrip(context.Context, string) ([]domainmember.TripMember, error) {
	return nil, nil
}

func itemServiceForTest() (*ItineraryItemService, *itemRepoStub) {
	itemRepo := &itemRepoStub{}
	return NewItineraryItemService(
		&itemMemberRepoStub{member: domainmember.TripMember{Id: "member-1", Role: serviceshared.TripRoleEditor}},
		&itemDayRepoStub{day: domainday.ItineraryDay{Id: "day-1", TripId: "trip-1"}}, itemRepo,
	), itemRepo
}

func TestItineraryItemServiceCreateUpdateDelete(t *testing.T) {
	svc, repo := itemServiceForTest()
	lat, lng := -6.2, 106.8
	created, err := svc.CreateItem(context.Background(), "user-1", "day-1", dto.CreateItineraryItemRequest{Title: " Cafe ", Description: " Note ", LocationName: " Place ", Latitude: &lat, Longitude: &lng, StartTime: "09:00", EndTime: "10:00", CostEstimate: 50000})
	if err != nil || !repo.stored || created.Title != "Cafe" || created.SortOrder != 1 || repo.item.StartTime == nil {
		t.Fatalf("create: %+v err=%v", created, err)
	}

	updated, err := svc.UpdateItem(context.Background(), "user-1", "item-1", dto.UpdateItineraryItemRequest{Title: " Updated ", StartTime: "11:00", EndTime: "12:00", CostEstimate: 70000, SortOrder: 2})
	if err != nil || !repo.updated || updated.Title != "Updated" || repo.item.SortOrder != 2 {
		t.Fatalf("update: %+v err=%v", updated, err)
	}
	if err := svc.DeleteItem(context.Background(), "user-1", "item-1"); err != nil || !repo.deleted {
		t.Fatalf("delete: err=%v deleted=%v", err, repo.deleted)
	}
}

func TestItineraryItemServiceCreateAndUpdateValidation(t *testing.T) {
	svc, repo := itemServiceForTest()
	lat, lng := 1.0, 2.0
	for _, tc := range []struct {
		name string
		req  dto.CreateItineraryItemRequest
		want error
	}{
		{"partial coords", dto.CreateItineraryItemRequest{Latitude: &lat}, ErrInvalidCoordinates},
		{"bad time", dto.CreateItineraryItemRequest{StartTime: "25:00"}, ErrInvalidTime},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.CreateItem(context.Background(), "user-1", "day-1", tc.req)
			if !errors.Is(err, tc.want) {
				t.Fatalf("err=%v want=%v", err, tc.want)
			}
		})
	}
	repo.err = errors.New("store failed")
	if _, err := svc.CreateItem(context.Background(), "user-1", "day-1", dto.CreateItineraryItemRequest{Latitude: &lat, Longitude: &lng}); !errors.Is(err, repo.err) {
		t.Fatalf("store err=%v", err)
	}
	repo.err = nil
	repo.item = domainitem.ItineraryItem{Id: "item-1", DayId: "day-1"}
	repo.err = errors.New("update failed")
	if _, err := svc.UpdateItem(context.Background(), "user-1", "item-1", dto.UpdateItineraryItemRequest{Title: "x"}); !errors.Is(err, repo.err) {
		t.Fatalf("update err=%v", err)
	}
}

func TestItineraryItemServiceAccessAndReorder(t *testing.T) {
	svc, repo := itemServiceForTest()
	if err := svc.ReorderItems(context.Background(), "user-1", "day-1", dto.ReorderItineraryItemsRequest{}); !errors.Is(err, ErrReorderEmpty) {
		t.Fatalf("empty reorder err=%v", err)
	}
	repo.items = []domainitem.ItineraryItem{{Id: "item-1", DayId: "day-1"}, {Id: "item-2", DayId: "day-1"}}
	if err := svc.ReorderItems(context.Background(), "user-1", "day-1", dto.ReorderItineraryItemsRequest{ItemIds: []string{"item-2", "item-1"}}); err != nil || !repo.reordered {
		t.Fatalf("reorder: err=%v reordered=%v", err, repo.reordered)
	}
	repo.items = []domainitem.ItineraryItem{{Id: "item-1", DayId: "other-day"}}
	if err := svc.ReorderItems(context.Background(), "user-1", "day-1", dto.ReorderItineraryItemsRequest{ItemIds: []string{"item-1"}}); !errors.Is(err, ErrReorderDifferentDay) {
		t.Fatalf("different day err=%v", err)
	}

	dayRepo := &itemDayRepoStub{err: errors.New("missing")}
	svc = NewItineraryItemService(&itemMemberRepoStub{}, dayRepo, repo)
	if _, err := svc.CreateItem(context.Background(), "user-1", "day-1", dto.CreateItineraryItemRequest{}); !errors.Is(err, serviceshared.ErrDayNotFound) {
		t.Fatalf("missing day err=%v", err)
	}
	dayRepo.err = nil
	memberRepo := &itemMemberRepoStub{err: errors.New("no member")}
	svc = NewItineraryItemService(memberRepo, &itemDayRepoStub{day: domainday.ItineraryDay{TripId: "trip-1"}}, repo)
	if _, err := svc.CreateItem(context.Background(), "user-1", "day-1", dto.CreateItineraryItemRequest{}); !errors.Is(err, serviceshared.ErrNotMember) {
		t.Fatalf("member err=%v", err)
	}
}

func TestNewItineraryItemService(t *testing.T) {
	svc := NewItineraryItemService(nil, nil, nil)
	if svc == nil {
		t.Fatal("NewItineraryItemService returned nil")
	}
}

func TestItemErrorsDistinct(t *testing.T) {
	errors := [...]error{
		ErrItemNotFound, serviceshared.ErrDayNotFound, serviceshared.ErrTripNotFound, ErrInvalidTime,
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
