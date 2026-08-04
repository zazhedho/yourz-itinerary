package serviceweather

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	domainitineraryitem "yourz-itinerary/internal/domain/itineraryitem"
	domainweather "yourz-itinerary/internal/domain/weather"
	"yourz-itinerary/internal/dto"
	interfaceitineraryday "yourz-itinerary/internal/interfaces/itineraryday"
	interfaceitineraryitem "yourz-itinerary/internal/interfaces/itineraryitem"
	interfacetripmember "yourz-itinerary/internal/interfaces/tripmember"
	interfaceweather "yourz-itinerary/internal/interfaces/weather"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/pkg/config"

	"gorm.io/gorm"
)

var ErrInternal = errors.New("weather internal failure")

type WeatherService struct {
	DayRepo    interfaceitineraryday.RepoItineraryDayInterface
	ItemRepo   interfaceitineraryitem.RepoItineraryItemInterface
	MemberRepo interfacetripmember.RepoTripMemberInterface
	Provider   interfaceweather.Provider
	Usage      interfaceweather.Usage
	Cache      interfaceweather.Cache
	Enabled    bool
	Now        func() time.Time
}

func NewWeatherService(
	dayRepo interfaceitineraryday.RepoItineraryDayInterface,
	itemRepo interfaceitineraryitem.RepoItineraryItemInterface,
	memberRepo interfacetripmember.RepoTripMemberInterface,
	provider interfaceweather.Provider,
	usage interfaceweather.Usage,
	configs ...config.WeatherConfig,
) *WeatherService {
	enabled := true
	if len(configs) > 0 {
		enabled = configs[0].Enabled
	}
	return &WeatherService{DayRepo: dayRepo, ItemRepo: itemRepo, MemberRepo: memberRepo, Provider: provider, Usage: usage, Enabled: enabled, Now: time.Now}
}

func (s *WeatherService) GetByDay(ctx context.Context, userID, dayID string) (dto.WeatherDayResponse, error) {
	day, err := s.DayRepo.GetByID(ctx, dayID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.WeatherDayResponse{}, serviceshared.ErrDayNotFound
		}
		return dto.WeatherDayResponse{}, fmt.Errorf("%w: load day", ErrInternal)
	}

	member, err := s.MemberRepo.GetActiveByTripAndUser(ctx, day.TripId, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.WeatherDayResponse{}, serviceshared.ErrNotMember
		}
		return dto.WeatherDayResponse{}, fmt.Errorf("%w: check membership", ErrInternal)
	}
	if member.Id == "" {
		return dto.WeatherDayResponse{}, serviceshared.ErrNotMember
	}
	if !serviceshared.CanViewTrip(member.Role) {
		return dto.WeatherDayResponse{}, serviceshared.ErrAccessDenied
	}

	items, err := s.ItemRepo.GetByDay(ctx, dayID)
	if err != nil {
		return dto.WeatherDayResponse{}, fmt.Errorf("%w: load items", ErrInternal)
	}
	result := dto.WeatherDayResponse{DayID: day.Id, Status: string(domainweather.StatusMissingCoordinates), Items: make([]dto.WeatherItemResponse, len(items))}
	for i, item := range items {
		result.Items[i] = dto.WeatherItemResponse{ItemID: item.Id, Status: string(domainweather.StatusMissingCoordinates)}
	}

	coordinates := make(map[string][]int)
	for i, item := range items {
		if item.Latitude == nil || item.Longitude == nil {
			continue
		}
		coordinates[coordinateKey(*item.Latitude, *item.Longitude)] = append(coordinates[coordinateKey(*item.Latitude, *item.Longitude)], i)
	}
	if len(coordinates) == 0 {
		return result, nil
	}

	if day.Date == nil {
		result.Status = string(domainweather.StatusOutOfRange)
		for _, indexes := range coordinates {
			for _, index := range indexes {
				result.Items[index].Status = string(domainweather.StatusOutOfRange)
			}
		}
		return result, nil
	}
	targetDate := dateOnly(*day.Date)
	now := s.currentTime()
	today := dateOnly(now)
	dayOffset := int(targetDate.Sub(today).Hours() / 24)
	if dayOffset < -1 {
		result.Status = string(domainweather.StatusPastDate)
		for _, indexes := range coordinates {
			for _, index := range indexes {
				result.Items[index].Status = string(domainweather.StatusPastDate)
			}
		}
		return result, nil
	}
	if dayOffset > 10 {
		result.Status = string(domainweather.StatusOutOfRange)
		for _, indexes := range coordinates {
			for _, index := range indexes {
				result.Items[index].Status = string(domainweather.StatusOutOfRange)
			}
		}
		return result, nil
	}
	result.Date = targetDate.Format("2006-01-02")

	if !s.Enabled || s.Provider == nil || s.Usage == nil {
		result.Status = string(domainweather.StatusProviderUnavailable)
		for _, indexes := range coordinates {
			for _, index := range indexes {
				result.Items[index].Status = string(domainweather.StatusProviderUnavailable)
			}
		}
		return result, nil
	}

	coordinateKeys := make([]string, 0, len(coordinates))
	for key := range coordinates {
		coordinateKeys = append(coordinateKeys, key)
	}
	sort.Strings(coordinateKeys)
	for _, key := range coordinateKeys {
		indexes := coordinates[key]
		latitude, longitude := roundedCoordinates(items[indexes[0]])
		cacheKey := dailyCacheKey(result.Date, latitude, longitude)
		forecast, err := s.loadDailyForecast(ctx, cacheKey, targetDate, latitude, longitude, now)
		if errors.Is(err, interfaceweather.ErrCacheUnavailable) {
			setWeatherStatus(result.Items, coordinates, domainweather.StatusProviderUnavailable)
			continue
		}
		if err != nil {
			status := domainweather.StatusProviderUnavailable
			if errors.Is(err, domainweather.ErrOutOfRange) {
				status = domainweather.StatusOutOfRange
			} else if errors.Is(err, domainweather.ErrUsageLimit) {
				status = domainweather.StatusLimitReached
			}
			if status == domainweather.StatusLimitReached {
				setWeatherStatus(result.Items, coordinates, status)
				break
			}
			for _, index := range indexes {
				result.Items[index].Status = string(status)
			}
			continue
		}

		for _, index := range indexes {
			result.Items[index] = mapWeatherItem(items[index], forecast)
		}
		s.enrichHourly(ctx, items, indexes, targetDate, now, forecast, result.Items)
	}
	result.Status = string(aggregateStatus(result.Items))
	return result, nil
}

