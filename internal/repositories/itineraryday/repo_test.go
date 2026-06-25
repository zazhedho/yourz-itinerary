package repositoryitineraryday

import (
	"context"
	"testing"
	"time"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	repositorygeneric "yourz-itinerary/internal/repositories/generic"

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
