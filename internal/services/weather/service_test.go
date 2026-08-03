package serviceweather

import (
	"context"
	"errors"
	"testing"
	"time"

	domainday "yourz-itinerary/internal/domain/itineraryday"
	domainitem "yourz-itinerary/internal/domain/itineraryitem"
	domainmember "yourz-itinerary/internal/domain/tripmember"
	domainweather "yourz-itinerary/internal/domain/weather"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/pkg/filter"

	"gorm.io/gorm"
)

type weatherDayRepoFake struct {
	day domainday.ItineraryDay
	err error
}

func (f *weatherDayRepoFake) Store(context.Context, domainday.ItineraryDay) error { return nil }
func (f *weatherDayRepoFake) GetByID(context.Context, string) (domainday.ItineraryDay, error) {
	return f.day, f.err
}
func (f *weatherDayRepoFake) GetAll(context.Context, filter.BaseParams) ([]domainday.ItineraryDay, int64, error) {
	return nil, 0, nil
}
func (f *weatherDayRepoFake) Update(context.Context, domainday.ItineraryDay) error { return nil }
func (f *weatherDayRepoFake) Delete(context.Context, string) error                 { return nil }
func (f *weatherDayRepoFake) SoftDelete(context.Context, string, string) error     { return nil }
func (f *weatherDayRepoFake) ListByTrip(context.Context, string) ([]domainday.ItineraryDay, error) {
	return nil, nil
}

type weatherItemRepoFake struct {
	items []domainitem.ItineraryItem
	err   error
}

func (f *weatherItemRepoFake) Store(context.Context, domainitem.ItineraryItem) error { return nil }
func (f *weatherItemRepoFake) GetByID(context.Context, string) (domainitem.ItineraryItem, error) {
	return domainitem.ItineraryItem{}, nil
}
func (f *weatherItemRepoFake) GetAll(context.Context, filter.BaseParams) ([]domainitem.ItineraryItem, int64, error) {
	return nil, 0, nil
}
func (f *weatherItemRepoFake) Update(context.Context, domainitem.ItineraryItem) error { return nil }
func (f *weatherItemRepoFake) Delete(context.Context, string) error                   { return nil }
func (f *weatherItemRepoFake) SoftDelete(context.Context, string, string) error       { return nil }
func (f *weatherItemRepoFake) GetByDay(context.Context, string) ([]domainitem.ItineraryItem, error) {
	return f.items, f.err
}
func (f *weatherItemRepoFake) GetByIDs(context.Context, []string) ([]domainitem.ItineraryItem, error) {
	return nil, nil
}
func (f *weatherItemRepoFake) Reorder(context.Context, string, []domainitem.ItineraryItem) error {
	return nil
}

type weatherMemberRepoFake struct {
	member domainmember.TripMember
	err    error
}

func (f *weatherMemberRepoFake) Store(context.Context, domainmember.TripMember) error { return nil }
func (f *weatherMemberRepoFake) GetByID(context.Context, string) (domainmember.TripMember, error) {
	return domainmember.TripMember{}, nil
}
func (f *weatherMemberRepoFake) GetAll(context.Context, filter.BaseParams) ([]domainmember.TripMember, int64, error) {
	return nil, 0, nil
}
func (f *weatherMemberRepoFake) Update(context.Context, domainmember.TripMember) error { return nil }
func (f *weatherMemberRepoFake) Delete(context.Context, string) error                  { return nil }
func (f *weatherMemberRepoFake) SoftDelete(context.Context, string, string) error      { return nil }
func (f *weatherMemberRepoFake) GetByTripAndUser(context.Context, string, string) (domainmember.TripMember, error) {
	return f.member, f.err
}
func (f *weatherMemberRepoFake) GetActiveByTripAndUser(context.Context, string, string) (domainmember.TripMember, error) {
	return f.member, f.err
}
func (f *weatherMemberRepoFake) ListByTrip(context.Context, string) ([]domainmember.TripMember, error) {
	return nil, nil
}

type weatherProviderFake struct {
	calls            int
	err              error
	errorsByLatitude map[float64]error
}

func (f *weatherProviderFake) GetDailyForecast(_ context.Context, latitude, _ float64, targetDate time.Time) (domainweather.Forecast, error) {
	f.calls++
	if err := f.errorsByLatitude[latitude]; err != nil {
		return domainweather.Forecast{}, err
	}
	if f.err != nil {
		return domainweather.Forecast{}, f.err
	}
	return domainweather.Forecast{ForecastDate: targetDate.UTC().Format("2006-01-02"), MinTemperatureC: 24, MaxTemperatureC: 31}, nil
}

type weatherUsageFake struct {
	allowed  bool
	sequence []bool
	calls    int
}

func (f *weatherUsageFake) Reserve(context.Context) (bool, int64, error) {
	index := f.calls
	f.calls++
	if index < len(f.sequence) {
		return f.sequence[index], int64(f.calls), nil
	}
	return f.allowed, int64(f.calls), nil
}

