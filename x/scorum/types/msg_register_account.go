package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRegisterAccount = "register_account"

var _ sdk.Msg = &MsgRegisterAccount{}

func NewMsgRegisterAccount(supervisor, address string) *MsgRegisterAccount {
	return &MsgRegisterAccount{
		Supervisor: supervisor,
		Address:    address,
	}
}

func (msg *MsgRegisterAccount) Route() string {
	return RouterKey
}

func (msg *MsgRegisterAccount) Type() string {
	return TypeMsgRegisterAccount
}

func (msg *MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Supervisor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRegisterAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterAccount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Supervisor); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid supervisor address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address for registration (%s)", err)
	}

	return nil
}
