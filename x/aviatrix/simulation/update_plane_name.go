package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/scorum/cosmos-network/x/aviatrix/keeper"
	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

func SimulateMsgUpdatePlaneName(
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgUpdatePlaneName{
			Supervisor: simAccount.Address.String(),
		}

		// TODO: Handling the MsgUpdatePlaneName simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "MsgUpdatePlaneName simulation not implemented"), nil, nil
	}
}
