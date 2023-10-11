package ante

import (
	"reflect"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
	"golang.org/x/exp/slices"
)

// CheckAddressesDecorator checks if all addresses are registered
type CheckAddressesDecorator struct {
	ak AccountKeeper
	sk ScorumKeeper

	ignore map[reflect.Type]struct{}
}

func NewCheckAddressesDecorator(ak AccountKeeper, sk ScorumKeeper) CheckAddressesDecorator {
	return CheckAddressesDecorator{
		ak: ak,
		sk: sk,

		ignore: map[reflect.Type]struct{}{
			reflect.TypeOf(&scorumtypes.MsgRegisterAccount{}): {},
		},
	}
}

func (d CheckAddressesDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		if _, ok := d.ignore[reflect.TypeOf(msg)]; ok || slices.Contains(strings.Split(sdk.MsgTypeURL(msg), "."), "ibc") {
			continue
		}

		for _, addr := range extractAddresses(msg) {
			if d.sk.IsSupervisor(ctx, addr.String()) {
				continue
			}

			if d.ak.HasAccount(ctx, addr) {
				continue
			}

			return sdk.Context{}, errorsmod.Wrap(sdkerrors.ErrUnknownAddress, "address is not registered in ScorumModule")
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
