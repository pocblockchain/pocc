package types

const (
	EventTypeExecuteTokenParamsChangeProposal = "execute_token_params_change_proposal"
	EventTypeExecuteDisableTokenProposal      = "execute_disable_token_proposal"
	EventTypeNewToken                         = "new_token"
	EventTypeBurnToken                        = "burn_token"
	EventTypeInflateToken                     = "inflate_token"

	AttributeKeyTokenParam      = "param"
	AttributeKeyTokenParamValue = "value"
	AttributeKeyToken           = "token"
	AttributeKeyIssuer          = "issuer"
	AttributeKeyRecipient       = "recipient"
	AttributeKeySymbol          = "symbol"
	AttributeKeyAmount          = "amount"
	AttributeKeyIssueFee        = "issue_fee"
	AttributeKeyBurner          = "burner"

	AttributeValueCategory = ModuleName
)
