package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

func (m msgServer) MintGas(goCtx context.Context, msg *types.MsgMintGas) (*types.MsgMintGasResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !m.IsSupervisor(ctx, msg.Supervisor) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "mint-gas is only allowed to supervisors")
	}

	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address: %s", err)
	}

	if err := m.Keeper.Mint(ctx, address, sdk.NewCoin(types.GasDenom, msg.Amount.Int)); err != nil {
		return nil, err
	}

	return &types.MsgMintGasResponse{}, nil
}
