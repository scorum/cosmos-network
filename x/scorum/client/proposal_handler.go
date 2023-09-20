package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/scorum/cosmos-network/x/scorum/client/cli"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.CmdSubmitMintProposal)
