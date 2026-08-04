package repositoryweather

import (
	"context"
	"errors"
	"testing"
	"time"

	domainweather "yourz-itinerary/internal/domain/weather"
	interfaceweather "yourz-itinerary/internal/interfaces/weather"

	redismock "github.com/go-redis/redismock/v9"
)

func TestNewWeatherCache(t *testing.T) {
	if NewWeatherCache(nil) == nil {
		t.Fatal("NewWeatherCache returned nil")
	}
}

func TestWeatherCacheDailyRoundTripAndTTL(t *testing.T) {
	client, mock := redismock.NewClientMock()
	cache := &WeatherCache{Redis: client}
	forecast := domainweather.Forecast{ForecastType: domainweather.ForecastTypeDaily, ForecastDate: "2026-08-03", TimeZone: "Asia/Jakarta", MinTemperatureC: 24}
	payload := `{"ForecastType":"daily","ForecastDate":"2026-08-03","ForecastTime":null,"TimeZone":"Asia/Jakarta","ConditionCode":"","ConditionDescription":"","IconURI":"","TemperatureC":0,"MinTemperatureC":24,"MaxTemperatureC":0,"FeelsLikeC":null,"FeelsLikeMinC":null,"FeelsLikeMaxC":null,"PrecipitationProbability":0,"HumidityPercent":0,"WindSpeedKPH":0}`
	mock.ExpectSet("weather:daily:2026-08-03:-6.1000:106.8000", payload, time.Hour).SetVal("OK")
	mock.ExpectGet("weather:daily:2026-08-03:-6.1000:106.8000").SetVal(payload)

	key := "weather:daily:2026-08-03:-6.1000:106.8000"
	if err := cache.SetDaily(context.Background(), key, forecast, time.Hour); err != nil {
		t.Fatalf("set daily: %v", err)
	}
	got, err := cache.GetDaily(context.Background(), key)
	if err != nil || got.ForecastDate != forecast.ForecastDate || got.MinTemperatureC != forecast.MinTemperatureC {
		t.Fatalf("get daily: %+v err=%v", got, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations: %v", err)
	}
}

func TestWeatherCacheMissMalformedAndUnavailable(t *testing.T) {
	client, mock := redismock.NewClientMock()
	cache := &WeatherCache{Redis: client}
	mock.ExpectGet("weather:daily:bad").SetVal("{")
	if _, err := cache.GetDaily(context.Background(), "weather:daily:bad"); !errors.Is(err, interfaceweather.ErrCacheMiss) {
		t.Fatalf("expected malformed cache miss, got %v", err)
	}
	mock.ExpectGet("weather:daily:empty").SetVal("{}")
	if _, err := cache.GetDaily(context.Background(), "weather:daily:empty"); !errors.Is(err, interfaceweather.ErrCacheMiss) {
		t.Fatalf("expected incomplete cache miss, got %v", err)
	}
	if _, err := (&WeatherCache{}).GetDaily(context.Background(), "weather:daily:any"); !errors.Is(err, interfaceweather.ErrCacheUnavailable) {
		t.Fatalf("expected unavailable cache, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations: %v", err)
	}
}

func TestWeatherCacheHourlyRoundTrip(t *testing.T) {
	client, mock := redismock.NewClientMock()
	cache := &WeatherCache{Redis: client}
	forecastTime := time.Date(2026, time.August, 3, 7, 0, 0, 0, time.UTC)
	forecasts := []domainweather.HourlyForecast{{ForecastTime: forecastTime, TimeZone: "Asia/Jakarta", TemperatureC: 30}}
	payload := `[{"ForecastTime":"2026-08-03T07:00:00Z","TimeZone":"Asia/Jakarta","ConditionCode":"","ConditionDescription":"","IconURI":"","TemperatureC":30,"FeelsLikeC":null,"PrecipitationProbability":0,"HumidityPercent":0,"WindSpeedKPH":0}]`
	key := "weather:hourly:20260803T07Z:-6.1000:106.8000"
	mock.ExpectSet(key, payload, time.Hour).SetVal("OK")
	mock.ExpectGet(key).SetVal(payload)
	if err := cache.SetHourly(context.Background(), key, forecasts, time.Hour); err != nil {
		t.Fatalf("set hourly: %v", err)
	}
	got, err := cache.GetHourly(context.Background(), key)
	if err != nil || len(got) != 1 || got[0].TemperatureC != 30 {
		t.Fatalf("get hourly: %+v err=%v", got, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations: %v", err)
	}
}
