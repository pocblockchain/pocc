package token

import (
	"github.com/pocblockchain/pocc/x/distribution"
	"testing"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/token/types"
	"github.com/stretchr/testify/assert"
)

var TestNewTokenFee = sdk.NewInt(1000000)

func TestHandleMsgNewTokenSuccess(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	tk := input.tokenKeeper
	ak := input.accountKeeper
	dk := input.distrKeeper
	bk := input.bankKeeper
	supplyKeeper := input.supplyKeeper

	params := tk.GetParams(ctx)
	params.NewTokenFee = TestNewTokenFee
	tk.SetParams(ctx, params)

	gotParams := tk.GetParams(ctx)
	assert.Equal(t, TestNewTokenFee, gotParams.NewTokenFee)

	fromAddr, err := sdk.AccAddressFromBech32("poc1kkkjfhv7t4swrurnftedsx0lngvc523cfa00s3")
	assert.Nil(t, err)
	fromAcc := ak.GetOrNewAccount(ctx, fromAddr)
	assert.Equal(t, sdk.Coins(nil), fromAcc.GetCoins())
	fromAcc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, TestNewTokenFee.MulRaw(5))))
	ak.SetAccount(ctx, fromAcc)

	fromAcc = ak.GetAccount(ctx, fromAddr)
	assert.Equal(t, TestNewTokenFee.MulRaw(5), fromAcc.GetCoins().AmountOf(sdk.NativeToken))

	toAddr, err := sdk.AccAddressFromBech32("poc12jwptqcnzkk4d7yupmwlnkzjkm9hvp0rr0chxr")
	assert.Nil(t, err)
	toAcc := ak.GetOrNewAccount(ctx, toAddr)
	assert.Equal(t, sdk.Coins(nil), toAcc.GetCoins())

	openFee := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, TestNewTokenFee))
	totalAmt, ok := sdk.NewIntFromString("10000000000000000000000")
	assert.True(t, ok)

	msg := types.NewMsgNewToken(fromAddr, toAddr, "bhd", 18, totalAmt)

	res := handleMsgNewToken(ctx, tk, msg)
	assert.Equal(t, sdk.CodeOK, res.Code)
	assert.Equal(t, 1, len(res.Events))
	assert.Equal(t, types.EventTypeNewToken, res.Events[0].Type)

	//check tokenInfo
	ti := tk.GetTokenInfo(ctx, "bhd")
	assert.Equal(t, "bhd", ti.Symbol.String())
	assert.Equal(t, uint64(18), ti.Decimals)
	assert.Equal(t, fromAddr.String(), ti.Issuer)
	assert.True(t, ti.IsSendEnabled)

	//check coins
	toAcc = ak.GetAccount(ctx, toAddr)
	assert.Equal(t, totalAmt, toAcc.GetCoins().AmountOf("bhd"))
	fromAcc = ak.GetAccount(ctx, fromAddr)
	assert.Equal(t, TestNewTokenFee.MulRaw(4), fromAcc.GetCoins().AmountOf(sdk.NativeToken))

	//check supply
	supply := supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("bhd")
	assert.Equal(t, totalAmt, supply)
	assert.Equal(t, totalAmt, tk.GetTotalSupply(ctx, "bhd"))

	//check feepool and distribution Account
	feePool := dk.GetFeePool(ctx)
	assert.Equal(t, openFee.AmountOf(sdk.NativeToken), feePool.CommunityPool.AmountOf(sdk.NativeToken).TruncateInt())
	assert.Equal(t, openFee.AmountOf(sdk.NativeToken), supplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins().AmountOf(sdk.NativeToken))

	//sendcoins back to fromCUAddr
	sendAmt := sdk.NewInt(2000000000)
	err = bk.SendCoins(ctx, toAddr, fromAddr, sdk.NewCoins(sdk.NewCoin("bhd", sendAmt)))
	toAcc = ak.GetAccount(ctx, toAddr)
	assert.Equal(t, totalAmt.Sub(sendAmt), toAcc.GetCoins().AmountOf("bhd"))

	fromAcc = ak.GetAccount(ctx, fromAddr)
	assert.Equal(t, TestNewTokenFee.MulRaw(4), fromAcc.GetCoins().AmountOf(sdk.NativeToken))
	assert.Equal(t, sendAmt, fromAcc.GetCoins().AmountOf("bhd"))

}

