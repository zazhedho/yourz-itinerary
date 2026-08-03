package domainweather

import (
	"errors"
	"time"
)

type Status string

const (
	StatusAvailable           Status = "available"
	StatusMissingCoordinates  Status = "missing_coordinates"
	StatusPastDate            Status = "past_date"
	StatusOutOfRange          Status = "out_of_range"
	StatusLimitReached        Status = "limit_reached"
	StatusProviderUnavailable Status = "provider_unavailable"
)

const (
	ForecastTypeDaily  = "daily"
	ForecastTypeHourly = "hourly"
)

var (
	ErrOutOfRange = errors.New("weather forecast date is outside provider range")
	ErrUsageLimit = errors.New("weather usage limit reached")
)

type Forecast struct {
	ForecastType             string
	ForecastDate             string
	ForecastTime             *time.Time
	TimeZone                 string
	ConditionCode            string
	ConditionDescription     string
	IconURI                  string
	TemperatureC             float64
	MinTemperatureC          float64
	MaxTemperatureC          float64
	FeelsLikeC               *float64
	FeelsLikeMinC            *float64
	FeelsLikeMaxC            *float64
	PrecipitationProbability int
	HumidityPercent          int
	WindSpeedKPH             float64
}

type HourlyForecast struct {
	ForecastTime             time.Time
	TimeZone                 string
	ConditionCode            string
	ConditionDescription     string
	IconURI                  string
	TemperatureC             float64
	FeelsLikeC               *float64
	PrecipitationProbability int
	HumidityPercent          int
	WindSpeedKPH             float64
}
