package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func (k Keeper) PlaneByName(goCtx context.Context, req *types.QueryPlaneByNameRequest) (*types.QueryPlaneByNameResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	id := string(k.planeNameIndex(ctx).Get([]byte(req.Name)))

	nft, ok := k.nftKeeper.GetNFT(ctx, types.NftClassID, id)
	if !ok {
		return nil, status.Error(codes.NotFound, "plane not found")
	}

	return &types.QueryPlaneByNameResponse{
		Nft: &nft,
	}, nil
}
