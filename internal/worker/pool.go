package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redisinfra "github.com/HARA-DID/did-queueing-engine/internal/infra/redis"

	"github.com/HARA-DID/did-queueing-engine/internal/config"
	"github.com/HARA-DID/did-queueing-engine/internal/domain"
	"github.com/HARA-DID/did-queueing-engine/internal/repository"
	"github.com/HARA-DID/did-queueing-engine/pkg"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// Pool is the Redis Streams consumer group worker.
type Pool struct {
	client   *redis.Client
	handler  *Handler
	jobRepo  repository.JobRepository
	cfg      config.WorkerConfig
	redisCfg config.RedisConfig
	metrics  *pkg.Metrics
	log      *logrus.Logger
}

func NewPool(
	client *redis.Client,
	handler *Handler,
	jobRepo repository.JobRepository,
	cfg config.WorkerConfig,
	redisCfg config.RedisConfig,
	metrics *pkg.Metrics,
	log *logrus.Logger,
) *Pool {
	return &Pool{
		client:   client,
		handler:  handler,
		jobRepo:  jobRepo,
		cfg:      cfg,
		redisCfg: redisCfg,
		metrics:  metrics,
		log:      log,
	}
}

func (p *Pool) Run(ctx context.Context) {
	p.log.WithFields(logrus.Fields{
		"stream":   p.redisCfg.StreamName,
		"group":    p.redisCfg.GroupName,
		"consumer": p.cfg.ConsumerName,
	}).Info("starting worker processing synchronously")

	for {
		select {
		case <-ctx.Done():
			p.log.Info("context cancelled; worker stopped")
			return
		default:
		}

		messages, err := p.readMessages(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			p.log.WithError(err).Error("XREADGROUP error; backing off")
			sleep(ctx, p.cfg.PollInterval*5)
			continue
		}

		if len(messages) == 0 {
			sleep(ctx, p.cfg.PollInterval)
			continue
		}

		for _, msg := range messages {
			p.processMessage(ctx, msg)
		}
	}
}

func (p *Pool) readMessages(ctx context.Context) ([]redis.XMessage, error) {
	streams, err := p.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    p.redisCfg.GroupName,
		Consumer: p.cfg.ConsumerName,
		Streams:  []string{p.redisCfg.StreamName, ">"},
		Count:    p.cfg.BatchSize,
		Block:    p.cfg.PollInterval,
		NoAck:    false,
	}).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	if len(streams) == 0 {
		return nil, nil
	}
	return streams[0].Messages, nil
}

func (p *Pool) processMessage(ctx context.Context, msg redis.XMessage) {
	log := p.log.WithField("msg_id", msg.ID)

	ack := p.handler.Handle(ctx, msg.ID, msg.Values)

	if ack {
		if err := p.ack(ctx, msg.ID); err != nil {
			log.WithError(err).Error("failed to ACK message")
		}
		return
	}

	// Handler returned false → push to DLQ (Redis + DB).
	payload, _ := json.Marshal(msg.Values)
	rawPayload := string(payload)

	if event, err := ParseEvent(msg.Values); err == nil {
		dlqErr := p.jobRepo.SaveToDLQ(ctx, &domain.DLQEvent{
			EventID:      event.ID,
			Type:         string(event.Type),
			Payload:      rawPayload,
			ErrorMessage: "failed after all retries",
		})
		if dlqErr != nil {
			log.WithError(dlqErr).Error("failed to save to DB DLQ")
		}
	}

	if err := redisinfra.PushToDLQ(ctx, p.client, p.redisCfg.DLQStreamName(), msg.ID, rawPayload); err != nil {
		log.WithError(err).Error("failed to push message to Redis DLQ")
	} else {
		p.metrics.EventsDLQ.Inc()
		log.Warn("message pushed to DLQ")
	}

	// ACK the original message so it doesn't stay in the PEL indefinitely.
	if err := p.ack(ctx, msg.ID); err != nil {
		log.WithError(err).Error("failed to ACK after DLQ push")
	}
}

// ack acknowledges a message in the consumer group.
func (p *Pool) ack(ctx context.Context, msgID string) error {
	if err := p.client.XAck(ctx, p.redisCfg.StreamName, p.redisCfg.GroupName, msgID).Err(); err != nil {
		return fmt.Errorf("XACK %s: %w", msgID, err)
	}
	return nil
}

func sleep(ctx context.Context, d time.Duration) {
	select {
	case <-time.After(d):
	case <-ctx.Done():
	}
}
