package token

import (
	"fmt"
	"github.com/pocblockchain/pocc/codec"
	sdk "github.com/pocblockchain/pocc/types"
	govtypes "github.com/pocblockchain/pocc/x/gov/types"
	"github.com/pocblockchain/pocc/x/token/types"
)

func processChangeParam(key, value string, ti *sdk.TokenInfoWithoutSupply, cdc *codec.Codec) error {
	switch key {
	case sdk.KeyIsSendEnabled:
		val := false
		err := cdc.UnmarshalJSON([]byte(value), &val)
		if err != nil {
			return err
		}
		ti.IsSendEnabled = val

	default:
		return fmt.Errorf("Unkonwn parameter:%v", key)
	}

	return nil
}

func handleTokenParamsChangeProposal(ctx sdk.Context, keeper Keeper, proposal types.TokenParamsChangeProposal) sdk.Result {
	ctx.Logger().Info("handleTokenParamsChangeProposal", "proposal", proposal)

	if proposal.Symbol == sdk.NativeToken {
		return sdk.ErrInvalidTx("Not allowed to change native token's params").Result()
	}

	ti := keeper.GetTokenInfoWithoutSupply(ctx, sdk.Symbol(proposal.Symbol))
	if ti == nil {
		return sdk.ErrUnSupportToken(fmt.Sprintf("token %s dose not exist", proposal.Symbol)).Result()
	}

	attr := []sdk.Attribute{}
	for _, pc := range proposal.Changes {
		err := processChangeParam(pc.Key, pc.Value, ti, keeper.cdc)
		if err != nil {
			return types.ErrInvalidParameter(pc.Key, pc.Value).Result()
		}
		attr = append(attr, sdk.NewAttribute(types.AttributeKeyTokenParam, pc.Key), sdk.NewAttribute(types.AttributeKeyTokenParamValue, pc.Value))
	}

	keeper.SetTokenInfoWithoutSupply(ctx, ti)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeExecuteTokenParamsChangeProposal, attr...),
	)
	return sdk.Result{}
}

func handleDisableTokenProposal(ctx sdk.Context, keeper Keeper, proposal types.DisableTokenProposal) sdk.Result {
	ctx.Logger().Info("handleDisableTokenProposal", "proposal", proposal)

	if proposal.Symbol == sdk.NativeToken {
		return sdk.ErrInvalidTx("Not allowed to disable native token").Result()
	}

	ti := keeper.GetTokenInfoWithoutSupply(ctx, sdk.Symbol(proposal.Symbol))
	if ti == nil {
		return sdk.ErrUnSupportToken(fmt.Sprintf("token %s does not exist", proposal.Symbol)).Result()
	}

	ti.IsSendEnabled = false
	keeper.SetTokenInfoWithoutSupply(ctx, ti)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeExecuteDisableTokenProposal,
			sdk.NewAttribute(types.AttributeKeyToken, proposal.Symbol),
		),
	)
	return sdk.Result{}
}

//NewTokenProposalHandler create handler for token's proposal
func NewTokenProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) sdk.Result {
		switch c := content.(type) {
		case types.TokenParamsChangeProposal:
			return handleTokenParamsChangeProposal(ctx, k, c)

		case types.DisableTokenProposal:
			return handleDisableTokenProposal(ctx, k, c)

		default:
			errMsg := fmt.Sprintf("unrecognized token proposal content type: %T", c)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
