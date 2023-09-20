package aviatrix

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/scorum/cosmos-network/testutil/sample"
	aviatrixsimulation "github.com/scorum/cosmos-network/x/aviatrix/simulation"
	"github.com/scorum/cosmos-network/x/aviatrix/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = aviatrixsimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

// nolint
const (
	opWeightMsgCreatePlane          = "op_weight_msg_create_plane"
	defaultWeightMsgCreatePlane int = 100

	opWeightMsgUpdatePlaneName          = "op_weight_msg_update_plane_name"
	defaultWeightMsgUpdatePlaneName int = 50
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	nftGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&nftGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {

	return []simtypes.ParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreatePlane int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreatePlane, &weightMsgCreatePlane, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePlane = defaultWeightMsgCreatePlane
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePlane,
		aviatrixsimulation.SimulateMsgCreatePlane(am.keeper),
	))

	var weightMsgUpdatePlaneName int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdatePlaneName, &weightMsgUpdatePlaneName, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePlaneName = defaultWeightMsgUpdatePlaneName
		},
	)

	return operations
}
