package types

import (
	"context"

	"cosmossdk.io/x/nft"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NftKeeper defines the expected interface needed to store nft.
type NftKeeper interface {
	HasClass(ctx context.Context, classID string) bool
	SaveClass(ctx context.Context, class nft.Class) error

	GetNFT(ctx context.Context, classID, nftID string) (nft.NFT, bool)
	Mint(ctx context.Context, token nft.NFT, receiver sdk.AccAddress) error
	Update(ctx context.Context, token nft.NFT) error

	GetOwner(ctx context.Context, classID string, nftID string) sdk.AccAddress

	Send(ctx context.Context, msg *nft.MsgSend) (*nft.MsgSendResponse, error)
}

// ScorumKeeper defines the expected interface needed to verify supervisors.
type ScorumKeeper interface {
	IsSupervisor(ctx sdk.Context, addr string) bool
}
