package rest

import (
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/types/rest"
	"github.com/pocblockchain/pocc/x/token/client/cli"
)

type (

	// TokenParamsChangeProposalReq defines a token params change request body.
	TokenParamsChangeProposalReq struct {
		BaseReq     rest.BaseReq         `json:"base_req" yaml:"base_req"`
		Title       string               `json:"title" yaml:"title"`
		Description string               `json:"description" yaml:"description"`
		Symbol      string               `json:"symbol" yaml:"symbol"`
		Changes     cli.ParamChangesJSON `json:"changes" yaml:"changes"`
		Deposit     sdk.Coins            `json:"deposit" yaml:"deposit"`
		Proposer    sdk.AccAddress       `json:"proposer" yaml:"proposer"`
	}

	// DisableTokenProposalReq defines a disable token request body.
	DisableTokenProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Symbol      string         `json:"symbol" yaml:"symbol"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
	}
)
