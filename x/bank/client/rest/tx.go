package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pocblockchain/pocc/client/context"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/types/rest"
	"github.com/pocblockchain/pocc/x/auth/client/utils"

	"github.com/pocblockchain/pocc/x/bank/internal/types"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/bank/accounts/{address}/transfers", SendRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/bank/balances/{address}", QueryBalancesRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bank/accounts/{address}/escrow", EscrowRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/bank/accounts/{address}/reclaim", ReclaimRequestHandlerFn(cliCtx)).Methods("POST")
}

// SendReq defines the properties of a send request's body.
type SendReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coins    `json:"amount" yaml:"amount"`
}

// SendRequestHandlerFn - http request handler to send coins to a address.
func SendRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		toAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req SendReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSend(fromAddr, toAddr, req.Amount)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// EscrowReq defines the properties of a send request's body.
type EscrowReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coins    `json:"amount" yaml:"amount"`
}

// EscrowRequestHandlerFn - http request handler to send coins to a address.
func EscrowRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		toAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req SendReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgEscrow(fromAddr, toAddr, req.Amount)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// ReclaimReq defines the properties of a send request's body.
type ReclaimReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coins    `json:"amount" yaml:"amount"`
}

// ReclaimRequestHandlerFn - http request handler to send coins to a address.
func ReclaimRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		toAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req SendReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgReclaim(fromAddr, toAddr, req.Amount)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
