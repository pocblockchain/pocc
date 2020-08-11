package keeper

import (
	"testing"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/mint/internal/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestNewQuerier(t *testing.T) {
	input := setupTestInput(t)
	querier := NewQuerier(input.mintKeeper)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	_, err := querier(input.ctx, []string{types.QueryParameters}, query)
	require.NoError(t, err)

	_, err = querier(input.ctx, []string{types.QueryAnnualProvisions}, query)
	require.NoError(t, err)

	_, err = querier(input.ctx, []string{"foo"}, query)
	require.Error(t, err)
}

func TestQueryParams(t *testing.T) {
	input := setupTestInput(t)

	var params types.Params

	res, sdkErr := queryParams(input.ctx, input.mintKeeper)
	require.NoError(t, sdkErr)

	err := input.cdc.UnmarshalJSON(res, &params)
	require.NoError(t, err)

	require.Equal(t, input.mintKeeper.GetParams(input.ctx), params)
}

func TestQueryAnnualProvisions(t *testing.T) {
	input := setupTestInput(t)

	var annualProvisions sdk.Dec

	res, sdkErr := queryAnnualProvisions(input.ctx, input.mintKeeper)
	require.NoError(t, sdkErr)

	err := input.cdc.UnmarshalJSON(res, &annualProvisions)
	require.NoError(t, err)

	require.Equal(t, input.mintKeeper.GetMinter(input.ctx).AnnualProvisions, annualProvisions)
}
