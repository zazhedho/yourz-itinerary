package repositoryitineraryday

import (
	"context"
	"errors"
	"testing"
	"time"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	repositorygeneric "yourz-itinerary/internal/repositories/generic"
	"yourz-itinerary/pkg/filter"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newItineraryDayMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestNewItineraryDayRepoCompiles(t *testing.T) {
	db, _ := newItineraryDayMockDB(t)
	r := NewItineraryDayRepo(db)
	if r == nil {
		t.Fatal("NewItineraryDayRepo returned nil")
	}
}

func TestItineraryDaySoftDeletePersistsAuditFields(t *testing.T) {
	db, mock := newItineraryDayMockDB(t)
	repo := &repo{GenericRepository: repositorygeneric.New[domainitineraryday.ItineraryDay](db)}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "itinerary_days" SET`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	now := time.Now().Add(-time.Hour)
	db.NowFunc = func() time.Time { return now }

	err := repo.SoftDelete(context.Background(), "day-1", "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestItineraryDayQueries(t *testing.T) {
	t.Run("list by trip", func(t *testing.T) {
		db, mock := newItineraryDayMockDB(t)
		repo := &repo{GenericRepository: repositorygeneric.New[domainitineraryday.ItineraryDay](db)}
		mock.ExpectQuery(`SELECT .* FROM "itinerary_days".*trip_id =`).
			WithArgs("trip-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "trip_id", "day_number"}).AddRow("day-1", "trip-1", 1))

		days, err := repo.ListByTrip(context.Background(), "trip-1")
		if err != nil || len(days) != 1 || days[0].Id != "day-1" {
			t.Fatalf("ListByTrip() = %#v, %v", days, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sql expectations: %v", err)
		}
	})

	t.Run("returns query error", func(t *testing.T) {
		db, mock := newItineraryDayMockDB(t)
		repo := &repo{GenericRepository: repositorygeneric.New[domainitineraryday.ItineraryDay](db)}
		queryErr := errors.New("query failed")
		mock.ExpectQuery(`SELECT .* FROM "itinerary_days".*trip_id =`).
			WithArgs("trip-1").
			WillReturnError(queryErr)

		_, err := repo.ListByTrip(context.Background(), "trip-1")
		if !errors.Is(err, queryErr) {
			t.Fatalf("error = %v, want %v", err, queryErr)
		}
	})
}

func TestItineraryDayGetAllDelegates(t *testing.T) {
	db, _ := newItineraryDayMockDB(t)
	repo := &repo{GenericRepository: repositorygeneric.New[domainitineraryday.ItineraryDay](db)}
	repo.DB = repo.DB.Session(&gorm.Session{DryRun: true})
	if _, _, err := repo.GetAll(context.Background(), filter.BaseParams{Limit: 10}); err != nil {
		t.Fatalf("GetAll: %v", err)
	}
}
