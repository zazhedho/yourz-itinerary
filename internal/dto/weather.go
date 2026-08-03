package dto

type WeatherItemResponse struct {
	ItemID                   string   `json:"item_id"`
	Status                   string   `json:"status"`
	ForecastDate             string   `json:"forecast_date,omitempty"`
	TimeZone                 string   `json:"time_zone,omitempty"`
	ConditionCode            string   `json:"condition_code,omitempty"`
	ConditionDescription     string   `json:"condition_description,omitempty"`
	IconURI                  string   `json:"icon_uri,omitempty"`
	MinTemperatureC          float64  `json:"min_temperature_c,omitempty"`
	MaxTemperatureC          float64  `json:"max_temperature_c,omitempty"`
	FeelsLikeMinC            *float64 `json:"feels_like_min_c,omitempty"`
	FeelsLikeMaxC            *float64 `json:"feels_like_max_c,omitempty"`
	PrecipitationProbability int      `json:"precipitation_probability"`
	HumidityPercent          int      `json:"humidity_percent"`
	WindSpeedKPH             float64  `json:"wind_speed_kph"`
}

type WeatherDayResponse struct {
	DayID  string                `json:"day_id"`
	Date   string                `json:"date,omitempty"`
	Status string                `json:"status"`
	Items  []WeatherItemResponse `json:"items,omitempty"`
}
