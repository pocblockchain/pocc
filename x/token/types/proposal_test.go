package types

import (
	"testing"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/stretchr/testify/require"
)

func TestTokenParamsChangeProposal(t *testing.T) {
	expectedStr := "Change Token Param Proposal:\n Title:       Test\n Description: Description\n Symbol:      btc\n Changes:\nis_send_enabled: true\t"
	changes := []ParamChange{}
	changes = append(changes, NewParamChange(sdk.KeyIsSendEnabled, "true"))
	tpcp := NewTokenParamsChangeProposal("Test", "Description", "btc", changes)

	require.Equal(t, "Test", tpcp.GetTitle())
	require.Equal(t, "Description", tpcp.GetDescription())
	require.Equal(t, RouterKey, tpcp.ProposalRoute())
	require.Equal(t, ProposalTypeTokenParamsChange, tpcp.ProposalType())
	require.Nil(t, tpcp.ValidateBasic())
	require.Equal(t, expectedStr, tpcp.String())

	//symbol is illegal
	tpcp = NewTokenParamsChangeProposal("Test", "Description", "Btc", changes)
	err := tpcp.ValidateBasic()
	require.NotNil(t, err)

	//duplicated keys
	changes = append(changes, NewParamChange(sdk.KeyIsSendEnabled, "false"))
	tpcp = NewTokenParamsChangeProposal("Test", "Description", "btc", changes)
	err = tpcp.ValidateBasic()
	require.NotNil(t, err)

}

func TestDisableTokenProposal(t *testing.T) {
	expectedStr := "Disable Token Proposal:\n Title:       Test\n Description: Description\n Symbol:      btc\n"
	dtp := NewDisableTokenProposal("Test", "Description", "btc")

	require.Equal(t, "Test", dtp.GetTitle())
	require.Equal(t, "Description", dtp.GetDescription())
	require.Equal(t, RouterKey, dtp.ProposalRoute())
	require.Equal(t, ProposalTypeDisableToken, dtp.ProposalType())
	require.Nil(t, dtp.ValidateBasic())
	require.Equal(t, expectedStr, dtp.String())

	dtp = NewDisableTokenProposal("Test", "Description", "Btc")
	require.NotNil(t, dtp.ValidateBasic())
}
