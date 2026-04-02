package sdk

import (
	"context"
	"fmt"

	didfactory "github.com/HARA-DID/did-root-sdk/pkg/factory"
	haracontract "github.com/meQlause/hara-core-blockchain-lib/pkg/contract"
	harautils "github.com/meQlause/hara-core-blockchain-lib/utils"

	"github.com/myorg/worker-service/internal/config"
	"github.com/myorg/worker-service/internal/domain"
)

// DIDAdapter implements DID-related blockchain operations.
type DIDAdapter struct {
	provider *Provider
	factory  *didfactory.Factory
}

// NewDIDAdapter initializes the DID SDK factory with common resources.
func NewDIDAdapter(p *Provider, cfg config.BlockchainConfig) (*DIDAdapter, error) {
	initCtx := context.Background()
	var contract *haracontract.Contract
	var err error

	// ── Contract Resolution ──────────────────────────────────────────
	if cfg.HNSName != "" {
		contract, err = p.Chain.ContractWithHNS(initCtx, cfg.HNSName)
		if err != nil {
			return nil, fmt.Errorf("resolve contract via HNS %q: %w", cfg.HNSName, err)
		}
	} else if cfg.ContractAddress != "" && cfg.ContractABI != "" {
		parsedABI, parseErr := harautils.ParseABI(cfg.ContractABI)
		if parseErr != nil {
			return nil, fmt.Errorf("parse contract ABI: %w", parseErr)
		}
		contract = p.Chain.Contract(parsedABI, harautils.HexToAddress(cfg.ContractAddress))
	} else {
		return nil, fmt.Errorf("DID config requires either HNS_NAME or both CONTRACT_ADDRESS and CONTRACT_ABI")
	}

	// ── Factory ─────────────────────────────────────────────────────
	factory := didfactory.NewFactory(
		p.WalletAddr,
		harautils.ABI{}, 
		p.Chain,
		contract,
	)

	return &DIDAdapter{
		provider: p,
		factory:  factory,
	}, nil
}

// CreateDID registers a new DID on-chain.
func (a *DIDAdapter) CreateDID(ctx context.Context, p domain.CreateDIDPayload) (*domain.BlockchainResult, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	txHashes, err := a.factory.CreateDID(
		ctx,
		a.provider.Wallet,
		didfactory.CreateDIDParam{DID: p.DID},
		keyID,
		p.MultipleRPCCalls,
	)
	if err != nil {
		return nil, fmt.Errorf("factory.CreateDID: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

// AddKey adds a verification key to an existing DID document.
func (a *DIDAdapter) AddKey(ctx context.Context, p domain.AddKeyPayload) (*domain.BlockchainResult, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	txHashes, err := a.factory.AddKey(
		ctx,
		a.provider.Wallet,
		didfactory.StoreKeyParams{
			DIDIndex: p.DIDIndex,
			KeyType:  p.KeyType,
			Purpose:  p.Purpose,
		},
		keyID,
		p.MultipleRPCCalls,
	)
	if err != nil {
		return nil, fmt.Errorf("factory.AddKey: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

// AddClaim attaches a verifiable claim to a DID document.
func (a *DIDAdapter) AddClaim(ctx context.Context, p domain.AddClaimPayload) (*domain.BlockchainResult, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	txHashes, err := a.factory.AddClaim(
		ctx,
		a.provider.Wallet,
		didfactory.StoreClaimParams{
			DIDIndex: p.DIDIndex,
			Topic:    p.Topic,
			Data:     p.Data,
		},
		keyID,
		p.MultipleRPCCalls,
	)
	if err != nil {
		return nil, fmt.Errorf("factory.AddClaim: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

// StoreData stores arbitrary binary data linked to a DID property key.
func (a *DIDAdapter) StoreData(ctx context.Context, p domain.StoreDataPayload) (*domain.BlockchainResult, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	txHashes, err := a.factory.StoreData(
		ctx,
		a.provider.Wallet,
		didfactory.StoreDataParams{
			DIDIndex: p.DIDIndex,
		},
		keyID,
		p.MultipleRPCCalls,
	)
	if err != nil {
		return nil, fmt.Errorf("factory.StoreData: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

// ---------------------------------------------------------------------------
// DID Helpers
// ---------------------------------------------------------------------------

func resolveKeyIdentifier(provided string) (string, error) {
	if provided != "" {
		return provided, nil
	}
	keyID, err := didfactory.GenerateKeyIdentifier()
	if err != nil {
		return "", fmt.Errorf("generate key identifier: %w", err)
	}
	return keyID, nil
}
