package sdk

import (
	"context"
	"fmt"

	vcfact "github.com/HARA-DID/did-verifiable-credentials-sdk/pkg/vcfactory"
	vcstor "github.com/HARA-DID/did-verifiable-credentials-sdk/pkg/vcstorage"
	harautils "github.com/meQlause/hara-core-blockchain-lib/utils"

	"github.com/myorg/worker-service/internal/config"
	"github.com/myorg/worker-service/internal/domain"
)

// VCAdapter implements Verifiable Credentials related blockchain operations.
type VCAdapter struct {
	provider *Provider
	factory  *vcfact.VCFactory
	storage  *vcstor.VCStorage
}

// NewVCAdapter initializes the VC SDK components.
func NewVCAdapter(p *Provider, cfg config.BlockchainConfig) (*VCAdapter, error) {
	initCtx := context.Background()

	// ── Factory ────────────────────────────────────────────────────
	factory, err := vcfact.NewVCFactoryWithHNS(initCtx, cfg.VCFactoryHNS, p.Chain)
	if err != nil {
		if cfg.VCFactoryAddress != "" {
			factory = vcfact.NewVCFactory(
				harautils.HexToAddress(cfg.VCFactoryAddress),
				harautils.ABI{},
				p.Chain,
				nil,
			)
		} else {
			return nil, fmt.Errorf("resolve VCFactory via HNS %q: %w", cfg.VCFactoryHNS, err)
		}
	}

	// ── Storage ────────────────────────────────────────────────────
	storage, err := vcstor.NewVCStorageWithHNS(initCtx, cfg.VCStorageHNS, p.Chain)
	if err != nil {
		if cfg.VCStorageAddress != "" {
			storage = vcstor.NewVCStorage(
				harautils.HexToAddress(cfg.VCStorageAddress),
				harautils.ABI{},
				p.Chain,
				nil,
			)
		} else {
			return nil, fmt.Errorf("resolve VCStorage via HNS %q: %w", cfg.VCStorageHNS, err)
		}
	}

	return &VCAdapter{
		provider: p,
		factory:  factory,
		storage:  storage,
	}, nil
}

// ── BlockchainService implementation for VC ──────────────────────

func (a *VCAdapter) IssueCredential(ctx context.Context, p domain.IssueCredentialPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.IssueCredential(ctx, a.provider.Wallet, vcfact.IssueCredentialParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("vcFactory.IssueCredential: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *VCAdapter) BurnCredential(ctx context.Context, p domain.BurnCredentialPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.BurnCredential(ctx, a.provider.Wallet, vcfact.BurnCredentialParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("vcFactory.BurnCredential: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *VCAdapter) UpdateMetadata(ctx context.Context, p domain.UpdateMetadataPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.UpdateMetadata(ctx, a.provider.Wallet, vcfact.UpdateMetadataParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("vcFactory.UpdateMetadata: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *VCAdapter) RevokeCredential(ctx context.Context, p domain.RevokeCredentialPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.RevokeCredential(ctx, a.provider.Wallet, vcfact.RevokeCredentialParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("vcFactory.RevokeCredential: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *VCAdapter) ApproveCredentialOrg(ctx context.Context, p domain.ApproveCredentialOrgPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.ApproveCredentialOrg(ctx, a.provider.Wallet, vcfact.ApproveCredentialOrgParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("vcFactory.ApproveCredentialOrg: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *VCAdapter) ApproveCredential(ctx context.Context, p domain.ApproveCredentialPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.ApproveCredential(ctx, a.provider.Wallet, vcfact.ApproveCredentialParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("vcFactory.ApproveCredential: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *VCAdapter) SetDidRootStorage(ctx context.Context, p domain.SetAddressPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.storage.SetDidRootStorage(ctx, a.provider.Wallet, vcstor.SetAddressParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("vcStorage.SetDidRootStorage: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *VCAdapter) SetDidOrgStorage(ctx context.Context, p domain.SetAddressPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.storage.SetDidOrgStorage(ctx, a.provider.Wallet, vcstor.SetAddressParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("vcStorage.SetDidOrgStorage: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}
