package params

import (
	"fmt"
	"github.com/pocblockchain/pocc/x/params/types"

	sdk "github.com/pocblockchain/pocc/types"
	govtypes "github.com/pocblockchain/pocc/x/gov/types"
)

func NewParamChangeProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) sdk.Result {
		switch c := content.(type) {
		case ParameterChangeProposal:
			return handleParameterChangeProposal(ctx, k, c)

		default:
			errMsg := fmt.Sprintf("unrecognized param proposal content type: %T", c)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleParameterChangeProposal(ctx sdk.Context, k Keeper, p ParameterChangeProposal) sdk.Result {
	ctx.Logger().Info("handleParameterChangeProposal", "proposal", p)

	attr := []sdk.Attribute{}
	for _, c := range p.Changes {
		ss, ok := k.GetSubspace(c.Subspace)
		if !ok {
			return ErrUnknownSubspace(k.codespace, c.Subspace).Result()
		}

		var err error
		if len(c.Subkey) == 0 {
			k.Logger(ctx).Info(
				fmt.Sprintf("setting new parameter; key: %s, value: %s", c.Key, c.Value),
			)

			err = ss.Update(ctx, []byte(c.Key), []byte(c.Value))
		} else {
			k.Logger(ctx).Info(
				fmt.Sprintf("setting new parameter; key: %s, subkey: %s, value: %s", c.Key, c.Subspace, c.Value),
			)
			err = ss.UpdateWithSubkey(ctx, []byte(c.Key), []byte(c.Subkey), []byte(c.Value))
		}

		if err != nil {
			return ErrSettingParameter(k.codespace, c.Key, c.Subkey, c.Value, err.Error()).Result()
		}

		attr = append(attr, sdk.NewAttribute(types.AttributeKeyParam, c.Key), sdk.NewAttribute(types.AttributeKeyParamValue, c.Value))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeExecuteParamsChangeProposal, attr...),
	)

	return sdk.Result{}
}
