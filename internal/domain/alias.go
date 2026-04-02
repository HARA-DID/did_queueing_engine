package domain

// Alias Event Types
const (
	EventTypeRegisterTLD          EventType = "REGISTER_TLD"
	EventTypeRegisterDomain       EventType = "REGISTER_DOMAIN"
	EventTypeSetDIDAlias         EventType = "SET_DID_ALIAS"
	EventTypeSetDIDOrgAlias      EventType = "SET_DID_ORG_ALIAS"
	EventTypeExtendRegistration   EventType = "EXTEND_REGISTRATION"
	EventTypeRevokeAlias          EventType = "REVOKE_ALIAS"
	EventTypeUnrevokeAlias        EventType = "UNREVOKE_ALIAS"
	EventTypeRegisterSubdomain    EventType = "REGISTER_SUBDOMAIN"
	EventTypeTransferAliasOwnership EventType = "TRANSFER_ALIAS_OWNERSHIP"
	EventTypeTransferTLD          EventType = "TRANSFER_TLD"
	EventTypeSetAliasRootStorage  EventType = "SET_ALIAS_ROOT_STORAGE"
	EventTypeSetAliasOrgStorage   EventType = "SET_ALIAS_ORG_STORAGE"
)

// Alias Payloads

type RegisterTLDPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params mapped to aliasfactory.RegisterTLDParams
}

type RegisterDomainPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params mapped to aliasfactory.RegisterDomainParams
}

type SetDIDAliasPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params mapped to aliasfactory.SetDIDParams
}

type SetDIDOrgAliasPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params mapped to aliasfactory.SetDIDOrgParams
}

type ExtendRegistrationPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params mapped to aliasfactory.ExtendRegistrationParams
}

type RevokeAliasPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params mapped to aliasfactory.NodeOnlyParams
}

type UnrevokeAliasPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params mapped to aliasfactory.NodeOnlyParams
}

type RegisterSubdomainPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params mapped to aliasfactory.RegisterSubdomainParams
}

type TransferAliasOwnershipPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params mapped to aliasfactory.TransferAliasOwnershipParams
}

type TransferTLDPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params mapped to aliasfactory.TransferTLDParams
}
