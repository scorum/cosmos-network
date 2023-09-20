package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

func (m msgServer) ConvertSCR2SP(goCtx context.Context, msg *types.MsgConvertSCR2SP) (*types.MsgConvertSCR2SPResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	coin := sdk.NewCoin(types.SCRDenom, msg.Amount.Int)
	if m.bankKeeper.GetBalance(ctx, owner, types.SCRDenom).IsLT(coin) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "balance is lower than requested exchange")
	}

	if err := m.Keeper.Burn(ctx, owner, coin); err != nil {
		return nil, err
	}
	coin.Denom = types.SPDenom
	if err := m.Keeper.Mint(ctx, owner, coin); err != nil {
		return nil, err
	}

	return &types.MsgConvertSCR2SPResponse{}, nil
}
