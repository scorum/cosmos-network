package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgWithdrawSP = "withdraw_sp"

var _ sdk.Msg = &MsgWithdrawSP{}

func NewMsgWithdrawSP(owner, recipient string, amount sdkmath.Int) *MsgWithdrawSP {
	return &MsgWithdrawSP{
		Owner:     owner,
		Recipient: recipient,
		Amount:    sdk.IntProto{Int: amount},
	}
}

func (msg *MsgWithdrawSP) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawSP) Type() string {
	return TypeMsgWithdrawSP
}

func (msg *MsgWithdrawSP) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawSP) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawSP) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}

	if !msg.Amount.Int.IsPositive() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid amount (must be positive)")
	}

	return nil
}
