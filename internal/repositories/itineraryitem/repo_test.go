package repositoryitineraryitem

import (
	"context"
	"testing"
	"time"

	domainitineraryitem "yourz-itinerary/internal/domain/itineraryitem"
	repositorygeneric "yourz-itinerary/internal/repositories/generic"

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
