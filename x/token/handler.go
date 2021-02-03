package token

import (
	"fmt"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/token/types"
)

//NewHandler create a token message handler
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgNewToken:
			return handleMsgNewToken(ctx, keeper, msg)

		case types.MsgInflateToken:
			return handleMsgInflateToken(ctx, keeper, msg)

		case types.MsgBurnToken:
			return handleMsgBurnToken(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("Unrecognized token Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgNewToken(ctx sdk.Context, keeper Keeper, msg types.MsgNewToken) sdk.Result {
	ctx.Logger().Info("handleMsgNewToken", "msg", msg)

	if ti := keeper.GetTokenInfoWithoutSupply(ctx, msg.Symbol); ti != nil {
		return sdk.ErrAlreadyExitSymbol(fmt.Sprintf("token %s already exist", msg.Symbol)).Result()
	}
	reserved := keeper.GetParams(ctx).ReservedSymbols
	for _, r := range reserved {
		if r == msg.Symbol.String() {
			return types.ErrSymbolReserved(fmt.Sprintf("%v is reserved", msg.Symbol)).Result()
		}
	}

	issueFee := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, keeper.GetParams(ctx).NewTokenFee))

	//transfer openFee to communityPool
	err := keeper.dk.AddCoinsFromAccountToFeePool(ctx, msg.From, issueFee)
	if err != nil {
		return err.Result()
	}

	//Set TokenInfo
	totalSupply := msg.TotalSupply
	keeper.SetTokenInfoWithoutSupply(ctx, &sdk.TokenInfoWithoutSupply{
		Symbol:        msg.Symbol,
		Issuer:        msg.From.String(),
		IsSendEnabled: true,
		Decimals:      msg.Decimals,
	})

	//minted newCoins
	mintedCoins := sdk.NewCoins(sdk.NewCoin(msg.Symbol.String(), msg.TotalSupply))
	err = keeper.sk.MintCoins(ctx, types.ModuleName, mintedCoins)
	if err != nil {
		return err.Result()
	}

	err = keeper.sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, msg.To, mintedCoins)
	if err != nil {
		return err.Result()
	}

	//ignore events in SendCoinsFromAccountToModule and SendCoinsFromModuleToAccount
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNewToken,
			sdk.NewAttribute(types.AttributeKeyIssuer, msg.From.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.To.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, totalSupply.String()),
			sdk.NewAttribute(types.AttributeKeyIssueFee, issueFee.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgInflateToken(ctx sdk.Context, keeper Keeper, msg types.MsgInflateToken) sdk.Result {
	ctx.Logger().Info("handleMsgInflateToken", "msg", msg)


	symbol := sdk.Symbol(msg.Amount[0].Denom)
	if !keeper.IsSendEnabled(ctx, symbol) {
		return sdk.ErrTransactionIsNotEnabled(fmt.Sprintf("%v is not sendenable", symbol)).Result()
	}

	//Only token owner can inflate token
	if msg.From.String() != keeper.GetIssuer(ctx, symbol) {
		return sdk.ErrInvalidAccount(fmt.Sprintf("%s is not allowed to inflate %v", msg.From, symbol)).Result()
	}

	//transfer openFee to communityPool
	//since issuer have pay when new token, do not charge again when inflate token
	//issueFee := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, keeper.GetParams(ctx).NewTokenFee))
	//err := keeper.dk.AddCoinsFromAccountToFeePool(ctx, msg.From, issueFee)
	//if err != nil {
	//	return err.Result()
	//}

	//minted inflatedCoins
	inflatedCoins := sdk.NewCoins(msg.Amount[0])
	err := keeper.sk.MintCoins(ctx, types.ModuleName, inflatedCoins)
	if err != nil {
		return err.Result()
	}

	err = keeper.sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, msg.To, inflatedCoins)
	if err != nil {
		return err.Result()
	}

	//ignore events in SendCoinsFromAccountToModule and SendCoinsFromModuleToAccount
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeInflateToken,
			sdk.NewAttribute(types.AttributeKeyIssuer, msg.From.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.To.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount[0].Amount.String()),
	//		sdk.NewAttribute(types.AttributeKeyIssueFee, issueFee.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBurnToken(ctx sdk.Context, keeper Keeper, msg types.MsgBurnToken) sdk.Result {
	ctx.Logger().Info("handleMsgBurnToken", "msg", msg)

	for _, coin := range msg.Amount {
		if !keeper.IsSendEnabled(ctx, sdk.Symbol(coin.Denom)) {
			return sdk.ErrTransactionIsNotEnabled(fmt.Sprintf("%v is not sendenable", coin.Denom)).Result()
		}
	}

	//send the burned coins to token module
	err := keeper.sk.SendCoinsFromAccountToModule(ctx, msg.From, types.ModuleName, msg.Amount)
	if err != nil {
		return err.Result()
	}

	//burn from token module
	err = keeper.sk.BurnCoins(ctx, types.ModuleName, msg.Amount)
	if err != nil {
		return err.Result()
	}

	//ignore events in SendCoinsFromAccountToModule and SendCoinsFromModuleToAccount
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBurnToken,
			sdk.NewAttribute(types.AttributeKeyBurner, msg.From.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
