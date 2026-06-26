package repositorytrip

import (
	"context"
	"testing"
	"time"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	domaintrip "yourz-itinerary/internal/domain/trip"
	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	repositorygeneric "yourz-itinerary/internal/repositories/generic"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newTripMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("open gorm: %v", err)
	}

	return db, mock
}

func TestNewTripRepoCompiles(t *testing.T) {
	db, _ := newTripMockDB(t)
	r := NewTripRepo(db)
	if r == nil {
		t.Fatal("NewTripRepo returned nil")
	}
}

func TestTripSoftDeleteUsesProvidedTripID(t *testing.T) {
	db, mock := newTripMockDB(t)
	repo := &repo{GenericRepository: repositorygeneric.New[domaintrip.Trip](db)}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "trips" SET`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	now := time.Now().Add(-time.Hour)
	db.NowFunc = func() time.Time { return now }

	deletedBy := "user-1"
	err := repo.SoftDelete(context.Background(), "trip-1", deletedBy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestTripCreateTripOpensTransaction(t *testing.T) {
	db, mock := newTripMockDB(t)
	repo := &repo{GenericRepository: repositorygeneric.New[domaintrip.Trip](db)}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "trips"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO "trip_members"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	_, err := repo.CreateTrip(context.Background(),
		domaintrip.Trip{Id: "trip-1", OwnerId: "user-1", Title: "Test"},
		domaintripmember.TripMember{Id: "member-1", UserId: "user-1", Role: "owner"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestTripCreateTripPersistsGeneratedDaysInTransaction(t *testing.T) {
	db, mock := newTripMockDB(t)
	repo := &repo{GenericRepository: repositorygeneric.New[domaintrip.Trip](db)}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "trips"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO "trip_members"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO "itinerary_days"`).
		WillReturnResult(sqlmock.NewResult(1, 2))
	mock.ExpectCommit()

	_, err := repo.CreateTrip(context.Background(),
		domaintrip.Trip{Id: "trip-1", OwnerId: "user-1", Title: "Test"},
		domaintripmember.TripMember{Id: "member-1", UserId: "user-1", Role: "owner"},
		domainitineraryday.ItineraryDay{Id: "day-1", DayNumber: 1},
		domainitineraryday.ItineraryDay{Id: "day-2", DayNumber: 2},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
