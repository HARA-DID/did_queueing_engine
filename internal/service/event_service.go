package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/HARA-DID/did-queueing-engine/internal/domain"
	"github.com/HARA-DID/did-queueing-engine/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type eventHandler func(context.Context, json.RawMessage) (*domain.BlockchainResult, error)

// EventService orchestrates event processing: idempotency → DB → blockchain → DB.
type EventService struct {
	jobRepo    repository.JobRepository
	blockchain BlockchainService
	log        *logrus.Logger
	handlers   map[domain.EventType]eventHandler
}

func NewEventService(
	jobRepo repository.JobRepository,
	blockchain BlockchainService,
	log *logrus.Logger,
) *EventService {
	s := &EventService{
		jobRepo:    jobRepo,
		blockchain: blockchain,
		log:        log,
	}
	s.initHandlers()
	return s
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
		_ = s.jobRepo.UpdateStatus(ctx, job.ID, domain.JobStatusFailed, nil, errMsg)
		log.WithError(bcErr).Error("blockchain operation failed")
		return &domain.ErrBlockchain{Op: string(event.Type), Err: bcErr}
	}

	if err := s.jobRepo.UpdateStatus(ctx, job.ID, domain.JobStatusSuccess, result.TxHashes, ""); err != nil {
		log.WithError(err).Error("failed to update job status to success")
	}

	log.WithField("tx_hashes", result.TxHashes).Info("event processed successfully")
	return nil
}

func (s *EventService) initHandlers() {
	s.handlers = make(map[domain.EventType]eventHandler)

	// --- DID ---
	s.register(domain.EventTypeCreateDID, s.blockchain.CreateDID)
	s.register(domain.EventTypeAddKey, s.blockchain.AddKey)
	s.register(domain.EventTypeAddClaim, s.blockchain.AddClaim)
	s.register(domain.EventTypeStoreData, s.blockchain.StoreData)
	s.register(domain.EventTypeUpdateDID, s.blockchain.UpdateDID)
	s.register(domain.EventTypeDeactivateDID, s.blockchain.DeactivateDID)
	s.register(domain.EventTypeReactivateDID, s.blockchain.ReactivateDID)
	s.register(domain.EventTypeTransferDID, s.blockchain.TransferDIDOwner)
	s.register(domain.EventTypeDeleteData, s.blockchain.DeleteData)
	s.register(domain.EventTypeRemoveKey, s.blockchain.RemoveKey)
	s.register(domain.EventTypeRemoveClaim, s.blockchain.RemoveClaim)

	// --- Org ---
	s.register(domain.EventTypeCreateOrg, s.blockchain.CreateOrg)
	s.register(domain.EventTypeAddMember, s.blockchain.AddMember)
	s.register(domain.EventTypeRemoveMember, s.blockchain.RemoveMember)
	s.register(domain.EventTypeUpdateMember, s.blockchain.UpdateMember)
	s.register(domain.EventTypeDeactivateOrg, s.blockchain.DeactivateOrg)
	s.register(domain.EventTypeReactivateOrg, s.blockchain.ReactivateOrg)
	s.register(domain.EventTypeTransferOrgOwner, s.blockchain.TransferOrgOwner)

	// --- AA ---
	s.register(domain.EventTypeHandleOps, s.blockchain.HandleOps)
	s.register(domain.EventTypeDeployWallet, s.blockchain.DeployWallet)

	// --- VC ---
	s.register(domain.EventTypeIssueCredential, s.blockchain.IssueCredential)
	s.register(domain.EventTypeBurnCredential, s.blockchain.BurnCredential)
	s.register(domain.EventTypeUpdateMetadata, s.blockchain.UpdateMetadata)
	s.register(domain.EventTypeRevokeCredential, s.blockchain.RevokeCredential)
	s.register(domain.EventTypeApproveCredentialOrg, s.blockchain.ApproveCredentialOrg)
	s.register(domain.EventTypeApproveCredential, s.blockchain.ApproveCredential)
	s.register(domain.EventTypeSetDidRootStorage, s.blockchain.SetDidRootStorage)
	s.register(domain.EventTypeSetDidOrgStorage, s.blockchain.SetDidOrgStorage)

	// --- Alias ---
	s.register(domain.EventTypeRegisterTLD, s.blockchain.RegisterTLD)
	s.register(domain.EventTypeRegisterDomain, s.blockchain.RegisterDomain)
	s.register(domain.EventTypeSetDIDAlias, s.blockchain.SetDIDAlias)
	s.register(domain.EventTypeSetDIDOrgAlias, s.blockchain.SetDIDOrgAlias)
	s.register(domain.EventTypeExtendRegistration, s.blockchain.ExtendRegistration)
	s.register(domain.EventTypeRevokeAlias, s.blockchain.RevokeAlias)
	s.register(domain.EventTypeUnrevokeAlias, s.blockchain.UnrevokeAlias)
	s.register(domain.EventTypeRegisterSubdomain, s.blockchain.RegisterSubdomain)
	s.register(domain.EventTypeTransferAliasOwnership, s.blockchain.TransferAliasOwnership)
	s.register(domain.EventTypeTransferTLD, s.blockchain.TransferTLD)
	s.register(domain.EventTypeSetAliasRootStorage, s.blockchain.SetAliasRootStorage)
	s.register(domain.EventTypeSetAliasOrgStorage, s.blockchain.SetAliasOrgStorage)
	s.register(domain.EventTypeSetFactoryContract, s.blockchain.SetFactoryContract)
}

func sRegister[P any](
	s *EventService,
	eventType domain.EventType,
	fn func(context.Context, P) (*domain.BlockchainResult, error),
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
}

func (s *EventService) register(eventType domain.EventType, fn any) {
	switch f := fn.(type) {
	case func(context.Context, domain.CreateDIDPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.AddKeyPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.AddClaimPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.StoreDataPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.UpdateDIDPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.DIDLifecyclePayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.TransferDIDOwnerPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.DeleteDataPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.RemoveKeyPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.RemoveClaimPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.CreateOrgPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.OrgMemberPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.OrgLifecyclePayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.OrgTransferPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.HandleOpsPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.DeployWalletPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.IssueCredentialPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.BurnCredentialPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.UpdateMetadataPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.RevokeCredentialPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.ApproveCredentialOrgPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.ApproveCredentialPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.RegisterTLDPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.RegisterDomainPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.SetDIDAliasPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.SetDIDOrgAliasPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.ExtendRegistrationPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.RevokeAliasPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.UnrevokeAliasPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.RegisterSubdomainPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.TransferAliasOwnershipPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.TransferTLDPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.SetAddressPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.SetAliasAddressPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
	case func(context.Context, domain.SetFactoryContractPayload) (*domain.BlockchainResult, error):
		sRegister(s, eventType, f)
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
