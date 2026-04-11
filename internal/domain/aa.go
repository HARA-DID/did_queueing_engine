package domain

import "fmt"

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
	Target            string `json:"target"`
	Value             string `json:"value"` // hex string for big.Int
	Data              []byte `json:"data"`
	UserNonce         string `json:"user_nonce,omitempty"` // optional, fetch from EP if empty
	Signature         []byte `json:"signature"`
	ClientBlockNumber string `json:"client_block_number,omitempty"`
	MultipleRPCCalls  bool   `json:"multiple_rpc_calls,omitempty"`
}

func (p HandleOpsPayload) Validate() error {
	if p.Target == "" {
		return fmt.Errorf("target is required")
	}
	if len(p.Signature) == 0 {
		return fmt.Errorf("signature is required")
	}
	return nil
}

type BatchHandleOpsPayload struct {
	Wallet           string             `json:"wallet"`
	UserOps          []HandleOpsPayload `json:"user_ops"`
	MultipleRPCCalls bool               `json:"multiple_rpc_calls,omitempty"`
}

func (p BatchHandleOpsPayload) Validate() error {
	if p.Wallet == "" {
		return fmt.Errorf("wallet is required")
	}
	if len(p.UserOps) == 0 {
		return fmt.Errorf("user_ops cannot be empty")
	}
	for i, op := range p.UserOps {
		if err := op.Validate(); err != nil {
			return fmt.Errorf("user_ops[%d]: %w", i, err)
		}
	}
	return nil
}

type InitializeWalletPayload struct {
	Owner1            string `json:"owner_1"`
	EntryPointAddress string `json:"entry_point_address"`
	FactoryAddress    string `json:"factory_address"`
	GuardianAdmin     string `json:"guardian_admin"`
	Version           string `json:"version"`
	MultipleRPCCalls  bool   `json:"multiple_rpc_calls,omitempty"`
}

func (p InitializeWalletPayload) Validate() error {
	if p.Owner1 == "" {
		return fmt.Errorf("owner_1 is required")
	}
	if p.EntryPointAddress == "" {
		return fmt.Errorf("entry_point_address is required")
	}
	if p.FactoryAddress == "" {
		return fmt.Errorf("factory_address is required")
	}
	if p.GuardianAdmin == "" {
		return fmt.Errorf("guardian_admin is required")
	}
	if p.Version == "" {
		return fmt.Errorf("version is required")
	}
	return nil
}

type AddOwnerPayload struct {
	NewOwner         string `json:"new_owner"`
	MultipleRPCCalls bool   `json:"multiple_rpc_calls,omitempty"`
}

func (p AddOwnerPayload) Validate() error {
	if p.NewOwner == "" {
		return fmt.Errorf("new_owner is required")
	}
	return nil
}

type TransferERC20Payload struct {
	Token            string `json:"token"`
	To               string `json:"to"`
	Amount           string `json:"amount"` // Hex string for big.Int
	MultipleRPCCalls bool   `json:"multiple_rpc_calls,omitempty"`
}

func (p TransferERC20Payload) Validate() error {
	if p.Token == "" {
		return fmt.Errorf("token is required")
	}
	if p.To == "" {
		return fmt.Errorf("to is required")
	}
	if p.Amount == "" {
		return fmt.Errorf("amount is required")
	}
	return nil
}

type InitiateRecoveryPayload struct {
	OldOwner         string `json:"old_owner"`
	NewOwner         string `json:"new_owner"`
	MultipleRPCCalls bool   `json:"multiple_rpc_calls,omitempty"`
}

func (p InitiateRecoveryPayload) Validate() error {
	if p.OldOwner == "" {
		return fmt.Errorf("old_owner is required")
	}
	if p.NewOwner == "" {
		return fmt.Errorf("new_owner is required")
	}
	return nil
}

type ApproveRecoveryPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
}

func (p ApproveRecoveryPayload) Validate() error {
	return nil
}

type ExecuteRecoveryPayload struct {
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
}

func (p ExecuteRecoveryPayload) Validate() error {
	return nil
}

type DeployWalletPayload struct {
	Owner            string `json:"owner"`
	Salt             string `json:"salt"`
	MultipleRPCCalls bool   `json:"multiple_rpc_calls,omitempty"`
}

func (p DeployWalletPayload) Validate() error {
	if p.Owner == "" {
		return fmt.Errorf("owner is required")
	}
	if p.Salt == "" {
		return fmt.Errorf("salt is required")
	}
	return nil
}

type AddFactoryPayload struct {
	Factory          string `json:"factory"`
	MultipleRPCCalls bool   `json:"multiple_rpc_calls,omitempty"`
}

func (p AddFactoryPayload) Validate() error {
	if p.Factory == "" {
		return fmt.Errorf("factory is required")
	}
	return nil
}

type RemoveFactoryPayload struct {
	Factory          string `json:"factory"`
	MultipleRPCCalls bool   `json:"multiple_rpc_calls,omitempty"`
}

func (p RemoveFactoryPayload) Validate() error {
	if p.Factory == "" {
		return fmt.Errorf("factory is required")
	}
	return nil
}

type SetGasManagerPayload struct {
	GasManager       string `json:"gas_manager"`
	MultipleRPCCalls bool   `json:"multiple_rpc_calls,omitempty"`
}

func (p SetGasManagerPayload) Validate() error {
	if p.GasManager == "" {
		return fmt.Errorf("gas_manager is required")
	}
	return nil
}

type SetIsFreePayload struct {
	IsFree           bool `json:"is_free"`
	MultipleRPCCalls bool `json:"multiple_rpc_calls,omitempty"`
}

func (p SetIsFreePayload) Validate() error {
	return nil
}

type WithdrawPayload struct {
	Amount           string `json:"amount"` // Hex string for big.Int
	To               string `json:"to"`
	MultipleRPCCalls bool   `json:"multiple_rpc_calls,omitempty"`
}

func (p WithdrawPayload) Validate() error {
	if p.Amount == "" {
		return fmt.Errorf("amount is required")
	}
	if p.To == "" {
		return fmt.Errorf("to is required")
	}
	return nil
}
