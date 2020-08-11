/*
 * *******************************************************************
 * @项目名称: types
 * @文件名称: querier.go
 * @Date: 2019/06/05
 * @Author: Keep
 * @Copyright（C）: 2019 BlueHelix Inc.   All rights reserved.
 * 注意：本内容仅限于内部传阅，禁止外泄以及用于其他的商业目的.
 * *******************************************************************
 */
package types

import (
	"strings"

	sdk "github.com/pocblockchain/pocc/types"
)

//QueryTokenInfo
type QueryTokenInfo struct {
	Symbol string `json:"symbol"`
}

func NewQueryTokenInfo(symbol string) QueryTokenInfo {
	return QueryTokenInfo{Symbol: symbol}
}

//QueryResToken
type QueryResToken sdk.TokenInfo

func (qs QueryResToken) String() string {
	return sdk.TokenInfo(qs).String()
}

//QueryResTokens
type QueryResTokens []sdk.TokenInfo

func (qs QueryResTokens) String() string {
	if len(qs) == 0 {
		return ""
	}

	var b strings.Builder
	for _, ti := range qs {
		b.WriteString(ti.String())
		b.WriteString("\n")
	}

	return b.String()
}

//QueryResSymbols
type QueryResSymbols []string

func (qs QueryResSymbols) String() string {
	return strings.Join(qs, ",")
}
