package sdk

import (
	"context"
	"fmt"
	"math/big"

	harachain "github.com/meQlause/hara-core-blockchain-lib/pkg/blockchain"
	haranetwork "github.com/meQlause/hara-core-blockchain-lib/pkg/network"
	harawallet "github.com/meQlause/hara-core-blockchain-lib/pkg/wallet"
	harautils "github.com/meQlause/hara-core-blockchain-lib/utils"

	"github.com/myorg/worker-service/internal/config"
)

// Provider holds the shared blockchain infrastructure components.
type Provider struct {
	Network    *haranetwork.Network
	Wallet     *harawallet.Wallet
	Chain      *harachain.Blockchain
	WalletAddr harautils.Address
}

// NewProvider initializes the network, wallet, and blockchain layers.
func NewProvider(cfg config.BlockchainConfig) (*Provider, error) {
	// ── 1. Network ─────────────────────────────────────────────────────────
	network := haranetwork.NewNetwork(
		cfg.RPCURLs,
		"1.0",
		0,
		harautils.LogConfig{},
	)

	initCtx := context.Background()
	if !network.IsOnline(initCtx) {
		return nil, fmt.Errorf("all configured RPC endpoints are unreachable: %v", cfg.RPCURLs)
	}

	// ── 2. Chain ID ───────────────────────────────────────────────────────
	chainID, err := network.ChainID(initCtx)
	if err != nil {
		return nil, fmt.Errorf("fetch chain id: %w", err)
	}

	// ── 3. Wallet ──────────────────────────────────────────────────────────
	wallet := harawallet.NewWallet(cfg.PrivateKey)
	walletAddr, err := wallet.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("derive wallet address from private key: %w", err)
	}

	// ── 4. Blockchain ──────────────────────────────────────────────────────
	chain := harachain.NewBlockchain(network, new(big.Int).SetUint64(chainID))

	return &Provider{
		Network:    network,
		Wallet:     wallet,
		Chain:      chain,
		WalletAddr: walletAddr,
	}, nil
}
