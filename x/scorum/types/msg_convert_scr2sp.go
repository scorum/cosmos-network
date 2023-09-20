package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgConvertSCR2SP = "convert_scr2sp"

var _ sdk.Msg = &MsgConvertSCR2SP{}

func NewMsgConvertSCR2SP(owner string, amount sdkmath.Int) *MsgConvertSCR2SP {
	return &MsgConvertSCR2SP{
		Owner:  owner,
		Amount: sdk.IntProto{Int: amount},
	}
}

func (msg *MsgConvertSCR2SP) Route() string {
	return RouterKey
}

func (msg *MsgConvertSCR2SP) Type() string {
	return TypeMsgConvertSCR2SP
}

func (msg *MsgConvertSCR2SP) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgConvertSCR2SP) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgConvertSCR2SP) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if !msg.Amount.Int.IsPositive() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid amount (must be positive)")
	}

	return nil
}
