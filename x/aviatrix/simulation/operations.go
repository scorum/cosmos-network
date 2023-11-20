package simulation

import (
	"math/rand"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/google/uuid"
	"github.com/scorum/cosmos-network/x/aviatrix/types"
	scorumkeeper "github.com/scorum/cosmos-network/x/scorum/keeper"
)

// nolint
const (
	opWeightMsgCreatePlane          = "op_weight_msg_create_plane"
	defaultWeightMsgCreatePlane int = 100

	opWeightMsgUpdatePlaneExperience          = "op_weight_msg_update_plane_experience"
	defaultWeightMsgUpdatePlaneExperience int = 20

	opWeightMsgAdjustPlaneExperience          = "op_weight_msg_adjust_plane_experience"
	defaultWeightMsgAdjustPlaneExperience int = 100
)

// WeightedOperations returns the all the gov module operations with their respective weights.
func WeightedOperations(
	simState module.SimulationState,
	sk scorumkeeper.Keeper,
	nk nftkeeper.Keeper,
	ak accountkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
) []simtypes.WeightedOperation {
	var (
		weightMsgCreatePlane           int
		weightMsgUpdatePlaneExperience int
		weightMsgAdjustPlaneExperience int
	)

	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreatePlane, &weightMsgCreatePlane, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePlane = defaultWeightMsgCreatePlane
		},
	)
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdatePlaneExperience, &weightMsgUpdatePlaneExperience, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePlaneExperience = defaultWeightMsgUpdatePlaneExperience
		},
	)
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAdjustPlaneExperience, &weightMsgAdjustPlaneExperience, nil,
		func(_ *rand.Rand) {
			weightMsgAdjustPlaneExperience = defaultWeightMsgAdjustPlaneExperience
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreatePlane,
			SimulateMsgCreatePlane(sk, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgUpdatePlaneExperience,
			SimulateMsgUpdatePlaneExperience(sk, nk, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgAdjustPlaneExperience,
			SimulateMsgAdjustPlaneExperience(sk, nk, ak, bk),
		),
	}
}

func SimulateMsgCreatePlane(
	sk scorumkeeper.Keeper,
	ak accountkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		supervisors := sk.GetParams(ctx).Supervisors
		if len(supervisors) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePlane, "empty supervisors"), nil, nil
		}

		supervisor := supervisors[simtypes.RandIntBetween(r, 0, len(supervisors))]
		supervisorAcc, ok := simtypes.FindAccount(accs, sdk.MustAccAddressFromBech32(supervisor))
		if !ok {
			panic("account not found")
		}

		exp, err := simtypes.RandPositiveInt(r, math.NewInt(10000000))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePlane, "unable to generate positive amount"), nil, err
		}

		id, err := uuid.NewRandomFromReader(r)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePlane, "unable to generate uuid"), nil, err
		}

		owner, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgCreatePlane{
			Id:         id.String(),
			Supervisor: supervisor,
			Owner:      owner.Address.String(),
			Meta: &types.PlaneMeta{
				Experience: exp.Uint64(),
			},
		}

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         simapp.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    supervisorAcc,
			AccountKeeper: ak,
			Bankkeeper:    bk,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgUpdatePlaneExperience(
	sk scorumkeeper.Keeper,
	nk nftkeeper.Keeper,
	ak accountkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		supervisors := sk.GetParams(ctx).Supervisors
		if len(supervisors) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdatePlaneExperience, "empty supervisors"), nil, nil
		}

		supervisor := supervisors[simtypes.RandIntBetween(r, 0, len(supervisors))]
		supervisorAcc, ok := simtypes.FindAccount(accs, sdk.MustAccAddressFromBech32(supervisor))
		if !ok {
			panic("account not found")
		}

		nfts := nk.GetNFTsOfClass(ctx, types.NftClassID)
		if len(nfts) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdatePlaneExperience, "empty nfts"), nil, nil
		}

		nft := nfts[0]
		if len(nfts) > 1 {
			nft = nfts[simtypes.RandIntBetween(r, 0, len(nfts)-1)]
		}

		amount, err := simtypes.RandPositiveInt(r, math.NewInt(1000))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdatePlaneExperience, "failed to random amount"), nil, nil
		}

		msg := &types.MsgUpdatePlaneExperience{
			Supervisor: supervisor,
			Id:         nft.Id,
			Amount:     amount.Uint64(),
		}

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         simapp.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    supervisorAcc,
			AccountKeeper: ak,
			Bankkeeper:    bk,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgAdjustPlaneExperience(
	sk scorumkeeper.Keeper,
	nk nftkeeper.Keeper,
	ak accountkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		supervisors := sk.GetParams(ctx).Supervisors
		if len(supervisors) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAdjustPlaneExperience, "empty supervisors"), nil, nil
		}

		supervisor := supervisors[simtypes.RandIntBetween(r, 0, len(supervisors))]
		supervisorAcc, ok := simtypes.FindAccount(accs, sdk.MustAccAddressFromBech32(supervisor))
		if !ok {
			panic("account not found")
		}

		nfts := nk.GetNFTsOfClass(ctx, types.NftClassID)
		if len(nfts) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAdjustPlaneExperience, "empty nfts"), nil, nil
		}

		nft := nfts[0]
		if len(nfts) > 1 {
			nft = nfts[simtypes.RandIntBetween(r, 0, len(nfts)-1)]
		}

		amount := simtypes.RandIntBetween(r, -1000, 1000)
		if amount == 0 {
			amount = 1
		}

		msg := &types.MsgAdjustPlaneExperience{
			Supervisor: supervisor,
			Id:         nft.Id,
			Amount:     int64(amount),
		}

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         simapp.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    supervisorAcc,
			AccountKeeper: ak,
			Bankkeeper:    bk,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
