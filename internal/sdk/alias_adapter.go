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

func NewAliasAdapter(p *Provider, cfg config.BlockchainConfig) (*AliasAdapter, error) {
	initCtx := context.Background()

	factory, err := aliasfact.NewAliasFactoryWithHNS(initCtx, cfg.AliasFactoryHNS, p.Chain)
	if err != nil {
		return nil, fmt.Errorf("resolve AliasFactory via HNS %q: %w", cfg.AliasFactoryHNS, err)
	}

	storage, err := aliasstor.NewAliasStorageWithHNS(initCtx, cfg.AliasStorageHNS, p.Chain)
	if err != nil {
		return nil, fmt.Errorf("resolve AliasStorage via HNS %q: %w", cfg.AliasStorageHNS, err)
	}

	return &AliasAdapter{
		provider: p,
		factory:  factory,
		storage:  storage,
	}, nil
}

// ── Encode Methods ───────────────────────────────────────────────

func (a *AliasAdapter) EncodeRegisterTLD(p domain.RegisterTLDPayload) ([]byte, error) {
	params := aliasfact.RegisterTLDParams{
		TLD:   p.TLD,
		Owner: p.Owner,
	}
	return a.encode("registerTLD", params)
}

func (a *AliasAdapter) EncodeRegisterDomain(p domain.RegisterDomainPayload) ([]byte, error) {
	params := aliasfact.RegisterDomainParams{
		Label:  p.Label,
		TLD:    p.TLD,
		Period: aliasfact.RegistrationPeriod(p.Period),
	}
	return a.encode("registerDomain", params)
}

func (a *AliasAdapter) EncodeSetDIDAlias(p domain.SetDIDAliasPayload) ([]byte, error) {
	params := aliasfact.SetDIDParams{
		Name: p.Name,
		DID:  decodeHash(p.DID),
	}
	return a.encode("setDID", params)
}

func (a *AliasAdapter) EncodeSetDIDOrgAlias(p domain.SetDIDOrgAliasPayload) ([]byte, error) {
	params := aliasfact.SetDIDOrgParams{
		Name:        p.Name,
		OrgDIDHash:  decodeHash(p.OrgDIDHash),
		UserDIDHash: decodeHash(p.UserDIDHash),
	}
	return a.encode("setDIDOrg", params)
}

func (a *AliasAdapter) EncodeExtendRegistration(p domain.ExtendRegistrationPayload) ([]byte, error) {
	params := aliasfact.ExtendRegistrationParams{
		Node:   decodeHash(p.Node),
		Period: aliasfact.RegistrationPeriod(p.Period),
	}
	return a.encode("extendRegistration", params)
}

func (a *AliasAdapter) EncodeRevokeAlias(p domain.RevokeAliasPayload) ([]byte, error) {
	params := aliasfact.NodeOnlyParams{
		Node: decodeHash(p.Node),
	}
	return a.encode("revokeAlias", params)
}

func (a *AliasAdapter) EncodeUnrevokeAlias(p domain.UnrevokeAliasPayload) ([]byte, error) {
	params := aliasfact.NodeOnlyParams{
		Node: decodeHash(p.Node),
	}
	return a.encode("unrevokeAlias", params)
}

func (a *AliasAdapter) EncodeRegisterSubdomain(p domain.RegisterSubdomainPayload) ([]byte, error) {
	params := aliasfact.RegisterSubdomainParams{
		Label:        p.Label,
		ParentDomain: p.ParentDomain,
		Period:       aliasfact.RegistrationPeriod(p.Period),
	}
	return a.encode("registerSubdomain", params)
}

func (a *AliasAdapter) EncodeTransferAliasOwnership(p domain.TransferAliasOwnershipPayload) ([]byte, error) {
	params := aliasfact.TransferAliasOwnershipParams{
		Node:     decodeHash(p.Node),
		NewOwner: p.NewOwner,
	}
	return a.encode("transferAliasOwnership", params)
}

func (a *AliasAdapter) EncodeTransferTLD(p domain.TransferTLDPayload) ([]byte, error) {
	params := aliasfact.TransferTLDParams{
		TLD:      p.TLD,
		NewOwner: p.NewOwner,
	}
	return a.encode("transferTLD", params)
}

func (a *AliasAdapter) EncodeSetAliasRootStorage(p domain.SetAliasAddressPayload) ([]byte, error) {
	params := aliasfact.SetDIDRootStorageParams{
		DIDRootStorage: p.Address,
	}
	return a.encode("setDIDRootStorage", params)
}

func (a *AliasAdapter) EncodeSetAliasOrgStorage(p domain.SetAliasAddressPayload) ([]byte, error) {
	params := aliasfact.SetDIDOrgStorageParams{
		DIDOrgStorage: p.Address,
	}
	return a.encode("setDIDOrgStorage", params)
}

func (a *AliasAdapter) EncodeSetFactoryContract(p domain.SetFactoryContractPayload) ([]byte, error) {
	params := aliasstor.SetFactoryContractParams{
		FactoryContract: harautils.HexToAddress(p.FactoryContract),
	}
	return a.encodeStorage("setFactoryContract", params)
}

// ── Private Helpers ──────────────────────────────────────────────

func (a *AliasAdapter) encode(methodName string, params interface{ ToArgs() []any }) ([]byte, error) {
	method, ok := a.factory.ContractABI.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s not found in ABI", methodName)
	}

	inputs, err := method.Inputs.Pack(params.ToArgs()...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack %s arguments: %w", methodName, err)
	}

	return append(method.ID, inputs...), nil
}

func (a *AliasAdapter) encodeStorage(methodName string, params interface{ ToArgs() []any }) ([]byte, error) {
	method, ok := a.storage.ContractABI.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s not found in storage ABI", methodName)
	}

	inputs, err := method.Inputs.Pack(params.ToArgs()...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack %s arguments: %w", methodName, err)
	}

	return append(method.ID, inputs...), nil
}

func decodeHash(hexStr string) [32]byte {
	return harautils.HexToHash(hexStr)
}
