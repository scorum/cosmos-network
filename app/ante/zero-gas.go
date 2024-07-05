package ante

import (
	"fmt"
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

	for _, signer := range signers {
		v, err := ante.GetSignerAcc(ctx, d.ak, signer)
		if err != nil {
			return ctx, err
		}

		if d.sk.IsSupervisor(ctx.WithGasMeter(NewFixedGasMeter(0, ctx.GasMeter().Limit())), v.GetAddress().String()) {
			return next(
				ctx.
					WithGasMeter(NewFixedGasMeter(0, ctx.GasMeter().Limit())).
					WithMinGasPrices(sdk.NewDecCoins()),
				tx,
				simulate,
			)
		}
	}

	return next(ctx, tx, simulate)
}

type fixedGasMeter struct {
	limit    storetypes.Gas
	consumed storetypes.Gas
}

// NewFixedGasMeter returns a reference to a new basicGasMeter.
func NewFixedGasMeter(consumed, limit storetypes.Gas) storetypes.GasMeter {
	return &fixedGasMeter{
		limit:    limit,
		consumed: consumed,
	}
}

func (g *fixedGasMeter) GasConsumed() storetypes.Gas {
	return g.consumed
}

func (g *fixedGasMeter) Limit() storetypes.Gas {
	return g.limit
}

func (g *fixedGasMeter) GasRemaining() storetypes.Gas {
	return math.MaxUint64
}

func (g *fixedGasMeter) GasConsumedToLimit() storetypes.Gas {
	if g.IsPastLimit() {
		return g.limit
	}
	return g.consumed
}

func (g *fixedGasMeter) ConsumeGas(_ storetypes.Gas, _ string) {
}

func (g *fixedGasMeter) RefundGas(_ storetypes.Gas, _ string) {
}

func (g *fixedGasMeter) IsPastLimit() bool {
	return g.consumed > g.limit
}

func (g *fixedGasMeter) IsOutOfGas() bool {
	return g.consumed >= g.limit
}

func (g *fixedGasMeter) String() string {
	return fmt.Sprintf("FixedGasMeter:\n  consumed: %d", g.consumed)
}
