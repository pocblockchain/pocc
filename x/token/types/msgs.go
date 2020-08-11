package types

import (
	"fmt"
	sdk "github.com/pocblockchain/pocc/types"
)

var _ sdk.Msg = MsgNewToken{}

//MsgNewToken ...
type MsgNewToken struct {
	From        sdk.AccAddress `json:"from" yaml:"from"`
	To          sdk.AccAddress `json:"to" yaml:"to"`
	Symbol      sdk.Symbol     `json:"symbol" yaml:"symbol"`
	Decimals    uint64         `json:"decimals" yaml:"decimals"`
	TotalSupply sdk.Int        `json:"total_supply" yaml:"total_supply"`
}

//NewMsgNewToken is a constructor function for MsgTokenNew
func NewMsgNewToken(from, to sdk.AccAddress, symbol string, decimals uint64, totalSupply sdk.Int) MsgNewToken {
	return MsgNewToken{
		From:        from,
		To:          to,
		Symbol:      sdk.Symbol(symbol),
		Decimals:    decimals,
		TotalSupply: totalSupply,
	}
}

func (msg MsgNewToken) Route() string { return RouterKey }
func (msg MsgNewToken) Type() string  { return TypeMsgNewToken }

// ValidateBasic runs stateless checks on the message
func (msg MsgNewToken) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress(fmt.Sprintf("from address can not be empty:%v", msg.From))
	}

	if msg.To.Empty() {
		return sdk.ErrInvalidAddress(fmt.Sprintf("to address can not be empty%v", msg.To))
	}

	if !msg.Symbol.IsValidTokenName() {
		return sdk.ErrInvalidSymbol(fmt.Sprintf("symbol %v is invalid", msg.Symbol))
	}

	if msg.Decimals > sdk.Precision {
		return sdk.ErrTooMuchPrecision(fmt.Sprintf("maximum:%v, provided:%v", sdk.Precision, msg.Decimals))
	}

	if !msg.TotalSupply.IsPositive() {
		return sdk.ErrInvalidAmount(fmt.Sprintf("totalSupply %v is not positive", msg.TotalSupply))
	}

	return nil
}

//GetSignBytes ...
func (msg MsgNewToken) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

//GetSigners ...
func (msg MsgNewToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

var _ sdk.Msg = MsgInflateToken{}

//MsgInflateToken ...
type MsgInflateToken struct {
	From   sdk.AccAddress `json:"from" yaml:"from"`
	To     sdk.AccAddress `json:"to" yaml:"to"`
	Amount sdk.Coins      `json:"amount" yaml:"amount"`
}

//NewMsgInflateToken is a constructor function for MsgInflateToken
func NewMsgInflateToken(from, to sdk.AccAddress, amount sdk.Coins) MsgInflateToken {
	return MsgInflateToken{
		From:   from,
		To:     to,
		Amount: amount,
	}
}

func (msg MsgInflateToken) Route() string { return RouterKey }
func (msg MsgInflateToken) Type() string  { return TypeMsgInflateToken }

// ValidateBasic runs stateless checks on the message
func (msg MsgInflateToken) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress(fmt.Sprintf("from address can not be empty:%v", msg.From))
	}

	if msg.To.Empty() {
		return sdk.ErrInvalidAddress(fmt.Sprintf("to address can not be empty%v", msg.To))
	}

	if len(msg.Amount) != 1 {
		return sdk.ErrInvalidTx(fmt.Sprintf("inflate only ONE coin once"))
	}

	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidAmount(fmt.Sprintf("inflate %v is not positive", msg.Amount))
	}

	return nil
}

//GetSignBytes  implements Msg.
func (msg MsgInflateToken) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

//GetSigners  implements Msg.
func (msg MsgInflateToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

//__________________________________________________________
var _ sdk.Msg = MsgBurnToken{}

// MsgBurnToken burn coins from account
type MsgBurnToken struct {
	From   sdk.AccAddress `json:"from_address" yaml:"from_address"`
	Amount sdk.Coins      `json:"amount" yaml:"amount"`
}

// NewMsgBurnToken - burn coins msg
func NewMsgBurnToken(fromAddr sdk.AccAddress, amount sdk.Coins) MsgBurnToken {
	return MsgBurnToken{From: fromAddr, Amount: amount}
}

// Route Implements Msg.
func (msg MsgBurnToken) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgBurnToken) Type() string { return TypeMsgBurnToken }

// ValidateBasic Implements Msg.
func (msg MsgBurnToken) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}

	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("burn amount is invalid: " + msg.Amount.String())
	}
	if !msg.Amount.IsAllPositive() {
		return sdk.ErrInvalidCoins("burn amount must be positive")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBurnToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgBurnToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}
