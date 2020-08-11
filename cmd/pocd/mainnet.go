package main

// DONTCOVER

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/pocblockchain/pocc/client"
	"github.com/pocblockchain/pocc/client/keys"
	"github.com/pocblockchain/pocc/codec"
	"github.com/pocblockchain/pocc/server"
	srvconfig "github.com/pocblockchain/pocc/server/config"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/types/module"
	"github.com/pocblockchain/pocc/x/auth"
	"github.com/pocblockchain/pocc/x/genaccounts"
	"github.com/pocblockchain/pocc/x/genutil"
	genutiltypes "github.com/pocblockchain/pocc/x/genutil/types"
	"github.com/pocblockchain/pocc/x/staking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	"os"
	"path/filepath"
)

// get cmd to initialize all files for tendermint testnet and application
var genesisValidators = []string{
	"POCWorld",
	"Sheraton",
	"Rambo",
	"Diamond",
	"WeekEight",
	"NoellCapital",
	"POCFamily",
}

func mainnetCmd(ctx *server.Context, cdc *codec.Codec,
	mbm module.BasicManager, genAccIterator genutiltypes.GenesisAccountsIterator,
) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "mainnet",
		Short: "Initialize files for a poc chain mainnet",
		Long: `miannet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:
	pocd mainnet --output-dir ./output --starting-ip-address 192.168.0.2
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			config := ctx.Config

			outputDir := viper.GetString(flagOutputDir)
			chainID := viper.GetString(client.FlagChainID)
			minGasPrices := viper.GetString(server.FlagMinGasPrices)
			nodeDirPrefix := viper.GetString(flagNodeDirPrefix)
			nodeDaemonHome := viper.GetString(flagNodeDaemonHome)
			nodeCLIHome := viper.GetString(flagNodeCLIHome)
			startingIPAddress := viper.GetString(flagStartingIPAddress)
			initalAccountAmount := viper.GetString(flagInitalAccountAmount)
			initalStakingAmount := viper.GetString(flagInitalStakingAmount)

			if viper.GetBool(flagSameIPAddress) {
				panic(1)
			} else {
				return InitMainnet(cmd, config, cdc, mbm, genAccIterator, outputDir, chainID,
					minGasPrices, nodeDirPrefix, nodeDaemonHome, nodeCLIHome, startingIPAddress, initalAccountAmount, initalStakingAmount, len(genesisValidators))
			}
		},
	}

	cmd.Flags().StringP(flagOutputDir, "o", "./mainnet",
		"Directory to store initialization data for the main")
	cmd.Flags().String(flagNodeDirPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "pocd",
		"Home directory of the node's daemon configuration")
	cmd.Flags().String(flagNodeCLIHome, "poccli",
		"Home directory of the node's cli configuration")
	cmd.Flags().String(flagStartingIPAddress, "192.168.0.2",
		"Starting IP address (192.168.0.2 results in persistent peers list ID0@192.168.0.2:26656, ID1@192.168.0.3:26656, ...)")
	cmd.Flags().String(
		client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(
		server.FlagMinGasPrices, fmt.Sprintf("10000000000%s", sdk.DefaultBondDenom),
		"Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 1x10^-8poc)")
	cmd.Flags().Bool(flagSameIPAddress, false, "All nodes use same ip with different port")

	cmd.Flags().String(
		flagInitalAccountAmount, fmt.Sprintf("1500000000000000000000000"),
		"Each validator's inital account balance")
	cmd.Flags().String(
		flagInitalStakingAmount, fmt.Sprintf("90000000000000000000000"),
		"Each validator's staking in creating validator")

	return cmd
}

// Initialize the testnet
func InitMainnet(cmd *cobra.Command, config *tmconfig.Config, cdc *codec.Codec,
	mbm module.BasicManager, genAccIterator genutiltypes.GenesisAccountsIterator,
	outputDir, chainID, minGasPrices, nodeDirPrefix, nodeDaemonHome,
	nodeCLIHome, startingIPAddress string, initAccountAmount, initStakingAmount string, numValidators int) error {

	if chainID == "" {
		chainID = "chain-" + cmn.RandStr(6)
	}

	monikers := make([]string, numValidators)
	nodeIDs := make([]string, numValidators)
	valPubKeys := make([]crypto.PubKey, numValidators)

	bhConfig := srvconfig.DefaultConfig()
	bhConfig.MinGasPrices = minGasPrices

	var (
		accs     []genaccounts.GenesisAccount
		genFiles []string
	)

	// generate private keys, node IDs, and initial transactions
	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		//nodeDirName := genesisValidators[i]
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		clientDir := filepath.Join(outputDir, nodeDirName, nodeCLIHome)
		gentxsDir := filepath.Join(outputDir, "gentxs")

		config.SetRoot(nodeDir)
		config.RPC.ListenAddress = "tcp://0.0.0.0:26657"

		if err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		if err := os.MkdirAll(clientDir, nodeDirPerm); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		monikers = append(monikers, nodeDirName)
		config.Moniker = nodeDirName

		ip := ""
		var err error
		if startingIPAddress == "0.0.0.0" {
			buf := bufio.NewReader(cmd.InOrStdin())
			promt := fmt.Sprintf(
				"ip for validator '%s' (default %s):", nodeDirName, "0.0.0.0",
			)
			ip, err = client.GetString(promt, buf)
		} else {
			ip, err = getIP(i, startingIPAddress)
			if err != nil {
				_ = os.RemoveAll(outputDir)
				return err
			}
		}

		nodeIDs[i], valPubKeys[i], err = genutil.InitializeNodeValidatorFiles(config)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		memo := fmt.Sprintf("%s@%s:26656", nodeIDs[i], ip)
		genFiles = append(genFiles, config.GenesisFile())

		buf := bufio.NewReader(cmd.InOrStdin())
		prompt := fmt.Sprintf(
			"Password for account '%s' (default %s):", nodeDirName, client.DefaultKeyPass,
		)

		keyPass, err := client.GetPassword(prompt, buf)
		if err != nil && keyPass != "" {
			// An error was returned that either failed to read the password from
			// STDIN or the given password is not empty but failed to meet minimum
			// length requirements.
			return err
		}

		if keyPass == "" {
			keyPass = client.DefaultKeyPass
		}
		addr, secret, err := server.GenerateSaveCoinKey(clientDir, nodeDirName, keyPass, true)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		info := map[string]string{"secret": secret}

		cliPrint, err := json.Marshal(info)
		if err != nil {
			return err
		}

		// save private key seed words
		if err := writeFile(fmt.Sprintf("%v.json", "key_seed"), clientDir, cliPrint); err != nil {
			return err
		}

		accStakingTokens := sdk.TokensFromConsensusPower(90000)
		accStakingTokens, ok := sdk.NewIntFromString(initStakingAmount) //sdk.TokensFromConsensusPower(500)
		if !ok {
			return fmt.Errorf("fail to parse init staking amount:%v", err)
		}

		valTokens := sdk.TokensFromConsensusPower(1500000)
		valTokens, ok = sdk.NewIntFromString(initAccountAmount) // sdk.TokensFromConsensusPower(100)
		if !ok {
			return fmt.Errorf("fail to parse init account amount:%v", err)
		}

		if valTokens.LT(accStakingTokens) {
			return fmt.Errorf("insufficient coins to create validator")
		}

		accs = append(accs, genaccounts.GenesisAccount{
			Address: addr,
			Coins: sdk.Coins{
				sdk.NewCoin(sdk.DefaultBondDenom, valTokens),
			},
		})

		msg := staking.NewMsgCreateValidator(
			sdk.ValAddress(addr),
			valPubKeys[i],
			sdk.NewCoin(sdk.DefaultBondDenom, accStakingTokens),
			staking.NewDescription(genesisValidators[i], "", "", ""),
			staking.NewCommissionRates(sdk.NewDecWithPrec(5, 2), sdk.NewDecWithPrec(5, 2), sdk.ZeroDec()),
			sdk.OneInt(),
		)
		kb, err := keys.NewKeyBaseFromDir(clientDir)
		if err != nil {
			return err
		}
		tx := auth.NewStdTx([]sdk.Msg{msg}, auth.StdFee{}, []auth.StdSignature{}, memo)
		txBldr := auth.NewTxBuilderFromCLI().WithChainID(chainID).WithMemo(memo).WithKeybase(kb)

		signedTx, err := txBldr.SignStdTx(nodeDirName, client.DefaultKeyPass, tx, false)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		txBytes, err := cdc.MarshalJSON(signedTx)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		// gather gentxs folder
		if err := writeFile(fmt.Sprintf("%v.json", nodeDirName), gentxsDir, txBytes); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		bhConfigFilePath := filepath.Join(nodeDir, "config/pocchain.toml")
		srvconfig.WriteConfigFile(bhConfigFilePath, bhConfig)
	}

	if err := initGenFiles(cdc, mbm, chainID, accs, genFiles, numValidators); err != nil {
		return err
	}

	err := collectGenFiles(
		cdc, config, chainID, monikers, nodeIDs, valPubKeys, numValidators,
		outputDir, nodeDirPrefix, nodeDaemonHome, genAccIterator,
	)
	if err != nil {
		return err
	}

	cmd.PrintErrf("Successfully initialized %d node directories\n", numValidators)
	return nil

}
