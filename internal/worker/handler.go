package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/myorg/worker-service/internal/domain"
	"github.com/myorg/worker-service/internal/service"
	"github.com/myorg/worker-service/pkg"
	"github.com/sirupsen/logrus"
)

// Handler processes a single Redis stream message.
type Handler struct {
	eventSvc    *service.EventService
	retryCfg    pkg.RetryConfig
	metrics     *pkg.Metrics
	log         *logrus.Logger
}

func NewHandler(
	eventSvc *service.EventService,
	retryCfg pkg.RetryConfig,
	metrics *pkg.Metrics,
	log *logrus.Logger,
) *Handler {
	return &Handler{
		eventSvc: eventSvc,
		retryCfg: retryCfg,
		metrics:  metrics,
		log:      log,
	}
}

func (h *Handler) Handle(ctx context.Context, msgID string, values map[string]interface{}) (ack bool) {
	start := time.Now()
	h.metrics.EventsReceived.Inc()

	event, err := parseEvent(values)
	if err != nil {
		h.log.WithError(err).WithField("msg_id", msgID).Error("failed to parse event; will ACK to avoid poison pill")
		h.metrics.EventsProcessed.WithLabelValues("failed").Inc()
		return true 
	}

	log := h.log.WithFields(logrus.Fields{
		"msg_id":     msgID,
		"event_id":   event.ID,
		"event_type": string(event.Type),
	})

	if err := event.Validate(); err != nil {
		log.WithError(err).Error("invalid event; ACKing to avoid poison pill")
		h.metrics.EventsProcessed.WithLabelValues("failed").Inc()
		return true
	}

	// ── Retry loop ─────────────────────────────────────────────────────────
	var lastErr error
	retryErr := pkg.DoWithRetry(ctx, h.retryCfg, func(attempt int) error {
		if attempt > 1 {
			log.WithField("attempt", attempt).Warn("retrying event")
			h.metrics.EventsRetried.Inc()
		}

		processErr := h.eventSvc.Process(ctx, event)
		if processErr == nil {
			return nil
		}
		if service.IsAlreadyProcessed(processErr) {
			lastErr = processErr
			return nil 
		}
		lastErr = processErr
		return processErr
	})

	elapsed := time.Since(start).Seconds()
	h.metrics.ProcessDuration.Observe(elapsed)

	if service.IsAlreadyProcessed(lastErr) {
		log.Info("event already processed; ACKing")
		h.metrics.EventsProcessed.WithLabelValues("skipped").Inc()
		return true
	}

	if retryErr != nil {
		log.WithError(retryErr).Error("event failed after all retries")
		h.metrics.EventsProcessed.WithLabelValues("failed").Inc()
		return false 
	}

	h.metrics.EventsProcessed.WithLabelValues("success").Inc()
	return true
}

func parseEvent(values map[string]interface{}) (*domain.Event, error) {
	raw, ok := values["data"]
	if !ok {
		return nil, fmt.Errorf("stream message missing 'data' field")
	}

	rawStr, ok := raw.(string)
	if !ok {
		return nil, fmt.Errorf("stream 'data' field is not a string")
	}

	var event domain.Event
	if err := json.Unmarshal([]byte(rawStr), &event); err != nil {
		return nil, fmt.Errorf("unmarshal event JSON: %w", err)
	}
	return &event, nil
}
