package types

// distribution module event types
const (
	EventTypeSetWithdrawAddress                  = "set_withdraw_address"
	EventTypeRewards                             = "rewards"
	EventTypeCommission                          = "commission"
	EventTypeWithdrawRewards                     = "withdraw_rewards"
	EventTypeWithdrawCommission                  = "withdraw_commission"
	EventTypeProposerReward                      = "proposer_reward"
	EventTypeExecutionCommunityPoolSpendProposal = "exec_community_pool_spend_proposal"

	AttributeKeyWithdrawAddress = "withdraw_address"
	AttributeKeyValidator       = "validator"
	AttributeKeyRecipient       = "recipient"

	AttributeValueCategory = ModuleName
)
