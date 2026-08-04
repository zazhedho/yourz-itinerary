package handlercommon

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"yourz-itinerary/internal/authscope"
	"yourz-itinerary/internal/dto"
	serviceshared "yourz-itinerary/internal/services/shared"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func commonContext(method, body string) *gin.Context {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(method, "/", bytes.NewBufferString(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request = ctx.Request.WithContext(authscope.WithContext(context.Background(), authscope.NewFromClaims(map[string]interface{}{"user_id": "user-1"}, nil)))
	return ctx
}

func TestBindJSON(t *testing.T) {
	var request dto.CreateTripRequest
	if !BindJSON(commonContext(http.MethodPost, `{"title":"Trip"}`), uuid.New(), &request) || request.Title != "Trip" {
		t.Fatal("valid JSON failed")
	}
	if BindJSON(commonContext(http.MethodPost, `{`), uuid.New(), &request) {
		t.Fatal("invalid JSON accepted")
	}
}

func TestHandleTripAccessError(t *testing.T) {
	for _, err := range []error{serviceshared.ErrNotMember, serviceshared.ErrAccessDenied} {
		ctx := commonContext(http.MethodGet, "")
		if !HandleTripAccessError(ctx, uuid.New(), err) || ctx.Writer.Status() != http.StatusForbidden {
			t.Fatalf("access err=%v status=%d", err, ctx.Writer.Status())
		}
	}
	ctx := commonContext(http.MethodGet, "")
	if HandleTripAccessError(ctx, uuid.New(), errors.New("other")) {
		t.Fatal("unexpected access mapping")
	}
}

func TestHandleJSONMutation(t *testing.T) {
	ctx := commonContext(http.MethodPost, `{"title":"Trip"}`)
	var request dto.CreateTripRequest
	HandleJSONMutation(ctx, JSONMutation[dto.CreateTripRequest, dto.TripDetailResponse]{
		ID: "trip-1", Request: &request, LogID: uuid.New(), StatusCode: http.StatusOK, Message: "ok",
		ServiceCall: func(_ context.Context, userID, id string, req dto.CreateTripRequest) (dto.TripDetailResponse, error) {
			if userID != "user-1" || id != "trip-1" || req.Title != "Trip" {
				t.Errorf("mutation args: %s %s %+v", userID, id, req)
			}
			return dto.TripDetailResponse{Id: id}, nil
		},
		HandleError: func(*gin.Context, uuid.UUID, string, error, string) { t.Fatal("unexpected error handler") },
	})
	if ctx.Writer.Status() != http.StatusOK {
		t.Fatalf("success status=%d", ctx.Writer.Status())
	}

	ctx = commonContext(http.MethodPost, `{"title":"Trip"}`)
	HandleJSONMutation(ctx, JSONMutation[dto.CreateTripRequest, dto.TripDetailResponse]{
		Request: &request, LogID: uuid.New(), ServiceCall: func(context.Context, string, string, dto.CreateTripRequest) (dto.TripDetailResponse, error) {
			return dto.TripDetailResponse{}, errors.New("failed")
		},
		HandleError: func(_ *gin.Context, _ uuid.UUID, _ string, _ error, _ string) {},
	})
	if ctx.Writer.Status() != http.StatusOK {
		t.Fatalf("error handler status=%d", ctx.Writer.Status())
	}
}
