package cli

import (
	"fmt"

	"github.com/pocblockchain/pocc/client/context"
	"github.com/pocblockchain/pocc/client/flags"
	"github.com/pocblockchain/pocc/codec"
	"github.com/pocblockchain/pocc/x/token/types"
	"github.com/spf13/cobra"
)

//GetQueryCmd ...
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	tokenQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the token module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		//	RunE:                       utils.ValidateCmd,
	}
	tokenQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryToken(cdc),
		GetCmdQuerySymbols(cdc),
		GetCmdQueryTokens(cdc),
		GetCmdQueryParams(cdc),
	)...)
	return tokenQueryCmd
}

//GetCmdQueryToken ...
func GetCmdQueryToken(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "token [symbol]",
		Short: "token symbol",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			symbol := args[0]

			bz, err := cdc.MarshalJSON(types.QueryTokenInfo{symbol})
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryToken)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var out types.QueryResToken
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

//GetCmdQueryTokens ...
func GetCmdQueryTokens(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "tokens",
		Short: "tokens",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryTokens)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var out types.QueryResTokens
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

//GetCmdQuerySymbols return symbols regiserted in
func GetCmdQuerySymbols(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "symbols",
		Short: "symbols",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QuerySymbols)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var out types.QueryResSymbols
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdQueryParams implements a command to return the current token
// parameters.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current token parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryParameters)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := cdc.UnmarshalJSON(res, &params); err != nil {
				return err
			}

			return cliCtx.PrintOutput(params)
		},
	}
}
