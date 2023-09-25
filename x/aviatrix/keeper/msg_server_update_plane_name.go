package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func (k msgServer) UpdatePlaneName(goCtx context.Context, msg *types.MsgUpdatePlaneName) (*types.MsgUpdatePlaneNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	plane, err := k.GetPlane(ctx, msg.Id)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get plane")
	}

	if !k.scorumKeeper.IsSupervisor(ctx, msg.Supervisor) {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "updating plane is allowed only for supervisors")
	}

	plane.Meta.Name = msg.Name
	if err := k.UpdatePlane(ctx, plane.Id, plane.Meta); err != nil {
		return nil, err
	}

	return &types.MsgUpdatePlaneNameResponse{}, nil
}