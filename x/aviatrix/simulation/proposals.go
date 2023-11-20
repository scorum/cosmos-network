package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return []simtypes.WeightedProposalContent{}
}
