package dto

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestWeatherItemResponseKeepsZeroTemperature(t *testing.T) {
	zero := 0.0
	payload, err := json.Marshal(WeatherItemResponse{ForecastType: "hourly", TemperatureC: &zero})
	if err != nil || !strings.Contains(string(payload), `"temperature_c":0`) {
		t.Fatalf("zero temperature missing from response: %s err=%v", payload, err)
	}
}
