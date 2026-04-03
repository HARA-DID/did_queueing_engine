package mocks

import (
	"context"
	"sync"

	"github.com/myorg/worker-service/internal/domain"
	"github.com/myorg/worker-service/internal/service"
)

var _ service.BlockchainService = (*MockBlockchainService)(nil)

type Call struct {
	Method  string
	Payload interface{}
}

type MockBlockchainService struct {
	mu    sync.Mutex
	Calls []Call

	CreateDIDFn        func(ctx context.Context, p domain.CreateDIDPayload) (*domain.BlockchainResult, error)
	AddKeyFn           func(ctx context.Context, p domain.AddKeyPayload) (*domain.BlockchainResult, error)
	AddClaimFn         func(ctx context.Context, p domain.AddClaimPayload) (*domain.BlockchainResult, error)
	StoreDataFn        func(ctx context.Context, p domain.StoreDataPayload) (*domain.BlockchainResult, error)
	UpdateDIDFn        func(ctx context.Context, p domain.UpdateDIDPayload) (*domain.BlockchainResult, error)
	DeactivateDIDFn     func(ctx context.Context, p domain.DIDLifecyclePayload) (*domain.BlockchainResult, error)
	ReactivateDIDFn     func(ctx context.Context, p domain.DIDLifecyclePayload) (*domain.BlockchainResult, error)
	TransferDIDOwnerFn  func(ctx context.Context, p domain.TransferDIDOwnerPayload) (*domain.BlockchainResult, error)
	DeleteDataFn       func(ctx context.Context, p domain.DeleteDataPayload) (*domain.BlockchainResult, error)
	RemoveKeyFn        func(ctx context.Context, p domain.RemoveKeyPayload) (*domain.BlockchainResult, error)
	RemoveClaimFn      func(ctx context.Context, p domain.RemoveClaimPayload) (*domain.BlockchainResult, error)
	GeneralExecuteFn   func(ctx context.Context, p domain.GeneralExecutePayload) (*domain.BlockchainResult, error)
	CreateOrgFn        func(ctx context.Context, p domain.CreateOrgPayload) (*domain.BlockchainResult, error)
	DeactivateOrgFn     func(ctx context.Context, p domain.OrgLifecyclePayload) (*domain.BlockchainResult, error)
	ReactivateOrgFn     func(ctx context.Context, p domain.OrgLifecyclePayload) (*domain.BlockchainResult, error)
	TransferOrgOwnerFn  func(ctx context.Context, p domain.OrgTransferPayload) (*domain.BlockchainResult, error)
	AddMemberFn        func(ctx context.Context, p domain.OrgMemberPayload) (*domain.BlockchainResult, error)
	RemoveMemberFn     func(ctx context.Context, p domain.OrgMemberPayload) (*domain.BlockchainResult, error)
	UpdateMemberFn     func(ctx context.Context, p domain.OrgMemberPayload) (*domain.BlockchainResult, error)

	HandleOpsFn    func(ctx context.Context, p domain.HandleOpsPayload) (*domain.BlockchainResult, error)
	DeployWalletFn func(ctx context.Context, p domain.DeployWalletPayload) (*domain.BlockchainResult, error)

	IssueCredentialFn      func(ctx context.Context, p domain.IssueCredentialPayload) (*domain.BlockchainResult, error)
	BurnCredentialFn       func(ctx context.Context, p domain.BurnCredentialPayload) (*domain.BlockchainResult, error)
	UpdateMetadataFn       func(ctx context.Context, p domain.UpdateMetadataPayload) (*domain.BlockchainResult, error)
	RevokeCredentialFn     func(ctx context.Context, p domain.RevokeCredentialPayload) (*domain.BlockchainResult, error)
	ApproveCredentialOrgFn  func(ctx context.Context, p domain.ApproveCredentialOrgPayload) (*domain.BlockchainResult, error)
	ApproveCredentialFn     func(ctx context.Context, p domain.ApproveCredentialPayload) (*domain.BlockchainResult, error)
	SetDidRootStorageFn    func(ctx context.Context, p domain.SetAddressPayload) (*domain.BlockchainResult, error)
	SetDidOrgStorageFn     func(ctx context.Context, p domain.SetAddressPayload) (*domain.BlockchainResult, error)

	RegisterTLDFn           func(ctx context.Context, p domain.RegisterTLDPayload) (*domain.BlockchainResult, error)
	RegisterDomainFn        func(ctx context.Context, p domain.RegisterDomainPayload) (*domain.BlockchainResult, error)
	SetDIDAliasFn          func(ctx context.Context, p domain.SetDIDAliasPayload) (*domain.BlockchainResult, error)
	SetDIDOrgAliasFn       func(ctx context.Context, p domain.SetDIDOrgAliasPayload) (*domain.BlockchainResult, error)
	ExtendRegistrationFn    func(ctx context.Context, p domain.ExtendRegistrationPayload) (*domain.BlockchainResult, error)
	RevokeAliasFn           func(ctx context.Context, p domain.RevokeAliasPayload) (*domain.BlockchainResult, error)
	UnrevokeAliasFn         func(ctx context.Context, p domain.UnrevokeAliasPayload) (*domain.BlockchainResult, error)
	RegisterSubdomainFn     func(ctx context.Context, p domain.RegisterSubdomainPayload) (*domain.BlockchainResult, error)
	TransferAliasOwnershipFn func(ctx context.Context, p domain.TransferAliasOwnershipPayload) (*domain.BlockchainResult, error)
	TransferTLDFn           func(ctx context.Context, p domain.TransferTLDPayload) (*domain.BlockchainResult, error)
	SetAliasRootStorageFn   func(ctx context.Context, p domain.SetAliasAddressPayload) (*domain.BlockchainResult, error)
	SetAliasOrgStorageFn    func(ctx context.Context, p domain.SetAliasAddressPayload) (*domain.BlockchainResult, error)
	SetFactoryContractFn     func(ctx context.Context, p domain.SetFactoryContractPayload) (*domain.BlockchainResult, error)
}