func TestHandleMsgNewTokenFail(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	tk := input.tokenKeeper
	ak := input.accountKeeper
	supplyKeeper := input.supplyKeeper

	params := tk.GetParams(ctx)
	params.NewTokenFee = TestNewTokenFee
	tk.SetParams(ctx, params)

	fromAddr, err := sdk.AccAddressFromBech32("poc1kkkjfhv7t4swrurnftedsx0lngvc523cfa00s3")
	assert.Nil(t, err)
	fromAcc := ak.GetOrNewAccount(ctx, fromAddr)
	assert.Equal(t, sdk.Coins(nil), fromAcc.GetCoins())
	fromAcc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, TestNewTokenFee.MulRaw(5))))
	ak.SetAccount(ctx, fromAcc)
	fromAcc = ak.GetAccount(ctx, fromAddr)
	assert.Equal(t, TestNewTokenFee.MulRaw(5), fromAcc.GetCoins().AmountOf(sdk.NativeToken))

	toAddr, err := sdk.AccAddressFromBech32("poc12jwptqcnzkk4d7yupmwlnkzjkm9hvp0rr0chxr")
	assert.Nil(t, err)
	toAcc := ak.GetOrNewAccount(ctx, toAddr)
	assert.Equal(t, sdk.Coins(nil), toAcc.GetCoins())

	ti := tk.GetTokenInfo(ctx, sdk.NativeToken)
	assert.NotNil(t, ti)

	btcTokenInfo := sdk.TokenInfo{
		Symbol:        "btc",
		Issuer:        "",
		Decimals:      8,
		IsSendEnabled: true,
		TotalSupply:   sdk.NewInt(2100000000000000),
	}
	tk.SetTokenInfo(ctx, &btcTokenInfo)

	got := tk.GetTokenInfo(ctx, "btc")
	assert.NotNil(t, got)
	assert.Equal(t, btcTokenInfo, *got)
	assert.Equal(t, sdk.NewInt(2100000000000000), supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("btc"))

	got = tk.GetTokenInfo(ctx, "eth")
	assert.Nil(t, got)

	//token already exist
	msg := types.NewMsgNewToken(fromAddr, toAddr, "btc", 8, sdk.NewInt(1000000))
	res := handleMsgNewToken(ctx, tk, msg)
	assert.Equal(t, sdk.CodeSymbolAlreadyExist, res.Code)

	//token is a reserved symbol
	msg = types.NewMsgNewToken(fromAddr, toAddr, "eos", 18, sdk.NewInt(1000000))
	res = handleMsgNewToken(ctx, tk, msg)
	assert.Equal(t, types.CodeSymbolReserved, res.Code)

	msg = types.NewMsgNewToken(fromAddr, toAddr, "bsv", 18, sdk.NewInt(1000000))
	res = handleMsgNewToken(ctx, tk, msg)
	assert.Equal(t, types.CodeSymbolReserved, res.Code)

	msg = types.NewMsgNewToken(fromAddr, toAddr, "bch", 18, sdk.NewInt(1000000))
	res = handleMsgNewToken(ctx, tk, msg)
	assert.Equal(t, types.CodeSymbolReserved, res.Code)

	//fromAccount does not exist
	nonExistAddr, err := sdk.AccAddressFromBech32("poc1fk7g27wg5aznua285jt2kplmfr6rtv0sxn42gh")
	msg = types.NewMsgNewToken(nonExistAddr, toAddr, "bhd", 18, sdk.NewInt(1000000))
	res = handleMsgNewToken(ctx, tk, msg)
	assert.Equal(t, sdk.CodeInsufficientCoins, res.Code)

	//insufficient openFee
	param := tk.GetParams(ctx)
	param.NewTokenFee = TestNewTokenFee.MulRaw(6)
	tk.SetParams(ctx, param)
	msg = types.NewMsgNewToken(fromAddr, toAddr, "bhd", 18, sdk.NewInt(1000000))
	res = handleMsgNewToken(ctx, tk, msg)
	assert.Equal(t, sdk.CodeInsufficientCoins, res.Code)

}

