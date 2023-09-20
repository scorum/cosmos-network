package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeMint = "Mint"
)

var (
	_ govv1beta1.Content = &MintProposal{}
)

func init() {
	govv1beta1.RegisterProposalType(ProposalTypeMint)
}

// NewMintProposal creates a new mint proposal.
func NewMintProposal(title, description string, recipient sdk.AccAddress, amount sdk.Coin) *MintProposal {
	return &MintProposal{
		Title:       title,
		Description: description,
		Recipient:   recipient.String(),
		Amount:      amount,
	}
}

func (cdp *MintProposal) ProposalRoute() string { return RouterKey }

func (cdp *MintProposal) ProposalType() string {
	return ProposalTypeMint
}

func (cdp *MintProposal) ValidateBasic() error {
	err := govv1beta1.ValidateAbstract(cdp)
	if err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(cdp.Recipient); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}

	return cdp.Amount.Validate()
}
