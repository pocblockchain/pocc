package bank_test

import (
	"testing"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/auth"
	"github.com/pocblockchain/pocc/x/bank/internal/types"
	"github.com/pocblockchain/pocc/x/mock"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type (
	expectedBalance struct {
		addr  sdk.AccAddress
		coins sdk.Coins
	}

	appTestCase struct {
		expSimPass       bool
		expPass          bool
		msgs             []sdk.Msg
		accNums          []uint64
		accSeqs          []uint64
		privKeys         []crypto.PrivKey
		expectedBalances []expectedBalance
	}
)

var (
	priv1 = secp256k1.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	priv2 = secp256k1.GenPrivKey()
	addr2 = sdk.AccAddress(priv2.PubKey().Address())
	addr3 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	priv4 = secp256k1.GenPrivKey()
	addr4 = sdk.AccAddress(priv4.PubKey().Address())

	coins     = sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}
	halfCoins = sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}
	manyCoins = sdk.Coins{sdk.NewInt64Coin("foocoin", 1), sdk.NewInt64Coin("barcoin", 1)}
	freeFee   = auth.NewStdFee(100000, sdk.Coins{sdk.NewInt64Coin("foocoin", 0)})

	sendMsg1 = types.NewMsgSend(addr1, addr2, coins)
	sendMsg2 = types.NewMsgSend(addr1, moduleAccAddr, coins)

	multiSendMsg1 = types.MsgMultiSend{
		Inputs:  []types.Input{types.NewInput(addr1, coins)},
		Outputs: []types.Output{types.NewOutput(addr2, coins)},
	}
	multiSendMsg2 = types.MsgMultiSend{
		Inputs: []types.Input{types.NewInput(addr1, coins)},
		Outputs: []types.Output{
			types.NewOutput(addr2, halfCoins),
			types.NewOutput(addr3, halfCoins),
		},
	}
	multiSendMsg3 = types.MsgMultiSend{
		Inputs: []types.Input{
			types.NewInput(addr1, coins),
			types.NewInput(addr4, coins),
		},
		Outputs: []types.Output{
			types.NewOutput(addr2, coins),
			types.NewOutput(addr3, coins),
		},
	}
	multiSendMsg4 = types.MsgMultiSend{
		Inputs: []types.Input{
			types.NewInput(addr2, coins),
		},
		Outputs: []types.Output{
			types.NewOutput(addr1, coins),
		},
	}
	multiSendMsg5 = types.MsgMultiSend{
		Inputs: []types.Input{
			types.NewInput(addr1, manyCoins),
		},
		Outputs: []types.Output{
			types.NewOutput(addr2, manyCoins),
		},
	}
	multiSendMsg6 = types.MsgMultiSend{
		Inputs: []types.Input{
			types.NewInput(addr1, coins),
		},
		Outputs: []types.Output{
			types.NewOutput(moduleAccAddr, coins),
		},
	}

	// BonusSend
	bonusSendMsg1 = types.MsgBonusSend{
		Inputs:  []types.Input{types.NewInput(addr1, coins)},
		Outputs: []types.Output{types.NewOutput(addr2, coins)},
	}

	bonusSendMsg2 = types.MsgBonusSend{
		Inputs: []types.Input{types.NewInput(addr1, coins)},
		Outputs: []types.Output{
			types.NewOutput(addr2, halfCoins),
			types.NewOutput(addr3, halfCoins),
		},
	}
	bonusSendMsg3 = types.MsgBonusSend{
		Inputs: []types.Input{
			types.NewInput(addr1, coins),
			types.NewInput(addr4, coins),
		},
		Outputs: []types.Output{
			types.NewOutput(addr2, coins),
			types.NewOutput(addr3, coins),
		},
	}
	bonusSendMsg4 = types.MsgBonusSend{
		Inputs: []types.Input{
			types.NewInput(addr2, coins),
		},
		Outputs: []types.Output{
			types.NewOutput(addr1, coins),
		},
	}
	bonusSendMsg5 = types.MsgBonusSend{
		Inputs: []types.Input{
			types.NewInput(addr1, manyCoins),
		},
		Outputs: []types.Output{
			types.NewOutput(addr2, manyCoins),
		},
	}
	bonusSendMsg6 = types.MsgBonusSend{
		Inputs: []types.Input{
			types.NewInput(addr1, coins),
		},
		Outputs: []types.Output{
			types.NewOutput(moduleAccAddr, coins),
		},
	}

	// ReclaimSend
	reclaimSendMsg1 = types.MsgReclaimSend{
		Inputs:  []types.Input{types.NewInput(addr1, coins)},
		Outputs: []types.Output{types.NewOutput(addr2, coins)},
	}

	reclaimSendMsg2 = types.MsgReclaimSend{
		Inputs: []types.Input{types.NewInput(addr1, coins)},
		Outputs: []types.Output{
			types.NewOutput(addr2, halfCoins),
			types.NewOutput(addr3, halfCoins),
		},
	}
	reclaimSendMsg3 = types.MsgReclaimSend{
		Inputs: []types.Input{
			types.NewInput(addr1, coins),
			types.NewInput(addr4, coins),
		},
		Outputs: []types.Output{
			types.NewOutput(addr2, coins),
			types.NewOutput(addr3, coins),
		},
	}
	reclaimSendMsg4 = types.MsgReclaimSend{
		Inputs: []types.Input{
			types.NewInput(addr2, coins),
		},
		Outputs: []types.Output{
			types.NewOutput(addr1, coins),
		},
	}
	reclaimSendMsg5 = types.MsgReclaimSend{
		Inputs: []types.Input{
			types.NewInput(addr1, manyCoins),
		},
		Outputs: []types.Output{
			types.NewOutput(addr2, manyCoins),
		},
	}
	reclaimSendMsg6 = types.MsgReclaimSend{
		Inputs: []types.Input{
			types.NewInput(addr1, coins),
		},
		Outputs: []types.Output{
			types.NewOutput(moduleAccAddr, coins),
		},
	}
)