func TestHandleMsgInflateTokenSuccess(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	tk := input.tokenKeeper
	ak := input.accountKeeper
	dk := input.distrKeeper
	bk := input.bankKeeper
	supplyKeeper := input.supplyKeeper

	params := tk.GetParams(ctx)
	params.NewTokenFee = TestNewTokenFee
	tk.SetParams(ctx, params)

	gotParams := tk.GetParams(ctx)
	assert.Equal(t, TestNewTokenFee, gotParams.NewTokenFee)

	fromAddr, err := sdk.AccAddressFromBech32("poc1kkkjfhv7t4swrurnftedsx0lngvc523cfa00s3")
	assert.Nil(t, err)
	fromAcc := ak.GetOrNewAccount(ctx, fromAddr)
	assert.Equal(t, sdk.Coins(nil), fromAcc.GetCoins())
	fromAcc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, TestNewTokenFee.MulRaw(5))))
	ak.SetAccount(ctx, fromAcc)

	fromAcc = ak.GetAccount(ctx, fromAddr)
	assert.Equal(t, TestNewTokenFee.MulRaw(5), fromAcc.GetCoins().AmountOf(sdk.NativeToken))

	toAddr, err := sdk.AccAddressFromBech32("poc12jwptqcnzkk4d7yupmwlnkzjkm9hvp0rr0chxr")
	assert.Nil(t, err)
	toAcc := ak.GetOrNewAccount(ctx, toAddr)
	assert.Equal(t, sdk.Coins(nil), toAcc.GetCoins())

	openFee := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, TestNewTokenFee))
	totalAmt, ok := sdk.NewIntFromString("10000000000000000000000")
	assert.True(t, ok)

	msg := types.NewMsgNewToken(fromAddr, toAddr, "bhd", 18, totalAmt)

	res := handleMsgNewToken(ctx, tk, msg)
	assert.Equal(t, sdk.CodeOK, res.Code)
	assert.Equal(t, 1, len(res.Events))
	assert.Equal(t, types.EventTypeNewToken, res.Events[0].Type)

	//check tokenInfo
	ti := tk.GetTokenInfo(ctx, "bhd")
	assert.Equal(t, "bhd", ti.Symbol.String())
	assert.Equal(t, uint64(18), ti.Decimals)
	assert.Equal(t, fromAddr.String(), ti.Issuer)
	assert.True(t, ti.IsSendEnabled)

	//check coins
	toAcc = ak.GetAccount(ctx, toAddr)
	assert.Equal(t, totalAmt, toAcc.GetCoins().AmountOf("bhd"))
	fromAcc = ak.GetAccount(ctx, fromAddr)
	assert.Equal(t, TestNewTokenFee.MulRaw(4), fromAcc.GetCoins().AmountOf(sdk.NativeToken))

	//check supply
	supply := supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("bhd")
	assert.Equal(t, totalAmt, supply)
	assert.Equal(t, totalAmt, tk.GetTotalSupply(ctx, "bhd"))

	//check feepool and distribution Account
	feePool := dk.GetFeePool(ctx)
	assert.Equal(t, openFee.AmountOf(sdk.NativeToken), feePool.CommunityPool.AmountOf(sdk.NativeToken).TruncateInt())
	assert.Equal(t, openFee.AmountOf(sdk.NativeToken), supplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins().AmountOf(sdk.NativeToken))

	//sendcoins back to fromCUAddr
	sendAmt := sdk.NewInt(2000000000)
	err = bk.SendCoins(ctx, toAddr, fromAddr, sdk.NewCoins(sdk.NewCoin("bhd", sendAmt)))
	toAcc = ak.GetAccount(ctx, toAddr)
	assert.Equal(t, totalAmt.Sub(sendAmt), toAcc.GetCoins().AmountOf("bhd"))

	fromAcc = ak.GetAccount(ctx, fromAddr)
	assert.Equal(t, TestNewTokenFee.MulRaw(4), fromAcc.GetCoins().AmountOf(sdk.NativeToken))
	assert.Equal(t, sendAmt, fromAcc.GetCoins().AmountOf("bhd"))

	//inflate bhd
	inflateAmt := sdk.NewInt(123456789)
	inflateCoins := sdk.NewCoins(sdk.NewCoin("bhd", inflateAmt))
	infateMsg := types.NewMsgInflateToken(fromAddr, toAddr, inflateCoins)

	res = handleMsgInflateToken(ctx, tk, infateMsg)
	assert.Equal(t, sdk.CodeOK, res.Code)
	assert.Equal(t, 1, len(res.Events))
	assert.Equal(t, types.EventTypeInflateToken, res.Events[0].Type)

	ti = tk.GetTokenInfo(ctx, "bhd")
	assert.Equal(t, "bhd", ti.Symbol.String())
	assert.Equal(t, uint64(18), ti.Decimals)
	assert.Equal(t, fromAddr.String(), ti.Issuer)
	assert.True(t, ti.IsSendEnabled)

	//check coins
	totalAmt = totalAmt.Add(inflateAmt)
	toAcc = ak.GetAccount(ctx, toAddr)
	assert.Equal(t, totalAmt.Sub(sendAmt), toAcc.GetCoins().AmountOf("bhd"))
	fromAcc = ak.GetAccount(ctx, fromAddr)
	assert.Equal(t, TestNewTokenFee.MulRaw(4), fromAcc.GetCoins().AmountOf(sdk.NativeToken))

	//check supply
	supply = supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("bhd")
	assert.Equal(t, totalAmt, supply)
	assert.Equal(t, totalAmt, tk.GetTotalSupply(ctx, "bhd"))

	//check feepool and distribution Account
	feePool = dk.GetFeePool(ctx)
	assert.Equal(t, openFee.AmountOf(sdk.NativeToken).MulRaw(1), feePool.CommunityPool.AmountOf(sdk.NativeToken).TruncateInt())
	assert.Equal(t, openFee.AmountOf(sdk.NativeToken).MulRaw(1), supplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins().AmountOf(sdk.NativeToken))

}

