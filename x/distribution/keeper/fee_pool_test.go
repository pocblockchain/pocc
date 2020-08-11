package keeper

import (
	"testing"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/distribution/types"
	"github.com/stretchr/testify/require"
)

func TestAddCoinsFromAccountToFeePool(t *testing.T) {
	ctx, ak, keeper, _, supplyKeeper := CreateTestInputDefault(t, false, 1000)
	require.Equal(t, sdk.ZeroInt(), supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins().AmountOf(sdk.NativeToken))
	require.Equal(t, sdk.DecCoins(nil), keeper.GetFeePool(ctx).CommunityPool)
	initCoins := ak.GetAccount(ctx, delAddr1).GetCoins()

	coins1 := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, sdk.NewInt(1)))
	coins2 := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, sdk.NewInt(2)))
	coins3 := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, sdk.NewInt(3)))
	coins4 := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, sdk.NewInt(4)))

	coins6 := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, sdk.NewInt(6)))
	coins10 := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, sdk.NewInt(10)))

	err := keeper.AddCoinsFromAccountToFeePool(ctx, delAddr1, coins1)
	require.Nil(t, err)
	require.Equal(t, coins1, supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins())
	require.Equal(t, sdk.NewDecCoins(coins1), keeper.GetFeePool(ctx).CommunityPool)
	require.Equal(t, initCoins.Sub(coins1), ak.GetAccount(ctx, delAddr1).GetCoins())

	err = keeper.AddCoinsFromAccountToFeePool(ctx, delAddr1, coins2)
	require.Nil(t, err)
	require.Equal(t, coins3, supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins())
	require.Equal(t, sdk.NewDecCoins(coins3), keeper.GetFeePool(ctx).CommunityPool)
	require.Equal(t, initCoins.Sub(coins3), ak.GetAccount(ctx, delAddr1).GetCoins())

	err = keeper.AddCoinsFromAccountToFeePool(ctx, delAddr1, coins3)
	require.Nil(t, err)
	require.Equal(t, coins6, supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins())
	require.Equal(t, sdk.NewDecCoins(coins6), keeper.GetFeePool(ctx).CommunityPool)
	require.Equal(t, initCoins.Sub(coins6), ak.GetAccount(ctx, delAddr1).GetCoins())

	err = keeper.AddCoinsFromAccountToFeePool(ctx, delAddr1, coins4)
	require.Nil(t, err)
	require.Equal(t, coins10, supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins())
	require.Equal(t, sdk.NewDecCoins(coins10), keeper.GetFeePool(ctx).CommunityPool)
	require.Equal(t, initCoins.Sub(coins10), ak.GetAccount(ctx, delAddr1).GetCoins())
}

func TestDistributeFromFeePool(t *testing.T) {
	ctx, ak, keeper, _, supplyKeeper := CreateTestInputDefault(t, false, 1000)
	require.Equal(t, sdk.ZeroInt(), supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins().AmountOf(sdk.NativeToken))
	require.Equal(t, sdk.DecCoins(nil), keeper.GetFeePool(ctx).CommunityPool)
	initCoins := ak.GetAccount(ctx, delAddr1).GetCoins()

	err := keeper.AddCoinsFromAccountToFeePool(ctx, delAddr1, initCoins)
	require.Nil(t, err)
	require.Equal(t, initCoins, supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins())
	require.Equal(t, sdk.NewDecCoins(initCoins), keeper.GetFeePool(ctx).CommunityPool)
	require.Equal(t, sdk.Coins(nil), ak.GetAccount(ctx, delAddr1).GetCoins())

	coins1 := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, sdk.NewInt(1)))
	coins2 := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, sdk.NewInt(2)))
	coins3 := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, sdk.NewInt(3)))

	coins6 := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, sdk.NewInt(6)))

	//distribute 1NativeCoin to delAddr1
	err = keeper.DistributeFromFeePool(ctx, coins1, delAddr1)
	require.Nil(t, err)
	require.Equal(t, initCoins.Sub(coins1), supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins())
	require.Equal(t, sdk.NewDecCoins(initCoins.Sub(coins1)), keeper.GetFeePool(ctx).CommunityPool)
	require.Equal(t, sdk.Coins(coins1), ak.GetAccount(ctx, delAddr1).GetCoins())

	//distribute 2NativeCoin to delAddr2
	err = keeper.DistributeFromFeePool(ctx, coins2, delAddr2)
	require.Nil(t, err)
	require.Equal(t, initCoins.Sub(coins3), supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins())
	require.Equal(t, sdk.NewDecCoins(initCoins.Sub(coins3)), keeper.GetFeePool(ctx).CommunityPool)
	require.Equal(t, sdk.Coins(coins1), ak.GetAccount(ctx, delAddr1).GetCoins())
	require.Equal(t, sdk.Coins(initCoins.Add(coins2)), ak.GetAccount(ctx, delAddr2).GetCoins())

	//distribute 3NativeToken to newAddr
	newAddr, _ := sdk.AccAddressFromBech32("poc12jwptqcnzkk4d7yupmwlnkzjkm9hvp0rr0chxr")
	ak.GetOrNewAccount(ctx, newAddr)
	err = keeper.DistributeFromFeePool(ctx, coins3, newAddr)
	require.Nil(t, err)
	require.Equal(t, initCoins.Sub(coins6), supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins())
	require.Equal(t, sdk.NewDecCoins(initCoins.Sub(coins6)), keeper.GetFeePool(ctx).CommunityPool)
	require.Equal(t, sdk.Coins(coins1), ak.GetAccount(ctx, delAddr1).GetCoins())
	require.Equal(t, sdk.Coins(initCoins.Add(coins2)), ak.GetAccount(ctx, delAddr2).GetCoins())
	require.Equal(t, sdk.Coins(coins3), ak.GetAccount(ctx, newAddr).GetCoins())
}
