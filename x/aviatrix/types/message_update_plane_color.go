package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
)

const TypeMsgUpdatePlaneColor = "msg_update_plane_color"

var _ sdk.Msg = &MsgUpdatePlaneColor{}

func NewMsgUpdatePlaneColor(supervisor string, id string, color string) *MsgUpdatePlaneColor {
	return &MsgUpdatePlaneColor{
		Supervisor: supervisor,
		Id:         id,
		Color:      color,
	}
}

func (msg *MsgUpdatePlaneColor) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePlaneColor) Type() string {
	return TypeMsgUpdatePlaneColor
}

func (msg *MsgUpdatePlaneColor) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Supervisor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePlaneColor) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePlaneColor) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Supervisor)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if msg.Id == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "empty id")
	}

	if _, err := uuid.Parse(msg.Id); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "id must be uuid")
	}

	if msg.Color == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "empty color")
	}

	return nil
}
