package repositorytrip

import (
	"context"
	"errors"
	"testing"
	"time"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	domaintrip "yourz-itinerary/internal/domain/trip"
	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	repositorygeneric "yourz-itinerary/internal/repositories/generic"
	"yourz-itinerary/pkg/filter"

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

func TestTripListByMember(t *testing.T) {
	db, mock := newTripMockDB(t)
	repo := &repo{GenericRepository: repositorygeneric.New[domaintrip.Trip](db)}
	mock.ExpectQuery(`SELECT count\(\*\) FROM "trips".*trip_id.*user_id =`).
		WithArgs("user-1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT .* FROM "trips".*trip_id.*user_id =`).
		WithArgs("user-1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner_id", "title"}).AddRow("trip-1", "user-1", "Test"))

	trips, total, err := repo.ListByMember(context.Background(), "user-1")
	if err != nil || total != 1 || len(trips) != 1 || trips[0].Id != "trip-1" {
		t.Fatalf("ListByMember() = %#v, %d, %v", trips, total, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestTripListByMemberReturnsCountError(t *testing.T) {
	db, mock := newTripMockDB(t)
	repo := &repo{GenericRepository: repositorygeneric.New[domaintrip.Trip](db)}
	queryErr := errors.New("count failed")
	mock.ExpectQuery(`SELECT count\(\*\) FROM "trips".*trip_id.*user_id =`).
		WithArgs("user-1").
		WillReturnError(queryErr)

	trips, total, err := repo.ListByMember(context.Background(), "user-1")
	if !errors.Is(err, queryErr) || trips != nil || total != 0 {
		t.Fatalf("ListByMember() = %#v, %d, %v", trips, total, err)
	}
}

func TestTripGetAllDelegates(t *testing.T) {
	db, _ := newTripMockDB(t)
	repo := &repo{GenericRepository: repositorygeneric.New[domaintrip.Trip](db)}
	repo.DB = repo.DB.Session(&gorm.Session{DryRun: true})
	if _, _, err := repo.GetAll(context.Background(), filter.BaseParams{Limit: 10}); err != nil {
		t.Fatalf("GetAll: %v", err)
	}
}
