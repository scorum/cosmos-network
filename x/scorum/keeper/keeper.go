package keeper

import (
	"fmt"

	"cosmossdk.io/x/nft"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/scorum/cosmos-network/x/scorum/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		paramstore paramtypes.Subspace

		bankKeeper    types.BankKeeper
		stakingKeeper types.StakingKeeper
		accountKeeper types.AccountKeeper

		feeCollectorName string
		authority        string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	feeCollectorName string,
	authority string,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	// ensure nft module account is set
	if addr := accountKeeper.GetModuleAddress(nft.ModuleName); addr == nil {
		panic("the nft module account has not been set")
	}

	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramstore: ps,

		bankKeeper:       bankKeeper,
		accountKeeper:    accountKeeper,
		stakingKeeper:    stakingKeeper,
		feeCollectorName: feeCollectorName,
		authority:        authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Mint(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) error {
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return fmt.Errorf("failed to mint: %w", err)
	}

	if k.bankKeeper.BlockedAddr(addr) {
		for name, acc := range k.accountKeeper.GetModulePermissions() {
			if acc.GetAddress().Equals(addr) {
				if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, name, sdk.NewCoins(coin)); err != nil {
					return fmt.Errorf("failed to send minted coins: %w", err)
				}

				return nil
			}
		}
		return fmt.Errorf("%s is not a module and restricted to receive coins", addr)
	} else {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(coin)); err != nil {
			return fmt.Errorf("failed to send minted coins: %w", err)
		}
	}

	return nil
}

func (k Keeper) Burn(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) error {
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return fmt.Errorf("failed to send coins to burn: %w", err)
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return fmt.Errorf("failed to burn coins: %w", err)
	}

	return nil
}
