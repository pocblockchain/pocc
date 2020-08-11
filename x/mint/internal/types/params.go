package types

import (
	"fmt"
	"strings"

	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/params"
)

var DefaultInitalInflationAmount = sdk.NewDecFromInt(sdk.NewIntWithDecimal(81000000, 18)) //81000000 in the 1st year

var DefaultInflationFactor = []sdk.Dec{
	sdk.NewDec(1),             //1st year
	sdk.NewDecWithPrec(85, 2), //2nd year
	sdk.NewDecWithPrec(85, 2), //3rd year
	sdk.NewDecWithPrec(85, 2), //4th year
	sdk.NewDecWithPrec(9, 1),  //5th~ years
}

// Parameter store keys
var (
	KeyMintDenom              = []byte("MintDenom")
	KeyInitalInflationAmount  = []byte("InitalInflationAmount")
	KeyInflationFactorPerYear = []byte("InflationFactorPerYear")
	KeyBlocksPerYear          = []byte("BlocksPerYear")
)

// mint parameters
type Params struct {
	MintDenom              string    `json:"mint_denom" yaml:"mint_denom"`                               // type of coin to mint
	InitalInflationAmount  sdk.Dec   `json:"initial_inflation_amount" yaml:"initial_inflation_amount"`   //initial amount of minted coins
	InflationFactorPerYear []sdk.Dec `json:"inflation_factor_per_year" yaml:"inflation_factor_per_year"` //
	BlocksPerYear          uint64    `json:"blocks_per_year" yaml:"blocks_per_year"`                     // expected blocks per year
}

// ParamTable for minting module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(mintDenom string, initalInflationAmount sdk.Dec, inflationFactorPerYear []sdk.Dec, blocksPerYear uint64) Params {

	return Params{
		MintDenom:              mintDenom,
		InitalInflationAmount:  initalInflationAmount,
		InflationFactorPerYear: inflationFactorPerYear,
		BlocksPerYear:          blocksPerYear,
	}
}

func DefaultParams() Params {
	return Params{
		MintDenom:              sdk.DefaultBondDenom,
		InitalInflationAmount:  DefaultInitalInflationAmount,
		InflationFactorPerYear: DefaultInflationFactor,
		BlocksPerYear:          uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
	}
}

// validate params
func ValidateParams(params Params) error {
	if params.InitalInflationAmount.IsNegative() {
		return fmt.Errorf("mint parameter InflationAmountInCurrentYear should NOT be negative, is %s ", params.InitalInflationAmount.String())
	}
	for _, factor := range params.InflationFactorPerYear {
		if factor.IsNegative() {
			return fmt.Errorf("mint parameter InflationFactorPerYear should NOT be negative, is %s ", factor.String())
		}

	}
	if params.MintDenom == "" {
		return fmt.Errorf("mint parameter MintDenom can't be an empty string")
	}
	return nil
}

func (p Params) String() string {
	var b strings.Builder

	for i, factor := range p.InflationFactorPerYear {
		b.WriteString(fmt.Sprintf("%v:%s", i, factor.String()))
	}

	return fmt.Sprintf(`Minting Params:
  Mint Denom:%s
  Initial Inflation Amoun:%s
  Inflation Factors Per Year:%s
  Blocks Per Year:%d
`,
		p.MintDenom, p.InitalInflationAmount, b.String(), p.BlocksPerYear,
	)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyMintDenom, &p.MintDenom},
		{KeyInitalInflationAmount, &p.InitalInflationAmount},
		{KeyInflationFactorPerYear, &p.InflationFactorPerYear},
		{KeyBlocksPerYear, &p.BlocksPerYear},
	}
}
