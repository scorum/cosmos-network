package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type TrackGasConsumedDecorator struct {
	ak AccountKeeper
	bk BankKeeper
	sk ScorumKeeper
}

func NewTrackGasConsumedDecorator(ak AccountKeeper, bk BankKeeper, sk ScorumKeeper) TrackGasConsumedDecorator {
	return TrackGasConsumedDecorator{
		ak: ak,
		sk: sk,
		bk: bk,
	}
}

func (d TrackGasConsumedDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	for _, v := range tx.GetMsgs() {
		for _, addr := range v.GetSigners() {
			if !d.ak.HasAccount(ctx, addr) {
				return sdk.Context{}, errorsmod.Wrap(sdkerrors.ErrUnknownAddress, "address is not registered")
			}

			d.sk.SetAddressToRestoreGas(ctx, addr)
		}
	}

	return next(ctx, tx, simulate)
}
