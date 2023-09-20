package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

func HandleMintProposal(ctx sdk.Context, k Keeper, c *types.MintProposal) error {
	recipient, err := sdk.AccAddressFromBech32(c.Recipient)
	if err != nil {
		return err
	}

	if err := k.Mint(ctx, recipient, c.Amount); err != nil {
		return err
	}

	return nil
}
