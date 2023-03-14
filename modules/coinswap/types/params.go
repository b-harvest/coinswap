package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyFee                    = []byte("Fee")                    // fee key
	KeyPoolCreationFee        = []byte("PoolCreationFee")        // fee key
	KeyTaxRate                = []byte("TaxRate")                // fee key
	KeyStandardDenom          = []byte("StandardDenom")          // standard token denom key
	KeyMaxStandardCoinPerPool = []byte("MaxStandardCoinPerPool") // max standard coin amount per pool
	KeyWhitelistedDenoms      = []byte("WhitelistedDenoms")      // whitelisted denoms

	DefaultFee                    = sdk.NewDecWithPrec(0, 0)
	DefaultPoolCreationFee        = sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)
	DefaultTaxRate                = sdk.NewDecWithPrec(0, 1)
	DefaultMaxStandardCoinPerPool = sdk.NewInt(10_000_000_000)
	DefaultWhitelistedDenoms      = []string{}
)

// NewParams is the coinswap params constructor
func NewParams(fee, taxRate sdk.Dec, poolCreationFee sdk.Coin, maxStandardCoinPerPool sdk.Int, whitelistedDenoms []string) Params {
	return Params{
		Fee:                    fee,
		TaxRate:                taxRate,
		PoolCreationFee:        poolCreationFee,
		MaxStandardCoinPerPool: maxStandardCoinPerPool,
		WhitelistedDenoms:      whitelistedDenoms,
	}
}

// ParamKeyTable returns the TypeTable for coinswap module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements paramtypes.KeyValuePairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyFee, &p.Fee, validateFee),
		paramtypes.NewParamSetPair(KeyPoolCreationFee, &p.PoolCreationFee, validatePoolCreationFee),
		paramtypes.NewParamSetPair(KeyTaxRate, &p.TaxRate, validateTaxRate),
		paramtypes.NewParamSetPair(KeyMaxStandardCoinPerPool, &p.MaxStandardCoinPerPool, validateMaxStandardCoinPerPool),
		paramtypes.NewParamSetPair(KeyWhitelistedDenoms, &p.WhitelistedDenoms, validateWhitelistedDeoms),
	}
}

// DefaultParams returns the default coinswap module parameters
func DefaultParams() Params {
	return Params{
		Fee:                    DefaultFee,
		PoolCreationFee:        DefaultPoolCreationFee,
		TaxRate:                DefaultTaxRate,
		MaxStandardCoinPerPool: DefaultMaxStandardCoinPerPool,
		WhitelistedDenoms:      DefaultWhitelistedDenoms,
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate returns err if Params is invalid
func (p Params) Validate() error {
	if p.Fee.IsNegative() || !p.Fee.LT(sdk.OneDec()) {
		return fmt.Errorf("fee must be positive and less than 1: %s", p.Fee.String())
	}
	return nil
}

func validateFee(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() || !v.LT(sdk.OneDec()) {
		return fmt.Errorf("fee must be positive and less than 1: %s", v.String())
	}

	return nil
}

func validatePoolCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("poolCreationFee must be positive: %s", v.String())
	}
	return nil
}

func validateTaxRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() || !v.LT(sdk.OneDec()) {
		return fmt.Errorf("fee must be positive and less than 1: %s", v.String())
	}
	return nil
}

func validateMaxStandardCoinPerPool(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsPositive() {
		return fmt.Errorf("maxStandardCoinPerPool must be positive: %s", v.String())
	}
	return nil
}

func validateWhitelistedDeoms(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, denom := range v {
		if err := sdk.ValidateDenom(denom); err != nil {
			return err
		}
	}
	return nil
}
