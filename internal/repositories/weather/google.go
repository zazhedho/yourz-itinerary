package repositoryweather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	domainweather "yourz-itinerary/internal/domain/weather"
	interfaceweather "yourz-itinerary/internal/interfaces/weather"
	"yourz-itinerary/pkg/config"
	"yourz-itinerary/pkg/logger"
)

var ErrMalformedProviderResponse = errors.New("malformed weather provider response")

type GoogleWeatherProvider struct {
	Client  *http.Client
	APIKey  string
	BaseURL string
}

func NewGoogleWeatherProvider(apiKey, baseURL string, client *http.Client) interfaceweather.Provider {
	if client == nil {
		client = &http.Client{Timeout: config.DefaultWeatherRequestTimeout}
	}
	if strings.TrimSpace(baseURL) == "" {
		baseURL = config.DefaultWeatherBaseURL
	}
	return &GoogleWeatherProvider{
		Client:  client,
		APIKey:  strings.TrimSpace(apiKey),
		BaseURL: strings.TrimRight(baseURL, "/"),
	}
}

type googleForecastResponse struct {
	ForecastDays []googleForecastDay `json:"forecastDays"`
	TimeZone     struct {
		ID string `json:"id"`
	} `json:"timeZone"`
}

type googleForecastHoursResponse struct {
	ForecastHours []googleForecastHour `json:"forecastHours"`
	TimeZone      struct {
		ID string `json:"id"`
	} `json:"timeZone"`
}

type googleForecastDay struct {
	DisplayDate googleDate        `json:"displayDate"`
	Daytime     googleDayPart     `json:"daytimeForecast"`
	Min         googleTemperature `json:"minTemperature"`
	Max         googleTemperature `json:"maxTemperature"`
	FeelsMin    googleTemperature `json:"feelsLikeMinTemperature"`
	FeelsMax    googleTemperature `json:"feelsLikeMaxTemperature"`
}

type googleForecastHour struct {
	Interval struct {
		StartTime time.Time `json:"startTime"`
	} `json:"interval"`
	Condition struct {
		IconBaseURI string `json:"iconBaseUri"`
		Description struct {
			Text string `json:"text"`
		} `json:"description"`
		Type string `json:"type"`
	} `json:"weatherCondition"`
	Temperature          googleTemperature `json:"temperature"`
	FeelsLikeTemperature googleTemperature `json:"feelsLikeTemperature"`
	Precipitation        struct {
		Probability struct {
			Percent int `json:"percent"`
		} `json:"probability"`
	} `json:"precipitation"`
	RelativeHumidity int `json:"relativeHumidity"`
	Wind             struct {
		Speed struct {
			Unit  string  `json:"unit"`
			Value float64 `json:"value"`
		} `json:"speed"`
	} `json:"wind"`
}

type googleDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type googleTemperature struct {
	Degrees *float64 `json:"degrees"`
}

type googleDayPart struct {
	Condition struct {
		IconBaseURI string `json:"iconBaseUri"`
		Description struct {
			Text string `json:"text"`
		} `json:"description"`
		Type string `json:"type"`
	} `json:"weatherCondition"`
	Precipitation struct {
		Probability struct {
			Percent int `json:"percent"`
		} `json:"probability"`
	} `json:"precipitation"`
	RelativeHumidity int `json:"relativeHumidity"`
	Wind             struct {
		Speed struct {
			Unit  string  `json:"unit"`
			Value float64 `json:"value"`
		} `json:"speed"`
	} `json:"wind"`
}

