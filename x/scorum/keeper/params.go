package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"golang.org/x/exp/slices"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSetIfExists(ctx, &params)

	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

func (k Keeper) IsSupervisor(ctx sdk.Context, addr string) bool {
	return slices.Contains(k.GetParams(ctx).Supervisors, addr)
}
