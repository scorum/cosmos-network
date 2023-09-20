package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

func (m msgServer) RegisterAccount(goCtx context.Context, msg *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := sdk.AccAddressFromBech32(msg.Supervisor); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid supervisor address: %s", err)
	}

	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address for registration: %s", err)
	}

	if !m.IsSupervisor(ctx, msg.Supervisor) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "register_account is only allowed to supervisors")
	}

	if m.accountKeeper.HasAccount(ctx, addr) {
		return nil, errorsmod.Wrap(sdkerrors.ErrConflict, "address is already registered")
	}

	m.accountKeeper.NewAccountWithAddress(ctx, addr)

	if err := m.Keeper.Mint(ctx, addr, sdk.NewCoin(types.SPDenom, m.Keeper.GetParams(ctx).RegistrationSPDelegationAmount.Int)); err != nil {
		return nil, err
	}

	if err := m.Keeper.Mint(ctx, addr, sdk.NewCoin(types.GasDenom, m.Keeper.GetParams(ctx).GasLimit.Int)); err != nil {
		return nil, err
	}

	return &types.MsgRegisterAccountResponse{}, nil
}
