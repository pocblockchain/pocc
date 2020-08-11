/*
 * *******************************************************************
 * @项目名称: token
 * @文件名称: keeper.go
 * @Date: 2019/06/05
 * @Author: Keep
 * @Copyright（C）: 2019 BlueHelix Inc.   All rights reserved.
 * 注意：本内容仅限于内部传阅，禁止外泄以及用于其他的商业目的.
 * *******************************************************************
 */

package token

import (
	"bytes"
	"github.com/pocblockchain/pocc/codec"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/params"
	"github.com/pocblockchain/pocc/x/supply"
	"github.com/pocblockchain/pocc/x/token/types"
)

// TokenStoreKeyPrefix define prefix for storing tokeninfk
var TokenStoreKeyPrefix = []byte{0x01}

/*
 Note:
	TokenInfoWithoutSupply stored in token module and total supply stored in supply module.
*/

//TokenKeeper defines the token module interface
type TokenKeeper interface {
	SetTokenInfoWithoutSupply(ctx sdk.Context, tokenInfo *sdk.TokenInfoWithoutSupply)
	DeleteTokenInfoWithoutSupply(ctx sdk.Context, symbol string)

	GetTokenInfo(ctx sdk.Context, symbol sdk.Symbol) *sdk.TokenInfo
	GetIssuer(ctx sdk.Context, symbol sdk.Symbol) string
	GetDecimals(ctx sdk.Context, symbol sdk.Symbol) uint64
	GetTotalSupply(ctx sdk.Context, symbol sdk.Symbol) sdk.Int

	IsSendEnabled(ctx sdk.Context, symbol sdk.Symbol) bool
	EnableSend(ctx sdk.Context, symbol sdk.Symbol)
	DisableSend(ctx sdk.Context, symbol sdk.Symbol)
}

//Keeper ...
type Keeper struct {
	storeKey      sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc           *codec.Codec // The wire codec for binary encoding/decoding
	dk            types.DistrKeeper
	sk            types.SupplyKeeper
	paramSubSpace params.Subspace
}

//NewKeeper create token's Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, dk types.DistrKeeper, sk types.SupplyKeeper, paramSubSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		dk:            dk,
		sk:            sk,
		paramSubSpace: paramSubSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

func tokenStoreKey(symbol string) []byte {
	return append(TokenStoreKeyPrefix, []byte(symbol)...)
}

var _ TokenKeeper = (*Keeper)(nil)

//SetTokenInfoWithoutSupply sets TokenInfoWithoutSupply
func (k *Keeper) SetTokenInfoWithoutSupply(ctx sdk.Context, tokenInfo *sdk.TokenInfoWithoutSupply) {
	store := ctx.KVStore(k.storeKey)
	ti := sdk.TokenInfoWithoutSupply{
		Symbol:        tokenInfo.Symbol,
		Issuer:        tokenInfo.Issuer,
		Decimals:      tokenInfo.Decimals,
		IsSendEnabled: tokenInfo.IsSendEnabled,
	}
	store.Set(tokenStoreKey(tokenInfo.Symbol.String()), k.cdc.MustMarshalBinaryBare(ti))
}

//DeleteTokenInfoWithoutSupply delete TokenInfoWithoutSupply
func (k *Keeper) DeleteTokenInfoWithoutSupply(ctx sdk.Context, symbol string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(tokenStoreKey(symbol))
}

//GetAllTokenInfoWithoutSupply get all token's TokenInfoWithoutSupply
func (k *Keeper) GetAllTokenInfoWithoutSupply(ctx sdk.Context) []sdk.TokenInfoWithoutSupply {
	store := ctx.KVStore(k.storeKey)
	var tokens []sdk.TokenInfoWithoutSupply
	iter := sdk.KVStorePrefixIterator(store, TokenStoreKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var ti sdk.TokenInfoWithoutSupply
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &ti)
		tokens = append(tokens, ti)
	}
	return tokens
}

//GetTokenInfoWithoutSupply gets a specified token's GetTokenInfoWithoutSupply
func (k *Keeper) GetTokenInfoWithoutSupply(ctx sdk.Context, symbol sdk.Symbol) *sdk.TokenInfoWithoutSupply {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(tokenStoreKey(symbol.String())) {
		return nil
	}
	bz := store.Get(tokenStoreKey(symbol.String()))
	var ti sdk.TokenInfoWithoutSupply
	k.cdc.MustUnmarshalBinaryBare(bz, &ti)

	return &ti
}

//GetTokenInfo get a specified  tokeninfo, whose totalsupply is stored in supply module
func (k *Keeper) GetTokenInfo(ctx sdk.Context, symbol sdk.Symbol) *sdk.TokenInfo {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(tokenStoreKey(symbol.String())) {
		return nil
	}

	bz := store.Get(tokenStoreKey(symbol.String()))
	var tsi sdk.TokenInfoWithoutSupply
	k.cdc.MustUnmarshalBinaryBare(bz, &tsi)

	ti := castToTokenInfo(tsi)
	ti.TotalSupply = k.sk.GetSupply(ctx).GetTotal().AmountOf(ti.Symbol.String())
	return &ti
}

//GetAllTokenInfo get all tokeninfo, whose totalsupply is stored in supply module
func (k *Keeper) GetAllTokenInfo(ctx sdk.Context) []sdk.TokenInfo {
	store := ctx.KVStore(k.storeKey)
	var tokens []sdk.TokenInfo
	iter := sdk.KVStorePrefixIterator(store, TokenStoreKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var tsi sdk.TokenInfoWithoutSupply
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &tsi)

		ti := castToTokenInfo(tsi)
		ti.TotalSupply = k.sk.GetSupply(ctx).GetTotal().AmountOf(ti.Symbol.String())
		tokens = append(tokens, ti)
	}
	return tokens
}

