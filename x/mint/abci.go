package mint

import (
	"fmt"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/mint/internal/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k Keeper) {
	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	// recalculate inflation rate
	epoch := uint64(sdk.NewInt(ctx.BlockHeight()).QuoRaw(int64(params.BlocksPerYear)).Int64())

	if epoch > minter.CurrentYearIndex {
		//a new year begin, update minter
		minter.CurrentYearIndex = epoch
		idx := int(epoch)

		if int(epoch) >= len(params.InflationFactorPerYear) {
			idx = len(params.InflationFactorPerYear) - 1
		}
		minter.AnnualProvisions = minter.NextAnnualProvisions(params.InflationFactorPerYear[idx])
	}

	k.SetMinter(ctx, minter)

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyYear, fmt.Sprintf("%v", minter.CurrentYearIndex)),
			sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
}
