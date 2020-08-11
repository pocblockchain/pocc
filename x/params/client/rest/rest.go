package rest

import (
	"net/http"

	"github.com/pocblockchain/pocc/client/context"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/types/rest"
	"github.com/pocblockchain/pocc/x/auth/client/utils"
	"github.com/pocblockchain/pocc/x/gov"
	govrest "github.com/pocblockchain/pocc/x/gov/client/rest"
	"github.com/pocblockchain/pocc/x/params"
	paramscutils "github.com/pocblockchain/pocc/x/params/client/utils"
)

// ProposalRESTHandler returns a ProposalRESTHandler that exposes the param
// change REST handler with a given sub-route.
func ProposalRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "param_change",
		Handler:  postProposalHandlerFn(cliCtx),
	}
}

func postProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req paramscutils.ParamChangeProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := params.NewParameterChangeProposal(req.Title, req.Description, req.Changes.ToParamChanges())

		msg := gov.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
