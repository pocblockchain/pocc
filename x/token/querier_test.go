/*
 * *******************************************************************
 * @项目名称: token
 * @文件名称: querier_test.go
 * @Date: 2019/06/24
 * @Author: Keep
 * @Copyright（C）: 2019 BlueHelix Inc.   All rights reserved.
 * 注意：本内容仅限于内部传阅，禁止外泄以及用于其他的商业目的.
 * *******************************************************************
 */

package token

import (
	"fmt"
	"testing"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/token/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQuery(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	keeper := input.tokenKeeper
	cdc := input.cdc

	for _, ti := range TestTokenData {
		keeper.SetTokenInfo(ctx, &ti)
	}

	req := abci.RequestQuery{
		Path: fmt.Sprintf("token/%s/%s", types.RouterKey, types.QueryToken),
		Data: nil,
	}

	for _, ti := range TestTokenData {
		symbol := ti.Symbol

		bz, err := cdc.MarshalJSON(QueryTokenInfo{string(symbol)})
		assert.NoError(t, err)

		req.Data = bz
		bz, err = queryToken(ctx, req, keeper)
		assert.Nil(t, err)

		var res types.QueryResToken
		keeper.cdc.MustUnmarshalJSON(bz, &res)
		assert.Equal(t, sdk.TokenInfo(res), ti)
	}

	bz, err := cdc.MarshalJSON(QueryTokenInfo{"nonexist"})
	assert.NoError(t, err)
	req.Data = bz
	_, err = queryToken(ctx, req, keeper)
	assert.NotNil(t, err)

	//err
	bz, err = querySymbols(ctx, keeper)
	assert.Nil(t, err)
	var symbols []string
	keeper.cdc.MustUnmarshalJSON(bz, &symbols)
	for _, ti := range TestTokenData {
		assert.Contains(t, symbols, ti.Symbol.String())
	}
}

func TestQueryTokens(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	keeper := input.tokenKeeper

	// set token info
	for _, ti := range TestTokenData {
		keeper.SetTokenInfo(ctx, &ti)
	}

	//queryTokens
	bz, err := queryTokens(ctx, keeper)
	assert.Nil(t, err)

	// Unmarshal tokeninfos
	var res []sdk.TokenInfo
	keeper.cdc.MustUnmarshalJSON(bz, &res)

	for _, ti := range TestTokenData {
		for _, r := range res {
			if ti.Symbol == r.Symbol {
				assert.Equal(t, ti, r)
			}
		}
	}
}

func TestMarshalAndString(t *testing.T) {
	input := setupTestEnv(t)
	keeper := input.tokenKeeper

	prtintStr := "\n\tSymbol:\n\tIssuer:\n\tIsSendEnabled:false\n\tDecimals:8\n\tTotalSupply:2100\n\t"

	btcTokenInfo := sdk.TokenInfo{
		Issuer:        "",
		IsSendEnabled: false,

		Decimals:    8,
		TotalSupply: sdk.NewInt(2100),
	}

	bz, err := keeper.cdc.MarshalJSON(btcTokenInfo)
	assert.NoError(t, err)

	var res sdk.TokenInfo
	err = keeper.cdc.UnmarshalJSON(bz, &res)
	assert.NoError(t, err)
	assert.Equal(t, prtintStr, res.String())
}
