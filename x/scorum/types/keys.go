package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func init() {
	sdk.SetCoinDenomRegex(func() string { return `[a-zA-Z][a-zA-Z0-9/:._-]{1,127}` })
}

const (
	// ModuleName defines the module name
	ModuleName = "scorum"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var GasPrice = sdk.NewCoins(sdk.NewCoin(GasDenom, sdk.NewInt(1)))

func KeyPrefix(p string) []byte {
	return []byte(p)
}
