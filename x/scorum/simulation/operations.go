package simulation

import (
	"math/rand"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/mikluke/co-pilot/slice"
	"github.com/scorum/cosmos-network/x/scorum/keeper"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

// nolint
const (
	opWeightMsgBurn      = "op_weight_msg_burn"
	defaultWeightMsgBurn = 10

	opWeightMsgMintGas      = "op_weight_msg_mint_gas"
	defaultWeightMsgMintGas = 10

	opWeightMsgConvertSCR2SP      = "op_weight_msg_convert_scr_2_sp"
	defaultWeightMsgConvertSCR2SP = 10

	opWeightMsgWithdrawSP      = "op_weight_msg_withdraw_sp"
	defaultWeightMsgWithdrawSP = 90

	opWeightMsgStopSPWithdrawal      = "op_weight_msg_stop_sp_withdrawal"
	defaultWeightMsgStopSPWithdrawal = 10
)

// WeightedOperations returns the all the gov module operations with their respective weights.
func WeightedOperations(
	simState module.SimulationState, ak types.AccountKeeper, bk bankkeeper.Keeper, k keeper.Keeper,
) []simtypes.WeightedOperation {
	var (
		weightMsgBurn             int
		weightMsgMintGas          int
		weightMsgConvertSCR2SP    int
		weightMsgWithdrawSP       int
		weightMsgStopSPWithdrawal int
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
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgConvertSCR2SP, &weightMsgConvertSCR2SP, nil,
		func(_ *rand.Rand) {
			weightMsgConvertSCR2SP = defaultWeightMsgConvertSCR2SP
		},
	)
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgWithdrawSP, &weightMsgWithdrawSP, nil,
		func(_ *rand.Rand) {
			weightMsgWithdrawSP = defaultWeightMsgWithdrawSP
		},
	)
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgStopSPWithdrawal, &weightMsgStopSPWithdrawal, nil,
		func(_ *rand.Rand) {
			weightMsgStopSPWithdrawal = defaultWeightMsgStopSPWithdrawal
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
		simulation.NewWeightedOperation(
			weightMsgConvertSCR2SP,
			SimulateMsgConvertSCR2SP(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgConvertSCR2SP,
			SimulateMsgWithdrawSP(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgStopSPWithdrawal,
			SimulateMsgStopSPWithdrawal(k, ak, bk),
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
					TxGen:         simapp.MakeTestEncodingConfig().TxConfig,
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
			TxGen:         simapp.MakeTestEncodingConfig().TxConfig,
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

func SimulateMsgConvertSCR2SP(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk bankkeeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgConvertSCR2SP, "accounts are empty"), nil, nil
		}

		owner, _ := simtypes.RandomAcc(r, accs)
		amount, err := simtypes.RandPositiveInt(r, math.NewInt(10000))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgConvertSCR2SP, "failed to rand int"), nil, nil
		}

		if err := k.Mint(ctx, owner.Address, sdk.NewCoin(types.SCRDenom, amount)); err != nil {
			panic(err)
		}

		msg := &types.MsgConvertSCR2SP{
			Owner:  owner.Address.String(),
			Amount: sdk.IntProto{Int: amount},
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simapp.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      owner,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(sdk.NewCoin(types.SCRDenom, amount)),
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgWithdrawSP(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk bankkeeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawSP, "accounts are empty"), nil, nil
		}

		owner, _ := simtypes.RandomAcc(r, accs)
		recipient, _ := simtypes.RandomAcc(r, accs)
		amount, err := simtypes.RandPositiveInt(r, math.NewInt(10000))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawSP, "failed to rand int"), nil, nil
		}

		if err := k.Mint(ctx, owner.Address, sdk.NewCoin(types.SPDenom, amount)); err != nil {
			panic(err)
		}

		msg := &types.MsgWithdrawSP{
			Owner:     owner.Address.String(),
			Recipient: recipient.Address.String(),
			Amount:    sdk.IntProto{Int: amount},
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simapp.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      owner,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(sdk.NewCoin(types.SPDenom, amount)),
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgStopSPWithdrawal(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk bankkeeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStopSPWithdrawal, "accounts are empty"), nil, nil
		}

		withdrawals := slice.Filter(k.ListAllWithdrawals(ctx), func(v types.SPWithdrawal) bool {
			return v.IsActive
		})
		if len(withdrawals) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStopSPWithdrawal, "no withdrawals to stop"), nil, nil
		}
		withdrawalToStop := withdrawals[0]
		if len(withdrawals) > 1 {
			withdrawalToStop = withdrawals[simtypes.RandIntBetween(r, 0, len(withdrawals)-1)]
		}

		msg := &types.MsgStopSPWithdrawal{
			Owner: withdrawalToStop.From,
			Id:    withdrawalToStop.Id,
		}

		simAcc, ok := simtypes.FindAccount(accs, sdk.MustAccAddressFromBech32(withdrawalToStop.From))
		if !ok {
			panic("account not found")
		}

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         simapp.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    simAcc,
			AccountKeeper: ak,
			Bankkeeper:    bk,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
