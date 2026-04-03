package sdk

import (
	"context"
	"fmt"
	"math/big"

	didfactory "github.com/HARA-DID/did-root-sdk/pkg/factory"
	harautils "github.com/HARA-DID/hara-core-blockchain-lib/utils"

	"github.com/HARA-DID/did-queueing-engine/internal/config"
	"github.com/HARA-DID/did-queueing-engine/internal/domain"
)

type DIDAdapter struct {
	provider *Provider
	factory  *didfactory.Factory
}

func NewDIDAdapter(p *Provider, cfg config.BlockchainConfig) (*DIDAdapter, error) {
	initCtx := context.Background()

	factory, err := didfactory.NewFactoryWithHNS(initCtx, cfg.DIDRootFactoryHNS, p.Chain)
	if err != nil {
		return nil, fmt.Errorf("resolve DIDFactory via HNS %q: %w", cfg.DIDRootFactoryHNS, err)
	}

	return &DIDAdapter{
		provider: p,
		factory:  factory,
	}, nil
}

func (a *DIDAdapter) CreateDID(ctx context.Context, p domain.CreateDIDPayload) (*domain.BlockchainResult, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	txHashes, err := a.factory.CreateDID(ctx, a.provider.Wallet, didfactory.CreateDIDParam{DID: p.DID}, keyID, p.MultipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("factory.CreateDID: %w", err)
	}
	return &domain.BlockchainResult{TxHashes: txHashes}, nil
}

func (a *DIDAdapter) EncodeAddKey(p domain.AddKeyPayload) ([]byte, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	keyHashed := harautils.HexToHash(p.PublicKey)

	argBuilder := a.provider.Network.ArgBuilder().
		Type("uint256").Value(p.DIDIndex).
		Type("bytes32").Value(keyHashed).
		Type("string").Value("").
		Type("uint8").Value(p.Purpose).
		Type("uint8").Value(p.KeyType)
	data := harautils.EncodeArgs(argBuilder)

	return a.encodeDID(didfactory.TypeAddKey, data, keyID)
}

func (a *DIDAdapter) EncodeAddClaim(p domain.AddClaimPayload) ([]byte, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	argBuilder := a.provider.Network.ArgBuilder().
		Type("uint256").Value(p.DIDIndex).
		Type("uint8").Value(p.Topic).
		Type("bytes").Value(p.Data).
		Type("string").Value(p.URI).
		Type("bytes").Value(p.Signature)
	data := harautils.EncodeArgs(argBuilder)

	return a.encodeDID(didfactory.TypeAddClaim, data, keyID)
}

func (a *DIDAdapter) EncodeStoreData(p domain.StoreDataPayload) ([]byte, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	argBuilder := a.provider.Network.ArgBuilder().
		Type("uint256").Value(p.DIDIndex).
		Type("string").Value(p.PropertyKey).
		Type("string").Value(p.Data)
	data := harautils.EncodeArgs(argBuilder)

	return a.encodeDID(didfactory.TypeStoreData, data, keyID)
}

func (a *DIDAdapter) EncodeUpdateDID(p domain.UpdateDIDPayload) ([]byte, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	argBuilder := a.provider.Network.ArgBuilder().
		Type("uint256").Value(p.DIDIndex).
		Type("string").Value(p.URI)
	data := harautils.EncodeArgs(argBuilder)

	return a.encodeDID(didfactory.TypeUpdateDID, data, keyID)
}

func (a *DIDAdapter) EncodeDeactivateDID(p domain.DIDLifecyclePayload) ([]byte, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	argBuilder := a.provider.Network.ArgBuilder().
		Type("uint256").Value(p.DIDIndex)
	data := harautils.EncodeArgs(argBuilder)

	return a.encodeDID(didfactory.TypeDeactivateDID, data, keyID)
}

func (a *DIDAdapter) EncodeReactivateDID(p domain.DIDLifecyclePayload) ([]byte, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	argBuilder := a.provider.Network.ArgBuilder().
		Type("uint256").Value(p.DIDIndex)
	data := harautils.EncodeArgs(argBuilder)

	return a.encodeDID(didfactory.TypeReactivateDID, data, keyID)
}

func (a *DIDAdapter) EncodeTransferDIDOwner(p domain.TransferDIDOwnerPayload) ([]byte, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	argBuilder := a.provider.Network.ArgBuilder().
		Type("uint256").Value(p.DIDIndex).
		Type("address").Value(harautils.HexToAddress(p.NewOwner))
	data := harautils.EncodeArgs(argBuilder)

	return a.encodeDID(didfactory.TypeTransferDID, data, keyID)
}

