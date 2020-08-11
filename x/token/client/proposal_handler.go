package client

import (
	govclient "github.com/pocblockchain/pocc/x/gov/client"
	"github.com/pocblockchain/pocc/x/token/client/cli"
	"github.com/pocblockchain/pocc/x/token/client/rest"
)

// param change proposal handler
var (
	DisableTokenProposalHandler      = govclient.NewProposalHandler(cli.GetCmdDisableTokenProposal, rest.DisableTokenProposalRESTHandler)
	TokenParamsChangeProposalHandler = govclient.NewProposalHandler(cli.GetCmdTokenParamsChangeProposal, rest.TokenParamsChangeProposalRESTHandler)
)
