package types

import (
	"math/rand"
	"testing"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/stretchr/testify/require"
)

func TestNextAnnualProvisions(t *testing.T) {
	dp := DefaultParams()
	minter := DefaultInitialMinter(dp.InitalInflationAmount)

	require.Equal(t, DefaultInitalInflationAmount, minter.AnnualProvisions)
	require.Equal(t, DefaultInitalInflationAmount, minter.NextAnnualProvisions(dp.InflationFactorPerYear[0]))
	require.Equal(t, DefaultInitalInflationAmount.Mul(sdk.NewDecWithPrec(85, 2)), minter.NextAnnualProvisions(dp.InflationFactorPerYear[1]))
	require.Equal(t, DefaultInitalInflationAmount.Mul(sdk.NewDecWithPrec(9, 1)), minter.NextAnnualProvisions(dp.InflationFactorPerYear[4]))
}

func TestBlockProvision(t *testing.T) {
	minter := InitialMinter(0, sdk.NewDec(0))
	params := DefaultParams()

	secondsPerYear := int64(60 * 60 * 8766)

	tests := []struct {
		annualProvisions int64
		expProvisions    int64
	}{
		{secondsPerYear / 5, 1},
		{secondsPerYear/5 + 1, 1},
		{(secondsPerYear / 5) * 2, 2},
		{(secondsPerYear / 5) / 2, 0},
	}
	for i, tc := range tests {
		minter.AnnualProvisions = sdk.NewDec(tc.annualProvisions)
		provisions := minter.BlockProvision(params)

		expProvisions := sdk.NewCoin(params.MintDenom,
			sdk.NewInt(tc.expProvisions))

		require.True(t, expProvisions.IsEqual(provisions),
			"test: %v\n\tExp: %v\n\tGot: %v\n",
			i, tc.expProvisions, provisions)
	}
}

// Benchmarking :)
// previously using sdk.Int operations:
// BenchmarkBlockProvision-4 5000000 220 ns/op
//
// using sdk.Dec operations: (current implementation)
// BenchmarkBlockProvision-4 3000000 344 ns/op
func BenchmarkBlockProvision(b *testing.B) {
	s1 := rand.NewSource(100)
	r1 := rand.New(s1)

	annualAmount := sdk.NewDec(r1.Int63n(1000000))

	minter := InitialMinter(0, annualAmount)
	params := DefaultParams()
	params.InitalInflationAmount = annualAmount

	// run the BlockProvision function b.N times
	for n := 0; n < b.N; n++ {
		minter.BlockProvision(params)
	}
}

// Next annual provisions benchmarking
// BenchmarkNextAnnualProvisions-4 5000000 176 ns/op
func BenchmarkNextAnnualProvisions(b *testing.B) {
	s1 := rand.NewSource(100)
	r1 := rand.New(s1)

	annualAmount := sdk.NewDec(r1.Int63n(1000000))

	minter := InitialMinter(0, annualAmount)
	params := DefaultParams()
	params.InitalInflationAmount = annualAmount

	// run the NextAnnualProvisions function b.N times
	for n := 0; n < b.N; n++ {
		idx := n % len(params.InflationFactorPerYear)
		minter.NextAnnualProvisions(params.InflationFactorPerYear[idx])
	}

}
