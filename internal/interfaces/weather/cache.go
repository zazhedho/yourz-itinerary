package interfaceweather

import (
	"context"
	"errors"
	"time"

	domainweather "yourz-itinerary/internal/domain/weather"
)

var (
	ErrCacheMiss        = errors.New("weather cache miss")
	ErrCacheUnavailable = errors.New("weather cache unavailable")
)

type Cache interface {
	GetDaily(context.Context, string) (domainweather.Forecast, error)
	SetDaily(context.Context, string, domainweather.Forecast, time.Duration) error
	GetHourly(context.Context, string) ([]domainweather.HourlyForecast, error)
	SetHourly(context.Context, string, []domainweather.HourlyForecast, time.Duration) error
}
