package config

import (
	"errors"
	"strings"
	"time"

	"yourz-itinerary/pkg/configvalue"
	"yourz-itinerary/utils"
)

const (
	DefaultWeatherBaseURL        = "https://weather.googleapis.com"
	DefaultWeatherMonthlyLimit   = 9000
	DefaultWeatherRequestTimeout = 5 * time.Second
)

var (
	errWeatherMonthlyLimit   = errors.New("WEATHER_MONTHLY_LIMIT must be non-negative")
	errWeatherRequestTimeout = errors.New("WEATHER_REQUEST_TIMEOUT must be positive")
)

type WeatherConfig struct {
	Enabled        bool
	APIKey         string
	BaseURL        string
	MonthlyLimit   int
	RequestTimeout time.Duration
}

func LoadWeatherConfig() WeatherConfig {
	enabled, _ := configvalue.Bool(utils.GetEnv("WEATHER_ENABLED", "false"), false)
	limit, _ := configvalue.Int(utils.GetEnv("WEATHER_MONTHLY_LIMIT", "9000"), DefaultWeatherMonthlyLimit)
	timeout, _ := configvalue.Duration(utils.GetEnv("WEATHER_REQUEST_TIMEOUT", "5s"), DefaultWeatherRequestTimeout)

	return WeatherConfig{
		Enabled:        enabled,
		APIKey:         strings.TrimSpace(utils.GetEnv("GOOGLE_WEATHER_API_KEY", "")),
		BaseURL:        strings.TrimRight(configvalue.String(utils.GetEnv("GOOGLE_WEATHER_BASE_URL", ""), DefaultWeatherBaseURL), "/"),
		MonthlyLimit:   limit,
		RequestTimeout: timeout,
	}
}

func ValidateWeatherConfig(cfg WeatherConfig) error {
	if cfg.MonthlyLimit < 0 {
		return errWeatherMonthlyLimit
	}
	if cfg.RequestTimeout <= 0 {
		return errWeatherRequestTimeout
	}
	return nil
}
