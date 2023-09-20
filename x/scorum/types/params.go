package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeySupervisors                       = []byte("Supervisors")
	KeyGasLimit                          = []byte("GasLimit")
	KeyGasAdjustCoefficient              = []byte("GasAdjustCoefficient")
	KeyGasUnconditionedAmount            = []byte("GasUnconditionedAmount")
	KeyRegistrationSPDelegationAmount    = []byte("RegistrationSPDelegationAmount")
	KeySPWithdrawalTotalPeriods          = []byte("SPWithdrawalTotalPeriods")
	KeySPWithdrawalPeriodDurationSeconds = []byte("SPWithdrawalPeriodDurationSeconds")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	supervisors []string,
	gasLimit sdkmath.Int,
	gasUnconditionedAmount sdkmath.Int,
	gasAdjustCoefficient sdk.Dec,
	registrationSPDelegationAmount sdkmath.Int,
	spWithdrawalTotalPeriods uint,
	spWithdrawalPeriodDuration uint,
) Params {
	return Params{
		Supervisors:                       supervisors,
		GasLimit:                          sdk.IntProto{Int: gasLimit},
		GasUnconditionedAmount:            sdk.IntProto{Int: gasUnconditionedAmount},
		GasAdjustCoefficient:              sdk.DecProto{Dec: gasAdjustCoefficient},
		RegistrationSPDelegationAmount:    sdk.IntProto{Int: registrationSPDelegationAmount},
		SpWithdrawalTotalPeriods:          uint32(spWithdrawalTotalPeriods),
		SpWithdrawalPeriodDurationSeconds: uint32(spWithdrawalPeriodDuration),
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		nil,
		sdk.NewInt(1000000),
		sdk.NewInt(15000),
		sdk.NewDec(1),
		sdk.NewInt(5),
		52, 7*24*60*60, // 52 weeks
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeySupervisors, &p.Supervisors, validateSupervisors),
		paramtypes.NewParamSetPair(KeyGasLimit, &p.GasLimit, validateGasLimit),
		paramtypes.NewParamSetPair(KeyGasUnconditionedAmount, &p.GasUnconditionedAmount, validateGasUnconditionedAmount),
		paramtypes.NewParamSetPair(KeyGasAdjustCoefficient, &p.GasAdjustCoefficient, validateGasAdjustCoefficient),
		paramtypes.NewParamSetPair(KeyRegistrationSPDelegationAmount, &p.RegistrationSPDelegationAmount, validateRegistrationSPDelegationAmount),
		paramtypes.NewParamSetPair(KeySPWithdrawalTotalPeriods, &p.SpWithdrawalTotalPeriods, validateSPWithdrawalTotalPeriods),
		paramtypes.NewParamSetPair(KeySPWithdrawalPeriodDurationSeconds, &p.SpWithdrawalPeriodDurationSeconds, validateSPWithdrawalPeriodDurationSeconds),
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

	if err := validateRegistrationSPDelegationAmount(p.RegistrationSPDelegationAmount); err != nil {
		return fmt.Errorf("invalid registrationSPDelegationAmount: %w", err)
	}

	if err := validateSPWithdrawalTotalPeriods(p.SpWithdrawalTotalPeriods); err != nil {
		return fmt.Errorf("invalid spWithdrawalTotalPeriods: %w", err)
	}

	if err := validateSPWithdrawalPeriodDurationSeconds(p.SpWithdrawalPeriodDurationSeconds); err != nil {
		return fmt.Errorf("invalid spWithdrawalPeriodDurationSeconds: %w", err)
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateGasLimit(i interface{}) error {
	s, ok := i.(sdk.IntProto)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if s.Int.IsZero() || s.Int.IsNegative() {
		return fmt.Errorf("must be positive")
	}

	return nil
}

func validateGasUnconditionedAmount(i interface{}) error {
	s, ok := i.(sdk.IntProto)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if s.Int.IsNegative() {
		return fmt.Errorf("must not be negative")
	}

	return nil
}

func validateGasAdjustCoefficient(i interface{}) error {
	s, ok := i.(sdk.DecProto)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if s.Dec.IsZero() || s.Dec.IsNegative() {
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

func validateRegistrationSPDelegationAmount(i interface{}) error {
	s, ok := i.(sdk.IntProto)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if s.Int.IsNegative() {
		return fmt.Errorf("can not be negative")
	}

	return nil
}

func validateSPWithdrawalTotalPeriods(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("can not zero")
	}

	return nil
}

func validateSPWithdrawalPeriodDurationSeconds(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("can not zero")
	}

	return nil
}
