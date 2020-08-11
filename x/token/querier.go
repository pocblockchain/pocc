package token

import (
	"fmt"
	"github.com/pocblockchain/pocc/codec"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/token/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

//NewQuerier create a token's query handler
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryToken:
			return queryToken(ctx, req, keeper)
		case types.QuerySymbols:
			return querySymbols(ctx, keeper)
		case types.QueryTokens:
			return queryTokens(ctx, keeper)
		case types.QueryParameters:
			return queryParams(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest(path[0])
		}
	}
}

//=====
func queryToken(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var ti QueryTokenInfo
	if err := keeper.cdc.UnmarshalJSON(req.Data, &ti); err != nil {
		return nil, sdk.ErrJSONUnmarshal(fmt.Sprintf("%v", ti))
	}

	symbol := sdk.Symbol(ti.Symbol)
	token := keeper.GetTokenInfo(ctx, symbol)
	if token == nil {
		return nil, types.ErrNonExistSymbol(ti.Symbol)
	}

	bz, err := keeper.cdc.MarshalJSON(types.QueryResToken(*token))
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

//=====
func querySymbols(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	var symbols []string
	symbols = keeper.GetSymbols(ctx)
	bz, err := keeper.cdc.MarshalJSON(symbols)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryTokens(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	tokens := keeper.GetAllTokenInfo(ctx)
	bz, err := keeper.cdc.MarshalJSON(types.QueryResTokens(tokens))
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)
	res, err := codec.MarshalJSONIndent(k.cdc, params)
	if err != nil {
		return nil, sdk.ErrJSONUnmarshal(fmt.Sprintf("%v", params))
	}
	return res, nil
}
