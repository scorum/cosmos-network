package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
)

const TypeMsgUpdatePlaneExperience = "msg_update_plane_experience"

var _ sdk.Msg = &MsgUpdatePlaneExperience{}

func NewMsgUpdatePlaneExperience(supervisorAddr string, id string, amount uint64) *MsgUpdatePlaneExperience {
	return &MsgUpdatePlaneExperience{
		Supervisor: supervisorAddr,
		Id:         id,
		Amount:     amount,
	}
}

func (msg *MsgUpdatePlaneExperience) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePlaneExperience) Type() string {
	return TypeMsgUpdatePlaneExperience
}

func (msg *MsgUpdatePlaneExperience) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Supervisor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePlaneExperience) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePlaneExperience) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Supervisor)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid supervisor address (%s)", err)
	}

	if msg.Id == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "empty id")
	}

	if _, err := uuid.Parse(msg.Id); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "id must be uuid")
	}

	return nil
}
