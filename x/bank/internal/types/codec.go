package types

import (
	"github.com/pocblockchain/pocc/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "poc/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "poc/MsgMultiSend", nil)
	cdc.RegisterConcrete(MsgEscrow{}, "poc/MsgEscrow", nil)
	cdc.RegisterConcrete(MsgReclaim{}, "poc/MsgReclaim", nil)
	cdc.RegisterConcrete(MsgBonusSend{}, "poc/MsgBonusSend", nil)
	cdc.RegisterConcrete(MsgReclaimSend{}, "poc/MsgReclaimSend", nil)
}

// module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
