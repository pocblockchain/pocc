/*
 * *******************************************************************
 * @项目名称: token
 * @文件名称: genesis.go
 * @Date: 2019/06/14
 * @Author: Keep
 * @Copyright（C）: 2019 BlueHelix Inc.   All rights reserved.
 * 注意：本内容仅限于内部传阅，禁止外泄以及用于其他的商业目的.
 * *******************************************************************
 */

package token

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	sdk "github.com/pocblockchain/pocc/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

//GenesisState ...
type GenesisState struct {
	GenesisTokenInfos []sdk.TokenInfoWithoutSupply `json:"genesis_token_info"`
	Params            Params                       `json:"params"`
}

//NewGenesisState ...
func NewGenesisState(tokenInfos []sdk.TokenInfoWithoutSupply, params Params) GenesisState {
	sort.Sort(sortGenesisTokenInfoWithoutSupply(tokenInfos))
	return GenesisState{
		GenesisTokenInfos: tokenInfos,
		Params:            params,
	}
}

//ValidateGenesis ...
func ValidateGenesis(data GenesisState) error {
	tokenMap := make(map[string]struct{}, len(data.GenesisTokenInfos))

	for _, genInfo := range data.GenesisTokenInfos {
		if !genInfo.IsValid() {
			return fmt.Errorf("%v is not valid", genInfo)
		}

		symbol := genInfo.Symbol.String()

		if _, ok := tokenMap[symbol]; ok {
			return fmt.Errorf("invalid TokenInfoMap: Value: %v. duplicated token", genInfo)
		}

		tokenMap[symbol] = struct{}{}
	}

	return nil
}

//DefaultGenesisState ...
func DefaultGenesisState() GenesisState {
	genInfos := []sdk.TokenInfoWithoutSupply{
		{
			Symbol:        sdk.Symbol(sdk.NativeToken),
			Issuer:        "",
			IsSendEnabled: true,
			Decimals:      sdk.NativeTokenDecimal,
		},
	}
	sort.Sort(sortGenesisTokenInfoWithoutSupply(genInfos))
	return GenesisState{
		GenesisTokenInfos: genInfos,
		Params:            DefaultParams(),
	}
}

//InitGenesis ...
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, ti := range data.GenesisTokenInfos {
		k.SetTokenInfoWithoutSupply(ctx, &ti)
	}
	k.SetParams(ctx, data.Params)
	return []abci.ValidatorUpdate{}
}

//ExportGenesis ...
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	genTokenInfos := make([]sdk.TokenInfoWithoutSupply, 0)
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, TokenStoreKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		symbol := sdk.Symbol(bytes.TrimPrefix(iter.Key(), TokenStoreKeyPrefix))
		info := k.GetTokenInfoWithoutSupply(ctx, symbol)
		genTokenInfos = append(genTokenInfos, *info)
	}
	sort.Sort(sortGenesisTokenInfoWithoutSupply(genTokenInfos))

	params := k.GetParams(ctx)
	return GenesisState{GenesisTokenInfos: genTokenInfos, Params: params}
}

//AddTokenInfoWithoutSupplyIntoGenesis add a token into genesis
func (g *GenesisState) AddTokenInfoWithoutSupplyIntoGenesis(new sdk.TokenInfoWithoutSupply) error {
	genInfos := g.GenesisTokenInfos

	for _, info := range genInfos {
		if info.Symbol == new.Symbol {
			return fmt.Errorf("%v already exist", new.Symbol)
		}
	}

	genInfos = append(genInfos, new)
	sort.Sort(sortGenesisTokenInfoWithoutSupply(genInfos))
	g.GenesisTokenInfos = genInfos
	return nil
}

//Equal check whether 2 genesisstate equal
func (g GenesisState) Equal(g2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(g)
	b2 := ModuleCdc.MustMarshalBinaryBare(g2)
	return bytes.Equal(b1, b2)
}

//IsEmpty checks whether genesis state is empty
func (g GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return g.Equal(emptyGenState)
}

func (g GenesisState) String() string {
	var b strings.Builder

	for _, ti := range g.GenesisTokenInfos {
		b.WriteString(ti.String())
		b.WriteString("\n")
	}
	b.WriteString(g.Params.String())

	return b.String()
}

type sortGenesisTokenInfoWithoutSupply []sdk.TokenInfoWithoutSupply

func (s sortGenesisTokenInfoWithoutSupply) Len() int      { return len(s) }
func (s sortGenesisTokenInfoWithoutSupply) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortGenesisTokenInfoWithoutSupply) Less(i, j int) bool {
	return string(s[i].Symbol) < string(s[j].Symbol)
}
