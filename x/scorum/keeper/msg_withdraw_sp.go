package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

func (m msgServer) WithdrawSP(goCtx context.Context, msg *types.MsgWithdrawSP) (*types.MsgWithdrawSPResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address: %s", err)
	}

	reservedSP := m.getReservedSP(ctx, owner)
	if balance := m.bankKeeper.GetBalance(ctx, owner, types.SPDenom).Amount; reservedSP.Add(msg.Amount.Int).GT(balance) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"insufficient balance: reserved %s, total %s, requested %s",
			reservedSP, balance, msg.Amount,
		)
	}

	params := m.Keeper.GetParams(ctx)

	w := types.SPWithdrawal{
		Id:   uuid.NewSHA1(uuid.NameSpaceOID, ctx.TxBytes()).String(),
		From: msg.Owner,
		To:   msg.Recipient,

		Total: msg.Amount,

		PeriodDurationInSeconds: params.SpWithdrawalPeriodDurationSeconds,
		TotalPeriods:            params.SpWithdrawalTotalPeriods,

		ProcessedPeriod: 0,
		IsActive:        true,

		CreatedAt: uint64(ctx.BlockTime().Unix()),
	}
	m.Keeper.SetSPWithdrawal(ctx, w)

	return &types.MsgWithdrawSPResponse{
		WithdrawalId: w.Id,
	}, nil
}

func (m msgServer) getReservedSP(ctx sdk.Context, owner sdk.AccAddress) sdkmath.Int {
	total := sdk.ZeroInt()

	for _, w := range m.Keeper.ListWithdrawals(ctx, owner) {
		if !w.IsActive {
			continue
		}

		total = total.Add(w.Total.Int.Sub(w.WithdrownByPeriod(w.ProcessedPeriod)))
	}

	return total
}
