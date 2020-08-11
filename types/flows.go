package types

import (
	"errors"

	"github.com/pocblockchain/pocc/codec"
	abci "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// CategoryType indicates the category type that causes the receipt in the flow.
type CategoryType uint64

const (
	CategoryTypeTransfer      CategoryType = 0x1
	CategoryTypeMultiTransfer CategoryType = 0x2
)

// Receipt defines basic interface for all kind of receipts
type Receipt struct {
	// Category for the transaction that causes the receipt.
	Category CategoryType

	// Flows list of flows.
	Flows []Flow
}

// Flow defines the interface of the flow in the receipt
type Flow interface{}

// BalanceFlow for asset balance change
type BalanceFlow struct {
	// AccAddress the address for the custodian unit for the balance change
	AccAddress AccAddress

	// Symbol token symbol for the asset
	Symbol Symbol

	// PreviousBalance previous balance for the balance change
	PreviousBalance Int

	// BalanceChange the actual balance change
	BalanceChange Int
}

func GetReceiptFromResult(cdc *codec.Codec, result *Result) (Receipt, error) {
	if result.Data == nil {
		return Receipt{}, errors.New("invalid data")
	}
	rc := Receipt{}
	if err := cdc.UnmarshalBinaryLengthPrefixed(result.Data, &rc); err != nil {
		return rc, err
	}
	return rc, nil
}

func NewResultFromResultTx(res *ctypes.ResultTx) Result {
	return NewResultFromDeliverTx(&res.TxResult)
}

func NewResultFromDeliverTx(res *abci.ResponseDeliverTx) Result {
	return Result{
		Code:      CodeType(res.Code),
		Codespace: CodespaceType(res.Codespace),
		Data:      res.Data,
		Log:       res.Log,
		GasWanted: uint64(res.GasWanted),
		GasUsed:   uint64(res.GasUsed),
		//Tags:      res.Tags, FIXME(liyong.zhang): temp change
	}
}