func (a *DIDAdapter) EncodeDeleteData(p domain.DeleteDataPayload) ([]byte, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	argBuilder := a.provider.Network.ArgBuilder().
		Type("uint256").Value(p.DIDIndex).
		Type("string").Value(p.Key)
	data := harautils.EncodeArgs(argBuilder)

	return a.encodeDID(didfactory.TypeDeleteData, data, keyID)
}

func (a *DIDAdapter) EncodeRemoveKey(p domain.RemoveKeyPayload) ([]byte, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	argBuilder := a.provider.Network.ArgBuilder().
		Type("uint256").Value(p.DIDIndex).
		Type("bytes32").Value(harautils.HexToHash(p.KeyDataHashed))
	data := harautils.EncodeArgs(argBuilder)

	return a.encodeDID(didfactory.TypeRemoveKey, data, keyID)
}

func (a *DIDAdapter) EncodeRemoveClaim(p domain.RemoveClaimPayload) ([]byte, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}

	argBuilder := a.provider.Network.ArgBuilder().
		Type("uint256").Value(p.DIDIndex).
		Type("bytes32").Value(harautils.HexToHash(p.ClaimID))
	data := harautils.EncodeArgs(argBuilder)

	return a.encodeDID(didfactory.TypeRemoveClaim, data, keyID)
}

func (a *DIDAdapter) EncodeGeneralExecute(p domain.GeneralExecutePayload) ([]byte, error) {
	keyID, err := resolveKeyIdentifier(p.KeyIdentifier)
	if err != nil {
		return nil, err
	}
	return a.encodeDID(didfactory.TypeGeneralExecute, p.Data, keyID)
}

func (a *DIDAdapter) EncodeCreateOrg(p domain.CreateOrgPayload) ([]byte, error) {
	return a.encodeOrg(didfactory.TypeCreateOrgDID, p.Data, big.NewInt(0))
}

func (a *DIDAdapter) EncodeDeactivateOrg(p domain.OrgLifecyclePayload) ([]byte, error) {
	return a.encodeOrg(didfactory.TypeDeactivateOrgDID, nil, p.OrgDIDIndex)
}

func (a *DIDAdapter) EncodeReactivateOrg(p domain.OrgLifecyclePayload) ([]byte, error) {
	return a.encodeOrg(didfactory.TypeReactivateOrgDID, nil, p.OrgDIDIndex)
}

func (a *DIDAdapter) EncodeTransferOrgOwner(p domain.OrgTransferPayload) ([]byte, error) {
	return a.encodeOrg(didfactory.TypeTransferOrgDID, p.Data, p.OrgDIDIndex)
}

func (a *DIDAdapter) EncodeAddMember(p domain.OrgMemberPayload) ([]byte, error) {
	return a.encodeOrg(didfactory.TypeAddMember, p.Data, p.OrgDIDIndex)
}

func (a *DIDAdapter) EncodeRemoveMember(p domain.OrgMemberPayload) ([]byte, error) {
	return a.encodeOrg(didfactory.TypeRemoveMember, p.Data, p.OrgDIDIndex)
}

func (a *DIDAdapter) EncodeUpdateMember(p domain.OrgMemberPayload) ([]byte, error) {
	return a.encodeOrg(didfactory.TypeUpdateMember, p.Data, p.OrgDIDIndex)
}

func (a *DIDAdapter) encodeDID(txType uint8, data []byte, keyIdentifier string) ([]byte, error) {
	method, ok := a.factory.ContractABI.Methods["callExternalDID"]
	if !ok {
		return nil, fmt.Errorf("method callExternalDID not found in ABI")
	}

	inputs, err := method.Inputs.Pack(txType, data, keyIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to pack callExternalDID arguments: %w", err)
	}

	return append(method.ID, inputs...), nil
}

func (a *DIDAdapter) encodeOrg(txType uint8, data []byte, orgDIDIndex *big.Int) ([]byte, error) {
	method, ok := a.factory.ContractABI.Methods["callExternalOrg"]
	if !ok {
		return nil, fmt.Errorf("method callExternalOrg not found in ABI")
	}

	inputs, err := method.Inputs.Pack(txType, data, orgDIDIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to pack callExternalOrg arguments: %w", err)
	}

	return append(method.ID, inputs...), nil
}

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
