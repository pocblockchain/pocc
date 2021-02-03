package bank

import (
	"fmt"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/bank/internal/keeper"
	"github.com/pocblockchain/pocc/x/bank/internal/types"
	"strings"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgSend:
			return handleMsgSend(ctx, k, msg)

		case types.MsgMultiSend:
			return handleMsgMultiSend(ctx, k, msg)

		case types.MsgEscrow:
			return handleMsgEscrow(ctx, k, msg)

		case types.MsgReclaim:
			return handleMsgReclaim(ctx, k, msg)

		case types.MsgBonusSend:
			return handleMsgBonusSend(ctx, k, msg)

		case types.MsgReclaimSend:
			return handleMsgReclaimSend(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized bank message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, k keeper.Keeper, msg types.MsgSend) sdk.Result {
	if !k.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}

	if !k.IsCoinsSendEnabled(ctx, msg.Amount) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}

	if k.BlacklistedAddr(msg.ToAddress) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", msg.ToAddress)).Result()
	}

	err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

// Handle MsgMultiSend.
func handleMsgMultiSend(ctx sdk.Context, k keeper.Keeper, msg types.MsgMultiSend) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked
	if !k.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}

	for _, in := range msg.Inputs {
		if !k.IsCoinsSendEnabled(ctx, in.Coins) {
			return types.ErrSendDisabled(k.Codespace()).Result()
		}
	}

	for _, out := range msg.Outputs {
		if k.BlacklistedAddr(out.Address) {
			return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", out.Address)).Result()
		}

	}

	err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// Handle MsgEscrow.
func handleMsgEscrow(ctx sdk.Context, k keeper.Keeper, msg types.MsgEscrow) sdk.Result {
	if !k.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}

	if !k.IsCoinsSendEnabled(ctx, msg.Amount) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}

	if k.BlacklistedAddr(msg.ToAddress) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", msg.ToAddress)).Result()
	}

	err := k.EscrowCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// Handle MsgReclaim.
func handleMsgReclaim(ctx sdk.Context, k keeper.Keeper, msg types.MsgReclaim) sdk.Result {
	//just record the event,
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReclaim,
			sdk.NewAttribute(types.AttributeKeyReclaimFrom, msg.FromAddress.String()),
			sdk.NewAttribute(types.AttributeKeyReclaimTo, msg.ToAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(types.AttributeKeySender, msg.FromAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

// Handle MsgBonusSend.
func handleMsgBonusSend(ctx sdk.Context, k keeper.Keeper, msg types.MsgBonusSend) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked
	if !k.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}
	var inAddr = make([]string, len(msg.Inputs))
	for i, in := range msg.Inputs {
		if !k.IsCoinsSendEnabled(ctx, in.Coins) {
			return types.ErrSendDisabled(k.Codespace()).Result()
		}
		inAddr[i] = in.Address.String()
	}

	var outAddr = make([]string, len(msg.Outputs))
	for i, out := range msg.Outputs {
		if k.BlacklistedAddr(out.Address) {
			return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", out.Address)).Result()
		}

		outAddr[i] = out.Address.String()
	}

	err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBonusSend,
			sdk.NewAttribute(types.AttributeKeySender, strings.Join(inAddr, ",")),
			sdk.NewAttribute(types.AttributeKeyRecipient, strings.Join(outAddr, ",")),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// Handle MsgReclaimSend.
func handleMsgReclaimSend(ctx sdk.Context, k keeper.Keeper, msg types.MsgReclaimSend) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked
	if !k.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}
	var inAddr = make([]string, len(msg.Inputs))
	for i, in := range msg.Inputs {
		if !k.IsCoinsSendEnabled(ctx, in.Coins) {
			return types.ErrSendDisabled(k.Codespace()).Result()
		}
		inAddr[i] = in.Address.String()
	}

	var outAddr = make([]string, len(msg.Outputs))
	for i, out := range msg.Outputs {
		if k.BlacklistedAddr(out.Address) {
			return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", out.Address)).Result()
		}

		outAddr[i] = out.Address.String()
	}

	err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReclaimSend,
			sdk.NewAttribute(types.AttributeKeySender, strings.Join(inAddr, ",")),
			sdk.NewAttribute(types.AttributeKeyRecipient, strings.Join(outAddr, ",")),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}
