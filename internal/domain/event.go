package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

// Validator is implemented by any domain payload that requires custom validation.
type Validator interface {
	Validate() error
}

// EventType represents a distinct blockchain operation.
type EventType string

// IsValid checks that the event type is recognized by one of our modules.
func (e EventType) IsValid() bool {
	// We check for presence in any of our specialized event categories.
	// This will be expanded as we add more modules.
	switch e {
	case EventTypeCreateDID, EventTypeAddKey, EventTypeAddClaim, EventTypeStoreData,
		EventTypeHandleOps, EventTypeBatchHandleOps, EventTypeInitializeWallet,
		EventTypeAddOwner, EventTypeTransferERC20, EventTypeInitiateRecovery,
		EventTypeApproveRecovery, EventTypeExecuteRecovery, EventTypeDeployWallet,
		EventTypeAddFactory, EventTypeRemoveFactory, EventTypeSetGasManager,
		EventTypeSetIsFree, EventTypeWithdraw,
		EventTypeIssueCredential, EventTypeBurnCredential, EventTypeUpdateMetadata,
		EventTypeRevokeCredential, EventTypeApproveCredentialOrg, EventTypeApproveCredential,
		EventTypeSetDidRootStorage, EventTypeSetDidOrgStorage,
		EventTypeRegisterTLD, EventTypeRegisterDomain, EventTypeSetDIDAlias,
		EventTypeSetDIDOrgAlias, EventTypeExtendRegistration, EventTypeRevokeAlias,
		EventTypeUnrevokeAlias, EventTypeRegisterSubdomain, EventTypeTransferAliasOwnership,
		EventTypeTransferTLD, EventTypeSetAliasRootStorage, EventTypeSetAliasOrgStorage:
		return true
	}
	return false
}

type Event struct {
	ID        string          `json:"id"`
	Type      EventType       `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

func (e *Event) Validate() error {
	if e.ID == "" {
		return fmt.Errorf("event id is required")
	}
	if !e.Type.IsValid() {
		return fmt.Errorf("unsupported event type: %q", e.Type)
	}
	if len(e.Payload) == 0 {
		return fmt.Errorf("event payload must not be empty")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Job
// ---------------------------------------------------------------------------

type JobStatus string

const (
	JobStatusPending JobStatus = "pending"
	JobStatusSuccess JobStatus = "success"
	JobStatusFailed  JobStatus = "failed"
)

type Job struct {
	ID           string    `db:"id"`
	EventID      string    `db:"event_id"`
	Type         string    `db:"type"`
	Status       JobStatus `db:"status"`
	TxHashes     []string  `db:"tx_hashes"`
	RetryCount   int       `db:"retry_count"`
	ErrorMessage string    `db:"error_message"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// ---------------------------------------------------------------------------
// Blockchain result
// ---------------------------------------------------------------------------

type BlockchainResult struct {
	TxHashes []string
}

// ---------------------------------------------------------------------------
// Errors
// ---------------------------------------------------------------------------

var ErrAlreadyProcessed = fmt.Errorf("event already processed (idempotency check)")

type ErrValidation struct {
	Field   string
	Message string
}

func (e *ErrValidation) Error() string {
	return fmt.Sprintf("validation error on %q: %s", e.Field, e.Message)
}

type ErrBlockchain struct {
	Op  string
	Err error
}

func (e *ErrBlockchain) Error() string {
	return fmt.Sprintf("blockchain operation %q failed: %v", e.Op, e.Err)
}

func (e *ErrBlockchain) Unwrap() error { return e.Err }
