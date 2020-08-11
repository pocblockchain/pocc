package test

import (
	"github.com/pocblockchain/pocc/codec"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/token/client/cli"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTokenParamsChangeProposalJSON(t *testing.T) {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	cp, err := cli.ParseTokenParamsChangeProposalJSON(cdc, "token_params_change_proposal.json")
	assert.Nil(t, err)
	assert.Equal(t, "Token Parameters Change", cp.Title)
	assert.Equal(t, "token parameter change proposal", cp.Description)
	assert.Equal(t, "testtoken", cp.Symbol)
	assert.Equal(t, sdk.NewCoins(sdk.NewCoin("hbc", sdk.NewInt(10000))), cp.Deposit)
	changes := cp.Changes.ToParamChanges()
	assert.Equal(t, sdk.KeyIsSendEnabled, changes[0].Key)
	assert.Equal(t, "true", changes[0].Value)
}

func TestParseDisableTokenProposalJSON(t *testing.T) {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	dp, err := cli.ParseDisableTokenProposalJSON(cdc, "disable_token_proposal.json")
	assert.Nil(t, err)
	assert.Equal(t, "Disable Token", dp.Title)
	assert.Equal(t, "disable token proposal", dp.Description)
	assert.Equal(t, "usdt", dp.Symbol)
	assert.Equal(t, sdk.NewCoins(sdk.NewCoin("hbc", sdk.NewInt(10000))), dp.Deposit)
}
