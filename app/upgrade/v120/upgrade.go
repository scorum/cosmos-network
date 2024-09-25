package v120

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"

	"github.com/cosmos/cosmos-sdk/types"

	upgrade "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const Name = "v1.2.0"

//nolint:lll
type Params struct {
	Supervisors            []string               `protobuf:"bytes,1,rep,name=supervisors,proto3" json:"supervisors,omitempty"`
	GasLimit               types.IntProto         `protobuf:"bytes,2,opt,name=gas_limit,json=gasLimit,proto3" json:"gas_limit"`
	GasUnconditionedAmount types.IntProto         `protobuf:"bytes,3,opt,name=gas_unconditioned_amount,json=gasUnconditionedAmount,proto3" json:"gas_unconditioned_amount"`
	GasAdjustCoefficient   types.DecProto         `protobuf:"bytes,4,opt,name=gas_adjust_coefficient,json=gasAdjustCoefficient,proto3" json:"gas_adjust_coefficient"`
	ValidatorsReward       ValidatorsRewardParams `protobuf:"bytes,5,opt,name=validators_reward,json=validatorsReward,proto3" json:"validators_reward"`
}

type ValidatorsRewardParams struct {
	PoolAddress string     `protobuf:"bytes,1,opt,name=pool_address,json=poolAddress,proto3" json:"pool_address,omitempty"`
	BlockReward types.Coin `protobuf:"bytes,2,opt,name=block_reward,json=blockReward,proto3" json:"block_reward"`
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	noValidate := func(value interface{}) error { return nil }
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(scorumtypes.KeySupervisors, &p.Supervisors, noValidate),
		paramtypes.NewParamSetPair(scorumtypes.KeyGasLimit, &p.GasLimit, noValidate),
		paramtypes.NewParamSetPair(scorumtypes.KeyGasUnconditionedAmount, &p.GasUnconditionedAmount, noValidate),
		paramtypes.NewParamSetPair(scorumtypes.KeyGasAdjustCoefficient, &p.GasAdjustCoefficient, noValidate),
		paramtypes.NewParamSetPair(scorumtypes.KeyValidatorsRewardPoolAddress, &p.ValidatorsReward.PoolAddress, noValidate),
		paramtypes.NewParamSetPair(scorumtypes.KeyValidatorsRewardPoolBlockReward, &p.ValidatorsReward.BlockReward, noValidate),
	}
}

func Handler(
	cfg module.Configurator,
	mm *module.Manager,
	cdc *codec.LegacyAmino,
	ps paramtypes.Subspace,
) func(ctx context.Context, _ upgrade.Plan, _ module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, _ upgrade.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := types.UnwrapSDKContext(ctx)
		var old Params

		for _, pair := range old.ParamSetPairs() {
			if err := cdc.UnmarshalJSON(ps.GetRaw(sdkCtx, pair.Key), pair.Value); err != nil {
				return nil, fmt.Errorf("failed to get old scorum params: %w", err)
			}
		}

		ps.GetParamSetIfExists(sdkCtx, &old)
		ps.SetParamSet(sdkCtx, &scorumtypes.Params{
			Supervisors:            old.Supervisors,
			GasLimit:               old.GasLimit.Int,
			GasUnconditionedAmount: old.GasUnconditionedAmount.Int,
			GasAdjustCoefficient:   old.GasAdjustCoefficient.Dec,
			ValidatorsReward: scorumtypes.ValidatorsRewardParams{
				PoolAddress: old.ValidatorsReward.PoolAddress,
				BlockReward: old.ValidatorsReward.BlockReward,
			},
		})

		return mm.RunMigrations(ctx, cfg, fromVM)
	}
}
