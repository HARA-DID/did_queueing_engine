package sdk

import (
	"context"
	"fmt"

	aapkg "github.com/HARA-DID/account-abstraction-sdk/pkg/entrypoint"
	"github.com/HARA-DID/account-abstraction-sdk/pkg/gasmanager"
	"github.com/HARA-DID/account-abstraction-sdk/pkg/walletfactory"
	harautils "github.com/meQlause/hara-core-blockchain-lib/utils"

	"github.com/myorg/worker-service/internal/config"
	"github.com/myorg/worker-service/internal/domain"
)

// AAAdapter implements Account Abstraction related blockchain operations.
type AAAdapter struct {
	provider     *Provider
	entryPoint   *aapkg.EntryPoint
	gasManager   *gasmanager.GasManager
	walletFactory *walletfactory.WalletFactory
}

// NewAAAdapter initializes the AA SDK components.
func NewAAAdapter(p *Provider, cfg config.BlockchainConfig) (*AAAdapter, error) {
	initCtx := context.Background()

	// ── EntryPoint ──────────────────────────────────────────────────
	entryPoint, err := aapkg.NewEntryPointWithHNS(initCtx, cfg.EntryPointHNS, p.Chain)
	if err != nil {
		if cfg.EntryPointAddress != "" {
			entryPoint = aapkg.NewEntryPoint(
				harautils.HexToAddress(cfg.EntryPointAddress),
				harautils.ABI{}, 
				p.Chain,
				nil, 
			)
		} else {
			return nil, fmt.Errorf("resolve EntryPoint via HNS %q: %w", cfg.EntryPointHNS, err)
		}
	}

	// ── GasManager ─────────────────────────────────────────────────
	gasMgr := gasmanager.NewGasManager(
		harautils.HexToAddress(cfg.GasManagerAddress),
		harautils.ABI{},
		p.Chain,
		nil,
	)

	// ── WalletFactory ──────────────────────────────────────────────
	walletFact := walletfactory.NewWalletFactory(
		harautils.HexToAddress(cfg.FactoryAddress),
		harautils.ABI{},
		p.Chain,
		nil,
	)

	return &AAAdapter{
		provider:      p,
		entryPoint:    entryPoint,
		gasManager:    gasMgr,
		walletFactory: walletFact,
	}, nil
}

// ── BlockchainService implementation for AA ──────────────────────

func (a *AAAdapter) HandleOps(ctx context.Context, p domain.HandleOpsPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.entryPoint.HandleOps(ctx, a.provider.Wallet, aapkg.HandleOpsParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("entryPoint.HandleOps: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *AAAdapter) DeployWallet(ctx context.Context, p domain.DeployWalletPayload) (*domain.BlockchainResult, error) {
	txHashes, err := a.walletFactory.DeployWallet(ctx, a.provider.Wallet, walletfactory.DeployWalletParams{}, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("walletFactory.DeployWallet: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

// DID Stub methods (to be called via the Composite adapter, but implemented here for completeness if needed)
func (a *AAAdapter) CreateDID(ctx context.Context, p domain.CreateDIDPayload) (*domain.BlockchainResult, error) {
	return nil, fmt.Errorf("not implemented in AA adapter")
}
func (a *AAAdapter) AddKey(ctx context.Context, p domain.AddKeyPayload) (*domain.BlockchainResult, error) {
	return nil, fmt.Errorf("not implemented in AA adapter")
}
func (a *AAAdapter) AddClaim(ctx context.Context, p domain.AddClaimPayload) (*domain.BlockchainResult, error) {
	return nil, fmt.Errorf("not implemented in AA adapter")
}
func (a *AAAdapter) StoreData(ctx context.Context, p domain.StoreDataPayload) (*domain.BlockchainResult, error) {
	return nil, fmt.Errorf("not implemented in AA adapter")
}
