package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
)

const TypeMsgCreatePlane = "create_plane"

var _ sdk.Msg = &MsgCreatePlane{}

func NewMsgCreatePlane(supervisor, id, owner string, experience uint64) *MsgCreatePlane {
	return &MsgCreatePlane{
		Id:         id,
		Supervisor: supervisor,
		Owner:      owner,
		Meta: &PlaneMeta{
			Experience: experience,
		},
	}
}

func (msg *MsgCreatePlane) Route() string {
	return RouterKey
}

func (msg *MsgCreatePlane) Type() string {
	return TypeMsgCreatePlane
}

func (msg *MsgCreatePlane) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Supervisor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePlane) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePlane) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Supervisor); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid supervisor address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if _, err := uuid.Parse(msg.Id); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "id must be uuid")
	}

	if err := msg.Meta.ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "invalid meta")
	}

	return nil
}

func (msg *PlaneMeta) ValidateBasic() error {
	return nil
}
