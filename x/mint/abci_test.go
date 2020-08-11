package mint

import (
	"testing"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/auth"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestBegingBlocker(t *testing.T) {
	input := getMockApp(t)
	keeper := input.keeper
	supplyKeeper := input.supplyKeeper
	mapp := input.mApp

	//________in the first year________
	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})

	params := keeper.GetParams(ctx)
	annualAmt := params.InitalInflationAmount
	blockAmt := annualAmt.QuoInt64(int64(params.BlocksPerYear)).TruncateInt()
	collectFee := blockAmt
	require.Equal(t, DefaultParams(), params)
	minter := keeper.GetMinter(ctx)
	require.Equal(t, uint64(0), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	//height = params.BlocksPerYear - 1
	header = abci.Header{Height: int64(params.BlocksPerYear - 1)}
	collectFee = collectFee.Add(blockAmt)
	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	//________in the second year________
	// the 1st height of 2nd year
	epoch := uint64(2)
	header = abci.Header{Height: int64(params.BlocksPerYear * (epoch - 1))}
	annualAmt = annualAmt.Mul(params.InflationFactorPerYear[epoch-1])
	blockAmt = annualAmt.QuoInt64(int64(params.BlocksPerYear)).TruncateInt()
	require.Equal(t, DefaultParams(), params)

	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)
	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	// the last height of 2nd year
	header = abci.Header{Height: int64(params.BlocksPerYear*epoch - 1)}
	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)
	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	//________in the third year________
	// the 1st height of 3rd year
	epoch = uint64(3)
	header = abci.Header{Height: int64(params.BlocksPerYear * (epoch - 1))}
	annualAmt = annualAmt.Mul(params.InflationFactorPerYear[epoch-1])
	blockAmt = annualAmt.QuoInt64(int64(params.BlocksPerYear)).TruncateInt()
	require.Equal(t, DefaultParams(), params)

	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)
	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	// the last height of 3rd year
	header = abci.Header{Height: int64(params.BlocksPerYear*epoch - 1)}
	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)
	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	//________in the 4th year________
	// the 1st height of 4th year
	epoch = uint64(4)
	header = abci.Header{Height: int64(params.BlocksPerYear * (epoch - 1))}
	annualAmt = annualAmt.Mul(params.InflationFactorPerYear[epoch-1])
	blockAmt = annualAmt.QuoInt64(int64(params.BlocksPerYear)).TruncateInt()
	require.Equal(t, DefaultParams(), params)

	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)
	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	// the last height of 4rd year
	header = abci.Header{Height: int64(params.BlocksPerYear*epoch - 1)}
	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)
	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	//________in the 5th year________
	// the 1st height of 5th year
	epoch = uint64(5)
	header = abci.Header{Height: int64(params.BlocksPerYear * (epoch - 1))}
	annualAmt = annualAmt.Mul(params.InflationFactorPerYear[epoch-1])
	blockAmt = annualAmt.QuoInt64(int64(params.BlocksPerYear)).TruncateInt()
	require.Equal(t, DefaultParams(), params)

	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)
	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	// the last height of 5th year
	header = abci.Header{Height: int64(params.BlocksPerYear*epoch - 1)}
	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)
	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	//________in the 6th year________
	// the 1st height of 6th year
	epoch = uint64(6)
	header = abci.Header{Height: int64(params.BlocksPerYear * (epoch - 1))}
	annualAmt = annualAmt.Mul(params.InflationFactorPerYear[4])
	blockAmt = annualAmt.QuoInt64(int64(params.BlocksPerYear)).TruncateInt()
	require.Equal(t, DefaultParams(), params)

	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)

	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	// the last height of 6th year
	header = abci.Header{Height: int64(params.BlocksPerYear*epoch - 1)}
	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)

	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	//________in the 7th year________
	// the 1st height of 7th year
	epoch = uint64(7)
	header = abci.Header{Height: int64(params.BlocksPerYear * (epoch - 1))}
	annualAmt = annualAmt.Mul(params.InflationFactorPerYear[4])
	blockAmt = annualAmt.QuoInt64(int64(params.BlocksPerYear)).TruncateInt()
	require.Equal(t, DefaultParams(), params)

	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)

	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	// the last height of 7th year
	header = abci.Header{Height: int64(params.BlocksPerYear*epoch - 1)}
	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)

	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	//________in the 8th year________
	// the 1st height of 7th year
	epoch = uint64(8)
	header = abci.Header{Height: int64(params.BlocksPerYear * (epoch - 1))}
	annualAmt = annualAmt.Mul(params.InflationFactorPerYear[4])
	blockAmt = annualAmt.QuoInt64(int64(params.BlocksPerYear)).TruncateInt()
	require.Equal(t, DefaultParams(), params)

	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)

	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	// the last height of 8th year
	header = abci.Header{Height: int64(params.BlocksPerYear*epoch - 1)}
	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)

	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	//________in the 9th year________
	// the 1st height of 9th year
	epoch = uint64(9)
	header = abci.Header{Height: int64(params.BlocksPerYear * (epoch - 1))}
	annualAmt = annualAmt.Mul(params.InflationFactorPerYear[4])
	blockAmt = annualAmt.QuoInt64(int64(params.BlocksPerYear)).TruncateInt()
	require.Equal(t, DefaultParams(), params)

	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)

	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	// the last height of 9th year
	header = abci.Header{Height: int64(params.BlocksPerYear*epoch - 1)}
	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)

	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	//________in the 10th year________
	// the 1st height of 10th year
	epoch = uint64(10)
	header = abci.Header{Height: int64(params.BlocksPerYear * (epoch - 1))}
	annualAmt = annualAmt.Mul(params.InflationFactorPerYear[4])
	blockAmt = annualAmt.QuoInt64(int64(params.BlocksPerYear)).TruncateInt()
	require.Equal(t, DefaultParams(), params)

	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)

	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

	// the last height of 9th year
	header = abci.Header{Height: int64(params.BlocksPerYear*epoch - 1)}
	ctx = ctx.WithBlockHeader(header)
	BeginBlocker(ctx, keeper)

	minter = keeper.GetMinter(ctx)
	require.Equal(t, uint64(epoch-1), minter.CurrentYearIndex)
	require.Equal(t, annualAmt, minter.AnnualProvisions)
	require.Equal(t, blockAmt, minter.BlockProvision(params).Amount)

	collectFee = collectFee.Add(blockAmt)
	require.Equal(t, collectFee, supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins().AmountOf(sdk.DefaultBondDenom))

}
