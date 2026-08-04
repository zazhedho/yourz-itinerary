package repositoryitineraryitem

import (
	"context"
	"errors"
	"testing"
	"time"

	domainitineraryitem "yourz-itinerary/internal/domain/itineraryitem"
	repositorygeneric "yourz-itinerary/internal/repositories/generic"
	"yourz-itinerary/pkg/filter"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newItineraryItemMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestNewItineraryItemRepoCompiles(t *testing.T) {
	db, _ := newItineraryItemMockDB(t)
	r := NewItineraryItemRepo(db)
	if r == nil {
		t.Fatal("NewItineraryItemRepo returned nil")
	}
}

func TestReorderPersistsAuditFields(t *testing.T) {
	db, mock := newItineraryItemMockDB(t)
	repo := &repo{GenericRepository: repositorygeneric.New[domainitineraryitem.ItineraryItem](db)}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "itinerary_items" SET "sort_order"=.*WHERE \(id = .* AND day_id =`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`UPDATE "itinerary_items" SET .*"updated_at".*"updated_by".*WHERE id =`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	now := time.Now().Add(-time.Hour)
	db.NowFunc = func() time.Time { return now }

	err := repo.Reorder(context.Background(), "day-1", []domainitineraryitem.ItineraryItem{
		{Id: "item-1", DayId: "day-1", SortOrder: 1, UpdatedBy: "user-editor"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestItineraryItemQueries(t *testing.T) {
	t.Run("get by day", func(t *testing.T) {
		db, mock := newItineraryItemMockDB(t)
		repo := &repo{GenericRepository: repositorygeneric.New[domainitineraryitem.ItineraryItem](db)}
		mock.ExpectQuery(`SELECT .* FROM "itinerary_items".*day_id =`).
			WithArgs("day-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "day_id", "title"}).AddRow("item-1", "day-1", "Museum"))

		items, err := repo.GetByDay(context.Background(), "day-1")
		if err != nil || len(items) != 1 || items[0].Id != "item-1" {
			t.Fatalf("GetByDay() = %#v, %v", items, err)
		}
	})

	t.Run("get by ids", func(t *testing.T) {
		db, mock := newItineraryItemMockDB(t)
		repo := &repo{GenericRepository: repositorygeneric.New[domainitineraryitem.ItineraryItem](db)}
		mock.ExpectQuery(`SELECT .* FROM "itinerary_items".*id IN`).
			WithArgs("item-1", "item-2").
			WillReturnRows(sqlmock.NewRows([]string{"id", "day_id"}).AddRow("item-1", "day-1"))

		items, err := repo.GetByIDs(context.Background(), []string{"item-1", "item-2"})
		if err != nil || len(items) != 1 {
			t.Fatalf("GetByIDs() = %#v, %v", items, err)
		}
	})

	t.Run("returns query error", func(t *testing.T) {
		db, mock := newItineraryItemMockDB(t)
		repo := &repo{GenericRepository: repositorygeneric.New[domainitineraryitem.ItineraryItem](db)}
		queryErr := errors.New("query failed")
		mock.ExpectQuery(`SELECT .* FROM "itinerary_items".*day_id =`).
			WithArgs("day-1").
			WillReturnError(queryErr)

		_, err := repo.GetByDay(context.Background(), "day-1")
		if !errors.Is(err, queryErr) {
			t.Fatalf("error = %v, want %v", err, queryErr)
		}
	})
}

func TestItineraryItemGetAllDelegates(t *testing.T) {
	db, _ := newItineraryItemMockDB(t)
	repo := &repo{GenericRepository: repositorygeneric.New[domainitineraryitem.ItineraryItem](db)}
	repo.DB = repo.DB.Session(&gorm.Session{DryRun: true})
	if _, _, err := repo.GetAll(context.Background(), filter.BaseParams{Limit: 10}); err != nil {
		t.Fatalf("GetAll: %v", err)
	}
}
