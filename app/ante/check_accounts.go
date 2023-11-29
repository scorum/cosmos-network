package ante

import (
	"fmt"
	"reflect"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
)

// CheckAddressesDecorator checks if all addresses are registered
type CheckAddressesDecorator struct {
	ak AccountKeeper
	sk ScorumKeeper
}

func NewCheckAddressesDecorator(ak AccountKeeper, sk ScorumKeeper) CheckAddressesDecorator {
	return CheckAddressesDecorator{
		ak: ak,
		sk: sk,
	}
}

func (d CheckAddressesDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		for _, addr := range extractAddresses(msg) {
			if !d.ak.HasAccount(ctx, addr) {
				if err := d.sk.Mint(ctx, addr, sdk.NewCoin(scorumtypes.GasDenom, d.sk.GetParams(ctx).GasLimit.Int)); err != nil {
					return sdk.Context{}, errorsmod.Wrap(sdkerrors.ErrPanic, fmt.Sprintf("failed to mint gas to new account: %s", err.Error()))
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}

func extractAddresses(i interface{}) []sdk.AccAddress {
	if v, ok := i.(sdk.AccAddress); ok {
		return []sdk.AccAddress{v}
	}

	var out []sdk.AccAddress
	v := reflect.ValueOf(i)

	switch v.Kind() {
	case reflect.String:
		if addr, err := sdk.AccAddressFromBech32(v.String()); err == nil {
			out = []sdk.AccAddress{addr}
		}
	case reflect.Ptr, reflect.Interface:
		if !v.IsNil() {
			out = append(out, extractAddresses(v.Elem().Interface())...)
		}
	case reflect.Map:
		for iter := v.MapRange(); iter.Next(); {
			out = append(out, extractAddresses(iter.Key().Interface())...)
			out = append(out, extractAddresses(iter.Value().Interface())...)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			out = append(out, extractAddresses(v.Index(i).Interface())...)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			if !f.CanInterface() {
				continue
			}

			out = append(out, extractAddresses(f.Interface())...)
		}
	}

	return out
}
