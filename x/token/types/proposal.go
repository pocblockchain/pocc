package types

import (
	"fmt"
	"strings"

	sdk "github.com/pocblockchain/pocc/types"
	govtypes "github.com/pocblockchain/pocc/x/gov/types"
)

const (
	// ProposalTypeAddToken defines the type for a AddToken
	ProposalTypeTokenParamsChange = "TokenParamsChange"
	ProposalTypeDisableToken      = "DisableToken"
)

// Assert proposl implements govtypes.Content at compile-time
var _ govtypes.Content = TokenParamsChangeProposal{}
var _ govtypes.Content = DisableTokenProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeTokenParamsChange)
	govtypes.RegisterProposalTypeCodec(TokenParamsChangeProposal{}, "poc/token/TokenParamsChangeProposal")
	govtypes.RegisterProposalType(ProposalTypeDisableToken)
	govtypes.RegisterProposalTypeCodec(DisableTokenProposal{}, "poc/token/DisableTokenProposal")
}

type ParamChange struct {
	Key   string
	Value string
}

//NewParamChange create a paramchange(k,v)
func NewParamChange(key, value string) ParamChange {
	return ParamChange{key, value}
}

// TokenParamsChangeProposal modify a token's variable parameter
type TokenParamsChangeProposal struct {
	Title       string        `json:"title" yaml:"title"`
	Description string        `json:"description" yaml:"description"`
	Symbol      string        `json:"symbol" yaml:"symbol"`
	Changes     []ParamChange `json:"changes" yaml:"changes"`
}

// NewTokenParamsChangeProposal creates a new add token proposal.
func NewTokenParamsChangeProposal(title, description, sybmol string, changes []ParamChange) TokenParamsChangeProposal {
	return TokenParamsChangeProposal{
		Title:       title,
		Description: description,
		Symbol:      sybmol,
		Changes:     changes,
	}
}

// GetTitle returns the title of a token parameter change proposal.
func (ctpp TokenParamsChangeProposal) GetTitle() string { return ctpp.Title }

// GetDescription returns the description of a token parameter change proposal.
func (ctpp TokenParamsChangeProposal) GetDescription() string { return ctpp.Description }

// GetDescription returns the routing key of a token parameter change proposal.
func (ctpp TokenParamsChangeProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a token parameter change proposal.
func (ctpp TokenParamsChangeProposal) ProposalType() string { return ProposalTypeTokenParamsChange }

// ValidateBasic runs basic stateless validity checks
func (ctpp TokenParamsChangeProposal) ValidateBasic() sdk.Error {
	err := govtypes.ValidateAbstract(DefaultCodespace, ctpp)
	if err != nil {
		return err
	}

	if !sdk.Symbol(ctpp.Symbol).IsValidTokenName() {
		return ErrInvalidSymbol(ctpp.Symbol)
	}

	if ctpp.Symbol == sdk.NativeToken {
		return sdk.ErrInvalidTx("Not allowed to change native token's params")
	}

	//dectect duplicated keys if any
	keysMap := map[string]interface{}{}

	for _, pc := range ctpp.Changes {
		_, ok := keysMap[pc.Key]
		if !ok {
			keysMap[pc.Key] = nil
		} else {
			return ErrDuplicatedKey()
		}

		if len(pc.Key) == 0 {
			return ErrEmptyKey()
		}
		if len(pc.Value) == 0 {
			return ErrEmptyValue()
		}
	}

	return err
}

// String implements the Stringer interface.
func (ctpp TokenParamsChangeProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Change Token Param Proposal:
 Title:       %s
 Description: %s
 Symbol:      %s
 Changes:
`, ctpp.Title, ctpp.Description, ctpp.Symbol))

	for _, pc := range ctpp.Changes {
		b.WriteString(fmt.Sprintf("%s: %s\t", pc.Key, pc.Value))
	}
	return b.String()
}

// DisableTokenProposal modify a token's variable parameter
type DisableTokenProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Symbol      string `json:"symbol" yaml:"symbol"`
}

// NewDisableTokenProposal creates a new disable token proposal.
func NewDisableTokenProposal(title, description, sybmol string) DisableTokenProposal {
	return DisableTokenProposal{
		Title:       title,
		Description: description,
		Symbol:      sybmol,
	}
}

// GetTitle returns the title of a disable token proposal..
func (dtp DisableTokenProposal) GetTitle() string { return dtp.Title }

// GetDescription returns the description of a disable token proposal..
func (dtp DisableTokenProposal) GetDescription() string { return dtp.Description }

// GetDescription returns the routing key of a disable token proposal..
func (dtp DisableTokenProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a disable token proposal.
func (dtp DisableTokenProposal) ProposalType() string { return ProposalTypeDisableToken }

// ValidateBasic runs basic stateless validity checks
func (dtp DisableTokenProposal) ValidateBasic() sdk.Error {
	err := govtypes.ValidateAbstract(DefaultCodespace, dtp)
	if err != nil {
		return err
	}

	if !sdk.Symbol(dtp.Symbol).IsValidTokenName() {
		return ErrInvalidSymbol(dtp.Symbol)
	}
	if dtp.Symbol == sdk.NativeToken {
		return sdk.ErrInvalidTx("Not allowed to disable native token's params")
	}

	return err
}

// String implements the Stringer interface.
func (dtp DisableTokenProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Disable Token Proposal:
 Title:       %s
 Description: %s
 Symbol:      %s
`, dtp.Title, dtp.Description, dtp.Symbol))
	return b.String()
}
