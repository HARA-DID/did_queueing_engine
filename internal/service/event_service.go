package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/HARA-DID/did-queueing-engine/internal/callback"
	"github.com/HARA-DID/did-queueing-engine/internal/domain"
	"github.com/HARA-DID/did-queueing-engine/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type eventHandler func(context.Context, json.RawMessage) (*domain.BlockchainResult, error)

// EventService orchestrates event processing: idempotency → DB → blockchain → DB.
type EventService struct {
	jobRepo        repository.JobRepository
	blockchain     BlockchainService
	log            *logrus.Logger
	txCheckSvc     *TxCheckService
	handlers       map[domain.EventType]eventHandler
	EventCallbacks map[domain.EventType]callback.Func
}

func NewEventService(
	jobRepo repository.JobRepository,
	blockchainSvc BlockchainService,
	log *logrus.Logger,
) *EventService {
	s := &EventService{
		jobRepo:        jobRepo,
		blockchain:     blockchainSvc,
		log:            log,
		handlers:       make(map[domain.EventType]eventHandler),
		EventCallbacks: make(map[domain.EventType]callback.Func),
	}
	s.initHandlers()
	return s
}

func (s *EventService) SetTxCheckService(svc *TxCheckService) {
	s.txCheckSvc = svc
}

func (s *EventService) Process(ctx context.Context, event *domain.Event) error {
	log := s.log.WithFields(logrus.Fields{
		"event_id":   event.ID,
		"event_type": string(event.Type),
	})

	existing, err := s.jobRepo.FindByEventID(ctx, event.ID)
	if err != nil {
		return fmt.Errorf("idempotency check: %w", err)
	}
	if existing != nil {
		if existing.Status == domain.JobStatusSuccess {
			log.Info("event already processed successfully, skipping")
			return domain.ErrAlreadyProcessed
		}
		log.WithField("job_id", existing.ID).Warn("re-processing previously failed event")
	}

	job := &domain.Job{
		ID:      uuid.NewString(),
		EventID: event.ID,
		Type:    string(event.Type),
		Status:  domain.JobStatusPending,
	}
	if existing != nil {
		job = existing // reuse the row if it already exists
	} else {
		if err := s.jobRepo.Create(ctx, job); err != nil {
			return fmt.Errorf("create pending job: %w", err)
		}
	}

	log = log.WithField("job_id", job.ID)
	log.Info("processing event")

	result, bcErr := s.dispatch(ctx, event)

	if bcErr != nil {
		errMsg := bcErr.Error()
		if err := s.jobRepo.UpdateStatus(ctx, job.ID, domain.JobStatusFailed, nil, errMsg); err != nil {
			log.WithError(err).Error("failed to update job status to failed")
		}
		log.WithError(bcErr).Error("blockchain operation failed")

		s.triggerCallback(ctx, event, callback.Result{
			EventID:      event.ID,
			JobID:        job.ID,
			EventType:    string(event.Type),
			Success:      false,
			ErrorMessage: errMsg,
		})

		return &domain.ErrBlockchain{Op: string(event.Type), Err: bcErr}
	}

	// For successful dispatch, delegate confirmations and final updates to TxCheckService
	if len(result.TxHashes) > 0 {
		s.txCheckSvc.Enqueue(TxCheckTask{
			Event: event,
			Job:   job,
			Hash:  result.TxHashes[0],
		})
		log.WithField("tx_hash", result.TxHashes[0]).Info("event dispatched; confirmation queued in background")
	}

	return nil
}

