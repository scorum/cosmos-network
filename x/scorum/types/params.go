package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeySupervisors                     = []byte("Supervisors")
	KeyGasLimit                        = []byte("GasLimit")
	KeyGasAdjustCoefficient            = []byte("GasAdjustCoefficient")
	KeyGasUnconditionedAmount          = []byte("GasUnconditionedAmount")
	KeyValidatorsRewardPoolAddress     = []byte("ValidatorsRewardPoolAddress")
	KeyValidatorsRewardPoolBlockReward = []byte("ValidatorsRewardPoolBlockReward")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	supervisors []string,
	gasLimit math.Int,
	gasUnconditionedAmount math.Int,
	gasAdjustCoefficient math.LegacyDec,
) Params {
	return Params{
		Supervisors:            supervisors,
		GasLimit:               gasLimit,
		GasUnconditionedAmount: gasUnconditionedAmount,
		GasAdjustCoefficient:   gasAdjustCoefficient,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		Supervisors:            nil,
		GasLimit:               math.NewInt(1000000),
		GasUnconditionedAmount: math.NewInt(15000),
		GasAdjustCoefficient:   math.LegacyNewDec(1),
		ValidatorsReward: ValidatorsRewardParams{
			PoolAddress: "",
			BlockReward: sdk.Coin{
				Denom:  SCRDenom,
				Amount: math.ZeroInt(),
			},
		},
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeySupervisors, &p.Supervisors, validateSupervisors),
		paramtypes.NewParamSetPair(KeyGasLimit, &p.GasLimit, validateGasLimit),
		paramtypes.NewParamSetPair(KeyGasUnconditionedAmount, &p.GasUnconditionedAmount, validateGasUnconditionedAmount),
		paramtypes.NewParamSetPair(KeyGasAdjustCoefficient, &p.GasAdjustCoefficient, validateGasAdjustCoefficient),
		paramtypes.NewParamSetPair(KeyValidatorsRewardPoolAddress, &p.ValidatorsReward.PoolAddress, validateValidatorsRewardPoolAddress),
		paramtypes.NewParamSetPair(KeyValidatorsRewardPoolBlockReward, &p.ValidatorsReward.BlockReward, validateValidatorsRewardBlockReward),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateSupervisors(p.Supervisors); err != nil {
		return fmt.Errorf("invalid supervisors: %w", err)
	}

	if err := validateGasLimit(p.GasLimit); err != nil {
		return fmt.Errorf("invalid gasLimit: %w", err)
	}

	if err := validateGasAdjustCoefficient(p.GasAdjustCoefficient); err != nil {
		return fmt.Errorf("invalid gasAdjustCoefficient: %w", err)
	}

	if err := validateValidatorsRewardPoolAddress(p.ValidatorsReward.PoolAddress); err != nil {
		return fmt.Errorf("invalid validatorsReward.poolAddress: %w", err)
	}

	if err := validateValidatorsRewardBlockReward(p.ValidatorsReward.BlockReward); err != nil {
		return fmt.Errorf("invalid validatorsReward.blockReward: %w", err)
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateGasLimit(i interface{}) error {
	s, ok := i.(math.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if s.IsZero() || s.IsNegative() {
		return fmt.Errorf("must be positive")
	}

	return nil
}

func validateGasUnconditionedAmount(i interface{}) error {
	s, ok := i.(math.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if s.IsNegative() {
		return fmt.Errorf("must not be negative")
	}

	return nil
}

func validateGasAdjustCoefficient(i interface{}) error {
	s, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if s.IsZero() || s.IsNegative() {
		return fmt.Errorf("must be positive")
	}

	return nil
}

func validateSupervisors(i interface{}) error {
	s, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for i, v := range s {
		addr, err := sdk.AccAddressFromBech32(v)
		if err != nil {
			return fmt.Errorf("invalid address %d", i+1)
		}
		if err := sdk.VerifyAddressFormat(addr); err != nil {
			return fmt.Errorf("invalid address %d", i+1)
		}
	}

	return nil
}

func validateValidatorsRewardPoolAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return nil
	}

	if _, err := sdk.AccAddressFromBech32(v); err != nil {
		return err
	}

	return nil
}

func validateValidatorsRewardBlockReward(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsValid() {
		return nil
	}

	if v.IsNegative() {
		return fmt.Errorf("can not be negative")
	}

	return nil
}
