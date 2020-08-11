package types

import (
	"testing"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/stretchr/testify/require"
)

func TestParamsEqual(t *testing.T) {
	p1 := DefaultParams()
	p2 := DefaultParams()

	require.True(t, p1.Equal(p2))

	p2.TokenCacheSize--
	require.False(t, p1.Equal(p2))

	p2.TokenCacheSize++
	p2.ReservedSymbols = append(p2.ReservedSymbols, "mytoken")
	require.False(t, p1.Equal(p2))
}

func TestParamsString(t *testing.T) {
	expectedStr := "Params:TokenCacheSize:32\tNewTokenFee:1000000000000000000000\tReservedSymbols:btc,eth,eos,usdt,bch,bsv,ltc,bnb,xrp,okb,ht,dash,etc,neo,atom,zec,ont,doge,tusd,bat,qtum,vsys,iost,dcr,zrx,beam,grin\t"
	p := DefaultParams()

	require.Equal(t, expectedStr, p.String())
}

func TestParamsMarshalJSON(t *testing.T) {
	p := DefaultParams()
	bz, err := ModuleCdc.MarshalJSON(p)
	require.Nil(t, err)

	p1 := Params{}
	ModuleCdc.UnmarshalJSON(bz, &p1)
	require.True(t, p.Equal(p1))
}

func TestParamValidate(t *testing.T) {
	p := DefaultParams()
	require.Nil(t, p.Validate())

	p.NewTokenFee = sdk.NewInt(1)
	require.Nil(t, p.Validate())

	p.NewTokenFee = sdk.NewInt(-1)

	require.NotNil(t, p.Validate())

}