func (m *MockBlockchainService) CreateDID(ctx context.Context, p domain.CreateDIDPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "CreateDID", Payload: p})
	m.mu.Unlock()
	if m.CreateDIDFn != nil {
		return m.CreateDIDFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_CREATE_DID"}}, nil
}

func (m *MockBlockchainService) AddKey(ctx context.Context, p domain.AddKeyPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "AddKey", Payload: p})
	m.mu.Unlock()
	if m.AddKeyFn != nil {
		return m.AddKeyFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_ADD_KEY"}}, nil
}

func (m *MockBlockchainService) AddClaim(ctx context.Context, p domain.AddClaimPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "AddClaim", Payload: p})
	m.mu.Unlock()
	if m.AddClaimFn != nil {
		return m.AddClaimFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_ADD_CLAIM"}}, nil
}

func (m *MockBlockchainService) StoreData(ctx context.Context, p domain.StoreDataPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "StoreData", Payload: p})
	m.mu.Unlock()
	if m.StoreDataFn != nil {
		return m.StoreDataFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_STORE_DATA"}}, nil
}

func (m *MockBlockchainService) UpdateDID(ctx context.Context, p domain.UpdateDIDPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "UpdateDID", Payload: p})
	m.mu.Unlock()
	if m.UpdateDIDFn != nil {
		return m.UpdateDIDFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_UPDATE_DID"}}, nil
}

func (m *MockBlockchainService) DeactivateDID(ctx context.Context, p domain.DIDLifecyclePayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "DeactivateDID", Payload: p})
	m.mu.Unlock()
	if m.DeactivateDIDFn != nil {
		return m.DeactivateDIDFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_DEACTIVATE_DID"}}, nil
}

