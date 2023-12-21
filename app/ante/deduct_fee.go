package ante

import (
	"math"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// DeductFeeDecorator deducts fees from the first signer of the tx
// This decorator ALLOWS zero fees
// CONTRACT: Tx must implement FeeTx interface to use DeductFeeDecorator
type DeductFeeDecorator struct {
	ante.DeductFeeDecorator
}

func NewDeductFeeDecorator(ak ante.AccountKeeper, bk types.BankKeeper, fk ante.FeegrantKeeper, tfc ante.TxFeeChecker) DeductFeeDecorator {
	return DeductFeeDecorator{
		DeductFeeDecorator: ante.NewDeductFeeDecorator(ak, bk, fk, tfc),
	}
}

func (dfd DeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if feeTx.GetGas() == 0 {
		// allow gas-free and fee-free transactions
		return next(ctx.WithPriority(math.MaxInt64), tx, simulate)
	}

	return dfd.DeductFeeDecorator.AnteHandle(ctx, tx, simulate, next)
}