var weatherTestNow = time.Date(2026, time.August, 3, 12, 0, 0, 0, time.UTC)

func weatherServiceFixture(t *testing.T, items []domainitem.ItineraryItem, provider *weatherProviderFake, usage *weatherUsageFake) *WeatherService {
	t.Helper()
	return weatherServiceFixtureAt(t, items, provider, usage, weatherTestNow)
}

func weatherServiceFixtureAt(t *testing.T, items []domainitem.ItineraryItem, provider *weatherProviderFake, usage *weatherUsageFake, date time.Time) *WeatherService {
	t.Helper()
	service := NewWeatherService(
		&weatherDayRepoFake{day: domainday.ItineraryDay{Id: "day-1", TripId: "trip-1", Date: &date}},
		&weatherItemRepoFake{items: items},
		&weatherMemberRepoFake{member: domainmember.TripMember{Id: "member-1", Role: serviceshared.TripRoleViewer}},
		provider,
		usage,
	)
	service.Now = func() time.Time { return weatherTestNow }
	return service
}

func weatherItem(id string, latitude, longitude float64) domainitem.ItineraryItem {
	return domainitem.ItineraryItem{Id: id, Latitude: &latitude, Longitude: &longitude}
}

func TestWeatherServiceViewerCanFetchWeather(t *testing.T) {
	provider := &weatherProviderFake{}
	usage := &weatherUsageFake{allowed: true}
	got, err := weatherServiceFixture(t, []domainitem.ItineraryItem{weatherItem("item-1", -6.1, 106.8)}, provider, usage).GetByDay(context.Background(), "user-1", "day-1")
	if err != nil || got.Status != string(domainweather.StatusAvailable) || got.Items[0].Status != string(domainweather.StatusAvailable) {
		t.Fatalf("unexpected result: %+v err=%v", got, err)
	}
}

func TestWeatherServiceDeniesNonMemberAndMissingDay(t *testing.T) {
	service := weatherServiceFixture(t, nil, &weatherProviderFake{}, &weatherUsageFake{allowed: true})
	service.MemberRepo = &weatherMemberRepoFake{}
	if _, err := service.GetByDay(context.Background(), "user-1", "day-1"); !errors.Is(err, serviceshared.ErrNotMember) {
		t.Fatalf("expected not member, got %v", err)
	}
	service.DayRepo = &weatherDayRepoFake{err: gorm.ErrRecordNotFound}
	if _, err := service.GetByDay(context.Background(), "user-1", "day-1"); !errors.Is(err, serviceshared.ErrDayNotFound) {
		t.Fatalf("expected day not found, got %v", err)
	}
}

func TestWeatherServiceSkipsMissingCoordinatesAndDeduplicates(t *testing.T) {
	provider := &weatherProviderFake{}
	usage := &weatherUsageFake{allowed: true}
	items := []domainitem.ItineraryItem{
		{Id: "missing"},
		weatherItem("one", -6.12341, 106.82721),
		weatherItem("two", -6.12344, 106.82724),
	}
	got, err := weatherServiceFixture(t, items, provider, usage).GetByDay(context.Background(), "user-1", "day-1")
	if err != nil || provider.calls != 1 || usage.calls != 1 || got.Items[0].Status != string(domainweather.StatusMissingCoordinates) || got.Items[1].Status != string(domainweather.StatusAvailable) || got.Items[2].Status != string(domainweather.StatusAvailable) {
		t.Fatalf("unexpected dedupe result: %+v calls=%d usage=%d err=%v", got, provider.calls, usage.calls, err)
	}
}

func TestWeatherServiceOutOfRangeMakesNoProviderCall(t *testing.T) {
	for _, tc := range []struct {
		offset int
		status domainweather.Status
	}{
		{offset: -2, status: domainweather.StatusPastDate},
		{offset: 11, status: domainweather.StatusOutOfRange},
	} {
		provider := &weatherProviderFake{}
		usage := &weatherUsageFake{allowed: true}
		date := weatherTestNow.AddDate(0, 0, tc.offset)
		service := weatherServiceFixtureAt(t, []domainitem.ItineraryItem{weatherItem("item-1", 1, 2)}, provider, usage, date)
		got, err := service.GetByDay(context.Background(), "user-1", "day-1")
		if err != nil || provider.calls != 0 || usage.calls != 0 || got.Status != string(tc.status) || got.Items[0].Status != string(tc.status) {
			t.Fatalf("unexpected range result offset=%d: %+v provider=%d usage=%d err=%v", tc.offset, got, provider.calls, usage.calls, err)
		}
	}
}

func TestWeatherServiceAcceptsDatesAroundUTCBoundary(t *testing.T) {
	for _, offset := range []int{-1, 1} {
		provider := &weatherProviderFake{}
		usage := &weatherUsageFake{allowed: true}
		date := weatherTestNow.AddDate(0, 0, offset)
		service := weatherServiceFixtureAt(t, []domainitem.ItineraryItem{weatherItem("item-1", 1, 2)}, provider, usage, date)
		got, err := service.GetByDay(context.Background(), "user-1", "day-1")
		if err != nil || provider.calls != 1 || got.Status != string(domainweather.StatusAvailable) {
			t.Fatalf("timezone boundary offset=%d rejected: %+v provider=%d err=%v", offset, got, provider.calls, err)
		}
	}
}

