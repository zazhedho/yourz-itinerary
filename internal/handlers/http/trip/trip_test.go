package handlertrip

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"yourz-itinerary/internal/authscope"
	"yourz-itinerary/internal/dto"
	servicetrip "yourz-itinerary/internal/services/trip"
	"yourz-itinerary/utils"

	"github.com/gin-gonic/gin"
)

type tripServiceMock struct {
	trip  dto.TripDetailResponse
	trips []dto.TripListResponse
	err   error
}

func (m *tripServiceMock) CreateTrip(ctx context.Context, userId string, req dto.CreateTripRequest) (dto.TripDetailResponse, error) {
	return m.trip, m.err
}
func (m *tripServiceMock) GetTripDetail(ctx context.Context, userId, tripId string) (dto.TripDetailResponse, error) {
	return m.trip, m.err
}
func (m *tripServiceMock) ListTrips(ctx context.Context, userId string) ([]dto.TripListResponse, error) {
	return m.trips, m.err
}
func (m *tripServiceMock) UpdateTrip(ctx context.Context, userId, tripId string, req dto.UpdateTripRequest) (dto.TripDetailResponse, error) {
	return m.trip, m.err
}
func (m *tripServiceMock) DeleteTrip(ctx context.Context, userId, tripId string) error {
	return m.err
}

func performTripRequest(method, routePath, requestPath string, body interface{}, handler gin.HandlerFunc, authData map[string]interface{}) *httptest.ResponseRecorder {
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

func TestTripHandlerCreateTripSuccess(t *testing.T) {
	handler := NewTripHandler(&tripServiceMock{
		trip: dto.TripDetailResponse{Id: "t-1", Title: "My Trip"},
	})
	rec := performTripRequest(http.MethodPost, "/api/trips", "/api/trips",
		dto.CreateTripRequest{Title: "My Trip"},
		handler.CreateTrip,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestTripHandlerGetTripDetailNotFound(t *testing.T) {
	handler := NewTripHandler(&tripServiceMock{
		err: servicetrip.ErrTripNotFound,
	})
	rec := performTripRequest(http.MethodGet, "/api/trips/:id", "/api/trips/550e8400-e29b-41d4-a716-446655440000",
		nil, handler.GetTripDetail,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestTripHandlerListTripsSuccess(t *testing.T) {
	handler := NewTripHandler(&tripServiceMock{
		trips: []dto.TripListResponse{{Id: "t-1", Title: "Trip A"}, {Id: "t-2", Title: "Trip B"}},
	})
	rec := performTripRequest(http.MethodGet, "/api/trips", "/api/trips",
		nil, handler.ListTrips,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestTripHandlerInvalidUUID(t *testing.T) {
	handler := NewTripHandler(&tripServiceMock{})
	rec := performTripRequest(http.MethodGet, "/api/trips/:id", "/api/trips/not-a-uuid",
		nil, handler.GetTripDetail,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
