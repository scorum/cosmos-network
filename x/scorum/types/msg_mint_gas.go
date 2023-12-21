package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgMintGas = "mint_gas"

var _ sdk.Msg = &MsgMintGas{}

func NewMsgMintGas(supervisor string, address string, amount math.Int) *MsgMintGas {
	return &MsgMintGas{
		Supervisor: supervisor,
		Address:    address,
		Amount:     sdk.IntProto{Int: amount},
	}
}

func (msg *MsgMintGas) Route() string {
	return RouterKey
}

func (msg *MsgMintGas) Type() string {
	return TypeMsgMintGas
}

func (msg *MsgMintGas) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Supervisor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMintGas) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMintGas) ValidateBasic() error {
	if !msg.Amount.Int.IsPositive() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid amount (must be positive)")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Supervisor); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid supervisor address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address (%s)", err)
	}

	return nil
}