func (p *GoogleWeatherProvider) GetDailyForecast(ctx context.Context, latitude, longitude float64, targetDate time.Time) (domainweather.Forecast, error) {
	target := time.Date(targetDate.UTC().Year(), targetDate.UTC().Month(), targetDate.UTC().Day(), 0, 0, 0, 0, time.UTC)
	const days = 10

	query := url.Values{
		"key":                []string{p.APIKey},
		"location.latitude":  []string{strconv.FormatFloat(latitude, 'f', -1, 64)},
		"location.longitude": []string{strconv.FormatFloat(longitude, 'f', -1, 64)},
		"unitsSystem":        []string{"METRIC"},
		"languageCode":       []string{"id"},
		"days":               []string{strconv.Itoa(days)},
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, p.BaseURL+"/v1/forecast/days:lookup?"+query.Encode(), nil)
	if err != nil {
		return domainweather.Forecast{}, errors.New("failed to prepare weather provider request")
	}
	started := time.Now()
	response, err := p.Client.Do(request)
	if err != nil {
		logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("[GoogleWeather] latency_ms=%d outcome=provider_unavailable", time.Since(started).Milliseconds()))
		return domainweather.Forecast{}, errors.New("weather provider request failed")
	}
	defer response.Body.Close()
	latency := time.Since(started).Milliseconds()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("[GoogleWeather] latency_ms=%d status=%d outcome=provider_unavailable", latency, response.StatusCode))
		return domainweather.Forecast{}, fmt.Errorf("weather provider returned status %d", response.StatusCode)
	}

	var payload googleForecastResponse
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&payload); err != nil {
		logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("[GoogleWeather] latency_ms=%d status=%d outcome=provider_unavailable", latency, response.StatusCode))
		return domainweather.Forecast{}, ErrMalformedProviderResponse
	}
	for _, day := range payload.ForecastDays {
		if day.DisplayDate.Year != target.Year() || day.DisplayDate.Month != int(target.Month()) || day.DisplayDate.Day != target.Day() {
			continue
		}
		forecast := domainweather.Forecast{
			ForecastDate:             target.Format("2006-01-02"),
			TimeZone:                 payload.TimeZone.ID,
			ConditionCode:            day.Daytime.Condition.Type,
			ConditionDescription:     day.Daytime.Condition.Description.Text,
			IconURI:                  weatherIconURI(day.Daytime.Condition.IconBaseURI),
			PrecipitationProbability: day.Daytime.Precipitation.Probability.Percent,
			HumidityPercent:          day.Daytime.RelativeHumidity,
			WindSpeedKPH:             day.Daytime.Wind.Speed.Value,
		}
		if day.Min.Degrees != nil {
			forecast.MinTemperatureC = *day.Min.Degrees
		}
		if day.Max.Degrees != nil {
			forecast.MaxTemperatureC = *day.Max.Degrees
		}
		if day.FeelsMin.Degrees != nil {
			value := *day.FeelsMin.Degrees
			forecast.FeelsLikeMinC = &value
		}
		if day.FeelsMax.Degrees != nil {
			value := *day.FeelsMax.Degrees
			forecast.FeelsLikeMaxC = &value
		}
		if strings.EqualFold(day.Daytime.Wind.Speed.Unit, "MILES_PER_HOUR") {
			forecast.WindSpeedKPH *= 1.609344
		}
		logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("[GoogleWeather] latency_ms=%d status=%d outcome=available", latency, response.StatusCode))
		return forecast, nil
	}
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("[GoogleWeather] latency_ms=%d status=%d outcome=out_of_range", latency, response.StatusCode))
	return domainweather.Forecast{}, domainweather.ErrOutOfRange
}

func (p *GoogleWeatherProvider) GetHourlyForecast(ctx context.Context, latitude, longitude float64, _ time.Time) ([]domainweather.HourlyForecast, error) {
	query := url.Values{
		"key":                []string{p.APIKey},
		"location.latitude":  []string{strconv.FormatFloat(latitude, 'f', -1, 64)},
		"location.longitude": []string{strconv.FormatFloat(longitude, 'f', -1, 64)},
		"unitsSystem":        []string{"METRIC"},
		"languageCode":       []string{"id"},
		"hours":              []string{"24"},
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, p.BaseURL+"/v1/forecast/hours:lookup?"+query.Encode(), nil)
	if err != nil {
		return nil, errors.New("failed to prepare weather provider request")
	}
	started := time.Now()
	response, err := p.Client.Do(request)
	if err != nil {
		logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("[GoogleWeather] hourly latency_ms=%d outcome=provider_unavailable", time.Since(started).Milliseconds()))
		return nil, errors.New("weather provider hourly request failed")
	}
	defer response.Body.Close()
	latency := time.Since(started).Milliseconds()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("[GoogleWeather] hourly latency_ms=%d status=%d outcome=provider_unavailable", latency, response.StatusCode))
		return nil, fmt.Errorf("weather provider hourly returned status %d", response.StatusCode)
	}

	var payload googleForecastHoursResponse
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&payload); err != nil {
		logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("[GoogleWeather] hourly latency_ms=%d status=%d outcome=provider_unavailable", latency, response.StatusCode))
		return nil, ErrMalformedProviderResponse
	}
	forecasts := make([]domainweather.HourlyForecast, 0, len(payload.ForecastHours))
	for _, hour := range payload.ForecastHours {
		if hour.Interval.StartTime.IsZero() {
			continue
		}
		forecast := domainweather.HourlyForecast{
			ForecastTime:             hour.Interval.StartTime,
			TimeZone:                 payload.TimeZone.ID,
			ConditionCode:            hour.Condition.Type,
			ConditionDescription:     hour.Condition.Description.Text,
			IconURI:                  weatherIconURI(hour.Condition.IconBaseURI),
			PrecipitationProbability: hour.Precipitation.Probability.Percent,
			HumidityPercent:          hour.RelativeHumidity,
			WindSpeedKPH:             hour.Wind.Speed.Value,
		}
		if hour.Temperature.Degrees != nil {
			forecast.TemperatureC = *hour.Temperature.Degrees
		}
		if hour.FeelsLikeTemperature.Degrees != nil {
			value := *hour.FeelsLikeTemperature.Degrees
			forecast.FeelsLikeC = &value
		}
		if strings.EqualFold(hour.Wind.Speed.Unit, "MILES_PER_HOUR") {
			forecast.WindSpeedKPH *= 1.609344
		}
		forecasts = append(forecasts, forecast)
	}
	if len(forecasts) == 0 {
		logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("[GoogleWeather] hourly latency_ms=%d status=%d outcome=out_of_range", latency, response.StatusCode))
		return nil, domainweather.ErrOutOfRange
	}
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("[GoogleWeather] hourly latency_ms=%d status=%d outcome=available", latency, response.StatusCode))
	return forecasts, nil
}

func weatherIconURI(base string) string {
	base = strings.TrimSpace(base)
	if base == "" {
		return ""
	}
	if strings.HasSuffix(base, ".svg") || strings.HasSuffix(base, ".png") {
		return base
	}
	return base + ".svg"
}
