package v110

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
)

var (
	KeySPWithdrawalTotalPeriods          = []byte("SPWithdrawalTotalPeriods")
	KeySPWithdrawalPeriodDurationSeconds = []byte("SPWithdrawalPeriodDurationSeconds")
)

//nolint:lll
type Params struct {
	Supervisors                       []string                           `protobuf:"bytes,1,rep,name=supervisors,proto3" json:"supervisors,omitempty"`
	GasLimit                          sdk.IntProto                       `protobuf:"bytes,2,opt,name=gas_limit,json=gasLimit,proto3" json:"gas_limit"`
	GasUnconditionedAmount            sdk.IntProto                       `protobuf:"bytes,3,opt,name=gas_unconditioned_amount,json=gasUnconditionedAmount,proto3" json:"gas_unconditioned_amount"`
	GasAdjustCoefficient              sdk.DecProto                       `protobuf:"bytes,4,opt,name=gas_adjust_coefficient,json=gasAdjustCoefficient,proto3" json:"gas_adjust_coefficient"`
	SpWithdrawalTotalPeriods          uint32                             `protobuf:"varint,5,opt,name=sp_withdrawal_total_periods,json=spWithdrawalTotalPeriods,proto3" json:"sp_withdrawal_total_periods,omitempty"`
	SpWithdrawalPeriodDurationSeconds uint32                             `protobuf:"varint,6,opt,name=sp_withdrawal_period_duration_seconds,json=spWithdrawalPeriodDurationSeconds,proto3" json:"sp_withdrawal_period_duration_seconds,omitempty"`
	ValidatorsReward                  scorumtypes.ValidatorsRewardParams `protobuf:"bytes,7,opt,name=validators_reward,json=validatorsReward,proto3" json:"validators_reward"`
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	noValidate := func(value interface{}) error { return nil }
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(scorumtypes.KeySupervisors, &p.Supervisors, noValidate),
		paramtypes.NewParamSetPair(scorumtypes.KeyGasLimit, &p.GasLimit, noValidate),
		paramtypes.NewParamSetPair(scorumtypes.KeyGasUnconditionedAmount, &p.GasUnconditionedAmount, noValidate),
		paramtypes.NewParamSetPair(scorumtypes.KeyGasAdjustCoefficient, &p.GasAdjustCoefficient, noValidate),
		paramtypes.NewParamSetPair(KeySPWithdrawalTotalPeriods, &p.SpWithdrawalTotalPeriods, noValidate),
		paramtypes.NewParamSetPair(KeySPWithdrawalPeriodDurationSeconds, &p.SpWithdrawalPeriodDurationSeconds, noValidate),
		paramtypes.NewParamSetPair(scorumtypes.KeyValidatorsRewardPoolAddress, &p.ValidatorsReward.PoolAddress, noValidate),
		paramtypes.NewParamSetPair(scorumtypes.KeyValidatorsRewardPoolBlockReward, &p.ValidatorsReward.BlockReward, noValidate),
	}
}

func migrateScorumParams(
	ctx sdk.Context,
	cdc *codec.LegacyAmino,
	ps paramtypes.Subspace,
) error {
	var old Params

	for _, pair := range old.ParamSetPairs() {
		if err := cdc.UnmarshalJSON(ps.GetRaw(ctx, pair.Key), pair.Value); err != nil {
			return fmt.Errorf("failed to get old scorum params: %w", err)
		}
	}

	ps.SetParamSet(ctx, &scorumtypes.Params{
		Supervisors:            old.Supervisors,
		GasLimit:               old.GasLimit,
		GasUnconditionedAmount: old.GasUnconditionedAmount,
		GasAdjustCoefficient:   old.GasAdjustCoefficient,
		ValidatorsReward: scorumtypes.ValidatorsRewardParams{
			PoolAddress: old.ValidatorsReward.PoolAddress,
			BlockReward: old.ValidatorsReward.BlockReward,
		},
	})

	return nil
}
