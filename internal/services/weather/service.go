package serviceweather

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
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
	now := time.Now
	if s.Now != nil {
		now = s.Now
	}
	today := dateOnly(now())
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
		allowed, _, err := s.Usage.Reserve(ctx)
		if err != nil || !allowed {
			status := domainweather.StatusLimitReached
			if err != nil {
				status = domainweather.StatusProviderUnavailable
			}
			result.Status = string(status)
			for _, remainingIndexes := range coordinates {
				for _, index := range remainingIndexes {
					if result.Items[index].Status == string(domainweather.StatusMissingCoordinates) {
						result.Items[index].Status = string(status)
					}
				}
			}
			break
		}

		latitude, longitude := roundedCoordinates(items[indexes[0]])
		forecast, err := s.Provider.GetDailyForecast(ctx, latitude, longitude, targetDate)
		if err != nil {
			status := domainweather.StatusProviderUnavailable
			if errors.Is(err, domainweather.ErrOutOfRange) {
				status = domainweather.StatusOutOfRange
			}
			result.Status = string(status)
			for _, index := range indexes {
				result.Items[index].Status = string(status)
			}
			continue
		}

		for _, index := range indexes {
			result.Items[index] = mapWeatherItem(items[index], forecast)
		}
	}
	result.Status = string(aggregateStatus(result.Items))
	return result, nil
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
		ItemID: item.Id, Status: string(domainweather.StatusAvailable), ForecastDate: forecast.ForecastDate,
		TimeZone: forecast.TimeZone, ConditionCode: forecast.ConditionCode, ConditionDescription: forecast.ConditionDescription,
		IconURI: forecast.IconURI, MinTemperatureC: forecast.MinTemperatureC, MaxTemperatureC: forecast.MaxTemperatureC,
		FeelsLikeMinC: forecast.FeelsLikeMinC, FeelsLikeMaxC: forecast.FeelsLikeMaxC,
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
