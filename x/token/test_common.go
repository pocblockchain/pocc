package token

import (
	"testing"

	"github.com/pocblockchain/pocc/codec"
	"github.com/pocblockchain/pocc/store"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/auth"
	"github.com/pocblockchain/pocc/x/bank"
	distr "github.com/pocblockchain/pocc/x/distribution"
	"github.com/pocblockchain/pocc/x/gov"
	"github.com/pocblockchain/pocc/x/mint"
	"github.com/pocblockchain/pocc/x/params"
	"github.com/pocblockchain/pocc/x/slashing"
	"github.com/pocblockchain/pocc/x/staking"
	"github.com/pocblockchain/pocc/x/supply"
	"github.com/pocblockchain/pocc/x/token/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var (
	EthToken  = "eth"
	BtcToken  = "btc"
	UsdtToken = "usdt"
)

var TestTokenData = []sdk.TokenInfo{
	{
		Symbol:        sdk.Symbol(sdk.NativeToken),
		Issuer:        "",
		IsSendEnabled: true,
		Decimals:      8,
		TotalSupply:   sdk.NewIntWithDecimal(21, 15),
	},
	{
		Symbol:        sdk.Symbol(BtcToken),
		Issuer:        "",
		IsSendEnabled: true,
		Decimals:      8,
		TotalSupply:   sdk.NewIntWithDecimal(21, 15),
	},
	{
		Symbol:        sdk.Symbol(EthToken),
		Issuer:        "",
		IsSendEnabled: true,
		Decimals:      18,
		TotalSupply:   sdk.NewInt(0),
	},
	{
		Symbol:        sdk.Symbol(UsdtToken),
		Issuer:        "0xFF760fcB0fa4Ba68d9DD2e28fc7A3c593b5d2106",
		IsSendEnabled: true,
		Decimals:      18,
		TotalSupply:   sdk.NewIntWithDecimal(1, 28),
	},
}

type testEnv struct {
	cdc           *codec.Codec
	ctx           sdk.Context
	bankKeeper    bank.Keeper
	accountKeeper auth.AccountKeeper
	tokenKeeper   Keeper
	supplyKeeper  supply.Keeper
	paramsKeeper  params.Keeper
	mintKeeper    mint.Keeper
	distrKeeper   distr.Keeper
}

func setupTestEnv(t *testing.T) testEnv {
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyToken := sdk.NewKVStoreKey(types.StoreKey)
	keyGov := sdk.NewKVStoreKey(gov.StoreKey)
	keySlash := sdk.NewKVStoreKey(slashing.StoreKey)
	keyMint := sdk.NewKVStoreKey(mint.StoreKey)
	keyDistr := sdk.NewKVStoreKey(distr.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyToken, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyGov, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySlash, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMint, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	assert.Nil(t, err)

	//register cdc
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	staking.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	params.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	mint.RegisterCodec(cdc)
	RegisterCodec(cdc)

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)
	distrAcc := supply.NewEmptyModuleAccount(distr.ModuleName)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.String()] = true
	blacklistedAddrs[feeCollectorAcc.String()] = true
	blacklistedAddrs[notBondedPool.String()] = true
	blacklistedAddrs[bondPool.String()] = true
	blacklistedAddrs[distrAcc.String()] = true

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, nil, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, blacklistedAddrs)
	bankKeeper.SetSendEnabled(ctx, true)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		staking.NotBondedPoolName: []string{supply.Burner, supply.Staking},
		staking.BondedPoolName:    []string{supply.Burner, supply.Staking},
		mint.ModuleName:           []string{supply.Minter},
		types.ModuleName:          {supply.Minter, supply.Burner},
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)

	stakingKeeper := staking.NewKeeper(cdc, keyStaking, tkeyStaking, supplyKeeper, pk.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	stakingKeeper.SetParams(ctx, staking.DefaultParams())

	initPower := int64(10)
	numValidators := 4
	initTokens := sdk.TokensFromConsensusPower(initPower)
	totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, initTokens.MulRaw(int64(numValidators))))
	notBondedPool.SetCoins(totalSupply)

	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	supplyKeeper.SetModuleAccount(ctx, bondPool)
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)
	supplyKeeper.SetModuleAccount(ctx, distrAcc)
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	distrKeeper := distr.NewKeeper(cdc, keyDistr, pk.Subspace(distr.DefaultParamspace), stakingKeeper, supplyKeeper, distr.DefaultCodespace, auth.FeeCollectorName, nil)
	mintKeeper := mint.NewKeeper(cdc, keyMint, pk.Subspace(mint.DefaultParamspace), stakingKeeper, supplyKeeper, auth.FeeCollectorName)

	tokenKeeper := NewKeeper(cdc, keyToken, distrKeeper, supplyKeeper, pk.Subspace(DefaultParamspace))
	tokenKeeper.SetParams(ctx, DefaultParams())

	//set native token
	tokenKeeper.SetTokenInfoWithoutSupply(ctx, sdk.NewTokenInfoWithoutSupply(sdk.NativeToken, "", true, 8))

	// set the community pool to pay back the constant fee
	distrKeeper.SetFeePool(ctx, distr.InitialFeePool())
	distrKeeper.SetCommunityTax(ctx, sdk.NewDecWithPrec(2, 2))
	distrKeeper.SetBaseProposerReward(ctx, sdk.NewDecWithPrec(1, 2))
	distrKeeper.SetBonusProposerReward(ctx, sdk.NewDecWithPrec(4, 2))

	return testEnv{cdc: cdc, ctx: ctx, bankKeeper: bankKeeper, accountKeeper: accountKeeper, supplyKeeper: supplyKeeper, paramsKeeper: pk, mintKeeper: mintKeeper, distrKeeper: distrKeeper, tokenKeeper: tokenKeeper}
}
