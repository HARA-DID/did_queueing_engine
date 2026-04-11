package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/HARA-DID/did-queueing-engine/internal/callback"
	"github.com/HARA-DID/did-queueing-engine/internal/domain"
	"github.com/HARA-DID/did-queueing-engine/internal/mocks"
	"github.com/HARA-DID/did-queueing-engine/internal/service"
	"github.com/sirupsen/logrus"
)

func newLogger() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.PanicLevel)
	return log
}

func makeEvent(t *testing.T, id string, evType domain.EventType, payload interface{}) *domain.Event {
	t.Helper()
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	return &domain.Event{
		ID:        id,
		Type:      evType,
		Payload:   json.RawMessage(raw),
		CreatedAt: time.Now(),
	}
}

func TestEventService_Process_CreateDID_Success(t *testing.T) {
	repo := mocks.NewMockJobRepository()
	bc := &mocks.MockBlockchainService{}
	svc := service.NewEventService(repo, bc, newLogger())
	txCheckSvc := service.NewTxCheckService(repo, nil, callback.NewRegistry(), svc.EventCallbacks, newLogger(), 10)
	svc.SetTxCheckService(txCheckSvc)

	event := makeEvent(t, "evt-001", domain.EventTypeCreateDID, domain.CreateDIDPayload{
		DID: "did:hara:alice",
	})

	if err := svc.Process(context.Background(), event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := bc.CallCount("CreateDID"); got != 1 {
		t.Errorf("CreateDID called %d times, want 1", got)
	}

	job := repo.JobByEventID("evt-001")
	if job == nil {
		t.Fatal("job not persisted")
	}
	if job.Status != domain.JobStatusPending {
		t.Errorf("job status = %q, want pending (confirmation is background)", job.Status)
	}
	if len(job.TxHashes) == 0 {
		t.Error("expected at least one tx hash")
	}
}

func TestEventService_Process_Idempotency_AlreadySuccess(t *testing.T) {
	repo := mocks.NewMockJobRepository()
	bc := &mocks.MockBlockchainService{}
	svc := service.NewEventService(repo, bc, newLogger())
	txCheckSvc := service.NewTxCheckService(repo, nil, callback.NewRegistry(), svc.EventCallbacks, newLogger(), 10)
	svc.SetTxCheckService(txCheckSvc)

	event := makeEvent(t, "evt-dup", domain.EventTypeAddKey, domain.AddKeyPayload{
		DIDIndex: big.NewInt(42), KeyType: 0, PublicKey: "abc123", Purpose: 0,
	})

	if err := svc.Process(context.Background(), event); err != nil {
		t.Fatalf("first process: %v", err)
	}

	// Mark as success so idempotency fires on second call.
	job := repo.JobByEventID("evt-dup")
	_ = repo.UpdateStatus(context.Background(), job.ID, domain.JobStatusSuccess, job.TxHashes, "")

	err := svc.Process(context.Background(), event)
	if !errors.Is(err, domain.ErrAlreadyProcessed) {
		t.Errorf("expected ErrAlreadyProcessed, got: %v", err)
	}
	if got := bc.CallCount("AddKey"); got != 1 {
		t.Errorf("AddKey called %d times, want 1", got)
	}
}

func TestEventService_Process_BlockchainError(t *testing.T) {
	repo := mocks.NewMockJobRepository()
	bc := &mocks.MockBlockchainService{
		AddClaimFn: func(_ context.Context, _ domain.AddClaimPayload) (*domain.BlockchainResult, error) {
			return nil, errors.New("rpc timeout")
		},
	}
	svc := service.NewEventService(repo, bc, newLogger())
	txCheckSvc := service.NewTxCheckService(repo, nil, callback.NewRegistry(), svc.EventCallbacks, newLogger(), 10)
	svc.SetTxCheckService(txCheckSvc)

	event := makeEvent(t, "evt-bc-err", domain.EventTypeAddClaim, domain.AddClaimPayload{
		DIDIndex: big.NewInt(7), Topic: 1, IssuerAddress: "0xDeAdBeEf00000000000000000000000000000001",
		Data: []byte(`{"email":"carol@example.com"}`),
	})

	err := svc.Process(context.Background(), event)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var bcErr *domain.ErrBlockchain
	if !errors.As(err, &bcErr) {
		t.Errorf("expected *domain.ErrBlockchain, got %T: %v", err, err)
	}

	job := repo.JobByEventID("evt-bc-err")
	if job == nil {
		t.Fatal("job not found")
	}
	if job.Status != domain.JobStatusFailed {
		t.Errorf("job status = %q, want failed", job.Status)
	}
}

func TestEventService_Process_AllEventTypes(t *testing.T) {
	tests := []struct {
		name       string
		evType     domain.EventType
		payload    interface{}
		wantMethod string
	}{
		{
			"CREATE_DID", domain.EventTypeCreateDID,
			domain.CreateDIDPayload{DID: "did:hara:test-1"},
			"CreateDID",
		},
		{
			"ADD_KEY", domain.EventTypeAddKey,
			domain.AddKeyPayload{DIDIndex: big.NewInt(1), KeyType: 0, PublicKey: "pub", Purpose: 0},
			"AddKey",
		},
		{
			"ADD_CLAIM", domain.EventTypeAddClaim,
			domain.AddClaimPayload{DIDIndex: big.NewInt(1), Topic: 2, IssuerAddress: "0x0000000000000000000000000000000000000001", Data: []byte("claim")},
			"AddClaim",
		},
		{
			"STORE_DATA", domain.EventTypeStoreData,
			domain.StoreDataPayload{DIDIndex: big.NewInt(1), PropertyKey: "profile", Data: "{}"},
			"StoreData",
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewMockJobRepository()
			bc := &mocks.MockBlockchainService{}
			svc := service.NewEventService(repo, bc, newLogger())
			txCheckSvc := service.NewTxCheckService(repo, nil, callback.NewRegistry(), svc.EventCallbacks, newLogger(), 10)
			svc.SetTxCheckService(txCheckSvc)

			event := makeEvent(t, fmt.Sprintf("evt-%d", i), tc.evType, tc.payload)
			if err := svc.Process(context.Background(), event); err != nil {
				t.Fatalf("Process() error: %v", err)
			}
			if bc.CallCount(tc.wantMethod) != 1 {
				t.Errorf("%s: expected 1 call, got %d", tc.wantMethod, bc.CallCount(tc.wantMethod))
			}
		})
	}
}

func TestEventService_Process_CallbackInvocation(t *testing.T) {
	repo := mocks.NewMockJobRepository()
	bc := &mocks.MockBlockchainService{}
	registry := callback.NewRegistry()

	svc := service.NewEventService(repo, bc, newLogger())
	txCheckSvc := service.NewTxCheckService(repo, nil, registry, svc.EventCallbacks, newLogger(), 10)
	svc.SetTxCheckService(txCheckSvc)

	event := makeEvent(t, "evt-cb", domain.EventTypeCreateDID, domain.CreateDIDPayload{DID: "did:hara:cb"})

	if err := svc.Process(context.Background(), event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
