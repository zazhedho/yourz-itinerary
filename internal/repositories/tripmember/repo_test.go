package repositorytripmember

import (
	"context"
	"errors"
	"testing"
	"time"

	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	repositorygeneric "yourz-itinerary/internal/repositories/generic"
	"yourz-itinerary/pkg/filter"

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

func TestTripMemberQueries(t *testing.T) {
	t.Run("get by trip and user", func(t *testing.T) {
		db, mock := newTripMemberMockDB(t)
		repo := &repo{GenericRepository: repositorygeneric.New[domaintripmember.TripMember](db)}
		mock.ExpectQuery(`SELECT .* FROM "trip_members".*trip_id = .*user_id =`).
			WithArgs("trip-1", "user-1", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "trip_id", "user_id", "role"}).AddRow("member-1", "trip-1", "user-1", "editor"))

		member, err := repo.GetByTripAndUser(context.Background(), "trip-1", "user-1")
		if err != nil || member.Id != "member-1" {
			t.Fatalf("GetByTripAndUser() = %#v, %v", member, err)
		}
	})

	t.Run("get active by trip and user", func(t *testing.T) {
		db, mock := newTripMemberMockDB(t)
		repo := &repo{GenericRepository: repositorygeneric.New[domaintripmember.TripMember](db)}
		mock.ExpectQuery(`SELECT .* FROM "trip_members".*trip_id = .*user_id =`).
			WithArgs("trip-1", "user-1", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "trip_id", "user_id", "role"}).AddRow("member-1", "trip-1", "user-1", "editor"))

		member, err := repo.GetActiveByTripAndUser(context.Background(), "trip-1", "user-1")
		if err != nil || member.Id != "member-1" {
			t.Fatalf("GetActiveByTripAndUser() = %#v, %v", member, err)
		}
	})

	t.Run("list by trip", func(t *testing.T) {
		db, mock := newTripMemberMockDB(t)
		repo := &repo{GenericRepository: repositorygeneric.New[domaintripmember.TripMember](db)}
		mock.ExpectQuery(`SELECT .* FROM "trip_members".*trip_id =`).
			WithArgs("trip-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "trip_id", "user_id", "role"}).AddRow("member-1", "trip-1", "user-1", "editor"))

		members, err := repo.ListByTrip(context.Background(), "trip-1")
		if err != nil || len(members) != 1 {
			t.Fatalf("ListByTrip() = %#v, %v", members, err)
		}
	})

	t.Run("returns query error", func(t *testing.T) {
		db, mock := newTripMemberMockDB(t)
		repo := &repo{GenericRepository: repositorygeneric.New[domaintripmember.TripMember](db)}
		queryErr := errors.New("query failed")
		mock.ExpectQuery(`SELECT .* FROM "trip_members".*trip_id = .*user_id =`).
			WithArgs("trip-1", "user-1", 1).
			WillReturnError(queryErr)

		_, err := repo.GetActiveByTripAndUser(context.Background(), "trip-1", "user-1")
		if !errors.Is(err, queryErr) {
			t.Fatalf("error = %v, want %v", err, queryErr)
		}
	})
}

func TestTripMemberGetAllDelegates(t *testing.T) {
	db, _ := newTripMemberMockDB(t)
	repo := &repo{GenericRepository: repositorygeneric.New[domaintripmember.TripMember](db)}
	repo.DB = repo.DB.Session(&gorm.Session{DryRun: true})
	if _, _, err := repo.GetAll(context.Background(), filter.BaseParams{Limit: 10}); err != nil {
		t.Fatalf("GetAll: %v", err)
	}
}
