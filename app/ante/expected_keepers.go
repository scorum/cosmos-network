package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

type ScorumKeeper interface {
	GetParams(ctx sdk.Context) (params types.Params)
	IsSupervisor(ctx sdk.Context, addr string) bool

	SetAddressToRestoreGas(ctx sdk.Context, addr sdk.AccAddress)

	Mint(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) error
}

type AccountKeeper interface {
	ante.AccountKeeper

	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool
}

type BankKeeper interface {
	authtypes.BankKeeper

	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}
