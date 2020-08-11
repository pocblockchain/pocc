package keeper

import (
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/distribution/types"
)

// DistributeFromFeePool distributes funds from the distribution module account to
// a receiver address while updating the community pool
func (k Keeper) DistributeFromFeePool(ctx sdk.Context, amount sdk.Coins, receiveAddr sdk.AccAddress) sdk.Error {
	feePool := k.GetFeePool(ctx)

	// NOTE the community pool isn't a module account, however its coins
	// are held in the distribution module account. Thus the community pool
	// must be reduced separately from the SendCoinsFromModuleToAccount call
	newPool, negative := feePool.CommunityPool.SafeSub(sdk.NewDecCoins(amount))
	if negative {
		return types.ErrBadDistribution(k.codespace)
	}
	feePool.CommunityPool = newPool

	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiveAddr, amount)
	if err != nil {
		return err
	}

	k.SetFeePool(ctx, feePool)
	return nil
}

//AddCoinsFromAccountToFeePool updates communityPool and distribution's module account
func (k Keeper) AddCoinsFromAccountToFeePool(ctx sdk.Context, senderAddr sdk.AccAddress, amount sdk.Coins) sdk.Error {
	feePool := k.GetFeePool(ctx)

	// NOTE the community pool isn't a module account, however its coins
	// are held in the distribution module account. Thus the community pool
	// must be reduced separately from the SendCoinsFromModuleToAccount call
	newPool := feePool.CommunityPool.Add(sdk.NewDecCoins(amount))
	feePool.CommunityPool = newPool

	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, amount)
	if err != nil {
		return err
	}

	k.SetFeePool(ctx, feePool)
	return nil
}
