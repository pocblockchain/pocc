package types

import (
	"github.com/pocblockchain/pocc/codec"
	sdk "github.com/pocblockchain/pocc/types"
	authtypes "github.com/pocblockchain/pocc/x/auth/types"
	stakingtypes "github.com/pocblockchain/pocc/x/staking/types"
)

// ModuleCdc defines a generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

// TODO: abstract genesis transactions registration back to staking
// required for genesis transactions
func init() {
	ModuleCdc = codec.New()
	stakingtypes.RegisterCodec(ModuleCdc)
	authtypes.RegisterCodec(ModuleCdc)
	sdk.RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
