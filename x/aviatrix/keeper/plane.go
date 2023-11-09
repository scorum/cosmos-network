package keeper

import (
	errorsmod "cosmossdk.io/errors"
	codec "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/gogo/protobuf/proto"

	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func (k Keeper) GetPlane(ctx sdk.Context, id string) (*types.Plane, error) {
	nft, ok := k.nftKeeper.GetNFT(ctx, types.NftClassID, id)
	if !ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "plane doesn't exist")
	}

	var meta types.PlaneMeta
	if err := proto.Unmarshal(nft.Data.Value, &meta); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrPanic, "failed to unmarshal stored meta: %s", err)
	}

	return &types.Plane{
		Id:    nft.Id,
		Owner: k.nftKeeper.GetOwner(ctx, types.NftClassID, id).String(),
		Meta:  &meta,
	}, nil
}

func (k Keeper) UpdatePlane(ctx sdk.Context, id string, meta *types.PlaneMeta) error {
	data, err := codec.NewAnyWithValue(meta)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "failed to marshal meta: %s", err)
	}

	if err := k.nftKeeper.Update(ctx, nft.NFT{
		ClassId: types.NftClassID,
		Id:      id,
		Data:    data,
	}); err != nil {
		return errorsmod.Wrap(err, "failed to update plane")
	}

	return nil
}
