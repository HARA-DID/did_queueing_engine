package sdk

import (
	"context"
	"fmt"

	aliasfact "github.com/HARA-DID/alias-root-sdk/pkg/aliasfactory"
	aliasstor "github.com/HARA-DID/alias-root-sdk/pkg/aliasstorage"
	harautils "github.com/meQlause/hara-core-blockchain-lib/utils"

	"github.com/myorg/worker-service/internal/config"
	"github.com/myorg/worker-service/internal/domain"
)

// AliasAdapter implements Alias related blockchain operations.
type AliasAdapter struct {
	provider *Provider
	factory  *aliasfact.AliasFactory
	storage  *aliasstor.AliasStorage
}

// NewAliasAdapter initializes the Alias SDK components.
func NewAliasAdapter(p *Provider, cfg config.BlockchainConfig) (*AliasAdapter, error) {
	initCtx := context.Background()

	// ── Factory ────────────────────────────────────────────────────
	factory, err := aliasfact.NewAliasFactoryWithHNS(initCtx, cfg.AliasFactoryHNS, p.Chain)
	if err != nil {
		if cfg.AliasFactoryAddress != "" {
			factory = aliasfact.NewAliasFactory(
				harautils.HexToAddress(cfg.AliasFactoryAddress),
				harautils.ABI{},
				p.Chain,
				nil,
			)
		} else {
			return nil, fmt.Errorf("resolve AliasFactory via HNS %q: %w", cfg.AliasFactoryHNS, err)
		}
	}

	// ── Storage ────────────────────────────────────────────────────
	storage, err := aliasstor.NewAliasStorageWithHNS(initCtx, cfg.AliasStorageHNS, p.Chain)
	if err != nil {
		if cfg.AliasStorageAddress != "" {
			storage = aliasstor.NewAliasStorage(
				harautils.HexToAddress(cfg.AliasStorageAddress),
				harautils.ABI{},
				p.Chain,
				nil,
			)
		} else {
			return nil, fmt.Errorf("resolve AliasStorage via HNS %q: %w", cfg.AliasStorageHNS, err)
		}
	}

	return &AliasAdapter{
		provider: p,
		factory:  factory,
		storage:  storage,
	}, nil
}

// ── BlockchainService implementation for Alias ───────────────────

func (a *AliasAdapter) RegisterTLD(ctx context.Context, p domain.RegisterTLDPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.RegisterTLD(ctx, a.provider.Wallet, aliasfact.RegisterTLDParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("aliasFactory.RegisterTLD: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *AliasAdapter) RegisterDomain(ctx context.Context, p domain.RegisterDomainPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.RegisterDomain(ctx, a.provider.Wallet, aliasfact.RegisterDomainParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("aliasFactory.RegisterDomain: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *AliasAdapter) SetDIDAlias(ctx context.Context, p domain.SetDIDAliasPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.SetDID(ctx, a.provider.Wallet, aliasfact.SetDIDParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("aliasFactory.SetDID: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *AliasAdapter) SetDIDOrgAlias(ctx context.Context, p domain.SetDIDOrgAliasPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.SetDIDOrg(ctx, a.provider.Wallet, aliasfact.SetDIDOrgParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("aliasFactory.SetDIDOrg: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *AliasAdapter) ExtendRegistration(ctx context.Context, p domain.ExtendRegistrationPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.ExtendRegistration(ctx, a.provider.Wallet, aliasfact.ExtendRegistrationParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("aliasFactory.ExtendRegistration: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *AliasAdapter) RevokeAlias(ctx context.Context, p domain.RevokeAliasPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.RevokeAlias(ctx, a.provider.Wallet, aliasfact.NodeOnlyParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("aliasFactory.RevokeAlias: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *AliasAdapter) UnrevokeAlias(ctx context.Context, p domain.UnrevokeAliasPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.UnrevokeAlias(ctx, a.provider.Wallet, aliasfact.NodeOnlyParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("aliasFactory.UnrevokeAlias: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *AliasAdapter) RegisterSubdomain(ctx context.Context, p domain.RegisterSubdomainPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.RegisterSubdomain(ctx, a.provider.Wallet, aliasfact.RegisterSubdomainParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("aliasFactory.RegisterSubdomain: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *AliasAdapter) TransferAliasOwnership(ctx context.Context, p domain.TransferAliasOwnershipPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.TransferAliasOwnership(ctx, a.provider.Wallet, aliasfact.TransferAliasOwnershipParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("aliasFactory.TransferAliasOwnership: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *AliasAdapter) TransferTLD(ctx context.Context, p domain.TransferTLDPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.TransferTLD(ctx, a.provider.Wallet, aliasfact.TransferTLDParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("aliasFactory.TransferTLD: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *AliasAdapter) SetAliasRootStorage(ctx context.Context, p domain.SetAddressPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.SetDIDRootStorage(ctx, a.provider.Wallet, aliasfact.SetDIDRootStorageParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("aliasFactory.SetDIDRootStorage: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *AliasAdapter) SetAliasOrgStorage(ctx context.Context, p domain.SetAddressPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.factory.SetDIDOrgStorage(ctx, a.provider.Wallet, aliasfact.SetDIDOrgStorageParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("aliasFactory.SetDIDOrgStorage: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}
