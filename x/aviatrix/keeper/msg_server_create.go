package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	codec "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"

	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func (k msgServer) CreatePlane(goCtx context.Context, msg *types.MsgCreatePlane) (*types.MsgCreatePlaneResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.scorumKeeper.IsSupervisor(ctx, msg.Supervisor) {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "method is allowed only for supervisors")
	}

	data, err := codec.NewAnyWithValue(msg.Meta)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "failed to marshal meta: %s", err)
	}

	if !k.nftKeeper.HasClass(ctx, types.NftClassID) {
		if err := k.nftKeeper.SaveClass(ctx, nft.Class{
			Id:     types.NftClassID,
			Name:   types.NftClassID,
			Symbol: "ap",
		}); err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrPanic, "")
		}
	}

	_, exists := k.nftKeeper.GetNFT(ctx, types.NftClassID, msg.Id)
	if exists {
		ctx.Logger().With("owner", msg.Owner, "id", msg.Id).Info("attempt to create plane with existing id")
		return &types.MsgCreatePlaneResponse{}, nil
	}

	if err := k.nftKeeper.Mint(ctx, nft.NFT{
		ClassId: types.NftClassID,
		Id:      msg.Id,
		Data:    data,
	}, sdk.MustAccAddressFromBech32(msg.Owner)); err != nil {
		return nil, fmt.Errorf("failed to mint nft: %w", err)
	}

	return &types.MsgCreatePlaneResponse{}, nil
}