func TestSendNotEnoughBalance(t *testing.T) {
	input := getMockApp(t)
	mapp := input.mApp
	tk := input.tk
	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})
	acc := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin("foocoin", 67)},
	}

	mock.SetGenesis(mapp, []auth.Account{acc})

	res1 := mapp.AccountKeeper.GetAccount(ctxCheck, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*auth.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	require.True(t, tk.IsTokenSupported(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsSendEnabled(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsTokenSupported(ctxCheck, "foocoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "foocoin"))
	require.True(t, tk.IsTokenSupported(ctxCheck, "barcoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "barcoin"))
	sendMsg := types.NewMsgSend(addr1, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 100)})
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, header, []sdk.Msg{sendMsg}, []uint64{origAccNum}, []uint64{origSeq}, false, false, priv1)

	mock.CheckBalance(t, mapp, addr1, sdk.Coins{sdk.NewInt64Coin("foocoin", 67)})

	res2 := mapp.AccountKeeper.GetAccount(mapp.NewContext(true, abci.Header{}), addr1)
	require.NotNil(t, res2)

	require.True(t, res2.GetAccountNumber() == origAccNum)
	require.True(t, res2.GetSequence() == origSeq+1)
}

func TestSendToModuleAcc(t *testing.T) {
	input := getMockApp(t)
	mapp := input.mApp
	tk := input.tk
	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})
	acc := &auth.BaseAccount{
		Address: addr1,
		Coins:   coins,
	}

	macc := &auth.BaseAccount{
		Address: moduleAccAddr,
	}

	mock.SetGenesis(mapp, []auth.Account{acc, macc})

	res1 := mapp.AccountKeeper.GetAccount(ctxCheck, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*auth.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	require.True(t, tk.IsTokenSupported(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsSendEnabled(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsTokenSupported(ctxCheck, "foocoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "foocoin"))
	require.True(t, tk.IsTokenSupported(ctxCheck, "barcoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "barcoin"))
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, header, []sdk.Msg{sendMsg2}, []uint64{origAccNum}, []uint64{origSeq}, false, false, priv1)

	mock.CheckBalance(t, mapp, addr1, coins)
	mock.CheckBalance(t, mapp, moduleAccAddr, sdk.Coins(nil))

	res2 := mapp.AccountKeeper.GetAccount(mapp.NewContext(true, abci.Header{}), addr1)
	require.NotNil(t, res2)

	require.True(t, res2.GetAccountNumber() == origAccNum)
	require.True(t, res2.GetSequence() == origSeq+1)
}

func TestMsgMultiSendWithAccounts(t *testing.T) {
	input := getMockApp(t)
	mapp := input.mApp
	tk := input.tk
	acc := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin("foocoin", 67)},
	}

	macc := &auth.BaseAccount{
		Address: moduleAccAddr,
	}

	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})
	mock.SetGenesis(mapp, []auth.Account{acc, macc})

	res1 := mapp.AccountKeeper.GetAccount(ctxCheck, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*auth.BaseAccount))

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{multiSendMsg1},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin("foocoin", 57)}},
				{addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}},
			},
		},
		{
			msgs:       []sdk.Msg{multiSendMsg1, multiSendMsg2},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true, // doesn't check signature
			expPass:    false,
			privKeys:   []crypto.PrivKey{priv1},
		},
		{
			msgs:       []sdk.Msg{multiSendMsg6},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: false,
			expPass:    false,
			privKeys:   []crypto.PrivKey{priv1},
		},
	}
	require.True(t, tk.IsTokenSupported(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsSendEnabled(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsTokenSupported(ctxCheck, "foocoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "foocoin"))
	require.True(t, tk.IsTokenSupported(ctxCheck, "barcoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "barcoin"))

	for _, tc := range testCases {
		header := abci.Header{Height: mapp.LastBlockHeight() + 1}
		mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, header, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			mock.CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}
}

func TestMsgMultiSendMultipleOut(t *testing.T) {
	input := getMockApp(t)
	mapp := input.mApp
	tk := input.tk
	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})

	acc1 := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin("foocoin", 42)},
	}
	acc2 := &auth.BaseAccount{
		Address: addr2,
		Coins:   sdk.Coins{sdk.NewInt64Coin("foocoin", 42)},
	}

	mock.SetGenesis(mapp, []auth.Account{acc1, acc2})

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{multiSendMsg2},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin("foocoin", 32)}},
				{addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 47)}},
				{addr3, sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}},
			},
		},
	}

	require.True(t, tk.IsTokenSupported(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsSendEnabled(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsTokenSupported(ctxCheck, "foocoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "foocoin"))
	require.True(t, tk.IsTokenSupported(ctxCheck, "barcoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "barcoin"))
	for _, tc := range testCases {
		header := abci.Header{Height: mapp.LastBlockHeight() + 1}
		mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, header, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			mock.CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}
}

func TestMsgMultiSendMultipleInOut(t *testing.T) {
	input := getMockApp(t)
	mapp := input.mApp
	tk := input.tk
	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})

	acc1 := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin("foocoin", 42)},
	}
	acc2 := &auth.BaseAccount{
		Address: addr2,
		Coins:   sdk.Coins{sdk.NewInt64Coin("foocoin", 42)},
	}
	acc4 := &auth.BaseAccount{
		Address: addr4,
		Coins:   sdk.Coins{sdk.NewInt64Coin("foocoin", 42)},
	}

	mock.SetGenesis(mapp, []auth.Account{acc1, acc2, acc4})

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{multiSendMsg3},
			accNums:    []uint64{0, 2},
			accSeqs:    []uint64{0, 0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1, priv4},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin("foocoin", 32)}},
				{addr4, sdk.Coins{sdk.NewInt64Coin("foocoin", 32)}},
				{addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 52)}},
				{addr3, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}},
			},
		},
	}
	require.True(t, tk.IsTokenSupported(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsSendEnabled(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsTokenSupported(ctxCheck, "foocoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "foocoin"))
	require.True(t, tk.IsTokenSupported(ctxCheck, "barcoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "barcoin"))

	for _, tc := range testCases {
		header := abci.Header{Height: mapp.LastBlockHeight() + 1}
		mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, header, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			mock.CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}
}

func TestMsgMultiSendDependent(t *testing.T) {
	input := getMockApp(t)
	mapp := input.mApp
	tk := input.tk
	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})

	acc1 := auth.NewBaseAccountWithAddress(addr1)
	acc2 := auth.NewBaseAccountWithAddress(addr2)
	err := acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 42)))
	require.NoError(t, err)
	err = acc2.SetAccountNumber(1)
	require.NoError(t, err)

	mock.SetGenesis(mapp, []auth.Account{&acc1, &acc2})

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{multiSendMsg1},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin("foocoin", 32)}},
				{addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}},
			},
		},
		{
			msgs:       []sdk.Msg{multiSendMsg4},
			accNums:    []uint64{1},
			accSeqs:    []uint64{0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv2},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin("foocoin", 42)}},
			},
		},
	}

	require.True(t, tk.IsTokenSupported(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsSendEnabled(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsTokenSupported(ctxCheck, "foocoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "foocoin"))
	require.True(t, tk.IsTokenSupported(ctxCheck, "barcoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "barcoin"))
	for _, tc := range testCases {
		header := abci.Header{Height: mapp.LastBlockHeight() + 1}
		mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, header, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			mock.CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}
}

func TestMsgSendFailTokenNotSupported(t *testing.T) {
	input := getMockApp(t)
	mapp := input.mApp
	tk := input.tk
	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})

	acc1 := auth.NewBaseAccountWithAddress(addr1)
	acc2 := auth.NewBaseAccountWithAddress(addr2)
	err := acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("notsupport", 42)))
	require.NoError(t, err)
	err = acc2.SetAccountNumber(1)
	require.NoError(t, err)

	mock.SetGenesis(mapp, []auth.Account{&acc1, &acc2})

	msg := types.NewMsgSend(addr1, addr2, sdk.Coins{sdk.NewInt64Coin("notsupport", 10)})

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{msg},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: false,
			expPass:    false,
			privKeys:   []crypto.PrivKey{priv1},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin("notsupport", 42)}},
				{addr2, sdk.Coins(nil)},
			},
		},
	}

	require.True(t, tk.IsTokenSupported(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsSendEnabled(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsTokenSupported(ctxCheck, "foocoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "foocoin"))
	require.True(t, tk.IsTokenSupported(ctxCheck, "barcoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "barcoin"))
	require.False(t, tk.IsTokenSupported(ctxCheck, "notsupport"))
	require.False(t, tk.IsSendEnabled(ctxCheck, "notsupport"))

	for _, tc := range testCases {
		header := abci.Header{Height: mapp.LastBlockHeight() + 1}
		mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, header, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			mock.CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}

}

func TestMsgSendWithSetTokenKeeper(t *testing.T) {
	input := getMockAppWithSetTokenKeeper(t)
	mapp := input.mApp
	tk := input.tk
	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})

	acc1 := auth.NewBaseAccountWithAddress(addr1)
	acc2 := auth.NewBaseAccountWithAddress(addr2)
	err := acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 42)))
	require.NoError(t, err)
	err = acc2.SetAccountNumber(1)
	require.NoError(t, err)

	mock.SetGenesis(mapp, []auth.Account{&acc1, &acc2})

	msg := types.NewMsgSend(addr1, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)})

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{msg},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin("foocoin", 32)}},
				{addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}},
			},
		},
	}

	require.True(t, tk.IsTokenSupported(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsSendEnabled(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsTokenSupported(ctxCheck, "foocoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "foocoin"))
	require.True(t, tk.IsTokenSupported(ctxCheck, "barcoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "barcoin"))

	for _, tc := range testCases {
		header := abci.Header{Height: mapp.LastBlockHeight() + 1}
		mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, header, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			mock.CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}
}

func TestEscrowNotEnoughBalance(t *testing.T) {
	input := getMockApp(t)
	mapp := input.mApp
	tk := input.tk
	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})
	acc := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin("foocoin", 67)},
	}

	mock.SetGenesis(mapp, []auth.Account{acc})

	res1 := mapp.AccountKeeper.GetAccount(ctxCheck, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*auth.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	require.True(t, tk.IsTokenSupported(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsSendEnabled(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsTokenSupported(ctxCheck, "foocoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "foocoin"))
	require.True(t, tk.IsTokenSupported(ctxCheck, "barcoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "barcoin"))
	escrowMsg := types.NewMsgEscrow(addr1, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 100)})
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, header, []sdk.Msg{escrowMsg}, []uint64{origAccNum}, []uint64{origSeq}, false, false, priv1)

	mock.CheckBalance(t, mapp, addr1, sdk.Coins{sdk.NewInt64Coin("foocoin", 67)})

	res2 := mapp.AccountKeeper.GetAccount(mapp.NewContext(true, abci.Header{}), addr1)
	require.NotNil(t, res2)

	require.True(t, res2.GetAccountNumber() == origAccNum)
	require.True(t, res2.GetSequence() == origSeq+1)
}

func TestEscrowSuccess(t *testing.T) {
	input := getMockApp(t)
	mapp := input.mApp
	tk := input.tk
	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})
	acc := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin("foocoin", 67)},
	}

	mock.SetGenesis(mapp, []auth.Account{acc})

	res1 := mapp.AccountKeeper.GetAccount(ctxCheck, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*auth.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	require.True(t, tk.IsTokenSupported(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsSendEnabled(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsTokenSupported(ctxCheck, "foocoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "foocoin"))
	require.True(t, tk.IsTokenSupported(ctxCheck, "barcoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "barcoin"))
	escrowMsg := types.NewMsgEscrow(addr1, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 30)})
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, header, []sdk.Msg{escrowMsg}, []uint64{origAccNum}, []uint64{origSeq}, true, true, priv1)

	mock.CheckBalance(t, mapp, addr1, sdk.Coins{sdk.NewInt64Coin("foocoin", 37)})
	mock.CheckBalance(t, mapp, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 30)})

	res2 := mapp.AccountKeeper.GetAccount(mapp.NewContext(true, abci.Header{}), addr1)
	require.NotNil(t, res2)

	require.True(t, res2.GetAccountNumber() == origAccNum)
	require.True(t, res2.GetSequence() == origSeq+1)
}

func TestReclaimSuccess(t *testing.T) {
	input := getMockApp(t)
	mapp := input.mApp
	tk := input.tk
	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})
	acc := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin("foocoin", 67)},
	}

	acc2 := &auth.BaseAccount{
		Address: addr2,
		Coins:   sdk.Coins{sdk.NewInt64Coin("foocoin", 1)},
	}

	mock.SetGenesis(mapp, []auth.Account{acc, acc2})

	res1 := mapp.AccountKeeper.GetAccount(ctxCheck, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*auth.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	require.True(t, tk.IsTokenSupported(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsSendEnabled(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsTokenSupported(ctxCheck, "foocoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "foocoin"))
	require.True(t, tk.IsTokenSupported(ctxCheck, "barcoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "barcoin"))
	relcaimMsg := types.NewMsgReclaim(addr1, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 100)})
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, header, []sdk.Msg{relcaimMsg}, []uint64{origAccNum}, []uint64{origSeq}, true, true, priv1)

	mock.CheckBalance(t, mapp, addr1, sdk.Coins{sdk.NewInt64Coin("foocoin", 67)})
	mock.CheckBalance(t, mapp, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 1)})

	res2 := mapp.AccountKeeper.GetAccount(mapp.NewContext(true, abci.Header{}), addr1)
	require.NotNil(t, res2)

	require.True(t, res2.GetAccountNumber() == origAccNum)
	require.True(t, res2.GetSequence() == origSeq+1)
}

func TestMsgBonusSend(t *testing.T) {
	input := getMockApp(t)
	mapp := input.mApp
	tk := input.tk
	acc := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin("foocoin", 67)},
	}

	macc := &auth.BaseAccount{
		Address: moduleAccAddr,
	}

	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})
	mock.SetGenesis(mapp, []auth.Account{acc, macc})

	res1 := mapp.AccountKeeper.GetAccount(ctxCheck, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*auth.BaseAccount))

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{bonusSendMsg1},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin("foocoin", 57)}},
				{addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}},
			},
		},
		{
			msgs:       []sdk.Msg{bonusSendMsg1, bonusSendMsg2},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true, // doesn't check signature
			expPass:    false,
			privKeys:   []crypto.PrivKey{priv1},
		},
		{
			msgs:       []sdk.Msg{bonusSendMsg6},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: false,
			expPass:    false,
			privKeys:   []crypto.PrivKey{priv1},
		},
	}
	require.True(t, tk.IsTokenSupported(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsSendEnabled(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsTokenSupported(ctxCheck, "foocoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "foocoin"))
	require.True(t, tk.IsTokenSupported(ctxCheck, "barcoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "barcoin"))

	for _, tc := range testCases {
		header := abci.Header{Height: mapp.LastBlockHeight() + 1}
		mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, header, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			mock.CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}
}

func TestMsgReclaimSend(t *testing.T) {
	input := getMockApp(t)
	mapp := input.mApp
	tk := input.tk
	acc := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin("foocoin", 67)},
	}

	macc := &auth.BaseAccount{
		Address: moduleAccAddr,
	}

	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})
	mock.SetGenesis(mapp, []auth.Account{acc, macc})

	res1 := mapp.AccountKeeper.GetAccount(ctxCheck, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*auth.BaseAccount))

	testCases := []appTestCase{
		{
			msgs:       []sdk.Msg{reclaimSendMsg1},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true,
			expPass:    true,
			privKeys:   []crypto.PrivKey{priv1},
			expectedBalances: []expectedBalance{
				{addr1, sdk.Coins{sdk.NewInt64Coin("foocoin", 57)}},
				{addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}},
			},
		},
		{
			msgs:       []sdk.Msg{reclaimSendMsg1, reclaimSendMsg2},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: true, // doesn't check signature
			expPass:    false,
			privKeys:   []crypto.PrivKey{priv1},
		},
		{
			msgs:       []sdk.Msg{reclaimSendMsg6},
			accNums:    []uint64{0},
			accSeqs:    []uint64{0},
			expSimPass: false,
			expPass:    false,
			privKeys:   []crypto.PrivKey{priv1},
		},
	}
	require.True(t, tk.IsTokenSupported(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsSendEnabled(ctxCheck, sdk.NativeToken))
	require.True(t, tk.IsTokenSupported(ctxCheck, "foocoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "foocoin"))
	require.True(t, tk.IsTokenSupported(ctxCheck, "barcoin"))
	require.True(t, tk.IsSendEnabled(ctxCheck, "barcoin"))

	for _, tc := range testCases {
		header := abci.Header{Height: mapp.LastBlockHeight() + 1}
		mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, header, tc.msgs, tc.accNums, tc.accSeqs, tc.expSimPass, tc.expPass, tc.privKeys...)

		for _, eb := range tc.expectedBalances {
			mock.CheckBalance(t, mapp, eb.addr, eb.coins)
		}
	}
}
