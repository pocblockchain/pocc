package client

import (
	"github.com/pocblockchain/pocc/x/distribution/client/cli"
	"github.com/pocblockchain/pocc/x/distribution/client/rest"
	govclient "github.com/pocblockchain/pocc/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
