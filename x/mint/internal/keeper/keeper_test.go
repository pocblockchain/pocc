package keeper

import (
	"testing"

	"github.com/pocblockchain/pocc/x/mint/internal/types"
	"github.com/stretchr/testify/require"
)

func TestMinterSetGet(t *testing.T) {
	input := setupTestInput(t)
	keeper := input.mintKeeper
	ctx := input.ctx

	minter := types.NewMinter(0, types.DefaultInitalInflationAmount)
	keeper.SetMinter(ctx, minter)

	require.Equal(t, minter, keeper.GetMinter(ctx))
}
