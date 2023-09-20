package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

func (m msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	supervisor, err := sdk.AccAddressFromBech32(msg.Supervisor)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid supervisor address: %s", err)
	}

	if !m.IsSupervisor(ctx, msg.Supervisor) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "burn is only allowed to supervisors")
	}

	if err := m.Keeper.Burn(ctx, supervisor, msg.Amount); err != nil {
		return nil, err
	}

	return &types.MsgBurnResponse{}, nil
}
