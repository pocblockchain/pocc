package token

import (
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/token/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func changeProposal(symbol string, changes []types.ParamChange) types.TokenParamsChangeProposal {
	return types.TokenParamsChangeProposal{
		Title:       "Test",
		Description: "description",
		Symbol:      symbol,
		Changes:     changes,
	}
}

func disableProposal(symbol string) types.DisableTokenProposal {
	return types.DisableTokenProposal{
		Title:       "Test",
		Description: "description",
		Symbol:      symbol,
	}
}

func TestTokenParamsChangeProposalPassed(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	keeper := input.tokenKeeper

	keeper.SetTokenInfo(ctx, &sdk.TokenInfo{
		Symbol:        "btc",
		Issuer:        "",
		Decimals:      18,
		TotalSupply:   sdk.ZeroInt(),
		IsSendEnabled: true,
	})

	ctx.WithBlockHeight(10)

	hdlr := NewTokenProposalHandler(keeper)
	changes := []types.ParamChange{}
	changes = append(changes, types.NewParamChange(sdk.KeyIsSendEnabled, "false"))

	cp := changeProposal("btc", changes)
	res := hdlr(ctx, cp)
	require.Equal(t, sdk.CodeOK, res.Code)
	events := ctx.EventManager().Events()
	require.Equal(t, 1, len(events))
	require.Equal(t, types.EventTypeExecuteTokenParamsChangeProposal, events[0].Type)
	require.Equal(t, 2, len(events[0].Attributes))
	require.Equal(t, sdk.KeyIsSendEnabled, string(events[0].Attributes[0].Value))

	ti := keeper.GetTokenInfo(ctx, "btc")
	require.Equal(t, false, ti.IsSendEnabled)

}

func TestTokenParamsChangeProposalFailed(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	keeper := input.tokenKeeper
	ctx.WithBlockHeight(10)

	btcTi := sdk.TokenInfoWithoutSupply{
		Symbol:        "btc",
		Issuer:        "btc",
		IsSendEnabled: true,
		Decimals:      8,
	}
	keeper.SetTokenInfoWithoutSupply(ctx, &btcTi)

	nativeTi := sdk.TokenInfoWithoutSupply{
		Symbol:        sdk.NativeToken,
		IsSendEnabled: true,
		Decimals:      8,
	}

	keeper.SetTokenInfoWithoutSupply(ctx, &nativeTi)

	ctx.WithBlockHeight(10)

	hdlr := NewTokenProposalHandler(keeper)
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	// try to change btc's non-exist collect params
	changes := []types.ParamChange{}
	changes = append(changes, types.NewParamChange("collect", `"21000000000000000"`))
	cp := changeProposal("btc", changes)
	res := hdlr(ctx, cp)
	require.NotEqual(t, sdk.CodeOK, res.Code)
	require.Equal(t, 0, len(res.Events))
	gotTi := keeper.GetTokenInfoWithoutSupply(ctx, "btc")
	require.Equal(t, btcTi, *gotTi)

	//try to change a non-exist token
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	changes = append(changes, types.NewParamChange(sdk.KeyIsSendEnabled, "false"))
	cp = changeProposal("ebt", changes)
	res = hdlr(ctx, cp)
	require.Equal(t, sdk.CodeUnsupportToken, res.Code)
	require.Contains(t, res.Log, "dose not exist")

	//try to change NativeToken
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	changes = append(changes, types.NewParamChange(sdk.KeyIsSendEnabled, "false"))
	cp = changeProposal(sdk.NativeToken, changes)
	res = hdlr(ctx, cp)
	require.NotEqual(t, sdk.CodeOK, res.Code)
	require.Contains(t, res.Log, "Not allowed to change native token's params")
	require.Equal(t, 0, len(res.Events))
	gotTi = keeper.GetTokenInfoWithoutSupply(ctx, sdk.NativeToken)
	require.Equal(t, nativeTi, *gotTi)

}

func TestDisableTokenProposalPassed(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	keeper := input.tokenKeeper
	supplyKeeper := input.supplyKeeper
	ctx.WithBlockHeight(10)

	btcAmt := sdk.NewInt(2100000000000000)
	keeper.SetTokenInfo(ctx, &sdk.TokenInfo{
		Symbol:        "btc",
		Issuer:        "btc",
		IsSendEnabled: true,
		Decimals:      8,
		TotalSupply:   btcAmt,
	})

	require.Equal(t, btcAmt, supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("btc"))
	require.Equal(t, btcAmt, keeper.GetTotalSupply(ctx, "btc"))

	symbols := keeper.GetSymbols(ctx)
	require.Contains(t, symbols, sdk.NativeToken)
	require.Contains(t, symbols, "btc")

	hdlr := NewTokenProposalHandler(keeper)
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	dp := disableProposal("btc")
	res := hdlr(ctx, dp)
	require.Equal(t, sdk.CodeOK, res.Code)
	events := ctx.EventManager().Events()
	require.Equal(t, 1, len(events))
	require.Equal(t, types.EventTypeExecuteDisableTokenProposal, events[0].Type)

	//after disable, token info still exist but isenabled = false
	tokenInfo := keeper.GetTokenInfoWithoutSupply(ctx, "btc")
	require.NotNil(t, tokenInfo)
	require.Equal(t, false, tokenInfo.IsSendEnabled)
	require.Equal(t, btcAmt, supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("btc"))
	require.Equal(t, btcAmt, keeper.GetTotalSupply(ctx, "btc"))

	symbols = keeper.GetSymbols(ctx)
	require.Contains(t, symbols, sdk.NativeToken)
	require.Contains(t, symbols, "btc")

}

func TestDisableTokenProposalFailed(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	keeper := input.tokenKeeper
	ctx.WithBlockHeight(10)

	//try to disable an non-exist token
	hdlr := NewTokenProposalHandler(keeper)
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	dp := disableProposal("ebt")
	res := hdlr(ctx, dp)
	require.Equal(t, sdk.CodeUnsupportToken, res.Code)
	require.Equal(t, 0, len(res.Events))
	require.Contains(t, res.Log, "does not exist")

	//try to disable native token
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	dp = disableProposal(sdk.NativeToken)
	res = hdlr(ctx, dp)
	require.Equal(t, sdk.CodeInvalidTx, res.Code)
	events := ctx.EventManager().Events()
	require.Equal(t, 0, len(events))
	require.Contains(t, res.Log, "Not allowed to disable native token")
}
