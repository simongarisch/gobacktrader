package broker

import (
	"errors"
	"gobacktrader/asset"
	"math"
)

// NoCharges applies zero charges to a trade.
type NoCharges struct{}

// NewNoCharges returns a new instance of NoCharges.
func NewNoCharges() NoCharges {
	return NoCharges{}
}

// Charge for the NoCharges strategy will apply no charges
// and leave portfolio cash unchanged.
func (c NoCharges) Charge(trade asset.ITrade) error {
	return nil
}

// FixedRatePlusPercentageCharges applies a fixed charge plus some
// percentage of the trade.
type FixedRatePlusPercentageCharges struct {
	fixedAmount  float64
	percentage   float64
	currencyCode string
}

// NewFixedRatePlusPercentageCharges returns a new instance of FixedRatePlusPercentageCharges.
func NewFixedRatePlusPercentageCharges(fixedAmount float64, percentage float64, currencyCode string) (FixedRatePlusPercentageCharges, error) {
	currencyCode, err := asset.ValidateCurrency(currencyCode)
	charges := FixedRatePlusPercentageCharges{
		fixedAmount:  fixedAmount,
		percentage:   percentage,
		currencyCode: currencyCode,
	}
	return charges, err
}

// GetFixedAmount returns the fixed amount portion of the charge.
func (c FixedRatePlusPercentageCharges) GetFixedAmount() float64 {
	return c.fixedAmount
}

// GetPercentage returns the percentage portion of the charge.
func (c FixedRatePlusPercentageCharges) GetPercentage() float64 {
	return c.percentage
}

// GetCurrencyCode returns the currency code for our charge.
func (c FixedRatePlusPercentageCharges) GetCurrencyCode() string {
	return c.currencyCode
}

// Charge for the FixedRatePlusPercentageCharges strategy will deduct
// some fixed amount plus a percentage of the trade from portfolio
// cash. You have the option to choose the currency code in which
// this charge is applied.
func (c FixedRatePlusPercentageCharges) Charge(trade asset.ITrade) error {
	portfolio, targetAsset := trade.GetPortfolio(), trade.GetAsset()
	tradeValue := trade.GetLocalCurrencyValue()
	if !tradeValue.Valid {
		return errors.New("cannot apply charges to a trade with invalid value")
	}

	// we'll apply charges in the chosen currency
	cash, err := asset.NewCash(c.currencyCode)
	if err != nil {
		return err
	}

	// apply the fixed charge
	portfolio.Transfer(cash, -math.Abs(c.fixedAmount))

	// apply the variable charge
	fxRates := portfolio.GetFxRates()
	fxPair := targetAsset.GetBaseCurrency() + c.currencyCode
	fxRate, _, err := fxRates.GetRate(fxPair)
	if err != nil {
		return err
	}
	portfolio.Transfer(cash, -math.Abs(tradeValue.Float64*fxRate*c.percentage))

	return nil
}
