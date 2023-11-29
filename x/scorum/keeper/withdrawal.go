package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/tendermint/tendermint/libs/json"
)

func (k Keeper) ListAllWithdrawals(ctx sdk.Context) []types.SPWithdrawal {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), withdrawalPrefix)

	out := make([]types.SPWithdrawal, 0)

	it := store.Iterator(nil, nil)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		var w types.SPWithdrawal
		k.cdc.MustUnmarshal(it.Value(), &w)

		out = append(out, w)
	}

	return out
}

func (k Keeper) ListWithdrawals(ctx sdk.Context, owner sdk.AccAddress) []types.SPWithdrawal {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(withdrawalPrefix, owner...))

	out := make([]types.SPWithdrawal, 0)

	it := store.Iterator(nil, nil)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		w, _ := k.GetSPWithdrawal(ctx, owner, string(it.Key()))
		out = append(out, w)
	}

	return out
}

func (k Keeper) GetSPWithdrawal(ctx sdk.Context, owner sdk.AccAddress, id string) (types.SPWithdrawal, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), withdrawalPrefix)

	bz := prefix.NewStore(store, owner).Get([]byte(id))
	if bz == nil {
		return types.SPWithdrawal{}, false
	}

	var w types.SPWithdrawal
	k.cdc.MustUnmarshal(bz, &w)

	return w, true
}

func (k Keeper) SetSPWithdrawal(ctx sdk.Context, w types.SPWithdrawal) {
	owner := sdk.MustAccAddressFromBech32(w.From)
	key := getWithdrawalKey(owner, w.Id)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(withdrawalPrefix, owner...))

	if old, has := k.GetSPWithdrawal(ctx, owner, w.Id); has {
		k.deleteWithdrawalFromIndex(ctx, old.PeriodTime(old.ProcessedPeriod+1), key)
	}

	if w.IsActive {
		k.addWithdrawalToIndex(ctx, w.PeriodTime(w.ProcessedPeriod+1), key)
	}

	store.Set([]byte(w.Id), k.cdc.MustMarshal(&w))
}

func (k Keeper) addWithdrawalToIndex(ctx sdk.Context, t uint64, key []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), withdrawalTimeIdxPrefix)

	var keys [][]byte
	bz := store.Get(uint64ToBytes(t))
	if bz != nil {
		if err := json.Unmarshal(bz, &keys); err != nil {
			panic(fmt.Errorf("failed to unmarshal withdrawals for index key %d: %w", t, err))
		}
	}
	for _, v := range keys {
		if bytes.Equal(v, key) {
			return
		}
	}
	keys = append(keys, key)

	bz, err := json.Marshal(keys)
	if err != nil {
		panic(fmt.Errorf("failed to marshal withdrawals for index key %d: %w", t, err))
	}
	store.Set(uint64ToBytes(t), bz)
}

func (k Keeper) deleteWithdrawalFromIndex(ctx sdk.Context, t uint64, key []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), withdrawalTimeIdxPrefix)

	var keys [][]byte
	bz := store.Get(uint64ToBytes(t))
	if bz != nil {
		if err := json.Unmarshal(bz, &keys); err != nil {
			panic(fmt.Errorf("failed to unmarshal withdrawals for index key %d: %w", t, err))
		}
	}

	out := make([][]byte, 0, len(keys))
	for _, v := range keys {
		if bytes.Equal(v, key) {
			continue
		}
		out = append(out, v)
	}

	if len(out) == 0 {
		store.Delete(uint64ToBytes(t))

		return
	}

	bz, err := json.Marshal(out)
	if err != nil {
		panic(fmt.Errorf("failed to marshal withdrawals for index key %d: %w", t, err))
	}
	store.Set(uint64ToBytes(t), bz)
}

func (k Keeper) listWithdrawalsByNextWithdrawalTime(ctx sdk.Context, t uint64) []types.SPWithdrawal {
	idx := prefix.NewStore(ctx.KVStore(k.storeKey), withdrawalTimeIdxPrefix)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), withdrawalPrefix)

	out := make([]types.SPWithdrawal, 0)

	it := idx.Iterator(nil, uint64ToBytes(t+1))
	defer it.Close()
	for ; it.Valid(); it.Next() {
		var keys [][]byte
		if err := json.Unmarshal(it.Value(), &keys); err != nil {
			panic(fmt.Errorf("failed to unmarshal keys array: %w", err))
		}

		for _, key := range keys {
			var toWithdraw types.SPWithdrawal
			k.cdc.MustUnmarshal(store.Get(key), &toWithdraw)
			out = append(out, toWithdraw)
		}
	}

	return out
}

func (k Keeper) WithdrawSP(ctx sdk.Context, timestamp uint64) {
	for _, w := range k.listWithdrawalsByNextWithdrawalTime(ctx, timestamp) {
		balance := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(w.From), types.SPDenom).Amount

		toWithdraw := w.ToWithdraw(timestamp)
		if balance.LT(toWithdraw) {
			toWithdraw = balance
		}

		from := sdk.MustAccAddressFromBech32(w.From)
		recipient := sdk.MustAccAddressFromBech32(w.To)
		if toWithdraw.IsPositive() {
			if err := k.Burn(ctx, from, sdk.NewCoin(types.SPDenom, toWithdraw)); err != nil {
				panic(fmt.Errorf("failed to burn(%s %s): %w", w.Id, toWithdraw, err))
			}
			if err := k.Mint(ctx, recipient, sdk.NewCoin(types.SCRDenom, toWithdraw)); err != nil {
				panic(fmt.Errorf("failed to mint(%s %s): %w", w.Id, toWithdraw, err))
			}

			w.ProcessedPeriod = w.CurrentPeriod(timestamp)
		}

		w.IsActive = balance.GT(toWithdraw) && w.PeriodTime(w.ProcessedPeriod+1) > 0
		k.SetSPWithdrawal(ctx, w)
	}
}

func uint64ToBytes(v uint64) []byte {
	out := make([]byte, 8)
	binary.BigEndian.PutUint64(out, v)

	return out
}

func getWithdrawalKey(owner sdk.AccAddress, id string) []byte {
	return append(owner, []byte(id)...)
}
