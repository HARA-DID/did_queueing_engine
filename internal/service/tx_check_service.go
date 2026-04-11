package service

import (
	"context"
	"time"

	"github.com/HARA-DID/did-queueing-engine/internal/callback"
	"github.com/HARA-DID/did-queueing-engine/internal/domain"
	"github.com/HARA-DID/did-queueing-engine/internal/repository"
	"github.com/HARA-DID/hara-core-blockchain-lib/pkg/blockchain"
	"github.com/sirupsen/logrus"
)

type TxCheckTask struct {
	Event *domain.Event
	Job   *domain.Job
	Hash  string
}

type TxCheckService struct {
	jobRepo        repository.JobRepository
	bcClient       *blockchain.Blockchain
	callbacks      *callback.Registry
	eventCallbacks map[domain.EventType]callback.Func
	log            *logrus.Logger
	tasks          chan TxCheckTask
}

func NewTxCheckService(
	jobRepo repository.JobRepository,
	bcClient *blockchain.Blockchain,
	callbacks *callback.Registry,
	eventCallbacks map[domain.EventType]callback.Func,
	log *logrus.Logger,
	bufferSize int,
) *TxCheckService {
	return &TxCheckService{
		jobRepo:        jobRepo,
		bcClient:       bcClient,
		callbacks:      callbacks,
		eventCallbacks: eventCallbacks,
		log:            log,
		tasks:          make(chan TxCheckTask, bufferSize),
	}
}

func (s *TxCheckService) Enqueue(task TxCheckTask) {
	s.tasks <- task
}

func (s *TxCheckService) Start(ctx context.Context) {
	s.log.Info("TxCheckService started")
	for {
		select {
		case <-ctx.Done():
			s.log.Info("TxCheckService stopping")
			return
		case task := <-s.tasks:
			batch := []TxCheckTask{task}

			drainLoop := true
			for drainLoop {
				select {
				case next := <-s.tasks:
					batch = append(batch, next)
				default:
					drainLoop = false
				}
			}

			s.processBatch(ctx, batch)
		}
	}
}

func (s *TxCheckService) processBatch(ctx context.Context, batch []TxCheckTask) {
	if len(batch) == 0 {
		return
	}

	hashes := make([]string, 0, len(batch))
	for _, t := range batch {
		hashes = append(hashes, t.Hash)
	}

	s.log.WithField("batch_size", len(batch)).Info("checking transaction batch")

	// TODO change the behavior in the future (currently only logs errors)
	_, errs := s.bcClient.CheckTxs(ctx, hashes)
	if len(errs) > 0 {
		s.log.WithField("check_txs_errors", errs).Warn("some transactions failed check")
	}

	// Finalize jobs and trigger callbacks
	for _, t := range batch {
		if err := s.jobRepo.UpdateStatus(ctx, t.Job.ID, domain.JobStatusSuccess, []string{t.Hash}, ""); err != nil {
			s.log.WithError(err).WithField("job_id", t.Job.ID).Error("failed to update job status to success")
		}

		s.triggerCallback(ctx, t.Event, callback.Result{
			EventID:   t.Event.ID,
			JobID:     t.Job.ID,
			EventType: string(t.Event.Type),
			Success:   true,
			TxHashes:  []string{t.Hash},
		})

		s.log.WithFields(logrus.Fields{
			"job_id":   t.Job.ID,
			"tx_hash":  t.Hash,
			"event_id": t.Event.ID,
		}).Info("event processed successfully in background")
	}
}

func (s *TxCheckService) triggerCallback(ctx context.Context, event *domain.Event, result callback.Result) {
	cb, ok := s.eventCallbacks[event.Type]
	if !ok {
		return
	}

	execCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := cb(execCtx, result); err != nil {
		s.log.WithFields(logrus.Fields{
			"event_id":   event.ID,
			"event_type": event.Type,
		}).WithError(err).Error("callback execution failed")
	}
}
