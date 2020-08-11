package types

// bank module event types
const (
	EventTypeTransfer      = "transfer"
	EventTypeEscrow        = "escrow"
	EventTypeReclaim       = "reclaim"
	EventTypeMultiTransfer = "multi_transfer"
	EventTypeBonusSend     = "bonus_send"
	EventTypeReclaimSend   = "relcaim_send"

	AttributeKeyRecipient   = "recipient"
	AttributeKeySender      = "sender"
	AttributeKeyReclaimTo   = "reclaim_to"
	AttributeKeyReclaimFrom = "reclaim_from"

	AttributeValueCategory = ModuleName
)
