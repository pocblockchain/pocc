package token

import (
	"github.com/pocblockchain/pocc/x/token/client"
	"github.com/pocblockchain/pocc/x/token/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	QueryTokens       = types.QueryTokens
	QueryToken        = types.QueryToken
	QuerySymbols      = types.QuerySymbols
	QueryParameters   = types.QueryParameters
	DefaultParamspace = types.DefaultParamspace
	DefaultCodespace  = types.DefaultCodespace
)

type (
	QueryTokenInfo = types.QueryTokenInfo
	Params         = types.Params
)

var (
	ModuleCdc                        = types.ModuleCdc
	RegisterCodec                    = types.RegisterCodec
	DefaultParams                    = types.DefaultParams
	KeyTokenCacheSize                = types.KeyTokenCacheSize
	DisableTokenProposalHandler      = client.DisableTokenProposalHandler
	TokenParamsChangeProposalHandler = client.TokenParamsChangeProposalHandler
	NewTokenParamsChangeProposal     = types.NewTokenParamsChangeProposal
	NewDisableTokenProposal          = types.NewDisableTokenProposal
)
