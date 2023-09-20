package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
)

const TypeMsgAdjustPlaneExperience = "msg_adjust_plane_experience"

var _ sdk.Msg = &MsgAdjustPlaneExperience{}

func NewMsgAdjustPlaneExperience(supervisorAddr string, id string, amount int64) *MsgAdjustPlaneExperience {
	return &MsgAdjustPlaneExperience{
		Supervisor: supervisorAddr,
		Id:         id,
		Amount:     amount,
	}
}

func (msg *MsgAdjustPlaneExperience) Route() string {
	return RouterKey
}

func (msg *MsgAdjustPlaneExperience) Type() string {
	return TypeMsgAdjustPlaneExperience
}

func (msg *MsgAdjustPlaneExperience) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Supervisor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAdjustPlaneExperience) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAdjustPlaneExperience) ValidateBasic() error {
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

	if msg.Amount == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "zero amount")
	}

	return nil
}
