package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewRoutesRegistersHealthcheck(t *testing.T) {
	routes := NewRoutes()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	routes.App.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestSessionRoutesSkipsWhenRedisUnavailable(t *testing.T) {
	routes := NewRoutes()
	routes.SessionRoutes()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/user/sessions", nil)
	routes.App.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected unregistered route to return 404, got %d", rec.Code)
	}
}

func newRouterDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()
	sqlDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}), &gorm.Config{
		DryRun:                 true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("open gorm: %v", err)
	}
	return db
}

func TestRouteGroupsRegisterWithDryRunDB(t *testing.T) {
	routes := NewRoutes()
	routes.DB = newRouterDryRunDB(t)

	routes.UserRoutes()
	routes.RoleRoutes()
	routes.PermissionRoutes()
	routes.MenuRoutes()
	routes.AppConfigRoutes()
	routes.AuditRoutes()
	routes.LocationRoutes()
	routes.TripRoutes()
	routes.TripMemberRoutes()
	routes.ItineraryDayRoutes()
	routes.ItineraryItemRoutes()
	routes.WeatherRoutes()

	registered := map[string]bool{}
	for _, route := range routes.App.Routes() {
		registered[route.Method+" "+route.Path] = true
	}

	for _, want := range []string{
		"POST /api/user/register",
		"POST /api/user/login",
		"GET /api/roles",
		"GET /api/permissions",
		"GET /api/menus",
		"GET /api/configs",
		"GET /api/audits",
		"GET /api/location/province",
		"POST /api/location/sync",
		"DELETE /api/trips/:id/leave",
		"DELETE /api/trips/:id/members/:member_id",
		"GET /api/itinerary-days/:id/weather",
	} {
		if !registered[want] {
			t.Fatalf("expected route %s to be registered", want)
		}
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/itinerary-days/550e8400-e29b-41d4-a716-446655440000/weather", nil)
	routes.App.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected weather route to require auth, got %d: %s", rec.Code, rec.Body.String())
	}
}
