package simulation

import (
	"math/rand"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/scorum/cosmos-network/x/scorum/keeper"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

// nolint
const (
	opWeightMsgBurn      = "op_weight_msg_burn"
	defaultWeightMsgBurn = 10

	opWeightMsgMintGas      = "op_weight_msg_mint_gas"
	defaultWeightMsgMintGas = 10
)

// WeightedOperations returns the all the gov module operations with their respective weights.
func WeightedOperations(
	simState module.SimulationState, ak types.AccountKeeper, bk bankkeeper.Keeper, k keeper.Keeper,
) []simtypes.WeightedOperation {
	var (
		weightMsgBurn    int
		weightMsgMintGas int
	)

	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgBurn, &weightMsgBurn, nil,
		func(_ *rand.Rand) {
			weightMsgBurn = defaultWeightMsgBurn
		},
	)
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgMintGas, &weightMsgMintGas, nil,
		func(_ *rand.Rand) {
			weightMsgMintGas = defaultWeightMsgMintGas
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgBurn,
			SimulateMsgBurn(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgMintGas,
			SimulateMsgMintGas(k, ak, bk),
		),
	}
}

func SimulateMsgBurn(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk bankkeeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBurn, "accounts are empty"), nil, nil
		}

		supervisor := accs[0]
		if !k.IsSupervisor(ctx, supervisor.Address.String()) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBurn, "first acc is not a supervisor"), nil, nil
		}

		balances := bk.GetAllBalances(ctx, supervisor.Address)
		if len(balances) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBurn, "empty balance"), nil, nil
		}

		for _, b := range balances {
			if b.Amount.GT(sdk.OneInt()) {
				msg := &types.MsgBurn{
					Supervisor: supervisor.Address.String(),
					Amount:     sdk.NewCoin(b.Denom, sdk.OneInt()),
				}

				txCtx := simulation.OperationInput{
					R:             r,
					App:           app,
					TxGen:         moduletestutil.MakeTestEncodingConfig().TxConfig,
					Cdc:           nil,
					Msg:           msg,
					MsgType:       msg.Type(),
					Context:       ctx,
					SimAccount:    supervisor,
					AccountKeeper: ak,
					Bankkeeper:    bk,
					ModuleName:    types.ModuleName,
				}

				return simulation.GenAndDeliverTxWithRandFees(txCtx)
			}
		}

		return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBurn, "empty balance"), nil, nil
	}
}

func SimulateMsgMintGas(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk bankkeeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBurn, "accounts are empty"), nil, nil
		}

		supervisor := accs[0]
		if !k.IsSupervisor(ctx, supervisor.Address.String()) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBurn, "first acc is not a supervisor"), nil, nil
		}

		addr, _ := simtypes.RandomAcc(r, accs)
		amount, err := simtypes.RandPositiveInt(r, math.NewInt(1000000))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgMintGas, "failed to rand int"), nil, nil
		}

		msg := &types.MsgMintGas{
			Supervisor: supervisor.Address.String(),
			Address:    addr.Address.String(),
			Amount:     sdk.IntProto{Int: amount},
		}

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    supervisor,
			AccountKeeper: ak,
			Bankkeeper:    bk,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
