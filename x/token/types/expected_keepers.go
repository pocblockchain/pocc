package types

import (
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/auth/exported"
	supplyexported "github.com/pocblockchain/pocc/x/supply/exported"
)

// AccountKeeper defines the account contract that must be fulfilled when
// creating a x/bank keeper.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	GetAllAccounts(ctx sdk.Context) []exported.Account
	SetAccount(ctx sdk.Context, acc exported.Account)
}

// SupplyKeeper defines the supply Keeper for module accounts
type SupplyKeeper interface {
	GetSupply(ctx sdk.Context) (supply supplyexported.SupplyI) //for get total supply from supply module
	SetSupply(ctx sdk.Context, supply supplyexported.SupplyI)
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI

	// TODO remove with genesis 2-phases refactor https://github.com/pocblockchain/pocc/issues/2862
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error
}

//DistrKeeper defines the distribution Keeper for module token
type DistrKeeper interface {
	AddCoinsFromAccountToFeePool(ctx sdk.Context, senderAddr sdk.AccAddress, amount sdk.Coins) sdk.Error
}
