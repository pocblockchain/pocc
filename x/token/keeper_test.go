/*
 * *******************************************************************
 * @项目名称: token
 * @文件名称: Keeper_test.go
 * @Date: 2019/06/06
 * @Author: Keep
 * @Copyright（C）: 2019 BlueHelix Inc.   All rights reserved.
 * 注意：本内容仅限于内部传阅，禁止外泄以及用于其他的商业目的.
 * *******************************************************************
 */

package token

import (
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/supply"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestSetTokenInfoWithoutSupply(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	keeper := input.tokenKeeper

	for _, ti := range TestTokenData {
		tsi := castToTokenInfoWithoutSupply(ti)
		keeper.SetTokenInfoWithoutSupply(ctx, &tsi)

		oldCoins := keeper.sk.GetSupply(ctx).GetTotal()
		oldAmt := sdk.NewCoins(sdk.NewCoin(ti.Symbol.String(), oldCoins.AmountOf(ti.Symbol.String())))
		oldCoins = oldCoins.Sub(oldAmt)
		coins := sdk.NewCoins(sdk.NewCoin(ti.Symbol.String(), ti.TotalSupply))
		newCoins := oldCoins.Add(coins)
		keeper.sk.SetSupply(ctx, supply.NewSupply(newCoins))
	}

	for _, ti := range TestTokenData {
		symbol := ti.Symbol
		assert.Equal(t, ti.Issuer, keeper.GetIssuer(ctx, symbol))
		assert.Equal(t, ti.IsSendEnabled, keeper.IsSendEnabled(ctx, symbol))
		assert.Equal(t, ti.Decimals, keeper.GetDecimals(ctx, symbol))
		assert.Equal(t, ti.TotalSupply, keeper.GetTotalSupply(ctx, symbol))
	}

	symbols := keeper.GetSymbols(ctx)
	assert.Contains(t, symbols, BtcToken)
	assert.Contains(t, symbols, EthToken)
	assert.Contains(t, symbols, UsdtToken)
}

func TestSetTokenInfo(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	keeper := input.tokenKeeper

	for _, ti := range TestTokenData {
		keeper.SetTokenInfo(ctx, &ti)
	}

	for _, ti := range TestTokenData {
		symbol := ti.Symbol
		assert.Equal(t, ti.Issuer, keeper.GetIssuer(ctx, symbol))
		assert.Equal(t, ti.IsSendEnabled, keeper.IsSendEnabled(ctx, symbol))
		assert.Equal(t, ti.Decimals, keeper.GetDecimals(ctx, symbol))
		assert.Equal(t, ti.TotalSupply, keeper.GetTotalSupply(ctx, symbol))
	}

	symbols := keeper.GetSymbols(ctx)
	assert.Contains(t, symbols, BtcToken)
	assert.Contains(t, symbols, EthToken)
	assert.Contains(t, symbols, UsdtToken)
}

func TestModifySend(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	keeper := input.tokenKeeper

	ti := TestTokenData[1]
	symbol := ti.Symbol

	tsi := castToTokenInfoWithoutSupply(ti)
	keeper.SetTokenInfoWithoutSupply(ctx, &tsi)
	assert.True(t, keeper.IsSendEnabled(ctx, symbol))

	keeper.EnableSend(ctx, symbol)
	assert.True(t, keeper.IsSendEnabled(ctx, symbol))

	keeper.DisableSend(ctx, symbol)
	assert.False(t, keeper.IsSendEnabled(ctx, symbol))
}

func TestModifySend1(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	keeper := input.tokenKeeper

	ti := TestTokenData[1]
	symbol := ti.Symbol

	keeper.SetTokenInfo(ctx, &ti)
	assert.True(t, keeper.IsSendEnabled(ctx, symbol))

	keeper.EnableSend(ctx, symbol)
	assert.True(t, keeper.IsSendEnabled(ctx, symbol))

	keeper.DisableSend(ctx, symbol)
	assert.False(t, keeper.IsSendEnabled(ctx, symbol))
}

func TestTokenInfoEncoding(t *testing.T) {
	input := setupTestEnv(t)
	keeper := input.tokenKeeper

	for _, ti := range TestTokenData {
		tsi := castToTokenInfoWithoutSupply(ti)

		bz := keeper.cdc.MustMarshalJSON(tsi)
		var got sdk.TokenInfoWithoutSupply
		keeper.cdc.MustUnmarshalJSON(bz, &got)
		//	t.Logf("tsi:%v, got:%v", tsi, got)
		assert.Equal(t, tsi, got)
	}
}
