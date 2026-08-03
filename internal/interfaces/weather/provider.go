package interfaceweather

import (
	"context"
	"time"

	domainweather "yourz-itinerary/internal/domain/weather"
)

type Provider interface {
	GetDailyForecast(context.Context, float64, float64, time.Time) (domainweather.Forecast, error)
}
