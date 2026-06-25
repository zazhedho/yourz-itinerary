package repositorytripmember

import (
	"context"
	"testing"
	"time"

	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	repositorygeneric "yourz-itinerary/internal/repositories/generic"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newTripMemberMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestNewTripMemberRepoCompiles(t *testing.T) {
	db, _ := newTripMemberMockDB(t)
	r := NewTripMemberRepo(db)
	if r == nil {
		t.Fatal("NewTripMemberRepo returned nil")
	}
}

func TestTripMemberSoftDeletePersistsAuditFields(t *testing.T) {
	db, mock := newTripMemberMockDB(t)
	repo := &repo{GenericRepository: repositorygeneric.New[domaintripmember.TripMember](db)}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "trip_members" SET`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	now := time.Now().Add(-time.Hour)
	db.NowFunc = func() time.Time { return now }

	err := repo.SoftDelete(context.Background(), "member-1", "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
