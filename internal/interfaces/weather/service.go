package interfaceweather

import (
	"context"

	"yourz-itinerary/internal/dto"
)

type Service interface {
	GetByDay(context.Context, string, string) (dto.WeatherDayResponse, error)
}
