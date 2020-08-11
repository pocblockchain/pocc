package types

import (
	"github.com/pocblockchain/pocc/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgWithdrawDelegatorReward{}, "poc/MsgWithdrawDelegationReward", nil)
	cdc.RegisterConcrete(MsgWithdrawValidatorCommission{}, "poc/MsgWithdrawValidatorCommission", nil)
	cdc.RegisterConcrete(MsgSetWithdrawAddress{}, "poc/MsgModifyWithdrawAddress", nil)
	cdc.RegisterConcrete(CommunityPoolSpendProposal{}, "poc/CommunityPoolSpendProposal", nil)
}

// generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
