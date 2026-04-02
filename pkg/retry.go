package pkg

import (
	"context"
	"math"
	"time"
)

type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

func DefaultRetryConfig(maxAttempts int, baseDelay time.Duration) RetryConfig {
	return RetryConfig{
		MaxAttempts: maxAttempts,
		BaseDelay:   baseDelay,
		MaxDelay:    30 * time.Second,
	}
}

// DoWithRetry executes fn with exponential backoff.
// It stops on context cancellation or when fn succeeds.
// The callback receives the current attempt number (1-based).
func DoWithRetry(ctx context.Context, cfg RetryConfig, fn func(attempt int) error) error {
	var lastErr error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = fn(attempt)
		if lastErr == nil {
			return nil
		}

		if attempt == cfg.MaxAttempts {
			break
		}

		delay := exponentialDelay(cfg.BaseDelay, cfg.MaxDelay, attempt)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return lastErr
}

// exponentialDelay computes delay = base * 2^(attempt-1), capped at max.
func exponentialDelay(base, max time.Duration, attempt int) time.Duration {
	factor := math.Pow(2, float64(attempt-1))
	delay := time.Duration(float64(base) * factor)
	if delay > max {
		return max
	}
	return delay
}
