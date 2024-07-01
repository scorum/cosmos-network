package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/scorum/cosmos-network/x/scorum/types"
)

func (m msgServer) Mint(ctx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	if m.authority != msg.Authority {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "expected %s got %s", m.authority, msg.Authority)
	}

	recipient, err := m.accountKeeper.AddressCodec().StringToBytes(msg.Recipient)
	if err != nil {
		return nil, err
	}

	if err := m.Keeper.Mint(sdk.UnwrapSDKContext(ctx), recipient, msg.Amount); err != nil {
		return nil, err
	}

	return &types.MsgMintResponse{}, nil
}
