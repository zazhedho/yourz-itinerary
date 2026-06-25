package handlertripmember

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"yourz-itinerary/internal/authscope"
	"yourz-itinerary/internal/dto"
	serviceshared "yourz-itinerary/internal/services/shared"
	servicetripmember "yourz-itinerary/internal/services/tripmember"
	"yourz-itinerary/utils"

	"github.com/gin-gonic/gin"
)

type tripMemberServiceMock struct {
	member dto.TripMemberResponse
	err    error
}

func (m *tripMemberServiceMock) AddMember(ctx context.Context, userId, tripId string, req dto.AddTripMemberRequest) (dto.TripMemberResponse, error) {
	return m.member, m.err
}
func (m *tripMemberServiceMock) UpdateMemberRole(ctx context.Context, userId, tripId, memberId string, req dto.UpdateTripMemberRoleRequest) (dto.TripMemberResponse, error) {
	return m.member, m.err
}
func (m *tripMemberServiceMock) RemoveMember(ctx context.Context, userId, tripId, memberId string) error {
	return m.err
}
func (m *tripMemberServiceMock) LeaveTrip(ctx context.Context, userId, tripId string) error {
	return m.err
}

func performTripMemberRequest(method, routePath, requestPath string, body interface{}, handler gin.HandlerFunc, authData map[string]interface{}) *httptest.ResponseRecorder {
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

func TestTripMemberAddMemberSuccess(t *testing.T) {
	handler := NewTripMemberHandler(&tripMemberServiceMock{
		member: dto.TripMemberResponse{Id: "m-1", UserId: "user-b", Role: "viewer"},
	})
	rec := performTripMemberRequest(http.MethodPost, "/api/trips/:id/members", "/api/trips/550e8400-e29b-41d4-a716-446655440000/members",
		dto.AddTripMemberRequest{Email: "b@test.com", Role: "viewer"},
		handler.AddMember,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestTripMemberAddMemberForbidden(t *testing.T) {
	handler := NewTripMemberHandler(&tripMemberServiceMock{
		err: serviceshared.ErrAccessDenied,
	})
	rec := performTripMemberRequest(http.MethodPost, "/api/trips/:id/members", "/api/trips/550e8400-e29b-41d4-a716-446655440000/members",
		dto.AddTripMemberRequest{Email: "b@test.com", Role: "viewer"},
		handler.AddMember,
		authClaims("u-1", "bob", "member"),
	)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestTripMemberOwnerLeaveBadRequest(t *testing.T) {
	handler := NewTripMemberHandler(&tripMemberServiceMock{
		err: servicetripmember.ErrOwnerLeave,
	})
	rec := performTripMemberRequest(http.MethodDelete, "/api/trips/:id/leave", "/api/trips/550e8400-e29b-41d4-a716-446655440000/leave",
		nil, handler.LeaveTrip,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestTripMemberDuplicateConflict(t *testing.T) {
	handler := NewTripMemberHandler(&tripMemberServiceMock{
		err: servicetripmember.ErrDuplicateMember,
	})
	rec := performTripMemberRequest(http.MethodPost, "/api/trips/:id/members", "/api/trips/550e8400-e29b-41d4-a716-446655440000/members",
		dto.AddTripMemberRequest{Email: "b@test.com", Role: "viewer"},
		handler.AddMember,
		authClaims("u-1", "alice", "member"),
	)
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
}
