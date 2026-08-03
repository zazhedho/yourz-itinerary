package handlerweather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"yourz-itinerary/internal/authscope"
	"yourz-itinerary/internal/dto"
	serviceshared "yourz-itinerary/internal/services/shared"

	"github.com/gin-gonic/gin"
)

type weatherServiceFake struct {
	data dto.WeatherDayResponse
	err  error
}

func (f *weatherServiceFake) GetByDay(context.Context, string, string) (dto.WeatherDayResponse, error) {
	return f.data, f.err
}

func TestWeatherHandlerMapsSuccessAndErrors(t *testing.T) {
	validPath := "/api/itinerary-days/550e8400-e29b-41d4-a716-446655440000/weather"
	for _, test := range []struct {
		name string
		err  error
		want int
	}{
		{name: "success", want: http.StatusOK},
		{name: "forbidden", err: serviceshared.ErrNotMember, want: http.StatusForbidden},
		{name: "not found", err: serviceshared.ErrDayNotFound, want: http.StatusNotFound},
	} {
		t.Run(test.name, func(t *testing.T) {
			service := &weatherServiceFake{data: dto.WeatherDayResponse{DayID: "day-1", Status: "available"}, err: test.err}
			handler := NewWeatherHandler(service)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.GET("/api/itinerary-days/:id/weather", func(ctx *gin.Context) {
				ctx.Request = ctx.Request.WithContext(authscope.WithContext(ctx.Request.Context(), authscope.New("user-1", "viewer", "viewer", nil)))
				handler.GetByDay(ctx)
			})
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, validPath, nil))
			if rec.Code != test.want {
				t.Fatalf("expected %d, got %d: %s", test.want, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestWeatherHandlerRejectsInvalidDayID(t *testing.T) {
	handler := NewWeatherHandler(&weatherServiceFake{})
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/itinerary-days/:id/weather", handler.GetByDay)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/itinerary-days/not-a-uuid/weather", nil))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}
