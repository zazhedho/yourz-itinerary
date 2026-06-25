package handleritineraryday

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"yourz-itinerary/internal/authscope"
	"yourz-itinerary/internal/dto"
	serviceitineraryday "yourz-itinerary/internal/services/itineraryday"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/utils"

	"github.com/gin-gonic/gin"
)

type itineraryDayServiceMock struct {
	day dto.ItineraryDayResponse
	err error
}

func (m *itineraryDayServiceMock) CreateDay(ctx context.Context, userId, tripId string, req dto.CreateItineraryDayRequest) (dto.ItineraryDayResponse, error) {
	return m.day, m.err
}
func (m *itineraryDayServiceMock) UpdateDay(ctx context.Context, userId, dayId string, req dto.UpdateItineraryDayRequest) (dto.ItineraryDayResponse, error) {
	return m.day, m.err
}
func (m *itineraryDayServiceMock) DeleteDay(ctx context.Context, userId, dayId string) error {
	return m.err
}

func performItineraryDayRequest(method, routePath, requestPath string, body interface{}, handler gin.HandlerFunc, authData map[string]interface{}) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Handle(method, routePath, func(ctx *gin.Context) {
		if authData != nil {
			ctx.Set(utils.CtxKeyAuthData, authData)
			ctx.Request = ctx.Request.WithContext(authscope.WithContext(ctx.Request.Context(), authscope.NewFromClaims(authData, nil)))
		}
		handler(ctx)
	})
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		raw, _ := json.Marshal(body)
		reader = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, requestPath, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func authClaims(userID, username, role string) map[string]interface{} {
	return map[string]interface{}{
		"user_id":  userID,
		"username": username,
		"role":     role,
	}
}

func TestItineraryDayCreateDaySuccess(t *testing.T) {
	handler := NewItineraryDayHandler(&itineraryDayServiceMock{
		day: dto.ItineraryDayResponse{Id: "d-1", DayNumber: 1},
	})
	rec := performItineraryDayRequest(http.MethodPost, "/api/trips/:id/days", "/api/trips/550e8400-e29b-41d4-a716-446655440000/days",
		dto.CreateItineraryDayRequest{DayNumber: 1},
		handler.CreateDay,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestItineraryDayDayNotFound(t *testing.T) {
	handler := NewItineraryDayHandler(&itineraryDayServiceMock{
		err: serviceitineraryday.ErrDayNotFound,
	})
	rec := performItineraryDayRequest(http.MethodPut, "/api/days/:id", "/api/days/550e8400-e29b-41d4-a716-446655440000",
		dto.UpdateItineraryDayRequest{DayNumber: 2},
		handler.UpdateDay,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestItineraryDayForbidden(t *testing.T) {
	handler := NewItineraryDayHandler(&itineraryDayServiceMock{
		err: serviceshared.ErrAccessDenied,
	})
	rec := performItineraryDayRequest(http.MethodDelete, "/api/days/:id", "/api/days/550e8400-e29b-41d4-a716-446655440000",
		nil, handler.DeleteDay,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", rec.Code, rec.Body.String())
	}
}
