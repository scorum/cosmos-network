package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func (k msgServer) UpdatePlaneExperience(
	goCtx context.Context,
	msg *types.MsgUpdatePlaneExperience,
) (*types.MsgUpdatePlaneExperienceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	plane, err := k.GetPlane(ctx, msg.Id)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get plane")
	}

	if isSupervisor := k.scorumKeeper.IsSupervisor(ctx, msg.Supervisor); !isSupervisor {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "adjusting experience is allowed only for supervisors")
	}

	plane.Meta.Experience = msg.Amount

	if err := k.UpdatePlane(ctx, plane.Id, plane.Meta); err != nil {
		return nil, err
	}

	return &types.MsgUpdatePlaneExperienceResponse{}, nil
}
