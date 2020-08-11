/*
 * *******************************************************************
 * @项目名称: token
 * @文件名称: genesis_test.go
 * @Date: 2019/06/14
 * @Author: Keep
 * @Copyright（C）: 2019 BlueHelix Inc.   All rights reserved.
 * 注意：本内容仅限于内部传阅，禁止外泄以及用于其他的商业目的.
 * *******************************************************************
 */
package token

import (
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGenesisState(t *testing.T) {
	tsis := make([]sdk.TokenInfoWithoutSupply, 0)

	for _, ti := range TestTokenData {
		tsi := castToTokenInfoWithoutSupply(ti)
		tsis = append(tsis, tsi)
	}

	genState := NewGenesisState(tsis, DefaultParams())
	err := ValidateGenesis(genState)
	assert.Nil(t, err)

	sorted := sortGenesisTokenInfoWithoutSupply(tsis)
	for i, info := range genState.GenesisTokenInfos {
		assert.Equal(t, sorted[i], info)
	}
	assert.Equal(t, DefaultParams(), genState.Params)
}

func TestDefaultGensisState(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	keeper := input.tokenKeeper

	InitGenesis(ctx, keeper, DefaultGenesisState())
	genState1 := ExportGenesis(ctx, keeper)
	assert.Equal(t, DefaultGenesisState(), genState1)

	got := keeper.GetTokenInfoWithoutSupply(ctx, DefaultGenesisState().GenesisTokenInfos[0].Symbol)
	assert.Equal(t, DefaultGenesisState().GenesisTokenInfos[0], *got)

}

func TestIllegalGensisState(t *testing.T) {
	//Illegal symbol
	ti := []sdk.TokenInfoWithoutSupply{
		{
			Symbol:   "Illegal",
			Decimals: 18,
		},
	}
	genState := NewGenesisState(ti, DefaultParams())
	err := ValidateGenesis(genState)
	assert.NotNil(t, err)

	//decimal > Precision
	ti = []sdk.TokenInfoWithoutSupply{
		{
			Symbol:   "poc",
			Decimals: 19,
		},
	}
	genState = NewGenesisState(ti, DefaultParams())
	err = ValidateGenesis(genState)
	assert.NotNil(t, err)

	//duplicated symbol
	ti = []sdk.TokenInfoWithoutSupply{
		{
			Symbol:   "poc",
			Decimals: 18,
		},
		{
			Symbol:   "btc",
			Decimals: 8,
		},
		{
			Symbol:   "poc",
			Decimals: 8,
		},
	}

	genState = NewGenesisState(ti, DefaultParams())
	err = ValidateGenesis(genState)
	assert.NotNil(t, err)

}

func TestSortGenTokenInfo(t *testing.T) {
	ti := []sdk.TokenInfoWithoutSupply{
		{
			Symbol:   "poc",
			Decimals: 18,
		},
		{
			Symbol:   "btc",
			Decimals: 8,
		},
		{
			Symbol:   "xrp",
			Decimals: 8,
		},
	}

	genState := NewGenesisState(ti, DefaultParams())
	err := ValidateGenesis(genState)
	assert.Nil(t, err)

	tis := genState.GenesisTokenInfos
	assert.Equal(t, "btc", tis[0].Symbol.String())
	assert.Equal(t, "poc", tis[1].Symbol.String())
	assert.Equal(t, "xrp", tis[2].Symbol.String())
}

func TestDefaultGenesisStateMarshal(t *testing.T) {
	defaulGenState := DefaultGenesisState()
	bz, err := ModuleCdc.MarshalJSON(defaulGenState)
	assert.Nil(t, err)

	var gotGenState GenesisState
	err = ModuleCdc.UnmarshalJSON(bz, &gotGenState)
	assert.Nil(t, err)
	assert.True(t, defaulGenState.Equal(gotGenState))
}

func TestAddTokenInfoWithoutSupplyIntoGenesis(t *testing.T) {
	defaulGenState := DefaultGenesisState()
	p := &defaulGenState

	beforeLen := len(p.GenesisTokenInfos)

	eos := sdk.TokenInfoWithoutSupply{
		Symbol:        "eos",
		Issuer:        "",
		IsSendEnabled: false,
		Decimals:      18,
	}

	err := p.AddTokenInfoWithoutSupplyIntoGenesis(eos)
	assert.Nil(t, err)
	assert.Equal(t, beforeLen+1, len(p.GenesisTokenInfos))

	//
	err = p.AddTokenInfoWithoutSupplyIntoGenesis(eos)
	assert.NotNil(t, err)
	assert.Equal(t, beforeLen+1, len(p.GenesisTokenInfos))
}
