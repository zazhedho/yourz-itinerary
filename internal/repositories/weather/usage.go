package repositoryweather

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	interfaceweather "yourz-itinerary/internal/interfaces/weather"
	"yourz-itinerary/pkg/logger"
)

var ErrUsageUnavailable = errors.New("weather usage counter unavailable")

type UsageRepository struct {
	Redis *redis.Client
	Limit int
	Now   func() time.Time
}

func NewUsageRepository(redisClient *redis.Client, limit int) interfaceweather.Usage {
	return &UsageRepository{Redis: redisClient, Limit: limit}
}

func (r *UsageRepository) Reserve(ctx context.Context) (bool, int64, error) {
	if r == nil || r.Limit < 1 {
		return false, 0, nil
	}
	if r.Redis == nil {
		return false, 0, ErrUsageUnavailable
	}

	now := time.Now().UTC()
	if r.Now != nil {
		now = r.Now().UTC()
	}
	key := fmt.Sprintf("weather:usage:%s", now.Format("2006-01"))
	count, err := r.Redis.Incr(ctx, key).Result()
	if err != nil {
		return false, 0, ErrUsageUnavailable
	}
	if count == 1 {
		expiresAt := time.Date(now.Year(), now.Month()+2, 1, 0, 0, 0, 0, time.UTC)
		if err := r.Redis.ExpireAt(ctx, key, expiresAt).Err(); err != nil {
			return false, count, ErrUsageUnavailable
		}
	}
	allowed := count <= int64(r.Limit)
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("[WeatherUsage] month=%s count=%d allowed=%t", now.Format("2006-01"), count, allowed))
	return allowed, count, nil
}