func TestHandleMsgInflateTokenFail(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	tk := input.tokenKeeper
	ak := input.accountKeeper
	dk := input.distrKeeper
	bk := input.bankKeeper
	supplyKeeper := input.supplyKeeper

	params := tk.GetParams(ctx)
	params.NewTokenFee = TestNewTokenFee
	tk.SetParams(ctx, params)

	gotParams := tk.GetParams(ctx)
	assert.Equal(t, TestNewTokenFee, gotParams.NewTokenFee)

	fromAddr, err := sdk.AccAddressFromBech32("poc1kkkjfhv7t4swrurnftedsx0lngvc523cfa00s3")
	assert.Nil(t, err)
	fromAcc := ak.GetOrNewAccount(ctx, fromAddr)
	assert.Equal(t, sdk.Coins(nil), fromAcc.GetCoins())
	fromAcc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, TestNewTokenFee.MulRaw(5))))
	ak.SetAccount(ctx, fromAcc)

	fromAcc = ak.GetAccount(ctx, fromAddr)
	assert.Equal(t, TestNewTokenFee.MulRaw(5), fromAcc.GetCoins().AmountOf(sdk.NativeToken))

	toAddr, err := sdk.AccAddressFromBech32("poc12jwptqcnzkk4d7yupmwlnkzjkm9hvp0rr0chxr")
	assert.Nil(t, err)
	toAcc := ak.GetOrNewAccount(ctx, toAddr)
	assert.Equal(t, sdk.Coins(nil), toAcc.GetCoins())

	openFee := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, TestNewTokenFee))
	totalAmt, ok := sdk.NewIntFromString("10000000000000000000000")
	assert.True(t, ok)

	msg := types.NewMsgNewToken(fromAddr, toAddr, "bhd", 18, totalAmt)

	res := handleMsgNewToken(ctx, tk, msg)
	assert.Equal(t, sdk.CodeOK, res.Code)
	assert.Equal(t, 1, len(res.Events))
	assert.Equal(t, types.EventTypeNewToken, res.Events[0].Type)

	//check tokenInfo
	ti := tk.GetTokenInfo(ctx, "bhd")
	assert.Equal(t, "bhd", ti.Symbol.String())
	assert.Equal(t, uint64(18), ti.Decimals)
	assert.Equal(t, fromAddr.String(), ti.Issuer)
	assert.True(t, ti.IsSendEnabled)

	//check coins
	toAcc = ak.GetAccount(ctx, toAddr)
	assert.Equal(t, totalAmt, toAcc.GetCoins().AmountOf("bhd"))
	fromAcc = ak.GetAccount(ctx, fromAddr)
	assert.Equal(t, TestNewTokenFee.MulRaw(4), fromAcc.GetCoins().AmountOf(sdk.NativeToken))

	//check supply
	supply := supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("bhd")
	assert.Equal(t, totalAmt, supply)
	assert.Equal(t, totalAmt, tk.GetTotalSupply(ctx, "bhd"))

	//check feepool and distribution Account
	feePool := dk.GetFeePool(ctx)
	assert.Equal(t, openFee.AmountOf(sdk.NativeToken), feePool.CommunityPool.AmountOf(sdk.NativeToken).TruncateInt())
	assert.Equal(t, openFee.AmountOf(sdk.NativeToken), supplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins().AmountOf(sdk.NativeToken))

	//sendcoins back to fromCUAddr
	sendAmt := sdk.NewInt(2000000000)
	err = bk.SendCoins(ctx, toAddr, fromAddr, sdk.NewCoins(sdk.NewCoin("bhd", sendAmt)))
	toAcc = ak.GetAccount(ctx, toAddr)
	assert.Equal(t, totalAmt.Sub(sendAmt), toAcc.GetCoins().AmountOf("bhd"))

	fromAcc = ak.GetAccount(ctx, fromAddr)
	assert.Equal(t, TestNewTokenFee.MulRaw(4), fromAcc.GetCoins().AmountOf(sdk.NativeToken))
	assert.Equal(t, sendAmt, fromAcc.GetCoins().AmountOf("bhd"))

	//fail because sendenable = false
	tk.DisableSend(ctx, "bhd")
	inflateAmt := sdk.NewInt(123456789)
	inflateCoins := sdk.NewCoins(sdk.NewCoin("bhd", inflateAmt))
	infateMsg := types.NewMsgInflateToken(fromAddr, toAddr, inflateCoins)
	res = handleMsgInflateToken(ctx, tk, infateMsg)
	assert.Equal(t, sdk.CodeTransactionIsNotEnabled, res.Code)

	//Obsolete, more than 1 coins is checked in MsgInflateToken's ValidateBasic()
	//fail because  more than 1 coins,

	//infateMsg = types.NewMsgInflateToken(fromAddr, toAddr, inflateCoins.Add(sdk.NewCoins(sdk.NewCoin("bht", sdk.NewInt(100)))))
	//res = handleMsgInflateToken(ctx, tk, infateMsg)
	//assert.Equal(t, sdk.CodeInvalidTx, res.Code)

	//fail because no authority
	tk.EnableSend(ctx, "bhd")
	infateMsg = types.NewMsgInflateToken(toAddr, toAddr, inflateCoins)
	res = handleMsgInflateToken(ctx, tk, infateMsg)
	assert.Equal(t, sdk.CodeInvalidAccount, res.Code)

	//insufficient openFee
	//param := tk.GetParams(ctx)
	//param.NewTokenFee = TestNewTokenFee.MulRaw(5)
	//tk.SetParams(ctx, param)
	//infateMsg = types.NewMsgInflateToken(fromAddr, toAddr, inflateCoins)
	//res = handleMsgInflateToken(ctx, tk, infateMsg)
	//assert.Equal(t, sdk.CodeInsufficientCoins, res.Code)

	ti = tk.GetTokenInfo(ctx, "bhd")
	assert.Equal(t, "bhd", ti.Symbol.String())
	assert.Equal(t, uint64(18), ti.Decimals)
	assert.Equal(t, fromAddr.String(), ti.Issuer)
	assert.True(t, ti.IsSendEnabled)

	//check coins
	toAcc = ak.GetAccount(ctx, toAddr)
	assert.Equal(t, totalAmt.Sub(sendAmt), toAcc.GetCoins().AmountOf("bhd"))
	fromAcc = ak.GetAccount(ctx, fromAddr)
	assert.Equal(t, TestNewTokenFee.MulRaw(4), fromAcc.GetCoins().AmountOf(sdk.NativeToken))

	//check supply
	supply = supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("bhd")
	assert.Equal(t, totalAmt, supply)
	assert.Equal(t, totalAmt, tk.GetTotalSupply(ctx, "bhd"))

	//check feepool and distribution Account
	feePool = dk.GetFeePool(ctx)
	assert.Equal(t, openFee.AmountOf(sdk.NativeToken), feePool.CommunityPool.AmountOf(sdk.NativeToken).TruncateInt())
	assert.Equal(t, openFee.AmountOf(sdk.NativeToken), supplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins().AmountOf(sdk.NativeToken))

}

