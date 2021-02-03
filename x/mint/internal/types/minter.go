package types

import (
	"fmt"

	sdk "github.com/pocblockchain/pocc/types"
)

// default minting module parameters
//AnnualProvisions in current year = AnnualProvisions in previous year * InflationFactor
/*
   year      AnnualProvisions
	1         81000000 *InflationFactor[0]
    2         81000000 *InflationFactor[0] * InflationFactor[1]
    3         81000000 *InflationFactor[0] * InflationFactor[1] * InflationFactor[2]
    4         81000000 *InflationFactor[0] * InflationFactor[1] * InflationFactor[2] * InflationFactor[3]
    5         81000000 *InflationFactor[0] * InflationFactor[1] * InflationFactor[2] * InflationFactor[3] * InflationFactor[4]
    6         81000000 *InflationFactor[0] * InflationFactor[1] * InflationFactor[2] * InflationFactor[3] * pow(InflationFactor[4],2)
   	.... ....
    n         81000000 *InflationFactor[0] * InflationFactor[1] * InflationFactor[2] * InflationFactor[3] * pow(InflationFactor[4], n-4)
*/

// Minter represents the minting state.
type Minter struct {
	CurrentYearIndex uint64  `json:"current_year_index" yaml:"current_year_index"` // current year index
	AnnualProvisions sdk.Dec `json:"annual_provisions" yaml:"annual_provisions"`   // current annual expected provisions
}

// NewMinter returns a new Minter object with the given inflation and annual
// provisions values.
func NewMinter(currentYearIndex uint64, annualProvisions sdk.Dec) Minter {
	return Minter{
		CurrentYearIndex: currentYearIndex,
		AnnualProvisions: annualProvisions,
	}
}

// InitialMinter returns an initial Minter object with a given inflation value.
func InitialMinter(currentYearIndex uint64, annualProvisions sdk.Dec) Minter {
	return NewMinter(currentYearIndex, annualProvisions)
}

//DefaultInitialMinter set the initial minter 
func DefaultInitialMinter(annualProvisions sdk.Dec) Minter {
	return InitialMinter(0, annualProvisions)
}

// validate minter
func ValidateMinter(minter Minter) error {
	if minter.AnnualProvisions.IsNegative() {
		return fmt.Errorf("mint parameter AnnualProvisions should Not be negative, is %s",
			minter.AnnualProvisions.String())
	}

	return nil
}

// NextAnnualProvisions returns the next annual provisions based
func (m Minter) NextAnnualProvisions(inflationFactor sdk.Dec) sdk.Dec {
	return m.AnnualProvisions.Mul(inflationFactor)
}

// BlockProvision returns the provisions for a block based on the annual
// provisions rate.
func (m Minter) BlockProvision(params Params) sdk.Coin {
	provisionAmt := m.AnnualProvisions.QuoInt(sdk.NewInt(int64(params.BlocksPerYear)))
	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt())
}

func (m Minter) String() string {
	return fmt.Sprintf(`Minter:
  CurrentYearIndex:%v
  AnnualProvisions:%s
`,
		m.CurrentYearIndex, m.AnnualProvisions.String(),
	)
}
