package service

import (
	"context"

	"github.com/myorg/worker-service/internal/domain"
)

// BlockchainService is the abstraction over the DID + Blockchain SDKs.
// All worker code calls this interface; no SDK import appears outside /sdk.
type BlockchainService interface {
	// ── DID Operations ──────────────────────────────────────────────
	CreateDID(ctx context.Context, payload domain.CreateDIDPayload) (*domain.BlockchainResult, error)
	AddKey(ctx context.Context, payload domain.AddKeyPayload) (*domain.BlockchainResult, error)
	AddClaim(ctx context.Context, payload domain.AddClaimPayload) (*domain.BlockchainResult, error)
	StoreData(ctx context.Context, payload domain.StoreDataPayload) (*domain.BlockchainResult, error)

	// ── Account Abstraction Operations ──────────────────────────────
	HandleOps(ctx context.Context, payload domain.HandleOpsPayload) (*domain.BlockchainResult, error)
	DeployWallet(ctx context.Context, payload domain.DeployWalletPayload) (*domain.BlockchainResult, error)

	// ── Verifiable Credentials Operations ───────────────────────────
	IssueCredential(ctx context.Context, payload domain.IssueCredentialPayload) (*domain.BlockchainResult, error)
	BurnCredential(ctx context.Context, payload domain.BurnCredentialPayload) (*domain.BlockchainResult, error)
	UpdateMetadata(ctx context.Context, payload domain.UpdateMetadataPayload) (*domain.BlockchainResult, error)
	RevokeCredential(ctx context.Context, payload domain.RevokeCredentialPayload) (*domain.BlockchainResult, error)
	ApproveCredentialOrg(ctx context.Context, payload domain.ApproveCredentialOrgPayload) (*domain.BlockchainResult, error)
	ApproveCredential(ctx context.Context, payload domain.ApproveCredentialPayload) (*domain.BlockchainResult, error)
	SetDidRootStorage(ctx context.Context, payload domain.SetAddressPayload) (*domain.BlockchainResult, error)
	SetDidOrgStorage(ctx context.Context, payload domain.SetAddressPayload) (*domain.BlockchainResult, error)

	// ── Alias Operations ──────────────────────────────────────────────
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
	SetAliasRootStorage(ctx context.Context, payload domain.SetAddressPayload) (*domain.BlockchainResult, error)
	SetAliasOrgStorage(ctx context.Context, payload domain.SetAddressPayload) (*domain.BlockchainResult, error)
}
