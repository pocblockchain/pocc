package cli

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pocblockchain/pocc/codec"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/token/types"
)

type (
	ParamChangesJSON []ParamChangeJSON

	// ParamChangeJSON defines a parameter change used in JSON input. This
	// allows values to be specified in raw JSON instead of being string encoded.
	ParamChangeJSON struct {
		Key   string          `json:"key" yaml:"key"`
		Value json.RawMessage `json:"value" yaml:"value"`
	}

	AddTokenProposalJSON struct {
		Title       string        `json:"title" yaml:"title"`
		Description string        `json:"description" yaml:"description"`
		TokenInfo   sdk.TokenInfo `json:"token_info" yaml:"token_info"`
		Deposit     sdk.Coins     `json:"deposit" yaml:"deposit"`
	}

	TokenParamsChangeProposalJSON struct {
		Title       string           `json:"title" yaml:"title"`
		Description string           `json:"description" yaml:"description"`
		Symbol      string           `json:"symbol" yaml:"symbol"`
		Changes     ParamChangesJSON `json:"changes" yaml:"changes"`
		Deposit     sdk.Coins        `json:"deposit" yaml:"deposit"`
	}

	DisableTokenProposalJSON struct {
		Title       string    `json:"title" yaml:"title"`
		Description string    `json:"description" yaml:"description"`
		Symbol      string    `json:"symbol" yaml:"symbol"`
		Deposit     sdk.Coins `json:"deposit" yaml:"deposit`
	}
)

// ToParamChange converts a ParamChangeJSON object to ParamChange.
func (pcj ParamChangeJSON) ToParamChange() types.ParamChange {
	return types.NewParamChange(pcj.Key, string(pcj.Value))
}

// ToParamChanges converts a slice of paramChangesJSON objects to a slice of
// ParamChange.
func (pcsj ParamChangesJSON) ToParamChanges() []types.ParamChange {
	res := make([]types.ParamChange, len(pcsj))
	for i, pc := range pcsj {
		res[i] = pc.ToParamChange()
	}
	return res
}

// ParseTokenParamsChangeProposalJSON reads and parses a tokenParamsChangeProposalJSON from a file.
func ParseTokenParamsChangeProposalJSON(cdc *codec.Codec, proposalFile string) (TokenParamsChangeProposalJSON, error) {
	proposal := TokenParamsChangeProposalJSON{}
	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

// ParseDisableTokenProposalJSON reads and parses a disableTokenProposalJSON from a file.
func ParseDisableTokenProposalJSON(cdc *codec.Codec, proposalFile string) (DisableTokenProposalJSON, error) {
	proposal := DisableTokenProposalJSON{}
	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
