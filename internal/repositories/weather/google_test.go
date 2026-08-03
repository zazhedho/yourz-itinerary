package repositoryweather

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	domainweather "yourz-itinerary/internal/domain/weather"
)

func TestGoogleWeatherProviderMapsDailyForecast(t *testing.T) {
	today := time.Date(2026, time.August, 2, 0, 0, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		for key, want := range map[string]string{
			"key": "secret", "unitsSystem": "METRIC", "languageCode": "id", "days": "10",
			"location.latitude": "-6.1754", "location.longitude": "106.8272",
		} {
			if query.Get(key) != want {
				t.Errorf("query %s=%q, want %q", key, query.Get(key), want)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
  "timeZone": {"id": "Asia/Jakarta"},
  "forecastDays": [{
    "displayDate": {"year": ` + strconv.Itoa(today.Year()) + `, "month": ` + strconv.Itoa(int(today.Month())) + `, "day": ` + strconv.Itoa(today.Day()) + `},
    "daytimeForecast": {
      "weatherCondition": {"type": "PARTLY_CLOUDY", "description": {"text": "Berawan sebagian"}, "iconBaseUri": "https://maps.gstatic.com/weather/v1/partly_cloudy"},
      "precipitation": {"probability": {"percent": 40}},
      "relativeHumidity": 78,
      "wind": {"speed": {"value": 12.4, "unit": "KILOMETERS_PER_HOUR"}}
    },
    "minTemperature": {"degrees": 24.1},
    "maxTemperature": {"degrees": 31.8},
    "feelsLikeMinTemperature": {"degrees": 25.0},
    "feelsLikeMaxTemperature": {"degrees": 35.2}
  }]
}`))
	}))
	defer server.Close()

	provider := NewGoogleWeatherProvider("secret", server.URL, server.Client())
	got, err := provider.GetDailyForecast(context.Background(), -6.1754, 106.8272, today)
	if err != nil {
		t.Fatalf("get forecast: %v", err)
	}
	if got.ConditionCode != "PARTLY_CLOUDY" || got.IconURI != "https://maps.gstatic.com/weather/v1/partly_cloudy.svg" || got.TimeZone != "Asia/Jakarta" || got.PrecipitationProbability != 40 || got.WindSpeedKPH != 12.4 {
		t.Fatalf("unexpected forecast: %+v", got)
	}
	if got.FeelsLikeMinC == nil || *got.FeelsLikeMinC != 25 || got.FeelsLikeMaxC == nil || *got.FeelsLikeMaxC != 35.2 {
		t.Fatalf("unexpected feels-like values: %+v", got)
	}
}

func TestGoogleWeatherProviderRejectsBadResponses(t *testing.T) {
	cases := []struct {
		name string
		body string
		code int
		want string
	}{
		{name: "status", body: `{"error":{"message":"secret"}}`, code: http.StatusBadGateway, want: "status 502"},
		{name: "malformed", body: "{", code: http.StatusOK, want: ErrMalformedProviderResponse.Error()},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.code)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()
			provider := NewGoogleWeatherProvider("secret", server.URL, server.Client())
			_, err := provider.GetDailyForecast(context.Background(), 1, 2, time.Date(2026, time.August, 2, 0, 0, 0, 0, time.UTC))
			if err == nil || !strings.Contains(err.Error(), tc.want) || strings.Contains(err.Error(), "secret") {
				t.Fatalf("error=%v, want %q without key", err, tc.want)
			}
		})
	}
}

func TestGoogleWeatherProviderRejectsMissingDateAndTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"forecastDays":[]}`))
	}))
	defer server.Close()
	provider := NewGoogleWeatherProvider("secret", server.URL, server.Client())
	target := time.Date(2026, time.August, 2, 0, 0, 0, 0, time.UTC)
	_, err := provider.GetDailyForecast(context.Background(), 1, 2, target)
	if !errors.Is(err, domainweather.ErrOutOfRange) {
		t.Fatalf("expected missing date out-of-range, got %v", err)
	}

	timeoutServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { time.Sleep(50 * time.Millisecond) }))
	defer timeoutServer.Close()
	client := timeoutServer.Client()
	client.Timeout = time.Millisecond
	provider = NewGoogleWeatherProvider("secret", timeoutServer.URL, client)
	_, err = provider.GetDailyForecast(context.Background(), 1, 2, target)
	if err == nil || !strings.Contains(err.Error(), "request failed") {
		t.Fatalf("expected timeout failure, got %v", err)
	}
}