func (s *WeatherService) loadDailyForecast(ctx context.Context, cacheKey string, targetDate time.Time, latitude, longitude float64, now time.Time) (domainweather.Forecast, error) {
	if s.Cache != nil {
		forecast, err := s.Cache.GetDaily(ctx, cacheKey)
		if err == nil {
			return forecast, nil
		}
		if !errors.Is(err, interfaceweather.ErrCacheMiss) {
			return domainweather.Forecast{}, err
		}
	}
	allowed, _, err := s.Usage.Reserve(ctx)
	if err != nil {
		return domainweather.Forecast{}, err
	}
	if !allowed {
		return domainweather.Forecast{}, domainweather.ErrUsageLimit
	}
	forecast, err := s.Provider.GetDailyForecast(ctx, latitude, longitude, targetDate)
	if err != nil {
		return domainweather.Forecast{}, err
	}
	if s.Cache != nil {
		_ = s.Cache.SetDaily(ctx, cacheKey, forecast, dailyCacheTTL(targetDate, now, forecast.TimeZone))
	}
	return forecast, nil
}

func (s *WeatherService) enrichHourly(ctx context.Context, items []domainitineraryitem.ItineraryItem, indexes []int, targetDate, now time.Time, daily domainweather.Forecast, responses []dto.WeatherItemResponse) {
	if daily.TimeZone == "" {
		return
	}
	location, err := time.LoadLocation(daily.TimeZone)
	if err != nil {
		return
	}
	eligible := make(map[int]time.Time)
	for _, index := range indexes {
		scheduled, ok := scheduledTime(targetDate, items[index].StartTime, location)
		if !ok || scheduled.Before(now) || scheduled.After(now.Add(24*time.Hour)) {
			continue
		}
		eligible[index] = scheduled
	}
	if len(eligible) == 0 {
		return
	}

	requestHour := now.UTC().Truncate(time.Hour)
	latitude, longitude := roundedCoordinates(items[indexes[0]])
	cacheKey := hourlyCacheKey(requestHour, latitude, longitude)
	var forecasts []domainweather.HourlyForecast
	if s.Cache != nil {
		forecasts, err = s.Cache.GetHourly(ctx, cacheKey)
		if err != nil && !errors.Is(err, interfaceweather.ErrCacheMiss) {
			return
		}
	}
	if forecasts == nil {
		allowed, _, reserveErr := s.Usage.Reserve(ctx)
		if reserveErr != nil || !allowed {
			return
		}
		forecasts, err = s.Provider.GetHourlyForecast(ctx, latitude, longitude, requestHour)
		if err != nil {
			return
		}
		if s.Cache != nil {
			_ = s.Cache.SetHourly(ctx, cacheKey, forecasts, time.Hour)
		}
	}

	for index, scheduled := range eligible {
		for _, forecast := range forecasts {
			start := forecast.ForecastTime.In(location)
			if scheduled.Before(start) || !scheduled.Before(start.Add(time.Hour)) {
				continue
			}
			responses[index] = mapHourlyItem(items[index], daily, forecast)
			break
		}
	}
}

func (s *WeatherService) currentTime() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
}

