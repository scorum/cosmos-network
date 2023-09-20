package keeper

import (
	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

var _ types.QueryServer = Keeper{}
