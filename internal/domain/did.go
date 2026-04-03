package domain

import (
	"math/big"
)

const (
	EventTypeCreateDID     EventType = "CREATE_DID"
	EventTypeAddKey        EventType = "ADD_KEY"
	EventTypeAddClaim      EventType = "ADD_CLAIM"
	EventTypeStoreData     EventType = "STORE_DATA"
	EventTypeUpdateDID    EventType = "UPDATE_DID"
	EventTypeDeactivateDID EventType = "DEACTIVATE_DID"
	EventTypeReactivateDID EventType = "REACTIVATE_DID"
	EventTypeTransferDID   EventType = "TRANSFER_DID"
	EventTypeDeleteData    EventType = "DELETE_DATA"
	EventTypeRemoveKey     EventType = "REMOVE_KEY"
	EventTypeRemoveClaim   EventType = "REMOVE_CLAIM"
	EventTypeCreateOrg     EventType = "CREATE_ORG"
	EventTypeAddMember     EventType = "ADD_MEMBER"
	EventTypeRemoveMember  EventType = "REMOVE_MEMBER"
	EventTypeUpdateMember  EventType = "UPDATE_MEMBER"
)

type CreateDIDPayload struct {
	TargetAddress    string `json:"target_address"`
	DID              string `json:"did"`
	KeyIdentifier    string `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool   `json:"multiple_rpc_calls,omitempty"`
}

type AddKeyPayload struct {
	TargetAddress    string   `json:"target_address"`
	DIDIndex         *big.Int `json:"did_index"`
	KeyType          uint8    `json:"key_type"`
	PublicKey        string   `json:"public_key"`
	Purpose          uint8    `json:"purpose"`
	KeyIdentifier    string   `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}

type AddClaimPayload struct {
	TargetAddress    string   `json:"target_address"`
	DIDIndex         *big.Int `json:"did_index"`
	Topic            uint8    `json:"topic"`
	IssuerAddress    string   `json:"issuer_address"`
	Data             []byte   `json:"data"`
	URI              string   `json:"uri"`
	Signature        []byte   `json:"signature"`
	KeyIdentifier    string   `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}

type StoreDataPayload struct {
	TargetAddress    string   `json:"target_address"`
	DIDIndex         *big.Int `json:"did_index"`
	PropertyKey      string   `json:"property_key"`
	Data             string   `json:"data"`
	KeyIdentifier    string   `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}

type UpdateDIDPayload struct {
	TargetAddress    string   `json:"target_address"`
	DIDIndex         *big.Int `json:"did_index"`
	URI              string   `json:"uri"`
	KeyIdentifier    string   `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}

type DIDLifecyclePayload struct {
	TargetAddress    string   `json:"target_address"`
	DIDIndex         *big.Int `json:"did_index"`
	KeyIdentifier    string   `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}

type TransferDIDOwnerPayload struct {
	TargetAddress    string   `json:"target_address"`
	DIDIndex         *big.Int `json:"did_index"`
	NewOwner         string   `json:"new_owner"`
	KeyIdentifier    string   `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}

type DeleteDataPayload struct {
	TargetAddress    string   `json:"target_address"`
	DIDIndex         *big.Int `json:"did_index"`
	Key              string   `json:"key"`
	KeyIdentifier    string   `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}

type RemoveKeyPayload struct {
	TargetAddress    string   `json:"target_address"`
	DIDIndex         *big.Int `json:"did_index"`
	KeyDataHashed    string   `json:"key_data_hashed"`
	KeyIdentifier    string   `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}

type RemoveClaimPayload struct {
	TargetAddress    string   `json:"target_address"`
	DIDIndex         *big.Int `json:"did_index"`
	ClaimID          string   `json:"claim_id"`
	KeyIdentifier    string   `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}

type CreateOrgPayload struct {
	TargetAddress    string `json:"target_address"`
	Data             []byte `json:"data"`
	MultipleRPCCalls bool   `json:"multiple_rpc_calls,omitempty"`
}

type OrgLifecyclePayload struct {
	TargetAddress    string   `json:"target_address"`
	OrgDIDIndex      *big.Int `json:"org_did_index"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}

type OrgTransferPayload struct {
	TargetAddress    string   `json:"target_address"`
	OrgDIDIndex      *big.Int `json:"org_did_index"`
	Data             []byte   `json:"data"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}

type OrgMemberPayload struct {
	TargetAddress    string   `json:"target_address"`
	OrgDIDIndex      *big.Int `json:"org_did_index"`
	Data             []byte   `json:"data"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}

type GeneralExecutePayload struct {
	TargetAddress    string `json:"target_address"`
	Data             []byte `json:"data"`
	KeyIdentifier    string `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool   `json:"multiple_rpc_calls,omitempty"`
}
