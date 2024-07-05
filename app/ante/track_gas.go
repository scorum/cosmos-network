package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
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
	sigTx, ok := tx.(authsigning.Tx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	signers, err := sigTx.GetSigners()
	if !ok {
		return ctx, err
	}

	for _, signer := range signers {
		addr, err := ante.GetSignerAcc(ctx, d.ak, signer)
		if err != nil {
			return ctx, err
		}

		if !d.ak.HasAccount(ctx, addr.GetAddress()) {
			return sdk.Context{}, errorsmod.Wrap(sdkerrors.ErrUnknownAddress, "address is not registered")
		}

		d.sk.SetAddressToRestoreGas(ctx, addr.GetAddress())
	}

	return next(ctx, tx, simulate)
}
