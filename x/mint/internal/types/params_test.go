package types

import (
	"testing"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/stretchr/testify/require"
)

func TestParamsString(t *testing.T) {
	expected := "Minting Params:\n  Mint Denom:poc\n  Initial Inflation Amoun:81000000000000000000000000.000000000000000000\n  Inflation Factors Per Year:0:1.0000000000000000001:0.8500000000000000002:0.8500000000000000003:0.8500000000000000004:0.900000000000000000\n  Blocks Per Year:6311520\n"
	require.Equal(t, expected, DefaultParams().String())
}

func TestValidateParams(t *testing.T) {
	params := NewParams("", sdk.NewDec(100), []sdk.Dec{sdk.NewDec(1), sdk.NewDec(2)}, 100)
	require.Error(t, ValidateParams(params))

	params = NewParams("poc", sdk.NewDec(-1), []sdk.Dec{sdk.NewDec(1), sdk.NewDec(2)}, 100)
	require.Error(t, ValidateParams(params))

	params = NewParams("poc", sdk.NewDec(100), []sdk.Dec{sdk.NewDec(1), sdk.NewDec(-2)}, 100)
	require.Error(t, ValidateParams(params))

}
