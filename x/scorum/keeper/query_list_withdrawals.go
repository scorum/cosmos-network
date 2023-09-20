package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv queryServer) ListWithdrawals(goCtx context.Context, req *types.QueryWithdrawalsRequest) (*types.QueryWithdrawalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid address")
	}

	return &types.QueryWithdrawalsResponse{
		Withdrawals: srv.Keeper.ListWithdrawals(ctx, owner),
	}, nil
}
