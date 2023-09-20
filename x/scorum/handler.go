package scorum

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/scorum/cosmos-network/x/scorum/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/scorum/cosmos-network/x/scorum/keeper"
)

func NewMintProposalHandler(k keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.MintProposal:
			return keeper.HandleMintProposal(ctx, k, c)

		default:
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized earn proposal content type: %T", c)
		}
	}
}
