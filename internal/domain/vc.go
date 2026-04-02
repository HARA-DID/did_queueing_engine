package domain

import (
	"math/big"
)

// Verifiable Credentials Event Types
const (
	EventTypeIssueCredential     EventType = "ISSUE_CREDENTIAL"
	EventTypeBurnCredential      EventType = "BURN_CREDENTIAL"
	EventTypeUpdateMetadata      EventType = "UPDATE_METADATA"
	EventTypeRevokeCredential    EventType = "REVOKE_CREDENTIAL"
	EventTypeApproveCredentialOrg EventType = "APPROVE_CREDENTIAL_ORG"
	EventTypeApproveCredential    EventType = "APPROVE_CREDENTIAL"
	EventTypeSetDidRootStorage    EventType = "SET_DID_ROOT_STORAGE"
	EventTypeSetDidOrgStorage     EventType = "SET_DID_ORG_STORAGE"
)

// VC Payloads

type IssueCredentialPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params would be mapped to vcfactory.IssueCredentialParams
}

type BurnCredentialPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params would be mapped to vcfactory.BurnCredentialParams
}

type UpdateMetadataPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params would be mapped to vcfactory.UpdateMetadataParams
}

type RevokeCredentialPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params would be mapped to vcfactory.RevokeCredentialParams
}

type ApproveCredentialOrgPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params would be mapped to vcfactory.ApproveCredentialOrgParams
}

type ApproveCredentialPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params would be mapped to vcfactory.ApproveCredentialParams
}

type SetAddressPayload struct {
	Address          string `json:"address"`
	MultipleRPCCalls bool   `json:"multiple_rpc_calls,omitempty"`
}

// Result types for queries (if needed by worker)
type TokenIdsResult struct {
	TokenIds []*big.Int
	Total    *big.Int
}
