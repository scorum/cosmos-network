package keeper

import (
	"github.com/scorum/cosmos-network/x/scorum/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServer returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
