package ante

import (
	"math"

	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"

	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ZeroGasTxDecorator struct {
	sk ScorumKeeper
	ak AccountKeeper
}

func NewZeroGasTxDecorator(ak AccountKeeper, sk ScorumKeeper) ZeroGasTxDecorator {
	return ZeroGasTxDecorator{
		ak: ak,
		sk: sk,
	}
}

func (d ZeroGasTxDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	sigTx, ok := tx.(authsigning.Tx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	signers, err := sigTx.GetSigners()
	if !ok {
		return ctx, err
	}

	zeroGasCtx := ctx.WithGasMeter(NewZeroGasMeter())

	for _, signer := range signers {
		v, err := ante.GetSignerAcc(zeroGasCtx, d.ak, signer)
		if err != nil {
			return ctx, err
		}

		if d.sk.IsSupervisor(zeroGasCtx, v.GetAddress().String()) {
			return next(
				ctx.
					WithGasMeter(NewZeroGasMeter()).
					WithMinGasPrices(sdk.NewDecCoins()),
				tx,
				simulate,
			)
		}
	}

	return next(ctx, tx, simulate)
}

type zeroGasMeter struct{}

// NewZeroGasMeter returns a reference to a new basicGasMeter.
func NewZeroGasMeter() storetypes.GasMeter {
	return &zeroGasMeter{}
}

func (g *zeroGasMeter) GasConsumed() storetypes.Gas {
	return 0
}

func (g *zeroGasMeter) Limit() storetypes.Gas {
	return 0
}

func (g *zeroGasMeter) GasRemaining() storetypes.Gas {
	return math.MaxUint64
}

func (g *zeroGasMeter) GasConsumedToLimit() storetypes.Gas {
	return 0
}

func (g *zeroGasMeter) ConsumeGas(_ storetypes.Gas, _ string) {
}

func (g *zeroGasMeter) RefundGas(_ storetypes.Gas, _ string) {
}

func (g *zeroGasMeter) IsPastLimit() bool {
	return false
}

func (g *zeroGasMeter) IsOutOfGas() bool {
	return false
}

func (g *zeroGasMeter) String() string {
	return "ZerGasMeter:\n  consumed: 0"
}
