package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
)

// NftKeeper defines the expected interface needed to store nft.
type NftKeeper interface {
	HasClass(ctx sdk.Context, classID string) bool
	SaveClass(ctx sdk.Context, class nft.Class) error

	GetNFT(ctx sdk.Context, classID, nftID string) (nft.NFT, bool)
	Mint(ctx sdk.Context, token nft.NFT, receiver sdk.AccAddress) error
	Update(ctx sdk.Context, token nft.NFT) error

	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress

	Send(goCtx context.Context, msg *nft.MsgSend) (*nft.MsgSendResponse, error)
}

// ScorumKeeper defines the expected interface needed to verify supervisors.
type ScorumKeeper interface {
	IsSupervisor(ctx sdk.Context, addr string) bool
}
