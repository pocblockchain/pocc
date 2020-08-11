package types

import (
	"github.com/pocblockchain/pocc/codec"
	"github.com/pocblockchain/pocc/x/supply/exported"
)

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.ModuleAccountI)(nil), nil)
	cdc.RegisterInterface((*exported.SupplyI)(nil), nil)
	cdc.RegisterConcrete(&ModuleAccount{}, "poc/ModuleAccount", nil)
	cdc.RegisterConcrete(&Supply{}, "poc/Supply", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}
