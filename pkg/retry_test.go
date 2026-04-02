package pkg_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/myorg/worker-service/pkg"
)

var errTransient = errors.New("transient error")

func TestDoWithRetry_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := pkg.DoWithRetry(context.Background(), pkg.RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   time.Millisecond,
		MaxDelay:    time.Second,
	}, func(attempt int) error {
		calls++
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Errorf("fn called %d times, want 1", calls)
	}
}

func TestDoWithRetry_SuccessAfterRetries(t *testing.T) {
	calls := 0
	err := pkg.DoWithRetry(context.Background(), pkg.RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   time.Millisecond,
		MaxDelay:    time.Second,
	}, func(attempt int) error {
		calls++
		if calls < 3 {
			return errTransient
		}
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Errorf("fn called %d times, want 3", calls)
	}
}

func TestDoWithRetry_ExhaustsMaxAttempts(t *testing.T) {
	calls := 0
	cfg := pkg.RetryConfig{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: time.Second}

	err := pkg.DoWithRetry(context.Background(), cfg, func(attempt int) error {
		calls++
		return errTransient
	})

	if !errors.Is(err, errTransient) {
		t.Errorf("expected errTransient, got: %v", err)
	}
	if calls != 3 {
		t.Errorf("fn called %d times, want 3", calls)
	}
}

func TestDoWithRetry_RespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	calls := 0
	err := pkg.DoWithRetry(ctx, pkg.RetryConfig{
		MaxAttempts: 10,
		BaseDelay:   50 * time.Millisecond,
		MaxDelay:    time.Second,
	}, func(attempt int) error {
		calls++
		if calls == 2 {
			cancel() // cancel after second attempt
		}
		return errTransient
	})

	if err == nil {
		t.Fatal("expected an error due to context cancellation")
	}
	// Should not have reached all 10 attempts.
	if calls >= 10 {
		t.Errorf("fn called %d times; should have been cut short by cancellation", calls)
	}
}

func TestDoWithRetry_SingleAttempt(t *testing.T) {
	calls := 0
	err := pkg.DoWithRetry(context.Background(), pkg.RetryConfig{
		MaxAttempts: 1,
		BaseDelay:   time.Millisecond,
		MaxDelay:    time.Second,
	}, func(attempt int) error {
		calls++
		return errTransient
	})

	if !errors.Is(err, errTransient) {
		t.Errorf("expected errTransient, got: %v", err)
	}
	if calls != 1 {
		t.Errorf("fn called %d times, want 1", calls)
	}
}