func scheduledTime(day time.Time, startTime *string, location *time.Location) (time.Time, bool) {
	if startTime == nil || strings.TrimSpace(*startTime) == "" {
		return time.Time{}, false
	}
	parsed, err := time.Parse("15:04", strings.TrimSpace(*startTime))
	if err != nil {
		parsed, err = time.Parse("15:04:05", strings.TrimSpace(*startTime))
		if err != nil {
			return time.Time{}, false
		}
	}
	return time.Date(day.Year(), day.Month(), day.Day(), parsed.Hour(), parsed.Minute(), parsed.Second(), 0, location), true
}

func dailyCacheKey(date string, latitude, longitude float64) string {
	return fmt.Sprintf("weather:daily:%s:%.4f:%.4f", date, latitude, longitude)
}

func hourlyCacheKey(requestHour time.Time, latitude, longitude float64) string {
	return fmt.Sprintf("weather:hourly:%s:%.4f:%.4f", requestHour.UTC().Format("20060102T15Z"), latitude, longitude)
}

func dailyCacheTTL(targetDate, now time.Time, timeZone string) time.Duration {
	location, err := time.LoadLocation(timeZone)
	if err == nil {
		now = now.In(location)
	}
	if targetDate.Year() == now.Year() && targetDate.YearDay() == now.YearDay() {
		return time.Hour
	}
	return 6 * time.Hour
}

func setWeatherStatus(items []dto.WeatherItemResponse, coordinates map[string][]int, status domainweather.Status) {
	for _, groupIndexes := range coordinates {
		for _, index := range groupIndexes {
			if items[index].Status == string(domainweather.StatusMissingCoordinates) {
				items[index].Status = string(status)
			}
		}
	}
}

func aggregateStatus(items []dto.WeatherItemResponse) domainweather.Status {
	seen := map[domainweather.Status]bool{}
	for _, item := range items {
		seen[domainweather.Status(item.Status)] = true
	}
	for _, status := range []domainweather.Status{
		domainweather.StatusAvailable,
		domainweather.StatusLimitReached,
		domainweather.StatusProviderUnavailable,
		domainweather.StatusOutOfRange,
		domainweather.StatusPastDate,
		domainweather.StatusMissingCoordinates,
	} {
		if seen[status] {
			return status
		}
	}
	return domainweather.StatusMissingCoordinates
}

func mapWeatherItem(item domainitineraryitem.ItineraryItem, forecast domainweather.Forecast) dto.WeatherItemResponse {
	return dto.WeatherItemResponse{
		ItemID: item.Id, Status: string(domainweather.StatusAvailable), ForecastType: domainweather.ForecastTypeDaily, ForecastDate: forecast.ForecastDate,
		TimeZone: forecast.TimeZone, ConditionCode: forecast.ConditionCode, ConditionDescription: forecast.ConditionDescription,
		IconURI: forecast.IconURI, MinTemperatureC: new(forecast.MinTemperatureC), MaxTemperatureC: new(forecast.MaxTemperatureC),
		FeelsLikeMinC: forecast.FeelsLikeMinC, FeelsLikeMaxC: forecast.FeelsLikeMaxC,
		PrecipitationProbability: forecast.PrecipitationProbability, HumidityPercent: forecast.HumidityPercent, WindSpeedKPH: forecast.WindSpeedKPH,
	}
}

func mapHourlyItem(item domainitineraryitem.ItineraryItem, daily domainweather.Forecast, forecast domainweather.HourlyForecast) dto.WeatherItemResponse {
	forecastTime := forecast.ForecastTime
	return dto.WeatherItemResponse{
		ItemID: item.Id, Status: string(domainweather.StatusAvailable), ForecastType: domainweather.ForecastTypeHourly,
		ForecastDate: daily.ForecastDate, ForecastTime: &forecastTime, TimeZone: forecast.TimeZone,
		ConditionCode: forecast.ConditionCode, ConditionDescription: forecast.ConditionDescription, IconURI: forecast.IconURI,
		TemperatureC: new(forecast.TemperatureC), FeelsLikeC: forecast.FeelsLikeC,
		PrecipitationProbability: forecast.PrecipitationProbability, HumidityPercent: forecast.HumidityPercent, WindSpeedKPH: forecast.WindSpeedKPH,
	}
}

func coordinateKey(latitude, longitude float64) string {
	return fmt.Sprintf("%.4f,%.4f", roundCoordinate(latitude), roundCoordinate(longitude))
}

func roundedCoordinates(item domainitineraryitem.ItineraryItem) (float64, float64) {
	return roundCoordinate(*item.Latitude), roundCoordinate(*item.Longitude)
}

func roundCoordinate(value float64) float64 {
	return math.Round(value*10000) / 10000
}

func dateOnly(value time.Time) time.Time {
	return time.Date(value.UTC().Year(), value.UTC().Month(), value.UTC().Day(), 0, 0, 0, 0, time.UTC)
}
