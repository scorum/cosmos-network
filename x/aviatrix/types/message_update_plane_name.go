package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
)

const TypeMsgUpdatePlaneName = "msg_update_plane_name"

var _ sdk.Msg = &MsgUpdatePlaneName{}

func NewMsgUpdatePlaneName(supervisor string, id string, name string) *MsgUpdatePlaneName {
	return &MsgUpdatePlaneName{
		Supervisor: supervisor,
		Id:         id,
		Name:       name,
	}
}

func (msg *MsgUpdatePlaneName) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePlaneName) Type() string {
	return TypeMsgUpdatePlaneName
}

func (msg *MsgUpdatePlaneName) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Supervisor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePlaneName) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePlaneName) ValidateBasic() error {
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

	if msg.Name == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "empty name")
	}

	if len(msg.Name) > 32 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "name is too long")
	}

	return nil
}
