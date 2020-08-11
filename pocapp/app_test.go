package pocapp

import (
	"github.com/pocblockchain/pocc/x/auth"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/pocblockchain/pocc/codec"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestPocAppExport(t *testing.T) {
	db := dbm.NewMemDB()
	app := NewPocApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)

	genesisState := NewDefaultGenesisState()
	stateBytes, err := codec.MarshalJSONIndent(app.cdc, genesisState)
	require.NoError(t, err)

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	app2 := NewPocApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	_, _, err = app2.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")

}

// ensure that black listed addresses are properly set in bank keeper
func TestBlackListedAddrs(t *testing.T) {
	db := dbm.NewMemDB()
	app := NewPocApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)

	for acc := range maccPerms {
		require.True(t, app.bankKeeper.BlacklistedAddr(app.supplyKeeper.GetModuleAddress(acc)))
	}
}

func TestPocAppGenesis(t *testing.T) {
	db := dbm.NewMemDB()
	app := NewPocApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)

	genesisState := NewDefaultGenesisState()
	stateBytes, err := codec.MarshalJSONIndent(app.cdc, genesisState)
	require.NoError(t, err)

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()

	header := abci.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := app.BaseApp.NewContext(false, abci.Header{})

	minter := app.mintKeeper.GetMinter(ctx)
	t.Logf("minter:%v", minter)
}

func TestStdSignedMarshal(t *testing.T) {
	db := dbm.NewMemDB()
	app := NewPocApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)

	bz := []byte(`{"type":"poc/StdTx","value":{"msg":[{"type":"poc/MsgSend","value":{"from_address":"poc1uhddf9ncrrxs6qfl64cqz8kj4ez8rdhgpn5a8j","to_address":"poc1mq3cnrrm6cluj06u05vjxgy2snwy8jeajncx5a","amount":[{"denom":"poc","amount":"1000000000000000000"}]}}],"fee":{"amount":[{"denom":"poc","amount":"2000000000000000000"}],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"AusYTPffSVKqhtaPoMoa0AeGZpfcWFVP88rvY074JNQI"},"signature":"B+2vc4HKp3bim5WxurSVfSFGDO3w7yarYSftIFGMFUZOVrDjwvRYb6Ml4Ecu8KKOQH6wMlMxzF1Q5ypblnuLuA=="}],"memo":""}}`)

	tx := auth.StdTx{}

	app.cdc.MustUnmarshalJSON(bz, &tx)
	t.Logf("tx :%v", tx)

	for _, msg := range tx.Msgs {
		t.Logf("msg:%+v", msg)
		err := msg.ValidateBasic()
		require.Nil(t, err)
	}
}

func TestStdSignedMarshal1(t *testing.T) {
	db := dbm.NewMemDB()
	app := NewPocApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)

	bz := []byte(`{"type":"poc/StdTx","value":{"msg":[{"type":"poc/MsgSend","value":{"from_address":"poc1uhddf9ncrrxs6qfl64cqz8kj4ez8rdhgpn5a8j","to_address":"poc1mq3cnrrm6cluj06u05vjxgy2snwy8jeajncx5a","amount":[{"denom":"poc","amount":"1000000000000000000"}]}}],"fee":{"amount":[{"denom":"poc","amount":"2000000000000000000"}],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"AusYTPffSVKqhtaPoMoa0AeGZpfcWFVP88rvY074JNQI"},"signature":"PFQkIFwfIHJTsbJxqPEsOB/sMLLsQdEl/zxj//mh36Z76unpMd7YA52/Mrt9CszV+aWuXSrrS4KbWEGS1Tv/lQ=="}],"memo":""}}`)

	tx := auth.StdTx{}

	app.cdc.MustUnmarshalJSON(bz, &tx)
	t.Logf("tx :%v", tx)

	for _, msg := range tx.Msgs {
		t.Logf("msg:%+v", msg)
		err := msg.ValidateBasic()
		require.Nil(t, err)
	}
}
