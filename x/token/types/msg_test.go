package types

import (
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgNewTokenRouteAndType(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	var msg = NewMsgNewToken(addr1, addr2, "btc", 8, sdk.NewInt(10000))

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), TypeMsgNewToken)
}

func TestMsgNewTokenValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))

	//from address is empty
	msg := NewMsgNewToken(nil, addr2, "btc", 8, sdk.NewInt(10000))
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	msg = NewMsgNewToken(sdk.AccAddress(nil), addr2, "btc", 8, sdk.NewInt(10000))
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	msg = NewMsgNewToken(sdk.AccAddress{}, addr2, "btc", 8, sdk.NewInt(10000))
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	//to address is empty
	msg = NewMsgNewToken(addr1, nil, "btc", 8, sdk.NewInt(10000))
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	msg = NewMsgNewToken(addr1, sdk.AccAddress(nil), "btc", 8, sdk.NewInt(10000))
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	msg = NewMsgNewToken(addr1, sdk.AccAddress{}, "btc", 8, sdk.NewInt(10000))
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	//too much precision
	msg = NewMsgNewToken(addr1, addr2, "btc", 19, sdk.NewInt(10000))
	require.Equal(t, sdk.CodeTooMuchPrecision, msg.ValidateBasic().Code())

	//total supply is negative
	msg = NewMsgNewToken(addr1, addr2, "btc", 8, sdk.NewInt(-10000))
	require.Equal(t, sdk.CodeInvalidAmount, msg.ValidateBasic().Code())

}

func TestMsgNewTokenGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	msg := NewMsgNewToken(addr1, addr2, "btc", 8, sdk.NewInt(10000))
	res := msg.GetSignBytes()

	expected := `{"type":"poc/token/MsgNewToken","value":{"decimals":"8","from":"poc1veex7mg3y9476","symbol":"btc","to":"poc1w3hssmamea","total_supply":"10000"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgNewTokenGetSigners(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	msg := NewMsgNewToken(addr1, addr2, "btc", 8, sdk.NewInt(10000))
	require.Equal(t, []sdk.AccAddress{addr1}, msg.GetSigners())

}

func TestMsgInflateTokenRouteAndType(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	coins := sdk.NewCoins(sdk.NewCoin("btc", sdk.NewInt(10000)))

	msg := NewMsgInflateToken(addr1, addr2, coins)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), TypeMsgInflateToken)
}

func TestMsgInflateValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	coins := sdk.NewCoins(sdk.NewCoin("btc", sdk.NewInt(10000)))

	//from address is empty
	msg := NewMsgInflateToken(nil, addr2, coins)
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	msg = NewMsgInflateToken(sdk.AccAddress(nil), addr2, coins)
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	msg = NewMsgInflateToken(sdk.AccAddress{}, addr2, coins)
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	//to address is empty
	msg = NewMsgInflateToken(addr1, nil, coins)
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	msg = NewMsgInflateToken(addr1, sdk.AccAddress(nil), coins)
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	msg = NewMsgInflateToken(addr1, sdk.AccAddress{}, coins)
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	msg = NewMsgInflateToken(addr1, addr2, coins.Add(sdk.NewCoins(sdk.NewCoin("eth", sdk.NewInt(1000)))))
	require.Equal(t, sdk.CodeInvalidTx, msg.ValidateBasic().Code())

	msg = NewMsgInflateToken(addr1, addr2, sdk.NewCoins(sdk.Coin{"eth", sdk.NewInt(0)}))
	require.Equal(t, sdk.CodeInvalidTx, msg.ValidateBasic().Code())

}

func TestMsgInflateTokenGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	coins := sdk.NewCoins(sdk.NewCoin("btc", sdk.NewInt(10000)))
	msg := NewMsgInflateToken(addr1, addr2, coins)

	res := msg.GetSignBytes()

	expected := `{"type":"poc/token/MsgInflateToken","value":{"amount":[{"amount":"10000","denom":"btc"}],"from":"poc1veex7mg3y9476","to":"poc1w3hssmamea"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgInflateTokenGetSigners(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	coins := sdk.NewCoins(sdk.NewCoin("btc", sdk.NewInt(10000)))
	msg := NewMsgInflateToken(addr1, addr2, coins)

	require.Equal(t, []sdk.AccAddress{addr1}, msg.GetSigners())
}

func TestMsgBurnTokenRouteAndType(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	coins := sdk.NewCoins(sdk.NewCoin("btc", sdk.NewInt(10000)))
	msg := NewMsgBurnToken(addr1, coins)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), TypeMsgBurnToken)
}

func TestMsgBurnTokenValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	coins := sdk.NewCoins(sdk.NewCoin("btc", sdk.NewInt(10000)))

	//from address is empty
	msg := NewMsgBurnToken(nil, coins)
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	msg = NewMsgBurnToken(sdk.AccAddress(nil), coins)
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	msg = NewMsgBurnToken(sdk.AccAddress{}, coins)
	require.Equal(t, sdk.CodeInvalidAddress, msg.ValidateBasic().Code())

	msg = NewMsgBurnToken(addr1, sdk.NewCoins(sdk.NewCoin("btc", sdk.NewInt(0))))
	require.Equal(t, sdk.CodeInvalidCoins, msg.ValidateBasic().Code())
}

func TestMsgBurnTokenGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	coins := sdk.NewCoins(sdk.NewCoin("btc", sdk.NewInt(10000)))
	msg := NewMsgBurnToken(addr1, coins)

	res := msg.GetSignBytes()

	expected := `{"type":"poc/token/MsgBurnToken","value":{"amount":[{"amount":"10000","denom":"btc"}],"from_address":"poc1veex7mg3y9476"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgBrunTokenGetSigners(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	coins := sdk.NewCoins(sdk.NewCoin("btc", sdk.NewInt(10000)))
	msg := NewMsgBurnToken(addr1, coins)

	require.Equal(t, []sdk.AccAddress{addr1}, msg.GetSigners())
}
