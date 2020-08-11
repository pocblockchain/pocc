package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/pocblockchain/pocc/types"
)

func TestMsgUnjailGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("abcd")
	msg := NewMsgUnjail(sdk.ValAddress(addr))
	bytes := msg.GetSignBytes()
	require.Equal(
		t,
		`{"type":"poc/MsgUnjail","value":{"address":"pocvaloper1v93xxeqjh3q7h"}}`,
		string(bytes),
	)
}
