package token

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/pocblockchain/pocc/client/context"
	"github.com/pocblockchain/pocc/codec"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/types/module"
	"github.com/pocblockchain/pocc/x/token/client/cli"
	"github.com/pocblockchain/pocc/x/token/client/rest"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the token module.
type AppModuleBasic struct{}

// Name returns the distribution module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the distribution module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the token
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the token module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	if err := ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %v", ModuleName, err)
	}

	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the token module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the token module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(StoreKey, cdc)
}

// GetQueryCmd returns the root query command for the token module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

//____________________________________________________________________________

// AppModule implements an application module for the token module.
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// Name returns the distribution module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the distribution module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// Route returns the message routing key for the token module.
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the token module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the distribution module's querier route name.
func (AppModule) QuerierRoute() string {
	return RouterKey
}

// NewQuerierHandler returns the distribution module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the token module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the distribution
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the distribution module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the distribution module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
