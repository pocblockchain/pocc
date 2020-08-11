package types

import (
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/auth/exported"
)

// AccountKeeper defines the account contract that must be fulfilled when
// creating a x/bank keeper.
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) exported.Account

	GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	GetAllAccounts(ctx sdk.Context) []exported.Account
	SetAccount(ctx sdk.Context, acc exported.Account)

	IterateAccounts(ctx sdk.Context, process func(exported.Account) bool)
}

type TokenKeeper interface {
	IsSendEnabled(ctx sdk.Context, symbol sdk.Symbol) bool
}
