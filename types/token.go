package types

import (
	"fmt"
)

const (
	NativeToken        = "poc"
	NativeTokenDecimal = 18
)

var (
	KeyIsSendEnabled = "is_send_enabled"
)

//TokenInfo defines information in token module
type TokenInfo struct {
	Symbol        Symbol `json:"symbol" yaml:"symbol"`
	Issuer        string `json:"issuer" yaml:"issuer"`                   //token's issuer
	IsSendEnabled bool   `json:"is_send_enabled" yaml:"is_send_enabled"` //whether send enabled or not
	Decimals      uint64 `json:"decimals" yaml:"decimals"`               //token's decimals, represents by the decimals's
	TotalSupply   Int    `json:"total_supply" yaml:"total_supply" `      //token's total supply
}

type TokenInfoWithoutSupply struct {
	Symbol        Symbol `json:"symbol" yaml:"symbol"`
	Issuer        string `json:"issuer" yaml:"issuer"`                   //token's issuer
	IsSendEnabled bool   `json:"is_send_enabled" yaml:"is_send_enabled"` //whether send enabled or not
	Decimals      uint64 `json:"decimals" yaml:"decimals"`               //token's decimals, represents by the decimals's
}

func NewTokenInfo(symbol Symbol, issuer string, isSendEnabled bool, decimals uint64, totalSupply Int) *TokenInfo {
	return &TokenInfo{
		Symbol:        symbol,
		Issuer:        issuer,
		IsSendEnabled: isSendEnabled,
		Decimals:      decimals,
		TotalSupply:   totalSupply,
	}
}

//NewTokenInfo create a TokenInfoWithoutSupply
func NewTokenInfoWithoutSupply(symbol Symbol, issuer string, isSendEnabled bool, decimals uint64) *TokenInfoWithoutSupply {
	return &TokenInfoWithoutSupply{
		Symbol:        symbol,
		Issuer:        issuer,
		IsSendEnabled: isSendEnabled,
		Decimals:      decimals,
	}
}

func (t TokenInfoWithoutSupply) String() string {
	return fmt.Sprintf(`
	Symbol:%s
	Issuer:%v
	IsSendEnabled:%v
	Decimals:%v
	`, t.Symbol, t.Issuer, t.IsSendEnabled, t.Decimals)
}

func (t TokenInfoWithoutSupply) IsValid() bool {
	if !t.Symbol.IsValidTokenName() {
		return false
	}

	if t.Decimals > Precision {
		return false
	}

	return true
}

func (t TokenInfo) String() string {
	return fmt.Sprintf(`
	Symbol:%s
	Issuer:%v
	IsSendEnabled:%v
	Decimals:%v
	TotalSupply:%v
	`, t.Symbol, t.Issuer, t.IsSendEnabled, t.Decimals, t.TotalSupply)
}

func (t TokenInfo) IsValid() bool {
	if !t.Symbol.IsValidTokenName() {
		return false
	}

	if t.Decimals > Precision || t.TotalSupply.IsNegative() {
		return false
	}

	return true
}

type Symbol string

// IsValidTokenName check token name.
// a valid token name must be a valid coin denom
func (s Symbol) IsValidTokenName() bool {
	// same as coin
	if reDnm.MatchString(string(s)) {
		return true
	}
	return false
}

func (s Symbol) ToDenomName() string {
	if s.IsValidTokenName() {
		return string(s)
	}
	return ""
}

func (s Symbol) String() string {
	if s.IsValidTokenName() {
		return string(s)
	}
	return ""
}
