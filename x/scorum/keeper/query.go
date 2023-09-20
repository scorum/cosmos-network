package keeper

import (
	"github.com/scorum/cosmos-network/x/scorum/types"
)

type queryServer struct {
	Keeper
}

// NewQueryServer returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServer(keeper Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

var _ types.QueryServer = queryServer{}
