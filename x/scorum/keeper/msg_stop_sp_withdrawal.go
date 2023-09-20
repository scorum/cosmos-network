package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/scorum/cosmos-network/x/scorum/types"
)

func (m msgServer) StopSPWithdrawal(goCtx context.Context, msg *types.MsgStopSPWithdrawal) (*types.MsgStopSPWithdrawalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	w, has := m.Keeper.GetSPWithdrawal(ctx, owner, msg.Id)
	if !has {
		return nil, errorsmod.Wrapf(sdkerrors.ErrNotFound, "withdrawal not found")
	}

	if !w.IsActive {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "withdrawal is not active")
	}

	w.IsActive = false
	m.Keeper.SetSPWithdrawal(ctx, w)

	return &types.MsgStopSPWithdrawalResponse{}, nil
}
