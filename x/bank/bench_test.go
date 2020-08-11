package bank_test

import (
	"github.com/pocblockchain/pocc/x/token"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/auth"
	"github.com/pocblockchain/pocc/x/bank"
	"github.com/pocblockchain/pocc/x/bank/internal/keeper"
	"github.com/pocblockchain/pocc/x/bank/internal/types"
	"github.com/pocblockchain/pocc/x/mock"
	"github.com/pocblockchain/pocc/x/supply"
)

var moduleAccAddr = sdk.AccAddress([]byte("moduleAcc"))

// initialize the mock application for this module
func getMockApp(t *testing.T) testInput {
	input, err := getBenchmarkMockApp()
	supply.RegisterCodec(input.mApp.Cdc)

	require.NoError(t, err)
	return input
}

// overwrite the mock init chainer
func getInitChainer(mapp *mock.App, keeper keeper.BaseKeeper, tk token.Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)
		bankGenesis := bank.DefaultGenesisState()
		bank.InitGenesis(ctx, keeper, bankGenesis)

		//init foocoinInfo and barcoinInfo
		infos := token.DefaultGenesisState()
		gs := &infos
		foocoinInfo := sdk.TokenInfoWithoutSupply{
			Symbol:        "foocoin",
			Issuer:        "",
			IsSendEnabled: true,
			Decimals:      8,
		}
		barcoinInfo := sdk.TokenInfoWithoutSupply{
			Symbol:        "barcoin",
			Issuer:        "",
			IsSendEnabled: true,
			Decimals:      8,
		}

		gs.AddTokenInfoWithoutSupplyIntoGenesis(foocoinInfo)
		gs.AddTokenInfoWithoutSupplyIntoGenesis(barcoinInfo)

		token.InitGenesis(ctx, tk, *gs)
		return abci.ResponseInitChain{}
	}
}

type testInput struct {
	mApp *mock.App
	tk   token.Keeper
}

// getBenchmarkMockApp initializes a mock application for this module, for purposes of benchmarking
// Any long term API support commitments do not apply to this function.
func getBenchmarkMockApp() (testInput, error) {
	mapp := mock.NewApp()
	types.RegisterCodec(mapp.Cdc)

	token.RegisterCodec(mapp.Cdc)
	keyToken := sdk.NewKVStoreKey(token.StoreKey)
	tokenKeeper := token.NewKeeper(mapp.Cdc, keyToken, nil, nil, mapp.ParamsKeeper.Subspace(token.DefaultParamspace))

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[moduleAccAddr.String()] = true

	bankKeeper := keeper.NewBaseKeeper(
		mapp.AccountKeeper, &tokenKeeper,
		mapp.ParamsKeeper.Subspace(types.DefaultParamspace),
		types.DefaultCodespace,
		blacklistedAddrs,
	)

	//	(&bankKeeper).SetTokenKeeper(&tokenKeeper)
	mapp.Router().AddRoute(types.RouterKey, bank.NewHandler(bankKeeper))
	mapp.SetInitChainer(getInitChainer(mapp, bankKeeper, tokenKeeper))

	err := mapp.CompleteSetup(keyToken)
	return testInput{mapp, tokenKeeper}, err
}

// initialize the mock application for this module
func getMockAppWithSetTokenKeeper(t *testing.T) testInput {
	input, err := getBenchmarkMockAppWithSetTokenKeeper()
	supply.RegisterCodec(input.mApp.Cdc)

	require.NoError(t, err)
	return input
}

// getBenchmarkMockAppWithSetTokenKeeper initializes a mock application for this module, token keeper is set
func getBenchmarkMockAppWithSetTokenKeeper() (testInput, error) {
	mapp := mock.NewApp()
	types.RegisterCodec(mapp.Cdc)

	token.RegisterCodec(mapp.Cdc)
	keyToken := sdk.NewKVStoreKey(token.StoreKey)
	tokenKeeper := token.NewKeeper(mapp.Cdc, keyToken, nil, nil, mapp.ParamsKeeper.Subspace(token.DefaultParamspace))

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[moduleAccAddr.String()] = true

	bankKeeper := keeper.NewBaseKeeper(
		mapp.AccountKeeper, nil,
		mapp.ParamsKeeper.Subspace(types.DefaultParamspace),
		types.DefaultCodespace,
		blacklistedAddrs,
	)

	bankKeeper.SetTokenKeeper(&tokenKeeper)
	mapp.Router().AddRoute(types.RouterKey, bank.NewHandler(bankKeeper))
	mapp.SetInitChainer(getInitChainer(mapp, bankKeeper, tokenKeeper))

	err := mapp.CompleteSetup(keyToken)
	return testInput{mapp, tokenKeeper}, err
}

func BenchmarkOneBankSendTxPerBlock(b *testing.B) {
	benchmarkApp, _ := getBenchmarkMockApp()

	// Add an account at genesis
	acc := &auth.BaseAccount{
		Address: addr1,
		// Some value conceivably higher than the benchmarks would ever go
		Coins: sdk.Coins{sdk.NewInt64Coin("foocoin", 100000000000)},
	}
	accs := []auth.Account{acc}

	// Construct genesis state
	mock.SetGenesis(benchmarkApp.mApp, accs)
	// Precompute all txs
	txs := mock.GenSequenceOfTxs([]sdk.Msg{sendMsg1}, []uint64{0}, []uint64{uint64(0)}, b.N, priv1)
	b.ResetTimer()
	// Run this with a profiler, so its easy to distinguish what time comes from
	// Committing, and what time comes from Check/Deliver Tx.
	for i := 0; i < b.N; i++ {
		benchmarkApp.mApp.BeginBlock(abci.RequestBeginBlock{})
		x := benchmarkApp.mApp.Check(txs[i])
		if !x.IsOK() {
			panic("something is broken in checking transaction")
		}
		benchmarkApp.mApp.Deliver(txs[i])
		benchmarkApp.mApp.EndBlock(abci.RequestEndBlock{})
		benchmarkApp.mApp.Commit()
	}
}

func BenchmarkOneBankMultiSendTxPerBlock(b *testing.B) {
	benchmarkApp, _ := getBenchmarkMockApp()

	// Add an account at genesis
	acc := &auth.BaseAccount{
		Address: addr1,
		// Some value conceivably higher than the benchmarks would ever go
		Coins: sdk.Coins{sdk.NewInt64Coin("foocoin", 100000000000)},
	}
	accs := []auth.Account{acc}

	// Construct genesis state
	mock.SetGenesis(benchmarkApp.mApp, accs)
	// Precompute all txs
	txs := mock.GenSequenceOfTxs([]sdk.Msg{multiSendMsg1}, []uint64{0}, []uint64{uint64(0)}, b.N, priv1)
	b.ResetTimer()
	// Run this with a profiler, so its easy to distinguish what time comes from
	// Committing, and what time comes from Check/Deliver Tx.
	for i := 0; i < b.N; i++ {
		benchmarkApp.mApp.BeginBlock(abci.RequestBeginBlock{})
		x := benchmarkApp.mApp.Check(txs[i])
		if !x.IsOK() {
			panic("something is broken in checking transaction")
		}
		benchmarkApp.mApp.Deliver(txs[i])
		benchmarkApp.mApp.EndBlock(abci.RequestEndBlock{})
		benchmarkApp.mApp.Commit()
	}
}
