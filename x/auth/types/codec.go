package types

import (
	"github.com/pocblockchain/pocc/codec"
	"github.com/pocblockchain/pocc/x/auth/exported"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "poc/Account", nil)
	cdc.RegisterInterface((*exported.VestingAccount)(nil), nil)
	cdc.RegisterConcrete(&BaseVestingAccount{}, "poc/BaseVestingAccount", nil)
	cdc.RegisterConcrete(&ContinuousVestingAccount{}, "poc/ContinuousVestingAccount", nil)
	cdc.RegisterConcrete(&DelayedVestingAccount{}, "poc/DelayedVestingAccount", nil)
	cdc.RegisterConcrete(StdTx{}, "poc/StdTx", nil)
}

// module wide codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
