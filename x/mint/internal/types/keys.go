package types

// the one key to use for the keeper store
var MinterKey = []byte{0x01}

// nolint
const (
	// module name
	ModuleName = "mint"

	// default paramspace for params keeper
	DefaultParamspace = ModuleName

	// StoreKey is the default store key for mint
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the minting querier
	QueryParameters       = "parameters"
	QueryAnnualProvisions = "annual_provisions"
)
