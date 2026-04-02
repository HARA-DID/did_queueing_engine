package domain

// Account Abstraction Event Types
const (
	EventTypeHandleOps        EventType = "HANDLE_OPS"
	EventTypeBatchHandleOps   EventType = "BATCH_HANDLE_OPS"
	EventTypeInitializeWallet EventType = "INITIALIZE_WALLET"
	EventTypeAddOwner         EventType = "ADD_OWNER"
	EventTypeTransferERC20    EventType = "TRANSFER_ERC20"
	EventTypeInitiateRecovery EventType = "INITIATE_RECOVERY"
	EventTypeApproveRecovery  EventType = "APPROVE_RECOVERY"
	EventTypeExecuteRecovery  EventType = "EXECUTE_RECOVERY"
	EventTypeDeployWallet     EventType = "DEPLOY_WALLET"
	EventTypeAddFactory       EventType = "ADD_FACTORY"
	EventTypeRemoveFactory    EventType = "REMOVE_FACTORY"
	EventTypeSetGasManager    EventType = "SET_GAS_MANAGER"
	EventTypeSetIsFree        EventType = "SET_IS_FREE"
	EventTypeWithdraw         EventType = "WITHDRAW"
)

// AA Payloads

type HandleOpsPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params would be mapped to SDK's HandleOpsParams
}

type BatchHandleOpsPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params would be mapped to SDK's BatchHandleOpsParams
}

type InitializeWalletPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
	// Params would be mapped to SDK's InitializeParams
}

type AddOwnerPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
}

type TransferERC20Payload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
}

type InitiateRecoveryPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
}

type DeployWalletPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
}

type AddFactoryPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
}

type RemoveFactoryPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
}

type SetGasManagerPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
}

type SetIsFreePayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
}

type WithdrawPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
}
