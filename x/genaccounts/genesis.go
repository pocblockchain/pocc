package genaccounts

import (
	"github.com/pocblockchain/pocc/codec"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/auth"
	authexported "github.com/pocblockchain/pocc/x/auth/exported"
	"github.com/pocblockchain/pocc/x/genaccounts/internal/types"
)

// InitGenesis initializes accounts and deliver genesis transactions
func InitGenesis(ctx sdk.Context, _ *codec.Codec, accountKeeper types.AccountKeeper,  supplyKeeper types.SupplyKeeper, genesisState GenesisState) {
	genesisState.Sanitize()

	// load the accounts
	for _, gacc := range genesisState {
		acc := gacc.ToAccount()
		acc = accountKeeper.NewAccount(ctx, acc) // set account number
		accountKeeper.SetAccount(ctx, acc)
	}
	supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
}

// ExportGenesis exports genesis for all accounts
func ExportGenesis(ctx sdk.Context, accountKeeper types.AccountKeeper) GenesisState {

	// iterate to get the accounts
	accounts := []GenesisAccount{}
	accountKeeper.IterateAccounts(ctx,
		func(acc authexported.Account) (stop bool) {
			account, err := NewGenesisAccountI(acc)
			if err != nil {
				panic(err)
			}
			accounts = append(accounts, account)
			return false
		},
	)

	return accounts
}
