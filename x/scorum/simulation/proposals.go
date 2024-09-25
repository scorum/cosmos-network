package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/address"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/scorum/cosmos-network/x/scorum/types"
)

//nolint:gosec
const (
	opWeightMintProposal = "op_weight_mint_proposal"
	defaultWeightMsgMint = 10
)

// ProposalMsgs defines the module weighted proposals' contents
func ProposalMsgs() []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMintProposal,
			defaultWeightMsgMint,
			SimulateMsgMint,
		),
	}
}

// SimulateMsgMint returns a random MsgMint
func SimulateMsgMint(r *rand.Rand, _ sdk.Context, accounts []simtypes.Account) sdk.Msg {
	// use the default gov module account address as authority
	var authority sdk.AccAddress = address.Module("gov")

	recipient, _ := simtypes.RandomAcc(r, accounts)
	amount, _ := simtypes.RandPositiveInt(r, math.NewInt(10000000))

	return &types.MsgMint{
		Authority: authority.String(),
		Recipient: recipient.Address.String(),
		Amount:    sdk.NewCoin(types.SCRDenom, amount),
	}
}
