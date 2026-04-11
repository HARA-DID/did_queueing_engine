package worker_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/HARA-DID/did-queueing-engine/internal/callback"
	"github.com/HARA-DID/did-queueing-engine/internal/domain"
	"github.com/HARA-DID/did-queueing-engine/internal/mocks"
	"github.com/HARA-DID/did-queueing-engine/internal/service"
	"github.com/HARA-DID/did-queueing-engine/internal/worker"
	"github.com/HARA-DID/did-queueing-engine/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func newTestLogger() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(logrus.PanicLevel)
	return l
}

func buildHandler(t *testing.T, bc *mocks.MockBlockchainService) (*worker.Handler, *mocks.MockJobRepository) {
	t.Helper()
	repo := mocks.NewMockJobRepository()
	log := newTestLogger()
	svc := service.NewEventService(repo, bc, log)
	txCheckSvc := service.NewTxCheckService(repo, nil, callback.NewRegistry(), svc.EventCallbacks, log, 10)
	svc.SetTxCheckService(txCheckSvc)
	retryCfg := pkg.RetryConfig{MaxAttempts: 2, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond}
	reg := prometheus.NewRegistry()
	metrics := pkg.NewMetrics(reg)
	return worker.NewHandler(svc, retryCfg, metrics, log), repo
}

func streamValues(t *testing.T, event *domain.Event) map[string]interface{} {
	t.Helper()
	raw, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}
	return map[string]interface{}{"data": string(raw)}
}

func TestHandler_Handle_ValidEvent_ReturnsTrue(t *testing.T) {
	bc := &mocks.MockBlockchainService{}
	h, repo := buildHandler(t, bc)

	event := &domain.Event{
		ID:        "test-evt-001",
		Type:      domain.EventTypeCreateDID,
		Payload:   mustMarshal(t, domain.CreateDIDPayload{DID: "did:hara:test"}),
		CreatedAt: time.Now(),
	}

	if ack := h.Handle(context.Background(), "redis-msg-1", streamValues(t, event)); !ack {
		t.Error("Handle() returned false, want true")
	}

	job := repo.JobByEventID("test-evt-001")
	if job == nil {
		t.Fatal("job not persisted")
	}
	if job.Status != domain.JobStatusPending {
		t.Errorf("job status = %q, want pending (confirmation is background)", job.Status)
	}
}

func TestHandler_Handle_MissingDataField_ACKs(t *testing.T) {
	bc := &mocks.MockBlockchainService{}
	h, _ := buildHandler(t, bc)

	ack := h.Handle(context.Background(), "msg-bad", map[string]interface{}{"wrong": "value"})
	if !ack {
		t.Error("malformed message should be ACKed (poison pill prevention)")
	}
}

func TestHandler_Handle_InvalidJSON_ACKs(t *testing.T) {
	bc := &mocks.MockBlockchainService{}
	h, _ := buildHandler(t, bc)

	ack := h.Handle(context.Background(), "msg-json", map[string]interface{}{"data": "not-json{{{"})
	if !ack {
		t.Error("invalid JSON should be ACKed")
	}
}

func TestHandler_Handle_InvalidEventType_ACKs(t *testing.T) {
	bc := &mocks.MockBlockchainService{}
	h, _ := buildHandler(t, bc)

	event := &domain.Event{
		ID: "bad-type", Type: domain.EventType("INVALID"), Payload: json.RawMessage(`{}`), CreatedAt: time.Now(),
	}
	if ack := h.Handle(context.Background(), "msg-type", streamValues(t, event)); !ack {
		t.Error("invalid event type should be ACKed")
	}
}

func TestHandler_Handle_BlockchainFailure_ReturnsFalse(t *testing.T) {
	bc := &mocks.MockBlockchainService{
		CreateDIDFn: func(_ context.Context, _ domain.CreateDIDPayload) (*domain.BlockchainResult, error) {
			return nil, errSentinel
		},
	}
	h, _ := buildHandler(t, bc)

	event := &domain.Event{
		ID:        "fail-evt-001",
		Type:      domain.EventTypeCreateDID,
		Payload:   mustMarshal(t, domain.CreateDIDPayload{DID: "did:hara:fail"}),
		CreatedAt: time.Now(),
	}
	if ack := h.Handle(context.Background(), "msg-fail", streamValues(t, event)); ack {
		t.Error("blockchain failure should return false (DLQ route)")
	}
}

func mustMarshal(t *testing.T, v interface{}) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("mustMarshal: %v", err)
	}
	return b
}
