package repositoryweather

import (
	"context"
	"errors"
	"testing"
	"time"

	redismock "github.com/go-redis/redismock/v9"
)

func TestNewUsageRepository(t *testing.T) {
	if NewUsageRepository(nil, 10) == nil {
		t.Fatal("NewUsageRepository returned nil")
	}
}

func TestUsageRepositoryReservesThroughLimit(t *testing.T) {
	client, mock := redismock.NewClientMock()
	now := time.Date(2026, time.August, 3, 12, 0, 0, 0, time.UTC)
	repo := &UsageRepository{Redis: client, Limit: 2, Now: func() time.Time { return now }}
	key := "weather:usage:2026-08"
	mock.ExpectIncr(key).SetVal(1)
	mock.ExpectExpireAt(key, time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC)).SetVal(true)
	mock.ExpectIncr(key).SetVal(2)
	mock.ExpectIncr(key).SetVal(3)

	for i, wantAllowed := range []bool{true, true, false} {
		allowed, count, err := repo.Reserve(context.Background())
		if err != nil || allowed != wantAllowed || count != int64(i+1) {
			t.Fatalf("reserve %d: allowed=%v count=%d err=%v", i, allowed, count, err)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations: %v", err)
	}
}

func TestUsageRepositoryUsesNewMonthKey(t *testing.T) {
	client, mock := redismock.NewClientMock()
	repo := &UsageRepository{Redis: client, Limit: 1, Now: func() time.Time {
		return time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)
	}}
	mock.ExpectIncr("weather:usage:2026-09").SetVal(1)
	mock.ExpectExpireAt("weather:usage:2026-09", time.Date(2026, time.November, 1, 0, 0, 0, 0, time.UTC)).SetVal(true)

	allowed, _, err := repo.Reserve(context.Background())
	if err != nil || !allowed {
		t.Fatalf("expected new month reservation, allowed=%v err=%v", allowed, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations: %v", err)
	}
}

func TestUsageRepositoryFailsClosed(t *testing.T) {
	client, mock := redismock.NewClientMock()
	repo := &UsageRepository{Redis: client, Limit: 1, Now: func() time.Time { return time.Now() }}
	mock.ExpectIncr("weather:usage:" + time.Now().UTC().Format("2006-01")).SetErr(errors.New("redis down"))
	allowed, _, err := repo.Reserve(context.Background())
	if allowed || !errors.Is(err, ErrUsageUnavailable) {
		t.Fatalf("expected closed failure, allowed=%v err=%v", allowed, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations: %v", err)
	}

	allowed, _, err = (&UsageRepository{Limit: 1}).Reserve(context.Background())
	if allowed || !errors.Is(err, ErrUsageUnavailable) {
		t.Fatalf("expected nil redis failure, allowed=%v err=%v", allowed, err)
	}
}

func TestUsageRepositoryDisablesZeroLimit(t *testing.T) {
	allowed, count, err := (&UsageRepository{Limit: 0}).Reserve(context.Background())
	if allowed || count != 0 || err != nil {
		t.Fatalf("expected disabled usage, allowed=%v count=%d err=%v", allowed, count, err)
	}
}
