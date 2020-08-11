package types

import (
	"github.com/pocblockchain/pocc/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "poc/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "poc/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "poc/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgUndelegate{}, "poc/MsgUndelegate", nil)
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "poc/MsgBeginRedelegate", nil)
}

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
