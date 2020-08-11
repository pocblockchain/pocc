package types

import (
	"github.com/pocblockchain/pocc/codec"
	sdk "github.com/pocblockchain/pocc/types"
)

var ModuleCdc *codec.Codec

//RegisterCodec register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(sdk.TokenInfo{}, "poc/token/TokenInfo", nil)
	cdc.RegisterConcrete(sdk.TokenInfoWithoutSupply{}, "poc/token/TokenInfoWithoutSupply", nil)
	cdc.RegisterConcrete(TokenParamsChangeProposal{}, "poc/token/TokenParamsChangeProposal", nil)
	cdc.RegisterConcrete(DisableTokenProposal{}, "poc/token/DisableTokenProposal", nil)
	cdc.RegisterConcrete(MsgNewToken{}, "poc/token/MsgNewToken", nil)
	cdc.RegisterConcrete(MsgBurnToken{}, "poc/token/MsgBurnToken", nil)
	cdc.RegisterConcrete(MsgInflateToken{}, "poc/token/MsgInflateToken", nil)

}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
