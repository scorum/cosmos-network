package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgBurn = "burn"

var _ sdk.Msg = &MsgBurn{}

func NewMsgBurn(supervisor string, amount sdk.Coin) *MsgBurn {
	return &MsgBurn{
		Supervisor: supervisor,
		Amount:     amount,
	}
}

func (msg *MsgBurn) Route() string {
	return RouterKey
}

func (msg *MsgBurn) Type() string {
	return TypeMsgBurn
}

func (msg *MsgBurn) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Supervisor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBurn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBurn) ValidateBasic() error {
	if err := msg.Amount.Validate(); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid amount (%s)", err)
	}

	if msg.Amount.IsZero() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid amount (zero)")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Supervisor); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid supervisor address (%s)", err)
	}

	return nil
}