func (s *EventService) initHandlers() {
	// --- DID ---
	s.register(domain.EventTypeCreateDID, s.blockchain.CreateDID, callback.NoOp)
	s.register(domain.EventTypeAddKey, s.blockchain.AddKey, callback.NoOp)
	s.register(domain.EventTypeAddClaim, s.blockchain.AddClaim, callback.NoOp)
	s.register(domain.EventTypeStoreData, s.blockchain.StoreData, callback.NoOp)
	s.register(domain.EventTypeUpdateDID, s.blockchain.UpdateDID, callback.NoOp)
	s.register(domain.EventTypeDeactivateDID, s.blockchain.DeactivateDID, callback.NoOp)
	s.register(domain.EventTypeReactivateDID, s.blockchain.ReactivateDID, callback.NoOp)
	s.register(domain.EventTypeTransferDID, s.blockchain.TransferDIDOwner, callback.NoOp)
	s.register(domain.EventTypeDeleteData, s.blockchain.DeleteData, callback.NoOp)
	s.register(domain.EventTypeRemoveKey, s.blockchain.RemoveKey, callback.NoOp)
	s.register(domain.EventTypeRemoveClaim, s.blockchain.RemoveClaim, callback.NoOp)

	// --- Org ---
	s.register(domain.EventTypeCreateOrg, s.blockchain.CreateOrg, callback.NoOp)
	s.register(domain.EventTypeAddMember, s.blockchain.AddMember, callback.NoOp)
	s.register(domain.EventTypeRemoveMember, s.blockchain.RemoveMember, callback.NoOp)
	s.register(domain.EventTypeUpdateMember, s.blockchain.UpdateMember, callback.NoOp)
	s.register(domain.EventTypeDeactivateOrg, s.blockchain.DeactivateOrg, callback.NoOp)
	s.register(domain.EventTypeReactivateOrg, s.blockchain.ReactivateOrg, callback.NoOp)
	s.register(domain.EventTypeTransferOrgOwner, s.blockchain.TransferOrgOwner, callback.NoOp)

	// --- AA ---
	s.register(domain.EventTypeHandleOps, s.blockchain.HandleOps, callback.NoOp)
	s.register(domain.EventTypeDeployWallet, s.blockchain.DeployWallet, callback.NoOp)

	// --- VC ---
	s.register(domain.EventTypeIssueCredential, s.blockchain.IssueCredential, callback.NoOp)
	s.register(domain.EventTypeBurnCredential, s.blockchain.BurnCredential, callback.NoOp)
	s.register(domain.EventTypeUpdateMetadata, s.blockchain.UpdateMetadata, callback.NoOp)
	s.register(domain.EventTypeRevokeCredential, s.blockchain.RevokeCredential, callback.NoOp)
	s.register(domain.EventTypeApproveCredentialOrg, s.blockchain.ApproveCredentialOrg, callback.NoOp)
	s.register(domain.EventTypeApproveCredential, s.blockchain.ApproveCredential, callback.NoOp)
	s.register(domain.EventTypeSetDidRootStorage, s.blockchain.SetDidRootStorage, callback.NoOp)
	s.register(domain.EventTypeSetDidOrgStorage, s.blockchain.SetDidOrgStorage, callback.NoOp)

	// --- Alias ---
	s.register(domain.EventTypeRegisterTLD, s.blockchain.RegisterTLD, callback.NoOp)
	s.register(domain.EventTypeRegisterDomain, s.blockchain.RegisterDomain, callback.NoOp)
	s.register(domain.EventTypeSetDIDAlias, s.blockchain.SetDIDAlias, callback.NoOp)
	s.register(domain.EventTypeSetDIDOrgAlias, s.blockchain.SetDIDOrgAlias, callback.NoOp)
	s.register(domain.EventTypeExtendRegistration, s.blockchain.ExtendRegistration, callback.NoOp)
	s.register(domain.EventTypeRevokeAlias, s.blockchain.RevokeAlias, callback.NoOp)
	s.register(domain.EventTypeUnrevokeAlias, s.blockchain.UnrevokeAlias, callback.NoOp)
	s.register(domain.EventTypeRegisterSubdomain, s.blockchain.RegisterSubdomain, callback.NoOp)
	s.register(domain.EventTypeTransferAliasOwnership, s.blockchain.TransferAliasOwnership, callback.NoOp)
	s.register(domain.EventTypeTransferTLD, s.blockchain.TransferTLD, callback.NoOp)
	s.register(domain.EventTypeSetAliasRootStorage, s.blockchain.SetAliasRootStorage, callback.NoOp)
	s.register(domain.EventTypeSetAliasOrgStorage, s.blockchain.SetAliasOrgStorage, callback.NoOp)
	s.register(domain.EventTypeSetFactoryContract, s.blockchain.SetFactoryContract, callback.NoOp)
}

func sRegister[P any](
	s *EventService,
	eventType domain.EventType,
	fn func(context.Context, P) (*domain.BlockchainResult, error),
	cb callback.Func,
) {
	s.handlers[eventType] = func(ctx context.Context, raw json.RawMessage) (*domain.BlockchainResult, error) {
		var p P
		if err := json.Unmarshal(raw, &p); err != nil {
			return nil, fmt.Errorf("unmarshal %T: %w", p, err)
		}
		if v, ok := any(&p).(domain.Validator); ok {
			if err := v.Validate(); err != nil {
				return nil, err
			}
		}
		return fn(ctx, p)
	}
	s.EventCallbacks[eventType] = cb
}

func (s *EventService) register(eventType domain.EventType, fn any, cb callback.Func) {
	switch f := fn.(type) {
	case func(context.Context, domain.CreateDIDPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.AddKeyPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.AddClaimPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.StoreDataPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.UpdateDIDPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.DIDLifecyclePayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.TransferDIDOwnerPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.DeleteDataPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.RemoveKeyPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.RemoveClaimPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.CreateOrgPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.OrgMemberPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.OrgLifecyclePayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.OrgTransferPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.HandleOpsPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.DeployWalletPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.IssueCredentialPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.BurnCredentialPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.UpdateMetadataPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.RevokeCredentialPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.ApproveCredentialOrgPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.ApproveCredentialPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.RegisterTLDPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.RegisterDomainPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.SetDIDAliasPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.SetDIDOrgAliasPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.ExtendRegistrationPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.RevokeAliasPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.UnrevokeAliasPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.RegisterSubdomainPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.TransferAliasOwnershipPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.TransferTLDPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.SetAddressPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.SetAliasAddressPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	case func(context.Context, domain.SetFactoryContractPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f, cb)
	default:
		panic(fmt.Sprintf("unsupported handler signature for event %q: %T", eventType, fn))
	}
}

func (s *EventService) dispatch(ctx context.Context, event *domain.Event) (*domain.BlockchainResult, error) {
	handler, ok := s.handlers[event.Type]
	if !ok {
		return nil, fmt.Errorf("unknown event type: %q", event.Type)
	}
	return handler(ctx, event.Payload)
}

func (s *EventService) RecordRetry(ctx context.Context, jobID, errMsg string) {
	if err := s.jobRepo.IncrementRetry(ctx, jobID, errMsg); err != nil {
		s.log.WithError(err).WithField("job_id", jobID).Error("failed to record retry")
	}
}

func IsAlreadyProcessed(err error) bool {
	return errors.Is(err, domain.ErrAlreadyProcessed)
}

func (s *EventService) triggerCallback(ctx context.Context, event *domain.Event, result callback.Result) {
	cb, ok := s.EventCallbacks[event.Type]
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := cb(ctx, result); err != nil {
		s.log.WithFields(logrus.Fields{
			"event_id":   event.ID,
			"event_type": event.Type,
		}).WithError(err).Error("callback execution failed")
	}
}
