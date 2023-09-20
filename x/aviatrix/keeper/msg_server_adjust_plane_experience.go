package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	codec "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"

	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func (k msgServer) AdjustPlaneExperience(
	goCtx context.Context,
	msg *types.MsgAdjustPlaneExperience,
) (*types.MsgAdjustPlaneExperienceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	plane, err := k.GetPlane(ctx, msg.Id)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get plane")
	}

	if isSupervisor := k.scorumKeeper.IsSupervisor(ctx, msg.Supervisor); !isSupervisor {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "adjusting experience is allowed only for supervisors")
	}

	experience := int64(plane.Meta.Experience) + msg.Amount
	if experience < 0 {
		experience = 0
	}
	plane.Meta.Experience = uint64(experience)

	data, err := codec.NewAnyWithValue(plane.Meta)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "failed to marshal meta: %s", err)
	}

	if err := k.nftKeeper.Update(ctx, nft.NFT{
		ClassId: types.NftClassID,
		Id:      msg.Id,
		Data:    data,
	}); err != nil {
		return nil, errorsmod.Wrap(err, "failed to update nft")
	}

	return &types.MsgAdjustPlaneExperienceResponse{}, nil
}
