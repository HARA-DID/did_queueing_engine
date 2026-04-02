package redisinfra

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/myorg/worker-service/internal/config"
)

// NewClient creates a Redis client from config.
func NewClient(cfg config.RedisConfig) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return client, nil
}

func EnsureConsumerGroup(ctx context.Context, client *redis.Client, stream, group string) error {
	err := client.XGroupCreateMkStream(ctx, stream, group, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return fmt.Errorf("create consumer group %q on stream %q: %w", group, stream, err)
	}
	return nil
}

func PushToDLQ(ctx context.Context, client *redis.Client, dlqStream, eventID, payload string) error {
	return client.XAdd(ctx, &redis.XAddArgs{
		Stream: dlqStream,
		Values: map[string]interface{}{
			"event_id": eventID,
			"payload":  payload,
		},
	}).Err()
}
