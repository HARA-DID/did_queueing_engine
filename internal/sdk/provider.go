package sdk

import (
	"context"
	"fmt"
	"math/big"

	harachain "github.com/HARA-DID/hara-core-blockchain-lib/pkg/blockchain"
	haranetwork "github.com/HARA-DID/hara-core-blockchain-lib/pkg/network"
	harawallet "github.com/HARA-DID/hara-core-blockchain-lib/pkg/wallet"
	harautils "github.com/HARA-DID/hara-core-blockchain-lib/utils"

	"github.com/HARA-DID/did-queueing-engine/internal/config"
)

// Provider holds the shared blockchain infrastructure components.
type Provider struct {
	Network    *haranetwork.Network
	Chain      *harachain.Blockchain
	Wallet     *harawallet.Wallet
	WalletAddr harautils.Address
}

func NewProvider(cfg config.BlockchainConfig, privateKey string) (*Provider, error) {
	initCtx := context.Background()

	network := haranetwork.NewNetwork(
		cfg.RPCURLs,
		"1.0",
		0,
		harautils.LogConfig{},
	)

	if !network.IsOnline(initCtx) {
		return nil, fmt.Errorf("all configured RPC endpoints are unreachable: %v", cfg.RPCURLs)
	}

	chainID, err := network.ChainID(initCtx)
	if err != nil {
		return nil, fmt.Errorf("fetch chain id: %w", err)
	}

	wallet := harawallet.NewWallet(privateKey)
	walletAddr, err := wallet.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("derive wallet address: %w", err)
	}

	chain := harachain.NewBlockchain(network, new(big.Int).SetUint64(chainID))

	return &Provider{
		Network:    network,
		Chain:      chain,
		Wallet:     wallet,
		WalletAddr: walletAddr,
	}, nil
}
