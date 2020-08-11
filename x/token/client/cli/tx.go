package cli

import (
	"fmt"
	"github.com/pocblockchain/pocc/client"
	"github.com/pocblockchain/pocc/x/auth"
	"github.com/pocblockchain/pocc/x/auth/client/utils"
	"strings"

	"github.com/pocblockchain/pocc/client/context"
	"github.com/pocblockchain/pocc/codec"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/version"
	govtype "github.com/pocblockchain/pocc/x/gov/types"
	"github.com/pocblockchain/pocc/x/token/types"
	"github.com/spf13/cobra"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "token subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(client.PostCommands(
		GetCmdNewToken(cdc),
		GetCmdInflateToken(cdc),
		GetCmdBurnToken(cdc),
	)...)

	return txCmd
}

//create a token in pocchain
func GetCmdNewToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [to][symbol][decimals][totalSupply]",
		Short: "new a token",
		Long:  ` Example: new-token poc1xxx bhetc 18 1000000000000000000000000000`,

		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			to, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			symbol := sdk.Symbol(args[1])
			if !symbol.IsValidTokenName() {
				return fmt.Errorf("%v is not a valid token name", args[1])
			}

			decimals, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("Fail to parse decimals:%v", args[2])
			}

			if decimals.IsNegative() || decimals.Int64() > sdk.Precision {
				return fmt.Errorf("invalid decimals:%v", args[2])
			}

			totalSupply, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("Fail to parse totalSupply:%v", args[3])
			}

			from := cliCtx.GetFromAddress()
			msg := types.NewMsgNewToken(from, to, symbol.String(), uint64(decimals.Int64()), totalSupply)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}

func GetCmdInflateToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inflate [to][amount]",
		Short: "inflate a token",
		Long:  ` Example: inflate-token poc1xxx 1000000000000000000000000000btc`,

		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			to, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// parse coins trying to be burn
			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			msg := types.NewMsgInflateToken(from, to, coins)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}

//create a token in pocchain
func GetCmdBurnToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn coins",
		Short: "burn some token",
		Long:  ` Example: burn 10000000000btc --from alice`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// parse coins trying to be burn
			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			msg := types.NewMsgBurnToken(from, coins)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}

// GetCmdTokenParamsChangeProposal implements the command to submit a TokenParamsChange proposal
func GetCmdTokenParamsChangeProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-params-change [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a token params change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a token params change proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal token-params-change <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Token Param Change",
  "description": "token param change proposal",
  "changes": [
    {
      "key": "is_send_enabled",
      "value": true
    },
  ],
  "deposit": [
    {
      "denom": "hbc",
      "amount": "10000"
    }
  ]
}
`, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := ParseTokenParamsChangeProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			changes := proposal.Changes.ToParamChanges()
			from := cliCtx.GetFromAddress()
			content := types.NewTokenParamsChangeProposal(proposal.Title, proposal.Description, proposal.Symbol, changes)

			msg := govtype.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdDisableTokenProposal implements the command to submit a DisableToken proposal
func GetCmdDisableTokenProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-token [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a disable token proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a disable token proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal disable-token <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Disable Token",
  "description": "disable token proposal",
  "symbol": "testtoken",
  "deposit": [
    {
      "denom": "hbc",
      "amount": "100000"
    }
  ]
}
`, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := ParseDisableTokenProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewDisableTokenProposal(proposal.Title, proposal.Description, proposal.Symbol)

			msg := govtype.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
