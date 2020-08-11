/*
 * *******************************************************************
 * @项目名称: token
 * @文件名称: params.go
 * @Date: 2019/08/07
 * @Author: Keep
 * @Copyright（C）: 2019 BlueHelix Inc.   All rights reserved.
 * 注意：本内容仅限于内部传阅，禁止外泄以及用于其他的商业目的.
 * *******************************************************************
 */
package token

import (
	"github.com/pocblockchain/pocc/x/token/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParamsEqual(t *testing.T) {
	p1 := types.DefaultParams()
	p2 := types.DefaultParams()
	assert.Equal(t, p1, p2)

	p1.TokenCacheSize += 10
	assert.NotEqual(t, p1, p2)
}

func TestGetSetParams(t *testing.T) {
	input := setupTestEnv(t)
	ctx := input.ctx
	keeper := input.tokenKeeper

	InitGenesis(ctx, keeper, DefaultGenesisState())
	genState1 := ExportGenesis(ctx, keeper)
	assert.Equal(t, DefaultGenesisState(), genState1)

	param1 := keeper.GetParams(ctx)
	assert.Equal(t, types.DefaultParams(), param1)

	//change tokenCacheSize
	param2 := param1
	param2.TokenCacheSize += 10

	keeper.SetParams(ctx, param2)
	param2Get := keeper.GetParams(ctx)
	assert.Equal(t, param2, param2Get)
	assert.True(t, param2.Equal(param2Get))
	assert.NotEqual(t, param1, param2Get)

	//append a new symbol
	param3 := param2
	param3.ReservedSymbols = append(param3.ReservedSymbols, "mytoken")
	keeper.SetParams(ctx, param3)
	param3Get := keeper.GetParams(ctx)
	assert.Equal(t, param3, param3Get)
	assert.True(t, param3.Equal(param3Get))
	assert.NotEqual(t, param2, param3Get)
	t.Logf("%v", param3.ReservedSymbols)

}