func (m *MockBlockchainService) ReactivateDID(ctx context.Context, p domain.DIDLifecyclePayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "ReactivateDID", Payload: p})
	m.mu.Unlock()
	if m.ReactivateDIDFn != nil {
		return m.ReactivateDIDFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_REACTIVATE_DID"}}, nil
}

func (m *MockBlockchainService) TransferDIDOwner(ctx context.Context, p domain.TransferDIDOwnerPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "TransferDIDOwner", Payload: p})
	m.mu.Unlock()
	if m.TransferDIDOwnerFn != nil {
		return m.TransferDIDOwnerFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_TRANSFER_DID"}}, nil
}

func (m *MockBlockchainService) DeleteData(ctx context.Context, p domain.DeleteDataPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "DeleteData", Payload: p})
	m.mu.Unlock()
	if m.DeleteDataFn != nil {
		return m.DeleteDataFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_DELETE_DATA"}}, nil
}

func (m *MockBlockchainService) RemoveKey(ctx context.Context, p domain.RemoveKeyPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "RemoveKey", Payload: p})
	m.mu.Unlock()
	if m.RemoveKeyFn != nil {
		return m.RemoveKeyFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_REMOVE_KEY"}}, nil
}

func (m *MockBlockchainService) RemoveClaim(ctx context.Context, p domain.RemoveClaimPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "RemoveClaim", Payload: p})
	m.mu.Unlock()
	if m.RemoveClaimFn != nil {
		return m.RemoveClaimFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_REMOVE_CLAIM"}}, nil
}

func (m *MockBlockchainService) GeneralExecute(ctx context.Context, p domain.GeneralExecutePayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "GeneralExecute", Payload: p})
	m.mu.Unlock()
	if m.GeneralExecuteFn != nil {
		return m.GeneralExecuteFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_GENERAL_EXECUTE"}}, nil
}

func (m *MockBlockchainService) CreateOrg(ctx context.Context, p domain.CreateOrgPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "CreateOrg", Payload: p})
	m.mu.Unlock()
	if m.CreateOrgFn != nil {
		return m.CreateOrgFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_CREATE_ORG"}}, nil
}

func (m *MockBlockchainService) DeactivateOrg(ctx context.Context, p domain.OrgLifecyclePayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "DeactivateOrg", Payload: p})
	m.mu.Unlock()
	if m.DeactivateOrgFn != nil {
		return m.DeactivateOrgFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_DEACTIVATE_ORG"}}, nil
}

func (m *MockBlockchainService) ReactivateOrg(ctx context.Context, p domain.OrgLifecyclePayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "ReactivateOrg", Payload: p})
	m.mu.Unlock()
	if m.ReactivateOrgFn != nil {
		return m.ReactivateOrgFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_REACTIVATE_ORG"}}, nil
}

func (m *MockBlockchainService) TransferOrgOwner(ctx context.Context, p domain.OrgTransferPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "TransferOrgOwner", Payload: p})
	m.mu.Unlock()
	if m.TransferOrgOwnerFn != nil {
		return m.TransferOrgOwnerFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_TRANSFER_ORG"}}, nil
}

func (m *MockBlockchainService) AddMember(ctx context.Context, p domain.OrgMemberPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "AddMember", Payload: p})
	m.mu.Unlock()
	if m.AddMemberFn != nil {
		return m.AddMemberFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_ADD_MEMBER"}}, nil
}

func (m *MockBlockchainService) RemoveMember(ctx context.Context, p domain.OrgMemberPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "RemoveMember", Payload: p})
	m.mu.Unlock()
	if m.RemoveMemberFn != nil {
		return m.RemoveMemberFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_REMOVE_MEMBER"}}, nil
}

func (m *MockBlockchainService) UpdateMember(ctx context.Context, p domain.OrgMemberPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "UpdateMember", Payload: p})
	m.mu.Unlock()
	if m.UpdateMemberFn != nil {
		return m.UpdateMemberFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_UPDATE_MEMBER"}}, nil
}

