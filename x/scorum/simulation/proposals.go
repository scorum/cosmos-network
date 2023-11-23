package simulation

import (
	"fmt"
	"math/rand"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

// nolint
const (
	OpWeightMintProposal = "op_weight_mint_proposal"
)

func ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	var (
		defaultWeightMintProposal = 5
	)

	return []simtypes.WeightedProposalContent{
		simulation.NewWeightedProposalContent(
			OpWeightMintProposal,
			defaultWeightMintProposal,
			SimulateMintProposalContent(),
		),
	}
}

func SimulateMintProposalContent() simtypes.ContentSimulatorFn {
	numProposals := 0

	return func(r *rand.Rand, _ sdk.Context, accounts []simtypes.Account) simtypes.Content {
		title := fmt.Sprintf("title from SimulateMintProposalContent-%d", numProposals)
		desc := fmt.Sprintf("desc from SimulateMintProposalContent-%d. Random short desc: %s",
			numProposals, simtypes.RandStringOfLength(r, 20))
		recipient, _ := simtypes.RandomAcc(r, accounts)
		amount := sdk.NewCoin(types.SCRDenom, simtypes.RandomAmount(r, math.NewInt(10000000)))
		numProposals++
		return types.NewMintProposal(title, desc, recipient.Address, amount)
	}
}
