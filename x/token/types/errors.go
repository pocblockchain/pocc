package types

import (
	"fmt"
	sdk "github.com/pocblockchain/pocc/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace   sdk.CodespaceType = "token"
	CodeInvalidInput   CodeType          = 103
	CodeEmptyData      CodeType          = 104
	CodeDuplicatedKey  CodeType          = 105
	CodeInvalidSymbol  CodeType          = 106
	CodeNonExistSymbol CodeType          = 107
	CodeSymbolReserved CodeType          = 108
)

// ErrEmptyKey returns an error for when an empty key is given.
func ErrEmptyKey() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeEmptyData, "parameter key is empty")
}

// ErrEmptyValue returns an error for when an empty key is given.
func ErrEmptyValue() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeEmptyData, "parameter value is empty")
}

// ErrEmptyValue returns an error for when an invalid parameter
func ErrInvalidParameter(key, value string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidInput, fmt.Sprintf("key:%v, value:%v", key, value))
}

// ErrDuplicatedKey returns an error for when duplicated keys appear
func ErrDuplicatedKey() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeDuplicatedKey, "parameter key is duplicated")
}

// ErrInvalidSymbol returns an error for when duplicated keys appear
func ErrInvalidSymbol(symbol string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidSymbol, "%v is invalid", symbol)
}

// ErrNonExistSymbol returns an error for when a symbol doesn't exist
func ErrNonExistSymbol(symbol string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNonExistSymbol, "%v does not exist", symbol)
}

// ErrSymbolReserved returns an error for when a symbol is reserved
func ErrSymbolReserved(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeSymbolReserved, msg)
}
