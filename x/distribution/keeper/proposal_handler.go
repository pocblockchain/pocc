package keeper

import (
	"fmt"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/distribution/types"
)

// HandleCommunityPoolSpendProposal is a handler for executing a passed community spend proposal
func HandleCommunityPoolSpendProposal(ctx sdk.Context, k Keeper, p types.CommunityPoolSpendProposal) sdk.Result {
	ctx.Logger().Info("handleCommunityPoolSpendProposal", "proposal", p)
	if k.blacklistedAddrs[p.Recipient.String()] {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is blacklisted from receiving external funds", p.Recipient)).Result()
	}

	err := k.DistributeFromFeePool(ctx, p.Amount, p.Recipient)
	if err != nil {
		return err.Result()
	}

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("transferred %s from the community pool to recipient %s", p.Amount, p.Recipient))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeExecutionCommunityPoolSpendProposal,
			sdk.NewAttribute(sdk.AttributeKeyAmount, p.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, p.Recipient.String()),
		),
	)
	return sdk.Result{}
}
