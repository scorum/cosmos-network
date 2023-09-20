package ante

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ZeroGasTxDecorator struct {
	sk ScorumKeeper
}

func NewZeroGasTxDecorator(sk ScorumKeeper) ZeroGasTxDecorator {
	return ZeroGasTxDecorator{
		sk: sk,
	}
}

func (d ZeroGasTxDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		for _, v := range msg.GetSigners() {
			if d.sk.IsSupervisor(ctx.WithGasMeter(NewFixedGasMeter(0, ctx.GasMeter().Limit())), v.String()) {
				return next(
					ctx.
						WithGasMeter(NewFixedGasMeter(0, ctx.GasMeter().Limit())).
						WithMinGasPrices(sdk.NewDecCoins()),
					tx,
					simulate,
				)
			}
		}
	}

	return next(ctx, tx, simulate)
}

type fixedGasMeter struct {
	limit    sdk.Gas
	consumed sdk.Gas
}

// NewFixedGasMeter returns a reference to a new basicGasMeter.
func NewFixedGasMeter(consumed, limit sdk.Gas) sdk.GasMeter {
	return &fixedGasMeter{
		limit:    limit,
		consumed: consumed,
	}
}

func (g *fixedGasMeter) GasConsumed() sdk.Gas {
	return g.consumed
}

func (g *fixedGasMeter) Limit() sdk.Gas {
	return g.limit
}

func (g *fixedGasMeter) GasRemaining() sdk.Gas {
	return math.MaxUint64
}

func (g *fixedGasMeter) GasConsumedToLimit() sdk.Gas {
	if g.IsPastLimit() {
		return g.limit
	}
	return g.consumed
}

func (g *fixedGasMeter) ConsumeGas(_ sdk.Gas, _ string) {
}

func (g *fixedGasMeter) RefundGas(_ sdk.Gas, _ string) {
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
