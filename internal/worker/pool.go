package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	redisinfra "github.com/myorg/worker-service/internal/infra/redis"
	"github.com/myorg/worker-service/internal/config"
	"github.com/myorg/worker-service/pkg"
	"github.com/sirupsen/logrus"
)

// Pool is the Redis Streams consumer group worker.
type Pool struct {
	client   *redis.Client
	handler  *Handler
	cfg      config.WorkerConfig
	redisCfg config.RedisConfig
	metrics  *pkg.Metrics
	log      *logrus.Logger
}

func NewPool(
	client *redis.Client,
	handler *Handler,
	cfg config.WorkerConfig,
	redisCfg config.RedisConfig,
	metrics *pkg.Metrics,
	log *logrus.Logger,
) *Pool {
	return &Pool{
		client:   client,
		handler:  handler,
		cfg:      cfg,
		redisCfg: redisCfg,
		metrics:  metrics,
		log:      log,
	}
}

func (p *Pool) Run(ctx context.Context) {
	p.log.WithFields(logrus.Fields{
		"stream":      p.redisCfg.StreamName,
		"group":       p.redisCfg.GroupName,
		"consumer":    p.cfg.ConsumerName,
		"concurrency": p.cfg.Concurrency,
	}).Info("starting worker pool")

	sem := make(chan struct{}, p.cfg.Concurrency)
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			p.log.Info("context cancelled; waiting for in-flight jobs")
			wg.Wait()
			p.log.Info("all in-flight jobs finished; worker pool stopped")
			return
		default:
		}

		messages, err := p.readMessages(ctx)
		if err != nil {
			if ctx.Err() != nil {
				// Context cancelled during read — clean shutdown.
				wg.Wait()
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
			// Acquire semaphore slot.
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				wg.Wait()
				return
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				defer func() { <-sem }()
				p.processMessage(ctx, msg)
			}()
		}
	}
}

// readMessages calls XREADGROUP with a blocking timeout equal to PollInterval.
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

// processMessage runs the handler for one message and ACKs or DLQs accordingly.
func (p *Pool) processMessage(ctx context.Context, msg redis.XMessage) {
	log := p.log.WithField("msg_id", msg.ID)

	ack := p.handler.Handle(ctx, msg.ID, msg.Values)

	if ack {
		if err := p.ack(ctx, msg.ID); err != nil {
			log.WithError(err).Error("failed to ACK message")
		}
		return
	}

	// Handler returned false → push to DLQ.
	payload, _ := json.Marshal(msg.Values)
	if err := redisinfra.PushToDLQ(ctx, p.client, p.redisCfg.DLQStreamName(), msg.ID, string(payload)); err != nil {
		log.WithError(err).Error("failed to push message to DLQ")
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