func (m *MockBlockchainService) HandleOps(ctx context.Context, p domain.HandleOpsPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "HandleOps", Payload: p})
	m.mu.Unlock()
	if m.HandleOpsFn != nil {
		return m.HandleOpsFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_HANDLE_OPS"}}, nil
}

func (m *MockBlockchainService) DeployWallet(ctx context.Context, p domain.DeployWalletPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "DeployWallet", Payload: p})
	m.mu.Unlock()
	if m.DeployWalletFn != nil {
		return m.DeployWalletFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_DEPLOY_WALLET"}}, nil
}

func (m *MockBlockchainService) IssueCredential(ctx context.Context, p domain.IssueCredentialPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "IssueCredential", Payload: p})
	m.mu.Unlock()
	if m.IssueCredentialFn != nil {
		return m.IssueCredentialFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_ISSUE_VC"}}, nil
}

func (m *MockBlockchainService) BurnCredential(ctx context.Context, p domain.BurnCredentialPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "BurnCredential", Payload: p})
	m.mu.Unlock()
	if m.BurnCredentialFn != nil {
		return m.BurnCredentialFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_BURN_VC"}}, nil
}

func (m *MockBlockchainService) UpdateMetadata(ctx context.Context, p domain.UpdateMetadataPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "UpdateMetadata", Payload: p})
	m.mu.Unlock()
	if m.UpdateMetadataFn != nil {
		return m.UpdateMetadataFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_UPDATE_VC"}}, nil
}

func (m *MockBlockchainService) RevokeCredential(ctx context.Context, p domain.RevokeCredentialPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "RevokeCredential", Payload: p})
	m.mu.Unlock()
	if m.RevokeCredentialFn != nil {
		return m.RevokeCredentialFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_REVOKE_VC"}}, nil
}

func (m *MockBlockchainService) ApproveCredentialOrg(ctx context.Context, p domain.ApproveCredentialOrgPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "ApproveCredentialOrg", Payload: p})
	m.mu.Unlock()
	if m.ApproveCredentialOrgFn != nil {
		return m.ApproveCredentialOrgFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_APPROVE_ORG_VC"}}, nil
}

func (m *MockBlockchainService) ApproveCredential(ctx context.Context, p domain.ApproveCredentialPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "ApproveCredential", Payload: p})
	m.mu.Unlock()
	if m.ApproveCredentialFn != nil {
		return m.ApproveCredentialFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_APPROVE_USER_VC"}}, nil
}

func (m *MockBlockchainService) SetDidRootStorage(ctx context.Context, p domain.SetAddressPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "SetDidRootStorage", Payload: p})
	m.mu.Unlock()
	if m.SetDidRootStorageFn != nil {
		return m.SetDidRootStorageFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_SET_ROOT_STORAGE"}}, nil
}

func (m *MockBlockchainService) SetDidOrgStorage(ctx context.Context, p domain.SetAddressPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "SetDidOrgStorage", Payload: p})
	m.mu.Unlock()
	if m.SetDidOrgStorageFn != nil {
		return m.SetDidOrgStorageFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_SET_ORG_STORAGE"}}, nil
}

func (m *MockBlockchainService) RegisterTLD(ctx context.Context, p domain.RegisterTLDPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "RegisterTLD", Payload: p})
	m.mu.Unlock()
	if m.RegisterTLDFn != nil {
		return m.RegisterTLDFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_REGISTER_TLD"}}, nil
}

func (m *MockBlockchainService) RegisterDomain(ctx context.Context, p domain.RegisterDomainPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "RegisterDomain", Payload: p})
	m.mu.Unlock()
	if m.RegisterDomainFn != nil {
		return m.RegisterDomainFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_REGISTER_DOMAIN"}}, nil
}

