package domain

import (
	"math/big"
)

// DID Event Types
const (
	EventTypeCreateDID EventType = "CREATE_DID"
	EventTypeAddKey    EventType = "ADD_KEY"
	EventTypeAddClaim  EventType = "ADD_CLAIM"
	EventTypeStoreData EventType = "STORE_DATA"
)

// CreateDIDPayload represents the payload for registering a new DID.
type CreateDIDPayload struct {
	DID              string `json:"did"`
	KeyIdentifier    string `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool   `json:"multiple_rpc_calls,omitempty"`
}

// AddKeyPayload represents the payload for adding a verification key.
type AddKeyPayload struct {
	DIDIndex         *big.Int `json:"did_index"`
	KeyType          uint8    `json:"key_type"`
	PublicKey        string   `json:"public_key"`
	Purpose          uint8    `json:"purpose"`
	KeyIdentifier    string   `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}

// AddClaimPayload represents the payload for attaching a claim.
type AddClaimPayload struct {
	DIDIndex         *big.Int `json:"did_index"`
	Topic            uint8    `json:"topic"`
	IssuerAddress    string   `json:"issuer_address"`
	Data             []byte   `json:"data"`
	KeyIdentifier    string   `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}

// StoreDataPayload represents the payload for storing arbitrary data.
type StoreDataPayload struct {
	DIDIndex         *big.Int `json:"did_index"`
	PropertyKey      string   `json:"property_key"`
	Data             []byte   `json:"data"`
	KeyIdentifier    string   `json:"key_identifier,omitempty"`
	MultipleRPCCalls bool     `json:"multiple_rpc_calls,omitempty"`
}