//GetIssuer ...
func (k *Keeper) GetIssuer(ctx sdk.Context, symbol sdk.Symbol) string {
	token := k.GetTokenInfoWithoutSupply(ctx, symbol)
	if token == nil {
		return ""
	}
	return token.Issuer
}

//IsTokenSupported ...
func (k *Keeper) IsTokenSupported(ctx sdk.Context, symbol sdk.Symbol) bool {
	store := ctx.KVStore(k.storeKey)
	if store.Has(tokenStoreKey(symbol.String())) {
		return true
	}
	return false
}

//IsSendEnabled ...
func (k *Keeper) IsSendEnabled(ctx sdk.Context, symbol sdk.Symbol) bool {
	token := k.GetTokenInfoWithoutSupply(ctx, symbol)
	if token == nil {
		return false
	}

	return token.IsSendEnabled
}

//GetDecimals ...
func (k *Keeper) GetDecimals(ctx sdk.Context, symbol sdk.Symbol) uint64 {
	token := k.GetTokenInfoWithoutSupply(ctx, symbol)
	if token == nil {
		return 0
	}
	return token.Decimals
}

//GetTotalSupply ...
func (k *Keeper) GetTotalSupply(ctx sdk.Context, symbol sdk.Symbol) sdk.Int {
	token := k.GetTokenInfoWithoutSupply(ctx, symbol)
	if token == nil {
		return sdk.NewInt(0)
	}
	return k.sk.GetSupply(ctx).GetTotal().AmountOf(string(symbol))
}

//EnableSend ...
func (k *Keeper) EnableSend(ctx sdk.Context, symbol sdk.Symbol) {
	token := k.GetTokenInfoWithoutSupply(ctx, symbol)
	if token == nil {
		return
	}
	token.IsSendEnabled = true
	k.SetTokenInfoWithoutSupply(ctx, token)
}

//DisableSend ...
func (k *Keeper) DisableSend(ctx sdk.Context, symbol sdk.Symbol) {
	token := k.GetTokenInfoWithoutSupply(ctx, symbol)
	if token == nil {
		return
	}
	token.IsSendEnabled = false
	k.SetTokenInfoWithoutSupply(ctx, token)
}

//GetSymbolIterator ...
func (k *Keeper) GetSymbolIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, TokenStoreKeyPrefix)
}

//GetSymbols ...
func (k *Keeper) GetSymbols(ctx sdk.Context) []string {
	var symbols []string
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, TokenStoreKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		symbols = append(symbols, string(bytes.TrimPrefix(iter.Key(), TokenStoreKeyPrefix)))
	}
	return symbols
}

// SetParams sets the token module's parameters.
func (k *Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSubSpace.SetParamSet(ctx, &params)
}

// GetParams gets the token module's parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params Params) {
	k.paramSubSpace.GetParamSet(ctx, &params)
	return
}

func castToTokenInfo(tsi sdk.TokenInfoWithoutSupply) sdk.TokenInfo {
	return sdk.TokenInfo{
		Symbol:        tsi.Symbol,
		Issuer:        tsi.Issuer,
		IsSendEnabled: tsi.IsSendEnabled,
		Decimals:      tsi.Decimals,
	}
}

func castToTokenInfoWithoutSupply(tsi sdk.TokenInfo) sdk.TokenInfoWithoutSupply {
	return sdk.TokenInfoWithoutSupply{
		Symbol:        tsi.Symbol,
		Issuer:        tsi.Issuer,
		IsSendEnabled: tsi.IsSendEnabled,
		Decimals:      tsi.Decimals,
	}
}

//__For Test only__
//The following methods will update total supply unconditionlly, shoul only be used in test
//SetTokenInfo set TokenInfo and total supply
func (k *Keeper) SetTokenInfo(ctx sdk.Context, tokenInfo *sdk.TokenInfo) {
	store := ctx.KVStore(k.storeKey)
	tsi := sdk.TokenInfoWithoutSupply{
		Symbol:        tokenInfo.Symbol,
		Issuer:        tokenInfo.Issuer,
		Decimals:      tokenInfo.Decimals,
		IsSendEnabled: tokenInfo.IsSendEnabled,
	}
	store.Set(tokenStoreKey(tokenInfo.Symbol.String()), k.cdc.MustMarshalBinaryBare(tsi))

	//update supply
	oldCoins := k.sk.GetSupply(ctx).GetTotal()
	oldAmt := sdk.NewCoins(sdk.NewCoin(tokenInfo.Symbol.String(), oldCoins.AmountOf(tokenInfo.Symbol.String())))
	oldCoins = oldCoins.Sub(oldAmt)
	coins := sdk.NewCoins(sdk.NewCoin(tokenInfo.Symbol.String(), tokenInfo.TotalSupply))
	newCoins := oldCoins.Add(coins)
	k.sk.SetSupply(ctx, supply.NewSupply(newCoins))

}

//DeleteTokenInfo delete TokenInfo and total supply
func (k *Keeper) DeleteTokenInfo(ctx sdk.Context, symbol string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(tokenStoreKey(symbol))

	//update supply
	oldCoins := k.sk.GetSupply(ctx).GetTotal()
	oldAmt := sdk.NewCoins(sdk.NewCoin(symbol, oldCoins.AmountOf(symbol)))
	newCoins := oldCoins.Sub(oldAmt)
	k.sk.SetSupply(ctx, supply.NewSupply(newCoins))
}
