package sdk

import (
	"context"
	"fmt"

	"github.com/HARA-DID/did-verifiable-credentials-sdk/pkg/nftbase"
	vcfact "github.com/HARA-DID/did-verifiable-credentials-sdk/pkg/vcfactory"
	vcstor "github.com/HARA-DID/did-verifiable-credentials-sdk/pkg/vcstorage"
	harautils "github.com/meQlause/hara-core-blockchain-lib/utils"

	"github.com/myorg/worker-service/internal/config"
	"github.com/myorg/worker-service/internal/domain"
)

// VCAdapter implements Verifiable Credentials related blockchain operations.
type VCAdapter struct {
	provider       *Provider
	factory        *vcfact.VCFactory
	storage        *vcstor.VCStorage
	IdentityNFT    *nftbase.NFTBase
	CertificateNFT *nftbase.NFTBase
}

func NewVCAdapter(p *Provider, cfg config.BlockchainConfig) (*VCAdapter, error) {
	initCtx := context.Background()

	factory, err := vcfact.NewVCFactoryWithHNS(initCtx, cfg.VCFactoryHNS, p.Chain)
	if err != nil {
		return nil, fmt.Errorf("resolve VCFactory via HNS %q: %w", cfg.VCFactoryHNS, err)
	}

	storage, err := vcstor.NewVCStorageWithHNS(initCtx, cfg.VCStorageHNS, p.Chain)
	if err != nil {
		return nil, fmt.Errorf("resolve VCStorage via HNS %q: %w", cfg.VCStorageHNS, err)
	}

	idNFT, err := nftbase.NewNFTBaseWithHNS(initCtx, cfg.VCIdentityNFTHNS, p.Chain)
	if err != nil {
		return nil, fmt.Errorf("resolve IdentityNFT via HNS %q: %w", cfg.VCIdentityNFTHNS, err)
	}

	certNFT, err := nftbase.NewNFTBaseWithHNS(initCtx, cfg.VCCertificateNFTHNS, p.Chain)
	if err != nil {
		return nil, fmt.Errorf("resolve CertificateNFT via HNS %q: %w", cfg.VCCertificateNFTHNS, err)
	}

	return &VCAdapter{
		provider:       p,
		factory:        factory,
		storage:        storage,
		IdentityNFT:    idNFT,
		CertificateNFT: certNFT,
	}, nil
}

// ── Address Getters ──────────────────────────────────────────────

func (a *VCAdapter) GetFactoryAddress() string {
	return a.factory.Address.Hex()
}

func (a *VCAdapter) GetStorageAddress() string {
	return a.storage.Address.Hex()
}

// ── Encode Methods ───────────────────────────────────────────────

func (a *VCAdapter) EncodeIssueCredential(p domain.IssueCredentialPayload) ([]byte, error) {
	params := vcfact.IssueCredentialParams{
		Option:         vcfact.Options(p.Option),
		DIDRecipient:   p.DIDRecipient,
		Issuer:         p.Issuer,
		ExpiredAt:      p.ExpiredAt,
		OffchainHash:   p.OffchainHash,
		MerkleTreeRoot: p.MerkleTreeRoot,
		PublicIdentity: p.PublicIdentity,
	}
	return a.encode("IssueCredential", params)
}

func (a *VCAdapter) EncodeBurnCredential(p domain.BurnCredentialPayload) ([]byte, error) {
	params := vcfact.BurnCredentialParams{
		Option:  vcfact.Options(p.Option),
		DID:     p.DID,
		TokenID: p.TokenID,
	}
	return a.encode("BurnCredential", params)
}

func (a *VCAdapter) EncodeUpdateMetadata(p domain.UpdateMetadataPayload) ([]byte, error) {
	params := vcfact.UpdateMetadataParams{
		Option:       vcfact.Options(p.Option),
		TokenID:      p.TokenID,
		ExpiredAt:    p.ExpiredAt,
		OffchainHash: p.OffchainHash,
	}
	return a.encode("UpdateMetadata", params)
}

func (a *VCAdapter) EncodeRevokeCredential(p domain.RevokeCredentialPayload) ([]byte, error) {
	params := vcfact.RevokeCredentialParams{
		Option:  vcfact.Options(p.Option),
		TokenID: p.TokenID,
	}
	return a.encode("RevokeCredential", params)
}

func (a *VCAdapter) EncodeApproveCredentialOrg(p domain.ApproveCredentialOrgPayload) ([]byte, error) {
	params := vcfact.ApproveCredentialOrgParams{
		Option:      vcfact.Options(p.Option),
		TokenID:     p.TokenID,
		OrgDIDHash:  p.OrgDIDHash,
		UserDIDHash: p.UserDIDHash,
		Signature:   p.Signature,
	}
	return a.encode("ApproveCredentialOrg", params)
}

func (a *VCAdapter) EncodeApproveCredential(p domain.ApproveCredentialPayload) ([]byte, error) {
	params := vcfact.ApproveCredentialParams{
		Option:    vcfact.Options(p.Option),
		TokenID:   p.TokenID,
		Signature: p.Signature,
	}
	return a.encode("ApproveCredential", params)
}

func (a *VCAdapter) EncodeSetDidRootStorage(p domain.SetAddressPayload) ([]byte, error) {
	params := vcstor.SetAddressParams{
		Address: harautils.HexToAddress(p.Address),
	}
	return a.encodeStorage("SetDidRootStorage", params)
}

func (a *VCAdapter) EncodeSetDidOrgStorage(p domain.SetAddressPayload) ([]byte, error) {
	params := vcstor.SetAddressParams{
		Address: harautils.HexToAddress(p.Address),
	}
	return a.encodeStorage("SetDidOrgStorage", params)
}

// ── Private Helpers ──────────────────────────────────────────────

func (a *VCAdapter) encode(methodName string, params interface{ ToArgs() []any }) ([]byte, error) {
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

func (a *VCAdapter) encodeStorage(methodName string, params interface{ ToArgs() []any }) ([]byte, error) {
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
