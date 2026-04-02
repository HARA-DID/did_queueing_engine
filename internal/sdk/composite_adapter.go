package sdk

import (
	"context"

	"github.com/myorg/worker-service/internal/domain"
	"github.com/myorg/worker-service/internal/service"
)

// Compile-time interface check.
var _ service.BlockchainService = (*CompositeAdapter)(nil)

// CompositeAdapter delegates calls to specialized adapters (DID, AA, etc.).
type CompositeAdapter struct {
	did   *DIDAdapter
	aa    *AAAdapter
	vc    *VCAdapter
	alias *AliasAdapter
}

// NewCompositeAdapter creates a new composite adapter.
func NewCompositeAdapter(did *DIDAdapter, aa *AAAdapter, vc *VCAdapter, alias *AliasAdapter) *CompositeAdapter {
	return &CompositeAdapter{
		did:   did,
		aa:    aa,
		vc:    vc,
		alias: alias,
	}
}

// ── DID Operations ──────────────────────────────────────────────

func (c *CompositeAdapter) CreateDID(ctx context.Context, p domain.CreateDIDPayload) (*domain.BlockchainResult, error) {
	return c.did.CreateDID(ctx, p)
}

func (c *CompositeAdapter) AddKey(ctx context.Context, p domain.AddKeyPayload) (*domain.BlockchainResult, error) {
	return c.did.AddKey(ctx, p)
}

func (c *CompositeAdapter) AddClaim(ctx context.Context, p domain.AddClaimPayload) (*domain.BlockchainResult, error) {
	return c.did.AddClaim(ctx, p)
}

func (c *CompositeAdapter) StoreData(ctx context.Context, p domain.StoreDataPayload) (*domain.BlockchainResult, error) {
	return c.did.StoreData(ctx, p)
}

// ── Account Abstraction Operations ──────────────────────────────

func (c *CompositeAdapter) HandleOps(ctx context.Context, p domain.HandleOpsPayload) (*domain.BlockchainResult, error) {
	return c.aa.HandleOps(ctx, p)
}

func (c *CompositeAdapter) DeployWallet(ctx context.Context, p domain.DeployWalletPayload) (*domain.BlockchainResult, error) {
	return c.aa.DeployWallet(ctx, p)
}

// ── Verifiable Credentials Operations ───────────────────────────

func (c *CompositeAdapter) IssueCredential(ctx context.Context, p domain.IssueCredentialPayload) (*domain.BlockchainResult, error) {
	return c.vc.IssueCredential(ctx, p)
}

func (c *CompositeAdapter) BurnCredential(ctx context.Context, p domain.BurnCredentialPayload) (*domain.BlockchainResult, error) {
	return c.vc.BurnCredential(ctx, p)
}

func (c *CompositeAdapter) UpdateMetadata(ctx context.Context, p domain.UpdateMetadataPayload) (*domain.BlockchainResult, error) {
	return c.vc.UpdateMetadata(ctx, p)
}

func (c *CompositeAdapter) RevokeCredential(ctx context.Context, p domain.RevokeCredentialPayload) (*domain.BlockchainResult, error) {
	return c.vc.RevokeCredential(ctx, p)
}

func (c *CompositeAdapter) ApproveCredentialOrg(ctx context.Context, p domain.ApproveCredentialOrgPayload) (*domain.BlockchainResult, error) {
	return c.vc.ApproveCredentialOrg(ctx, p)
}

func (c *CompositeAdapter) ApproveCredential(ctx context.Context, p domain.ApproveCredentialPayload) (*domain.BlockchainResult, error) {
	return c.vc.ApproveCredential(ctx, p)
}

func (c *CompositeAdapter) SetDidRootStorage(ctx context.Context, p domain.SetAddressPayload) (*domain.BlockchainResult, error) {
	return c.vc.SetDidRootStorage(ctx, p)
}

func (c *CompositeAdapter) SetDidOrgStorage(ctx context.Context, p domain.SetAddressPayload) (*domain.BlockchainResult, error) {
	return c.vc.SetDidOrgStorage(ctx, p)
}

// ── Alias Operations ──────────────────────────────────────────────

func (c *CompositeAdapter) RegisterTLD(ctx context.Context, p domain.RegisterTLDPayload) (*domain.BlockchainResult, error) {
	return c.alias.RegisterTLD(ctx, p)
}

func (c *CompositeAdapter) RegisterDomain(ctx context.Context, p domain.RegisterDomainPayload) (*domain.BlockchainResult, error) {
	return c.alias.RegisterDomain(ctx, p)
}

func (c *CompositeAdapter) SetDIDAlias(ctx context.Context, p domain.SetDIDAliasPayload) (*domain.BlockchainResult, error) {
	return c.alias.SetDIDAlias(ctx, p)
}

func (c *CompositeAdapter) SetDIDOrgAlias(ctx context.Context, p domain.SetDIDOrgAliasPayload) (*domain.BlockchainResult, error) {
	return c.alias.SetDIDOrgAlias(ctx, p)
}

func (c *CompositeAdapter) ExtendRegistration(ctx context.Context, p domain.ExtendRegistrationPayload) (*domain.BlockchainResult, error) {
	return c.alias.ExtendRegistration(ctx, p)
}

func (c *CompositeAdapter) RevokeAlias(ctx context.Context, p domain.RevokeAliasPayload) (*domain.BlockchainResult, error) {
	return c.alias.RevokeAlias(ctx, p)
}

func (c *CompositeAdapter) UnrevokeAlias(ctx context.Context, p domain.UnrevokeAliasPayload) (*domain.BlockchainResult, error) {
	return c.alias.UnrevokeAlias(ctx, p)
}

func (c *CompositeAdapter) RegisterSubdomain(ctx context.Context, p domain.RegisterSubdomainPayload) (*domain.BlockchainResult, error) {
	return c.alias.RegisterSubdomain(ctx, p)
}

func (c *CompositeAdapter) TransferAliasOwnership(ctx context.Context, p domain.TransferAliasOwnershipPayload) (*domain.BlockchainResult, error) {
	return c.alias.TransferAliasOwnership(ctx, p)
}

func (c *CompositeAdapter) TransferTLD(ctx context.Context, p domain.TransferTLDPayload) (*domain.BlockchainResult, error) {
	return c.alias.TransferTLD(ctx, p)
}

func (c *CompositeAdapter) SetAliasRootStorage(ctx context.Context, p domain.SetAddressPayload) (*domain.BlockchainResult, error) {
	return c.alias.SetAliasRootStorage(ctx, p)
}

func (c *CompositeAdapter) SetAliasOrgStorage(ctx context.Context, p domain.SetAddressPayload) (*domain.BlockchainResult, error) {
	return c.alias.SetAliasOrgStorage(ctx, p)
}
