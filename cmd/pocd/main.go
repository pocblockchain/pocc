package main

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/pocblockchain/pocc/baseapp"
	"github.com/pocblockchain/pocc/client"
	"github.com/pocblockchain/pocc/pocapp"
	"github.com/pocblockchain/pocc/server"
	"github.com/pocblockchain/pocc/store"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/genaccounts"
	genaccscli "github.com/pocblockchain/pocc/x/genaccounts/client/cli"
	genutilcli "github.com/pocblockchain/pocc/x/genutil/client/cli"
	"github.com/pocblockchain/pocc/x/staking"
)

// poc chain custom flags
const (
	flagInvCheckPeriod = "inv-check-period"
)

var invCheckPeriod uint

func main() {
	cdc := pocapp.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "pocd",
		Short:             "pocd Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(genutilcli.InitCmd(ctx, cdc, pocapp.ModuleBasics, pocapp.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, pocapp.DefaultNodeHome))
	//rootCmd.AddCommand(genutilcli.MigrateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(genutilcli.GenTxCmd(ctx, cdc, pocapp.ModuleBasics, staking.AppModuleBasic{},
		genaccounts.AppModuleBasic{}, pocapp.DefaultNodeHome, pocapp.DefaultCLIHome))
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc, pocapp.ModuleBasics))
	rootCmd.AddCommand(genaccscli.AddGenesisAccountCmd(ctx, cdc, pocapp.DefaultNodeHome, pocapp.DefaultCLIHome))
	rootCmd.AddCommand(client.NewCompletionCmd(rootCmd, true))
	rootCmd.AddCommand(testnetCmd(ctx, cdc, pocapp.ModuleBasics, genaccounts.AppModuleBasic{}))
	rootCmd.AddCommand(mainnetCmd(ctx, cdc, pocapp.ModuleBasics, genaccounts.AppModuleBasic{}))
	rootCmd.AddCommand(replayCmd())

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "BH", pocapp.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return pocapp.NewPocApp(
		logger, db, traceStore, true, invCheckPeriod,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(uint64(viper.GetInt(server.FlagHaltHeight))),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		gApp := pocapp.NewPocApp(logger, db, traceStore, false, uint(1))
		err := gApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return gApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}
	gApp := pocapp.NewPocApp(logger, db, traceStore, true, uint(1))
	return gApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
