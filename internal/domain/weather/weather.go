package domainweather

import "errors"

type Status string

const (
	StatusAvailable           Status = "available"
	StatusMissingCoordinates  Status = "missing_coordinates"
	StatusPastDate            Status = "past_date"
	StatusOutOfRange          Status = "out_of_range"
	StatusLimitReached        Status = "limit_reached"
	StatusProviderUnavailable Status = "provider_unavailable"
)

var ErrOutOfRange = errors.New("weather forecast date is outside provider range")

type Forecast struct {
	ForecastDate             string
	TimeZone                 string
	ConditionCode            string
	ConditionDescription     string
	IconURI                  string
	MinTemperatureC          float64
	MaxTemperatureC          float64
	FeelsLikeMinC            *float64
	FeelsLikeMaxC            *float64
	PrecipitationProbability int
	HumidityPercent          int
	WindSpeedKPH             float64
}
