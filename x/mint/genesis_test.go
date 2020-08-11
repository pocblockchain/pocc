package mint

import (
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenesisStateEqual(t *testing.T) {
	gs1 := DefaultGenesisState()
	gs2 := DefaultGenesisState()

	require.True(t, gs1.Equal(gs2))
	gs2.Params.InitalInflationAmount = sdk.NewDecWithPrec(1, 1)
	require.False(t, gs1.Equal(gs2))

	gs3 := NewGenesisState(gs1.Minter, gs1.Params)
	require.True(t, gs1.Equal(gs3))

	gs4 := NewGenesisState(gs2.Minter, gs2.Params)
	require.True(t, gs2.Equal(gs4))
}

func TestGenesisStateIsEmpty(t *testing.T) {
	gs1 := GenesisState{}
	require.True(t, gs1.IsEmpty())

	gs1.Params.MintDenom = sdk.DefaultBondDenom
	require.False(t, gs1.IsEmpty())

	gs1 = DefaultGenesisState()
	require.False(t, gs1.IsEmpty())
}

func TestGenesisStateString(t *testing.T) {
	expected := "Minter:\n  CurrentYearIndex:0\n  AnnualProvisions:81000000000000000000000000.000000000000000000\nMinting Params:\n  Mint Denom:poc\n  Initial Inflation Amoun:81000000000000000000000000.000000000000000000\n  Inflation Factors Per Year:0:1.0000000000000000001:0.8500000000000000002:0.8500000000000000003:0.8500000000000000004:0.900000000000000000\n  Blocks Per Year:6311520"
	gs1 := DefaultGenesisState()
	require.Equal(t, expected, gs1.String())
}
