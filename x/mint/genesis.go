package mint

import (
	"bytes"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/mint/internal/types"
	"strings"
)

// GenesisState - minter state
type GenesisState struct {
	Minter Minter `json:"minter" yaml:"minter"` // minter object
	Params Params `json:"params" yaml:"params"` // inflation params
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(minter Minter, params Params) GenesisState {
	return GenesisState{
		Minter: minter,
		Params: params,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	defaultParams := DefaultParams()
	return GenesisState{
		Minter: DefaultInitialMinter(defaultParams.InitalInflationAmount),
		Params: DefaultParams(),
	}
}

// InitGenesis new mint genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, supplyKeeper types.SupplyKeeper, data GenesisState) {
	keeper.SetMinter(ctx, data.Minter)
	keeper.SetParams(ctx, data.Params)
	supplyKeeper.GetModuleAccount(ctx, ModuleName)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	minter := keeper.GetMinter(ctx)
	params := keeper.GetParams(ctx)
	return NewGenesisState(minter, params)
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	err := ValidateParams(data.Params)
	if err != nil {
		return err
	}

	err = ValidateMinter(data.Minter)
	if err != nil {
		return err
	}

	return nil
}

func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}

func (data GenesisState) String() string {
	out := ""
	out += data.Minter.String()
	out += data.Params.String()
	return strings.TrimSpace(out)
}
