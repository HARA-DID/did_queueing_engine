package sdk

import (
	"context"

	"github.com/myorg/worker-service/internal/domain"
	"github.com/myorg/worker-service/internal/service"
)

var _ service.BlockchainService = (*CompositeAdapter)(nil)

type CompositeAdapter struct {
	did   *DIDAdapter
	aa    *AAAdapter
	vc    *VCAdapter
	alias *AliasAdapter
}

func NewCompositeAdapter(did *DIDAdapter, aa *AAAdapter, vc *VCAdapter, alias *AliasAdapter) *CompositeAdapter {
	return &CompositeAdapter{
		did:   did,
		aa:    aa,
		vc:    vc,
		alias: alias,
	}
}

func (c *CompositeAdapter) CreateDID(ctx context.Context, p domain.CreateDIDPayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeCreateDID(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) AddKey(ctx context.Context, p domain.AddKeyPayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeAddKey(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) AddClaim(ctx context.Context, p domain.AddClaimPayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeAddClaim(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) StoreData(ctx context.Context, p domain.StoreDataPayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeStoreData(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) UpdateDID(ctx context.Context, p domain.UpdateDIDPayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeUpdateDID(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) DeactivateDID(ctx context.Context, p domain.DIDLifecyclePayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeDeactivateDID(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) ReactivateDID(ctx context.Context, p domain.DIDLifecyclePayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeReactivateDID(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) TransferDIDOwner(ctx context.Context, p domain.TransferDIDOwnerPayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeTransferDIDOwner(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) DeleteData(ctx context.Context, p domain.DeleteDataPayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeDeleteData(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) RemoveKey(ctx context.Context, p domain.RemoveKeyPayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeRemoveKey(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) RemoveClaim(ctx context.Context, p domain.RemoveClaimPayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeRemoveClaim(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) GeneralExecute(ctx context.Context, p domain.GeneralExecutePayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeGeneralExecute(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) CreateOrg(ctx context.Context, p domain.CreateOrgPayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeCreateOrg(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) DeactivateOrg(ctx context.Context, p domain.OrgLifecyclePayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeDeactivateOrg(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) ReactivateOrg(ctx context.Context, p domain.OrgLifecyclePayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeReactivateOrg(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) TransferOrgOwner(ctx context.Context, p domain.OrgTransferPayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeTransferOrgOwner(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) AddMember(ctx context.Context, p domain.OrgMemberPayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeAddMember(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) RemoveMember(ctx context.Context, p domain.OrgMemberPayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeRemoveMember(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) UpdateMember(ctx context.Context, p domain.OrgMemberPayload) (*domain.BlockchainResult, error) {
	data, err := c.did.EncodeUpdateMember(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) HandleOps(ctx context.Context, p domain.HandleOpsPayload) (*domain.BlockchainResult, error) {
	return c.aa.HandleOps(ctx, p)
}

func (c *CompositeAdapter) DeployWallet(ctx context.Context, p domain.DeployWalletPayload) (*domain.BlockchainResult, error) {
	return c.aa.DeployWallet(ctx, p)
}

func (c *CompositeAdapter) IssueCredential(ctx context.Context, p domain.IssueCredentialPayload) (*domain.BlockchainResult, error) {
	data, err := c.vc.EncodeIssueCredential(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) BurnCredential(ctx context.Context, p domain.BurnCredentialPayload) (*domain.BlockchainResult, error) {
	data, err := c.vc.EncodeBurnCredential(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) UpdateMetadata(ctx context.Context, p domain.UpdateMetadataPayload) (*domain.BlockchainResult, error) {
	data, err := c.vc.EncodeUpdateMetadata(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) RevokeCredential(ctx context.Context, p domain.RevokeCredentialPayload) (*domain.BlockchainResult, error) {
	data, err := c.vc.EncodeRevokeCredential(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) ApproveCredentialOrg(ctx context.Context, p domain.ApproveCredentialOrgPayload) (*domain.BlockchainResult, error) {
	data, err := c.vc.EncodeApproveCredentialOrg(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) ApproveCredential(ctx context.Context, p domain.ApproveCredentialPayload) (*domain.BlockchainResult, error) {
	data, err := c.vc.EncodeApproveCredential(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) SetDidRootStorage(ctx context.Context, p domain.SetAddressPayload) (*domain.BlockchainResult, error) {
	data, err := c.vc.EncodeSetDidRootStorage(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) SetDidOrgStorage(ctx context.Context, p domain.SetAddressPayload) (*domain.BlockchainResult, error) {
	data, err := c.vc.EncodeSetDidOrgStorage(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) RegisterTLD(ctx context.Context, p domain.RegisterTLDPayload) (*domain.BlockchainResult, error) {
	data, err := c.alias.EncodeRegisterTLD(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) RegisterDomain(ctx context.Context, p domain.RegisterDomainPayload) (*domain.BlockchainResult, error) {
	data, err := c.alias.EncodeRegisterDomain(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) SetDIDAlias(ctx context.Context, p domain.SetDIDAliasPayload) (*domain.BlockchainResult, error) {
	data, err := c.alias.EncodeSetDIDAlias(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) SetDIDOrgAlias(ctx context.Context, p domain.SetDIDOrgAliasPayload) (*domain.BlockchainResult, error) {
	data, err := c.alias.EncodeSetDIDOrgAlias(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) ExtendRegistration(ctx context.Context, p domain.ExtendRegistrationPayload) (*domain.BlockchainResult, error) {
	data, err := c.alias.EncodeExtendRegistration(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) RevokeAlias(ctx context.Context, p domain.RevokeAliasPayload) (*domain.BlockchainResult, error) {
	data, err := c.alias.EncodeRevokeAlias(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) UnrevokeAlias(ctx context.Context, p domain.UnrevokeAliasPayload) (*domain.BlockchainResult, error) {
	data, err := c.alias.EncodeUnrevokeAlias(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) RegisterSubdomain(ctx context.Context, p domain.RegisterSubdomainPayload) (*domain.BlockchainResult, error) {
	data, err := c.alias.EncodeRegisterSubdomain(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) TransferAliasOwnership(ctx context.Context, p domain.TransferAliasOwnershipPayload) (*domain.BlockchainResult, error) {
	data, err := c.alias.EncodeTransferAliasOwnership(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) TransferTLD(ctx context.Context, p domain.TransferTLDPayload) (*domain.BlockchainResult, error) {
	data, err := c.alias.EncodeTransferTLD(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) SetAliasRootStorage(ctx context.Context, p domain.SetAliasAddressPayload) (*domain.BlockchainResult, error) {
	data, err := c.alias.EncodeSetAliasRootStorage(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) SetAliasOrgStorage(ctx context.Context, p domain.SetAliasAddressPayload) (*domain.BlockchainResult, error) {
	data, err := c.alias.EncodeSetAliasOrgStorage(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}

func (c *CompositeAdapter) SetFactoryContract(ctx context.Context, p domain.SetFactoryContractPayload) (*domain.BlockchainResult, error) {
	data, err := c.alias.EncodeSetFactoryContract(p)
	if err != nil {
		return nil, err
	}
	return c.aa.HandleOps(ctx, domain.HandleOpsPayload{
		Target:           p.TargetAddress,
		Data:             data,
		MultipleRPCCalls: p.MultipleRPCCalls,
	})
}
