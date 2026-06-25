package handleritineraryitem

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"yourz-itinerary/internal/authscope"
	"yourz-itinerary/internal/dto"
	serviceitineraryitem "yourz-itinerary/internal/services/itineraryitem"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/utils"

	"github.com/gin-gonic/gin"
)

type itineraryItemServiceMock struct {
	item dto.ItineraryItemResponse
	err  error
}

func (m *itineraryItemServiceMock) CreateItem(ctx context.Context, userId, dayId string, req dto.CreateItineraryItemRequest) (dto.ItineraryItemResponse, error) {
	return m.item, m.err
}
func (m *itineraryItemServiceMock) UpdateItem(ctx context.Context, userId, itemId string, req dto.UpdateItineraryItemRequest) (dto.ItineraryItemResponse, error) {
	return m.item, m.err
}
func (m *itineraryItemServiceMock) DeleteItem(ctx context.Context, userId, itemId string) error {
	return m.err
}
func (m *itineraryItemServiceMock) ReorderItems(ctx context.Context, userId, dayId string, req dto.ReorderItineraryItemsRequest) error {
	return m.err
}

func performItineraryItemRequest(method, routePath, requestPath string, body interface{}, handler gin.HandlerFunc, authData map[string]interface{}) *httptest.ResponseRecorder {
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

func TestItineraryItemCreateItemSuccess(t *testing.T) {
	handler := NewItineraryItemHandler(&itineraryItemServiceMock{
		item: dto.ItineraryItemResponse{Id: "i-1", Title: "Visit Museum"},
	})
	rec := performItineraryItemRequest(http.MethodPost, "/api/days/:id/items", "/api/days/550e8400-e29b-41d4-a716-446655440000/items",
		dto.CreateItineraryItemRequest{Title: "Visit Museum"},
		handler.CreateItem,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestItineraryItemItemNotFound(t *testing.T) {
	handler := NewItineraryItemHandler(&itineraryItemServiceMock{
		err: serviceitineraryitem.ErrItemNotFound,
	})
	rec := performItineraryItemRequest(http.MethodPut, "/api/items/:id", "/api/items/550e8400-e29b-41d4-a716-446655440000",
		dto.UpdateItineraryItemRequest{Title: "Updated"},
		handler.UpdateItem,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestItineraryItemInvalidTime(t *testing.T) {
	handler := NewItineraryItemHandler(&itineraryItemServiceMock{
		err: serviceitineraryitem.ErrInvalidTime,
	})
	rec := performItineraryItemRequest(http.MethodPost, "/api/days/:id/items", "/api/days/550e8400-e29b-41d4-a716-446655440000/items",
		dto.CreateItineraryItemRequest{Title: "Dinner"},
		handler.CreateItem,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestItineraryItemReorderEmpty(t *testing.T) {
	handler := NewItineraryItemHandler(&itineraryItemServiceMock{
		err: serviceitineraryitem.ErrReorderEmpty,
	})
	rec := performItineraryItemRequest(http.MethodPut, "/api/days/:id/items/reorder", "/api/days/550e8400-e29b-41d4-a716-446655440000/items/reorder",
		dto.ReorderItineraryItemsRequest{ItemIds: []string{}},
		handler.ReorderItems,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestItineraryItemForbidden(t *testing.T) {
	handler := NewItineraryItemHandler(&itineraryItemServiceMock{
		err: serviceshared.ErrAccessDenied,
	})
	rec := performItineraryItemRequest(http.MethodDelete, "/api/items/:id", "/api/items/550e8400-e29b-41d4-a716-446655440000",
		nil, handler.DeleteItem,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", rec.Code, rec.Body.String())
	}
}
