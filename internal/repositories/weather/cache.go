package repositoryweather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	domainweather "yourz-itinerary/internal/domain/weather"
	interfaceweather "yourz-itinerary/internal/interfaces/weather"

	"github.com/redis/go-redis/v9"
)

type WeatherCache struct {
	Redis *redis.Client
}

func NewWeatherCache(redisClient *redis.Client) interfaceweather.Cache {
	return &WeatherCache{Redis: redisClient}
}

func (c *WeatherCache) GetDaily(ctx context.Context, key string) (domainweather.Forecast, error) {
	raw, err := c.get(ctx, key)
	if err != nil {
		return domainweather.Forecast{}, err
	}
	var forecast domainweather.Forecast
	if err := json.Unmarshal([]byte(raw), &forecast); err != nil {
		return domainweather.Forecast{}, interfaceweather.ErrCacheMiss
	}
	if forecast.ForecastDate == "" || forecast.TimeZone == "" {
		return domainweather.Forecast{}, interfaceweather.ErrCacheMiss
	}
	return forecast, nil
}

func (c *WeatherCache) SetDaily(ctx context.Context, key string, forecast domainweather.Forecast, ttl time.Duration) error {
	return c.set(ctx, key, forecast, ttl)
}

func (c *WeatherCache) GetHourly(ctx context.Context, key string) ([]domainweather.HourlyForecast, error) {
	raw, err := c.get(ctx, key)
	if err != nil {
		return nil, err
	}
	var forecasts []domainweather.HourlyForecast
	if err := json.Unmarshal([]byte(raw), &forecasts); err != nil {
		return nil, interfaceweather.ErrCacheMiss
	}
	if len(forecasts) == 0 {
		return nil, interfaceweather.ErrCacheMiss
	}
	for _, forecast := range forecasts {
		if forecast.ForecastTime.IsZero() || forecast.TimeZone == "" {
			return nil, interfaceweather.ErrCacheMiss
		}
	}
	return forecasts, nil
}

func (c *WeatherCache) SetHourly(ctx context.Context, key string, forecasts []domainweather.HourlyForecast, ttl time.Duration) error {
	return c.set(ctx, key, forecasts, ttl)
}

func (c *WeatherCache) get(ctx context.Context, key string) (string, error) {
	if c == nil || c.Redis == nil {
		return "", interfaceweather.ErrCacheUnavailable
	}
	raw, err := c.Redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", interfaceweather.ErrCacheMiss
	}
	if err != nil {
		return "", fmt.Errorf("%w: %w", interfaceweather.ErrCacheUnavailable, err)
	}
	return raw, nil
}

func (c *WeatherCache) set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if c == nil || c.Redis == nil {
		return interfaceweather.ErrCacheUnavailable
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%w: marshal", interfaceweather.ErrCacheUnavailable)
	}
	if err := c.Redis.Set(ctx, key, string(payload), ttl).Err(); err != nil {
		return fmt.Errorf("%w: %w", interfaceweather.ErrCacheUnavailable, err)
	}
	return nil
}
