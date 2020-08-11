package types

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestIsValidTokenName(t *testing.T) {
	testdata := []struct {
		name  string
		valid bool
	}{
		{NativeToken, true},
		{"bh124", true},
		{"bh12345678901234", true},
		{"bhabc", true},
		{"bhabc123", true},
		{"bh123456789012345", false}, // length limit
		{"bhCABC", false},
		{" bh124", false},
		{"_bh124", false},
		{"bh123456789012345", false},
		{"BhT", false},
		{"bHT", false},
		{"bh 123", false},
		{"bh123 ", false},
		{"bh#123 ", false},
		{"bh123% ", false},
		{"HBC124", false},
		{"bh^124", false},
		{"bh 125", false},
		{"bh 125*", false},
		{"bh125 ", false},
		{"1bhC", false},
		{"B1HC", false},
		{"BTC", false},
		{"bhABCDEFGHIGKLMNOP", false},
	}

	for _, d := range testdata {

		assert.Equal(t, d.valid, Symbol(d.name).IsValidTokenName(), d.name)

		if Symbol(d.name).IsValidTokenName() {
			assert.Equal(t, strings.ToLower(d.name), Symbol(d.name).ToDenomName())
			assert.Equal(t, d.name, Symbol(d.name).String())
			// must be valid denom
			assert.Nil(t, validateDenom(d.name))
			// symbol == DenomName // TODO remove symbol ,use DenomName
			assert.Equal(t, d.name, Symbol(d.name).ToDenomName())
		} else {
			assert.Equal(t, "", Symbol(d.name).String())
		}
	}
}

func TestTokenInfoString(t *testing.T) {
	expected := "\n\tSymbol:btc\n\tIssuer:iss\n\tIsSendEnabled:false\n\tDecimals:8\n\tTotalSupply:1000000\n\t"
	d := TokenInfo{
		Symbol:        "btc",
		Issuer:        "iss",
		IsSendEnabled: false,
		Decimals:      8,
		TotalSupply:   NewInt(1000000),
	}
	assert.Equal(t, expected, d.String())
}

func TestTokenInfoIsValid(t *testing.T) {
	tokenInfo := TokenInfo{
		Symbol:        Symbol("btc"),
		Issuer:        "",
		IsSendEnabled: true,
		Decimals:      8,
		TotalSupply:   NewIntWithDecimal(21, 15),
	}

	assert.True(t, tokenInfo.IsValid())

	//symbol is illegal
	tokenInfo.Symbol = Symbol("Btc")
	assert.False(t, tokenInfo.IsValid())

	tokenInfo.Symbol = Symbol("btc")

	//Decimal is illegal
	for i := 1; i <= 18; i++ {
		tokenInfo.Decimals = uint64(i)
		assert.True(t, tokenInfo.IsValid())
	}
	tokenInfo.Decimals = 19
	assert.False(t, tokenInfo.IsValid())

	tokenInfo.Decimals = 18
	assert.True(t, tokenInfo.IsValid())

	//TotalSupply is -1
	tokenInfo.TotalSupply = NewInt(-1)
	assert.False(t, tokenInfo.IsValid())

	//TotalSupply is 0
	tokenInfo.TotalSupply = NewInt(0)
	assert.True(t, tokenInfo.IsValid())
}

func TestTokenInfoWithoutSupplyString(t *testing.T) {
	expected := "\n\tSymbol:btc\n\tIssuer:iss\n\tIsSendEnabled:false\n\tDecimals:8\n\t"
	d := TokenInfoWithoutSupply{
		Symbol:        "btc",
		Issuer:        "iss",
		IsSendEnabled: false,
		Decimals:      8,
	}
	assert.Equal(t, expected, d.String())
}

func TestTokenInfoWithoutSupplyIsValid(t *testing.T) {
	tokenInfo := TokenInfoWithoutSupply{
		Symbol:        Symbol("btc"),
		Issuer:        "",
		IsSendEnabled: true,
		Decimals:      8,
	}

	assert.True(t, tokenInfo.IsValid())
	//symbol is illegal
	tokenInfo.Symbol = Symbol("Btc")
	assert.False(t, tokenInfo.IsValid())

	tokenInfo.Symbol = Symbol("btc")

	//Decimal is illegal
	for i := 1; i <= 18; i++ {
		tokenInfo.Decimals = uint64(i)
		assert.True(t, tokenInfo.IsValid())
	}
	tokenInfo.Decimals = 19
	assert.False(t, tokenInfo.IsValid())

	tokenInfo.Decimals = 18
	assert.True(t, tokenInfo.IsValid())

}
