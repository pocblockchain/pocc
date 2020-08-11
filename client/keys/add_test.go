package keys

import (
	"github.com/pocblockchain/pocc/crypto/keys"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/pocblockchain/pocc/client/flags"
	"github.com/pocblockchain/pocc/tests"
)

func Test_runAddCmdBasic(t *testing.T) {
	cmd := addKeyCommand()
	assert.NotNil(t, cmd)
	mockIn, _, _ := tests.ApplyMockIO(cmd)

	kbHome, kbCleanUp := tests.NewTestCaseDir(t)
	assert.NotNil(t, kbHome)
	defer kbCleanUp()
	viper.Set(flags.FlagHome, kbHome)

	viper.Set(cli.OutputFlag, OutputFormatText)

	mockIn.Reset("test1234\ntest1234\n")
	err := runAddCmd(cmd, []string{"keyname1"})
	assert.NoError(t, err)

	viper.Set(cli.OutputFlag, OutputFormatText)

	mockIn.Reset("test1234\ntest1234\n")
	err = runAddCmd(cmd, []string{"keyname1"})
	assert.Error(t, err)

	viper.Set(cli.OutputFlag, OutputFormatText)

	mockIn.Reset("y\ntest1234\ntest1234\n")
	err = runAddCmd(cmd, []string{"keyname1"})
	assert.NoError(t, err)

	viper.Set(cli.OutputFlag, OutputFormatJSON)

	mockIn.Reset("test1234\ntest1234\n")
	err = runAddCmd(cmd, []string{"keyname2"})
	assert.NoError(t, err)
}

func TestExportPrivKey(t *testing.T) {
	kb := keys.NewInMemory()

	mnemonic := "lounge degree orphan snap fox prefer rail jealous rebuild flock mistake spell put know peace skate game laptop mixture amount unique fan sound gallery"

	info, err := kb.CreateAccount("alice", mnemonic, "", "12345678", 0, 0)
	assert.Nil(t, err)
	//t.Logf("info:%v", info)
	t.Logf("addr:%s", info.GetAddress())
}
