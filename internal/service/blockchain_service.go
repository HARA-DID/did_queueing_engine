package service

import (
	"context"

	"github.com/myorg/worker-service/internal/domain"
)

type BlockchainService interface {
	CreateDID(ctx context.Context, payload domain.CreateDIDPayload) (*domain.BlockchainResult, error)
	AddKey(ctx context.Context, payload domain.AddKeyPayload) (*domain.BlockchainResult, error)
	AddClaim(ctx context.Context, payload domain.AddClaimPayload) (*domain.BlockchainResult, error)
	StoreData(ctx context.Context, payload domain.StoreDataPayload) (*domain.BlockchainResult, error)
	UpdateDID(ctx context.Context, payload domain.UpdateDIDPayload) (*domain.BlockchainResult, error)
	DeactivateDID(ctx context.Context, payload domain.DIDLifecyclePayload) (*domain.BlockchainResult, error)
	ReactivateDID(ctx context.Context, payload domain.DIDLifecyclePayload) (*domain.BlockchainResult, error)
	TransferDIDOwner(ctx context.Context, payload domain.TransferDIDOwnerPayload) (*domain.BlockchainResult, error)
	DeleteData(ctx context.Context, payload domain.DeleteDataPayload) (*domain.BlockchainResult, error)
	RemoveKey(ctx context.Context, payload domain.RemoveKeyPayload) (*domain.BlockchainResult, error)
	RemoveClaim(ctx context.Context, payload domain.RemoveClaimPayload) (*domain.BlockchainResult, error)
	GeneralExecute(ctx context.Context, payload domain.GeneralExecutePayload) (*domain.BlockchainResult, error)

	CreateOrg(ctx context.Context, payload domain.CreateOrgPayload) (*domain.BlockchainResult, error)
	DeactivateOrg(ctx context.Context, payload domain.OrgLifecyclePayload) (*domain.BlockchainResult, error)
	ReactivateOrg(ctx context.Context, payload domain.OrgLifecyclePayload) (*domain.BlockchainResult, error)
	TransferOrgOwner(ctx context.Context, payload domain.OrgTransferPayload) (*domain.BlockchainResult, error)
	AddMember(ctx context.Context, payload domain.OrgMemberPayload) (*domain.BlockchainResult, error)
	RemoveMember(ctx context.Context, payload domain.OrgMemberPayload) (*domain.BlockchainResult, error)
	UpdateMember(ctx context.Context, payload domain.OrgMemberPayload) (*domain.BlockchainResult, error)

	HandleOps(ctx context.Context, payload domain.HandleOpsPayload) (*domain.BlockchainResult, error)
	DeployWallet(ctx context.Context, payload domain.DeployWalletPayload) (*domain.BlockchainResult, error)

	IssueCredential(ctx context.Context, payload domain.IssueCredentialPayload) (*domain.BlockchainResult, error)
	BurnCredential(ctx context.Context, payload domain.BurnCredentialPayload) (*domain.BlockchainResult, error)
	UpdateMetadata(ctx context.Context, payload domain.UpdateMetadataPayload) (*domain.BlockchainResult, error)
	RevokeCredential(ctx context.Context, payload domain.RevokeCredentialPayload) (*domain.BlockchainResult, error)
	ApproveCredentialOrg(ctx context.Context, payload domain.ApproveCredentialOrgPayload) (*domain.BlockchainResult, error)
	ApproveCredential(ctx context.Context, payload domain.ApproveCredentialPayload) (*domain.BlockchainResult, error)
	SetDidRootStorage(ctx context.Context, payload domain.SetAddressPayload) (*domain.BlockchainResult, error)
	SetDidOrgStorage(ctx context.Context, payload domain.SetAddressPayload) (*domain.BlockchainResult, error)

	RegisterTLD(ctx context.Context, payload domain.RegisterTLDPayload) (*domain.BlockchainResult, error)
	RegisterDomain(ctx context.Context, payload domain.RegisterDomainPayload) (*domain.BlockchainResult, error)
	SetDIDAlias(ctx context.Context, payload domain.SetDIDAliasPayload) (*domain.BlockchainResult, error)
	SetDIDOrgAlias(ctx context.Context, payload domain.SetDIDOrgAliasPayload) (*domain.BlockchainResult, error)
	ExtendRegistration(ctx context.Context, payload domain.ExtendRegistrationPayload) (*domain.BlockchainResult, error)
	RevokeAlias(ctx context.Context, payload domain.RevokeAliasPayload) (*domain.BlockchainResult, error)
	UnrevokeAlias(ctx context.Context, payload domain.UnrevokeAliasPayload) (*domain.BlockchainResult, error)
	RegisterSubdomain(ctx context.Context, payload domain.RegisterSubdomainPayload) (*domain.BlockchainResult, error)
	TransferAliasOwnership(ctx context.Context, payload domain.TransferAliasOwnershipPayload) (*domain.BlockchainResult, error)
	TransferTLD(ctx context.Context, payload domain.TransferTLDPayload) (*domain.BlockchainResult, error)
	SetAliasRootStorage(ctx context.Context, payload domain.SetAliasAddressPayload) (*domain.BlockchainResult, error)
	SetAliasOrgStorage(ctx context.Context, payload domain.SetAliasAddressPayload) (*domain.BlockchainResult, error)
	SetFactoryContract(ctx context.Context, payload domain.SetFactoryContractPayload) (*domain.BlockchainResult, error)
}