func TestWeatherServiceLimitAndProviderFailureDegrade(t *testing.T) {
	provider := &weatherProviderFake{}
	usage := &weatherUsageFake{allowed: false}
	got, err := weatherServiceFixture(t, []domainitem.ItineraryItem{weatherItem("item-1", 1, 2), weatherItem("item-2", 3, 4)}, provider, usage).GetByDay(context.Background(), "user-1", "day-1")
	if err != nil || provider.calls != 0 || got.Status != string(domainweather.StatusLimitReached) || got.Items[0].Status != string(domainweather.StatusLimitReached) {
		t.Fatalf("unexpected limit result: %+v provider=%d err=%v", got, provider.calls, err)
	}

	provider = &weatherProviderFake{err: errors.New("provider unavailable")}
	usage = &weatherUsageFake{allowed: true}
	got, err = weatherServiceFixture(t, []domainitem.ItineraryItem{weatherItem("item-1", 1, 2)}, provider, usage).GetByDay(context.Background(), "user-1", "day-1")
	if err != nil || got.Status != string(domainweather.StatusProviderUnavailable) || got.Items[0].Status != string(domainweather.StatusProviderUnavailable) {
		t.Fatalf("unexpected provider failure: %+v err=%v", got, err)
	}
}

func TestWeatherServiceAllowsAllViewRolesAndDisablesCallsWhenConfiguredOff(t *testing.T) {
	for _, role := range []string{serviceshared.TripRoleOwner, serviceshared.TripRoleEditor, serviceshared.TripRoleViewer} {
		provider := &weatherProviderFake{}
		usage := &weatherUsageFake{allowed: true}
		service := weatherServiceFixture(t, []domainitem.ItineraryItem{weatherItem("item-1", 1, 2)}, provider, usage)
		service.MemberRepo = &weatherMemberRepoFake{member: domainmember.TripMember{Id: "member-1", Role: role}}
		if _, err := service.GetByDay(context.Background(), "user-1", "day-1"); err != nil {
			t.Fatalf("role %s denied: %v", role, err)
		}
	}

	provider := &weatherProviderFake{}
	usage := &weatherUsageFake{allowed: true}
	service := weatherServiceFixture(t, []domainitem.ItineraryItem{weatherItem("item-1", 1, 2)}, provider, usage)
	service.Enabled = false
	got, err := service.GetByDay(context.Background(), "user-1", "day-1")
	if err != nil || provider.calls != 0 || usage.calls != 0 || got.Status != string(domainweather.StatusProviderUnavailable) {
		t.Fatalf("disabled weather called provider: result=%+v provider=%d usage=%d err=%v", got, provider.calls, usage.calls, err)
	}
}

func TestWeatherServiceStopsAfterMonthlyLimit(t *testing.T) {
	provider := &weatherProviderFake{}
	usage := &weatherUsageFake{sequence: []bool{true, false}}
	items := []domainitem.ItineraryItem{weatherItem("one", 1, 2), weatherItem("two", 3, 4)}
	got, err := weatherServiceFixture(t, items, provider, usage).GetByDay(context.Background(), "user-1", "day-1")
	if err != nil || provider.calls != 1 || usage.calls != 2 || got.Status != string(domainweather.StatusAvailable) {
		t.Fatalf("unexpected limit result: %+v provider=%d usage=%d err=%v", got, provider.calls, usage.calls, err)
	}
	available := 0
	limitReached := 0
	for _, item := range got.Items {
		if item.Status == string(domainweather.StatusAvailable) {
			available++
		}
		if item.Status == string(domainweather.StatusLimitReached) {
			limitReached++
		}
	}
	if available != 1 || limitReached != 1 {
		t.Fatalf("expected one available and one limited item, got %+v", got.Items)
	}
}

func TestWeatherServiceAggregatesMixedStatusesDeterministically(t *testing.T) {
	for i := 0; i < 100; i++ {
		provider := &weatherProviderFake{errorsByLatitude: map[float64]error{1: errors.New("provider unavailable")}}
		usage := &weatherUsageFake{allowed: true, sequence: []bool{true, false}}
		items := []domainitem.ItineraryItem{weatherItem("provider-error", 1, 2), weatherItem("limit", 3, 4)}
		got, err := weatherServiceFixture(t, items, provider, usage).GetByDay(context.Background(), "user-1", "day-1")
		if err != nil || got.Status != string(domainweather.StatusLimitReached) || got.Items[0].Status != string(domainweather.StatusProviderUnavailable) || got.Items[1].Status != string(domainweather.StatusLimitReached) {
			t.Fatalf("mixed status changed on iteration %d: %+v err=%v", i, got, err)
		}
	}
}
