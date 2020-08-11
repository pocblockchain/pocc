package client

import (
	govclient "github.com/pocblockchain/pocc/x/gov/client"
	"github.com/pocblockchain/pocc/x/params/client/cli"
	"github.com/pocblockchain/pocc/x/params/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
