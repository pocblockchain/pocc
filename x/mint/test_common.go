package mint

import (
	"testing"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/auth"
	"github.com/pocblockchain/pocc/x/bank"
	"github.com/pocblockchain/pocc/x/mint/internal/keeper"
	"github.com/pocblockchain/pocc/x/mock"
	"github.com/pocblockchain/pocc/x/staking"
	"github.com/pocblockchain/pocc/x/supply"
	supplyexported "github.com/pocblockchain/pocc/x/supply/exported"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	valTokens  = sdk.TokensFromConsensusPower(42)
	initTokens = sdk.TokensFromConsensusPower(100000)
	valCoins   = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, valTokens))
	initCoins  = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))
)

type testInput struct {
	mApp   *mock.App
	keeper Keeper
	//sk           staking.Keeper
	supplyKeeper supply.Keeper
}

func getMockApp(t *testing.T) testInput {
	mApp := mock.NewApp()

	RegisterCodec(mApp.Cdc)
	supply.RegisterCodec(mApp.Cdc)

	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyMint := sdk.NewKVStoreKey(StoreKey)

	feeCollector := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollector.String()] = true
	blacklistedAddrs[notBondedPool.String()] = true
	blacklistedAddrs[bondPool.String()] = true

	bankKeeper := bank.NewBaseKeeper(mApp.AccountKeeper, nil, mApp.ParamsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, blacklistedAddrs)
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.NotBondedPoolName: []string{supply.Burner, supply.Staking},
		staking.BondedPoolName:    []string{supply.Burner, supply.Staking},
		ModuleName:                []string{supply.Minter},
	}
	supplyKeeper := supply.NewKeeper(mApp.Cdc, keySupply, mApp.AccountKeeper, bankKeeper, maccPerms)
	stakingkeeper := staking.NewKeeper(mApp.Cdc, keyStaking, tkeyStaking, supplyKeeper, mApp.ParamsKeeper.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	mintKeeper := NewKeeper(mApp.Cdc, keyMint, mApp.ParamsKeeper.Subspace(DefaultParamspace), stakingkeeper, supplyKeeper, feeCollector.Name)

	mApp.SetBeginBlocker(getBeginBlocker(mintKeeper))
	mApp.SetInitChainer(getInitChainer(mApp, stakingkeeper, mApp.AccountKeeper, supplyKeeper, mintKeeper,
		[]supplyexported.ModuleAccountI{feeCollector, notBondedPool, bondPool}))

	require.NoError(t, mApp.CompleteSetup(keyStaking, tkeyStaking, keySupply, keyMint))

	mock.SetGenesis(mApp, nil)

	return testInput{mApp: mApp, keeper: mintKeeper, supplyKeeper: supplyKeeper}
}

//getBeginBlocker return mint beginblocker
func getBeginBlocker(keeper Keeper) sdk.BeginBlocker {
	return func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		BeginBlocker(ctx, keeper)
		return abci.ResponseBeginBlock{}
	}
}

// getInitChainer initializes the chainer of the mock app and sets the genesis
// state. It returns an empty ResponseInitChain.
func getInitChainer(mapp *mock.App, stakingkeeper staking.Keeper, accountKeeper auth.AccountKeeper, supplyKeeper supply.Keeper, keeper keeper.Keeper,
	blacklistedAddrs []supplyexported.ModuleAccountI) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}

		stakingGenesis := staking.DefaultGenesisState()
		totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens.MulRaw(int64(len(mapp.GenesisAccounts)))))
		supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

		validators := staking.InitGenesis(ctx, stakingkeeper, accountKeeper, supplyKeeper, stakingGenesis)

		//inital minter
		var mintGenesis = DefaultGenesisState()

		InitGenesis(ctx, keeper, mintGenesis)

		return abci.ResponseInitChain{
			Validators: validators,
		}

		return abci.ResponseInitChain{
			Validators: validators,
		}
	}
}