func TestHandleMsBurnTokenSuccess(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	tk := input.tokenKeeper
	ak := input.accountKeeper
	supplyKeeper := input.supplyKeeper

	params := tk.GetParams(ctx)
	params.NewTokenFee = TestNewTokenFee
	tk.SetParams(ctx, params)
	gotParams := tk.GetParams(ctx)
	assert.Equal(t, TestNewTokenFee, gotParams.NewTokenFee)

	fromAddr, err := sdk.AccAddressFromBech32("poc1kkkjfhv7t4swrurnftedsx0lngvc523cfa00s3")
	assert.Nil(t, err)
	fromAcc := ak.GetOrNewAccount(ctx, fromAddr)
	assert.Equal(t, sdk.Coins(nil), fromAcc.GetCoins())
	fromAcc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, TestNewTokenFee.MulRaw(5))))
	ak.SetAccount(ctx, fromAcc)

	toAddr := fromAddr
	totalSupply, ok := sdk.NewIntFromString("10000000000000000000000")
	assert.True(t, ok)

	newMsg := types.NewMsgNewToken(fromAddr, toAddr, "bhd", 18, totalSupply)
	res := handleMsgNewToken(ctx, tk, newMsg)
	assert.Equal(t, sdk.CodeOK, res.Code)
	assert.Equal(t, 1, len(res.Events))
	assert.Equal(t, types.EventTypeNewToken, res.Events[0].Type)

	burnAmt := totalSupply.QuoRaw(3)
	burnCoins := sdk.NewCoins(sdk.NewCoin("bhd", burnAmt))

	//burn token
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	burnMsg := types.NewMsgBurnToken(fromAddr, burnCoins)

	res = handleMsgBurnToken(ctx, tk, burnMsg)
	assert.Equal(t, sdk.CodeOK, res.Code)
	assert.Equal(t, 1, len(res.Events))
	assert.Equal(t, types.EventTypeBurnToken, res.Events[0].Type)

	//check  coins
	toAcc := ak.GetAccount(ctx, toAddr)
	assert.Equal(t, totalSupply.Sub(burnAmt), toAcc.GetCoins().AmountOf("bhd"))

	//check supply
	supply := supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("bhd")
	assert.Equal(t, totalSupply.Sub(burnAmt), supply)
	assert.Equal(t, totalSupply.Sub(burnAmt), tk.GetTotalSupply(ctx, "bhd"))

}

