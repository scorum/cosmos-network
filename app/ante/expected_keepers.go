package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type ScorumKeeper interface {
	IsSupervisor(ctx sdk.Context, addr string) bool

	SetAddressToRestoreGas(ctx sdk.Context, addr sdk.AccAddress)
}

type AccountKeeper interface {
	ante.AccountKeeper

	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool
}

type BankKeeper interface {
	authtypes.BankKeeper

	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}