func (m *MockBlockchainService) SetDIDAlias(ctx context.Context, p domain.SetDIDAliasPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "SetDIDAlias", Payload: p})
	m.mu.Unlock()
	if m.SetDIDAliasFn != nil {
		return m.SetDIDAliasFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_SET_DID_ALIAS"}}, nil
}

func (m *MockBlockchainService) SetDIDOrgAlias(ctx context.Context, p domain.SetDIDOrgAliasPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "SetDIDOrgAlias", Payload: p})
	m.mu.Unlock()
	if m.SetDIDOrgAliasFn != nil {
		return m.SetDIDOrgAliasFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_SET_DID_ORG_ALIAS"}}, nil
}

func (m *MockBlockchainService) ExtendRegistration(ctx context.Context, p domain.ExtendRegistrationPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "ExtendRegistration", Payload: p})
	m.mu.Unlock()
	if m.ExtendRegistrationFn != nil {
		return m.ExtendRegistrationFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_EXTEND_ALIAS"}}, nil
}

func (m *MockBlockchainService) RevokeAlias(ctx context.Context, p domain.RevokeAliasPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "RevokeAlias", Payload: p})
	m.mu.Unlock()
	if m.RevokeAliasFn != nil {
		return m.RevokeAliasFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_REVOKE_ALIAS"}}, nil
}

func (m *MockBlockchainService) UnrevokeAlias(ctx context.Context, p domain.UnrevokeAliasPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "UnrevokeAlias", Payload: p})
	m.mu.Unlock()
	if m.UnrevokeAliasFn != nil {
		return m.UnrevokeAliasFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_UNREVOKE_ALIAS"}}, nil
}

func (m *MockBlockchainService) RegisterSubdomain(ctx context.Context, p domain.RegisterSubdomainPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "RegisterSubdomain", Payload: p})
	m.mu.Unlock()
	if m.RegisterSubdomainFn != nil {
		return m.RegisterSubdomainFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_REGISTER_SUBDOMAIN"}}, nil
}

func (m *MockBlockchainService) TransferAliasOwnership(ctx context.Context, p domain.TransferAliasOwnershipPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "TransferAliasOwnership", Payload: p})
	m.mu.Unlock()
	if m.TransferAliasOwnershipFn != nil {
		return m.TransferAliasOwnershipFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_TRANSFER_ALIAS"}}, nil
}

func (m *MockBlockchainService) TransferTLD(ctx context.Context, p domain.TransferTLDPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "TransferTLD", Payload: p})
	m.mu.Unlock()
	if m.TransferTLDFn != nil {
		return m.TransferTLDFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_TRANSFER_TLD"}}, nil
}

func (m *MockBlockchainService) SetAliasRootStorage(ctx context.Context, p domain.SetAliasAddressPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "SetAliasRootStorage", Payload: p})
	m.mu.Unlock()
	if m.SetAliasRootStorageFn != nil {
		return m.SetAliasRootStorageFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_SET_ALIAS_ROOT_STORAGE"}}, nil
}

func (m *MockBlockchainService) SetAliasOrgStorage(ctx context.Context, p domain.SetAliasAddressPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "SetAliasOrgStorage", Payload: p})
	m.mu.Unlock()
	if m.SetAliasOrgStorageFn != nil {
		return m.SetAliasOrgStorageFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_SET_ALIAS_ORG_STORAGE"}}, nil
}

func (m *MockBlockchainService) SetFactoryContract(ctx context.Context, p domain.SetFactoryContractPayload) (*domain.BlockchainResult, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, Call{Method: "SetFactoryContract", Payload: p})
	m.mu.Unlock()
	if m.SetFactoryContractFn != nil {
		return m.SetFactoryContractFn(ctx, p)
	}
	return &domain.BlockchainResult{TxHashes: []string{"0xMOCK_SET_FACTORY_CONTRACT"}}, nil
}

func (m *MockBlockchainService) CallCount(method string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := 0
	for _, c := range m.Calls {
		if c.Method == method {
			n++
		}
	}
	return n
}
