package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
)

const TypeMsgStopSPWithdrawal = "stop_sp_withdrawal"

var _ sdk.Msg = &MsgStopSPWithdrawal{}

func NewMsgStopSPWithdrawal(owner, id string) *MsgStopSPWithdrawal {
	return &MsgStopSPWithdrawal{
		Owner: owner,
		Id:    id,
	}
}

func (msg *MsgStopSPWithdrawal) Route() string {
	return RouterKey
}

func (msg *MsgStopSPWithdrawal) Type() string {
	return TypeMsgStopSPWithdrawal
}

func (msg *MsgStopSPWithdrawal) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgStopSPWithdrawal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgStopSPWithdrawal) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if _, err := uuid.Parse(msg.Id); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid id (%s): must be uuid", err)
	}

	return nil
}