func TestHandleMsBurnTokenFail(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	tk := input.tokenKeeper
	ak := input.accountKeeper
	supplyKeeper := input.supplyKeeper

	params := tk.GetParams(ctx)
	params.NewTokenFee = TestNewTokenFee
	tk.SetParams(ctx, params)

	gotParams := tk.GetParams(ctx)
	assert.Equal(t, TestNewTokenFee, gotParams.NewTokenFee)

	fromAddr, err := sdk.AccAddressFromBech32("poc1kkkjfhv7t4swrurnftedsx0lngvc523cfa00s3")
	assert.Nil(t, err)
	fromAcc := ak.GetOrNewAccount(ctx, fromAddr)
	assert.Equal(t, sdk.Coins(nil), fromAcc.GetCoins())
	fromAcc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, TestNewTokenFee.MulRaw(5))))
	ak.SetAccount(ctx, fromAcc)

	toAddr := fromAddr
	totalSupply, ok := sdk.NewIntFromString("10000000000000000000000")
	assert.True(t, ok)

	newMsg := types.NewMsgNewToken(fromAddr, toAddr, "bhd", 18, totalSupply)
	res := handleMsgNewToken(ctx, tk, newMsg)
	assert.Equal(t, sdk.CodeOK, res.Code)
	assert.Equal(t, 1, len(res.Events))
	assert.Equal(t, types.EventTypeNewToken, res.Events[0].Type)

	burnAmt := totalSupply.AddRaw(1)
	burnCoins := sdk.NewCoins(sdk.NewCoin("bhd", burnAmt))

	//burn more than account have
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	msg := types.NewMsgBurnToken(fromAddr, burnCoins)
	res = handleMsgBurnToken(ctx, tk, msg)
	assert.Equal(t, sdk.CodeInsufficientCoins, res.Code)
	assert.Equal(t, 0, len(res.Events))

	//check supply
	supply := supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("bhd")
	assert.Equal(t, supply, totalSupply)

	//token sendenable is false
	tk.DisableSend(ctx, "bhd")
	burnAmt = totalSupply.QuoRaw(3)
	burnCoins = sdk.NewCoins(sdk.NewCoin("bhd", burnAmt))

	msg = types.NewMsgBurnToken(fromAddr, burnCoins)
	res = handleMsgBurnToken(ctx, tk, msg)
	assert.Equal(t, sdk.CodeTransactionIsNotEnabled, res.Code)

	supply = supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("bhd")
	assert.Equal(t, totalSupply, supply)
	assert.Equal(t, totalSupply, tk.GetTotalSupply(ctx, "bhd"))
}
