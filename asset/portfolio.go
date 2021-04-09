package asset

import (
	"database/sql"
	"errors"
	"fmt"
	"gobacktrader/btutil"
	"sort"
	"time"
)

// Weight represents an asset weight for a given portfolio.
type Weight sql.NullFloat64

var (
	nullWeight        = Weight{Float64: 0.0, Valid: false}
	errWrongPortfolio = errors.New("applying to the wrong portfolio")
)

// IComplianceRule defines the compliance rule interface.
type IComplianceRule interface {
	Passes(*Portfolio) (bool, error)
}

// PortfolioSnapshot takes a snapshot of portfolio value and weights
// for a specific timestamp.
type PortfolioSnapshot struct {
	timestamp time.Time
	value     Price
	weights   map[IAssetReadOnly]Weight
}

func newPortfolioSnapshot(timestamp time.Time, p *Portfolio) (PortfolioSnapshot, error) {
	portfolioValue, positionWeights, err := p.GetValueWeights()
	snap := PortfolioSnapshot{
		timestamp: timestamp,
		value:     portfolioValue,
		weights:   positionWeights,
	}
	return snap, err
}

// GetTime returns the timestamp for our snapshot.
func (s PortfolioSnapshot) GetTime() time.Time {
	return s.timestamp
}

// GetValue returns the portfolio value for our snapshot.
func (s PortfolioSnapshot) GetValue() Price {
	return s.value
}

// GetWeights returns the portfolio weights for our snapshot.
func (s PortfolioSnapshot) GetWeights() map[IAssetReadOnly]Weight {
	return s.weights
}

// Portfolio consists of a collection of positions.
type Portfolio struct {
	code            string
	baseCurrency    string
	positions       map[IAssetReadOnly]*Position
	fxRates         *FxRates
	history         map[time.Time]PortfolioSnapshot
	complianceRules []IComplianceRule
}

// Show will print and return a string representation
// of the portfolio positions.
func (p Portfolio) Show() string {
	maxTickerLen := 0
	for asset := range p.positions {
		ticker := asset.GetTicker()
		lenTicker := len(ticker)
		if lenTicker > maxTickerLen {
			maxTickerLen = lenTicker
		}
	}

	var positionList []string
	for asset, position := range p.positions {
		ticker := btutil.PadRight(asset.GetTicker(), " ", uint(maxTickerLen+5))
		units := fmt.Sprintf("%0.2f", position.GetUnits())
		positionList = append(positionList, ticker+units+"\n")
	}

	output := "---Portfolio('" + p.code + "')---\n"
	sort.Strings(positionList)
	for _, item := range positionList {
		output += item
	}

	fmt.Println(output)
	return output
}

// NewPortfolio returns a new instance of Portfolio.
func NewPortfolio(code string, baseCurrency string) (*Portfolio, error) {
	positions := make(map[IAssetReadOnly]*Position)
	history := make(map[time.Time]PortfolioSnapshot)
	baseCurrency, err := ValidateCurrency(baseCurrency)
	fxRates := &FxRates{}
	portfolio := Portfolio{
		code:         code,
		baseCurrency: baseCurrency,
		positions:    positions,
		fxRates:      fxRates,
		history:      history,
	}
	return &portfolio, err
}

// GetCode returns our portfolio code
func (p *Portfolio) GetCode() string {
	return p.code
}

// GetBaseCurrency returns the base currency for our portfolio.
func (p *Portfolio) GetBaseCurrency() string {
	return p.baseCurrency
}

// NumPositions returns the number of portfolio positions.
func (p *Portfolio) NumPositions() int {
	return len(p.positions)
}

// HasAsset returns a boolean of true if a portfolio
// contains some asset, false otherwise.
func (p *Portfolio) HasAsset(a IAssetReadOnly) bool {
	_, ok := p.positions[a]
	return ok
}

// GetUnits returns the units held for a given asset.
func (p *Portfolio) GetUnits(a IAssetReadOnly) float64 {
	position, ok := p.positions[a]
	if !ok {
		return 0.0
	}
	return position.GetUnits()
}

// GetWeight returns the portfolio weight in a given asset.
func (p *Portfolio) GetWeight(a IAssetReadOnly) (Weight, error) {
	_, positionWeights, err := p.GetValueWeights()
	if err != nil {
		return nullWeight, err
	}

	weight, ok := positionWeights[a]
	if !ok {
		return nullWeight, nil
	}

	return weight, nil
}

// ModifyPositions allows us to increment and decrement positions
// in the portfolio.
func (p *Portfolio) ModifyPositions(a IAssetReadOnly, units float64) {
	// if the asset is already held then modify its position
	if p.HasAsset(a) {
		p.positions[a].Increment(units)
	} else { // otherwise create a new position for this asset
		p.positions[a] = &Position{a, units}
	}
}

// Transfer has identical functionality to ModifyPositions
// and will increment or decrement some asset in the portfolio.
func (p *Portfolio) Transfer(a IAssetReadOnly, units float64) {
	p.ModifyPositions(a, units)
}

// GetValue returns our portfolio value.
func (p Portfolio) GetValue() (Price, error) {
	var totalValueFloat float64

	for _, position := range p.positions {
		value := position.GetValue()
		if !value.Valid { // no value, so value and weight invalid
			return nullValue, nil
		}

		assetBaseCurrency := position.GetBaseCurrency()
		pair := assetBaseCurrency + p.baseCurrency
		fxRate, ok, err := p.fxRates.GetRate(pair)
		if err != nil {
			return nullValue, err
		}
		if !ok { // no fx rate, so value and weight invalid
			return nullValue, nil
		}

		totalValueFloat += value.Float64 * fxRate
	}

	return Price{Float64: totalValueFloat, Valid: true}, nil
}

// GetValueWeights returns the portfolio value along with all position weights.
func (p *Portfolio) GetValueWeights() (Price, map[IAssetReadOnly]Weight, error) {
	var totalValueFloat float64
	portfolioValue := nullValue
	valid := true
	positionValues := make(map[IAssetReadOnly]Price)
	positionWeights := make(map[IAssetReadOnly]Weight)

	for _, position := range p.positions {
		units := position.GetUnits()
		if units == 0 {
			continue // nothing to value.
		}
		asset := position.GetAsset()
		value := position.GetValue()
		if !value.Valid { // no value, so value and weight invalid
			valid = false
			positionValues[asset] = nullValue
			positionWeights[asset] = nullWeight
			continue
		}

		assetBaseCurrency := position.GetBaseCurrency()
		pair := assetBaseCurrency + p.baseCurrency
		fxRate, ok, err := p.fxRates.GetRate(pair)
		if err != nil {
			return portfolioValue, positionWeights, err
		}
		if !ok { // no fx rate, so value and weight invalid
			valid = false
			positionValues[asset] = nullValue
			positionWeights[asset] = nullWeight
			continue
		}

		positionValueBaseCurrency := value.Float64 * fxRate
		positionValues[asset] = Price{Float64: positionValueBaseCurrency, Valid: true}
		totalValueFloat += positionValueBaseCurrency
	}

	if totalValueFloat == 0.0 {
		err := errors.New("cannot calculate weights for portfolio with zero value")
		return portfolioValue, positionWeights, err
	}

	if valid { // all assets have a valid value which we can use to derive portfolio value
		portfolioValue = Price{Float64: totalValueFloat, Valid: true}
		for asset, price := range positionValues {
			positionWeights[asset] = Weight{Float64: price.Float64 / totalValueFloat, Valid: true}
		}
	}

	return portfolioValue, positionWeights, nil
}

// SetFxRates sets the FX rates object to be used for this portfolio.
func (p *Portfolio) SetFxRates(fxRates *FxRates) {
	p.fxRates = fxRates
}

// GetFxRates returns the FX rates object for this portfolio
func (p *Portfolio) GetFxRates() *FxRates {
	return p.fxRates
}

// TakeSnapshot takes a portfolio snapshot for a given timestamp.
func (p *Portfolio) TakeSnapshot(timestamp time.Time) error {
	snap, err := newPortfolioSnapshot(timestamp, p)
	if err != nil {
		return err
	}

	p.history[timestamp] = snap
	return nil
}

// GetHistory returns the portfolio snapshot history.
func (p Portfolio) GetHistory() map[time.Time]PortfolioSnapshot {
	return p.history
}

// Copy returns a copy of the portfolio ex history.
func (p *Portfolio) Copy() (*Portfolio, error) {
	portfolioCopy, err := NewPortfolio(p.GetCode(), p.GetBaseCurrency())
	if err != nil {
		return portfolioCopy, err
	}

	// copy over positions
	for asset, position := range p.positions {
		portfolioCopy.Transfer(asset, position.GetUnits())
	}

	// fx rates
	portfolioCopy.SetFxRates(p.GetFxRates())

	// and any compliance rules
	for _, rule := range p.complianceRules {
		portfolioCopy.AddComplianceRule(rule)
	}

	return portfolioCopy, nil
}

// NumComplianceRules returns the number of compliance rules
// attached to this portfolio.
func (p *Portfolio) NumComplianceRules() int {
	return len(p.complianceRules)
}

// HasComplianceRule returns true if a portfolio has some
// compliance rule, false otherwise.
func (p *Portfolio) HasComplianceRule(rule IComplianceRule) bool {
	for _, portfolioRule := range p.complianceRules {
		if rule == portfolioRule {
			return true
		}
	}
	return false
}

// AddComplianceRule adds a compliance rule to the portfolio.
func (p *Portfolio) AddComplianceRule(rule IComplianceRule) error {
	if p.HasComplianceRule(rule) {
		return nil
	}
	p.complianceRules = append(p.complianceRules, rule)
	return nil
}

// RemoveComplianceRule removes a compliance rule from the portfolio.
func (p *Portfolio) RemoveComplianceRule(rule IComplianceRule) error {
	if !p.HasComplianceRule(rule) {
		return nil
	}
	var newRules []IComplianceRule
	for _, portfolioRule := range p.complianceRules {
		if portfolioRule != rule {
			newRules = append(newRules, portfolioRule)
		}
	}

	p.complianceRules = newRules
	return nil
}

// PassesCompliance returns true if all compliance rules pass,
// false otherwise.
func (p *Portfolio) PassesCompliance() (bool, error) {
	allPasses := true
	for _, rule := range p.complianceRules {
		rulePasses, err := rule.Passes(p)
		if err != nil {
			return false, err
		}
		if !rulePasses {
			allPasses = false
		}
	}
	return allPasses, nil
}

// Trade executes some trade for the portfolio.
func (p *Portfolio) Trade(asset IAssetReadOnly, units float64, consideration *float64) error {
	if consideration == nil { // calculate consideration
		assetValue := asset.GetValue()
		if !assetValue.Valid {
			return fmt.Errorf("'%s' cannot trade an asset with invalid value", asset.GetTicker())
		}
		considerationAmount := assetValue.Float64 * -units
		consideration = &considerationAmount
	}

	currencyCode := asset.GetBaseCurrency()
	cash, err := NewCash(currencyCode)
	if err != nil {
		return err
	}

	p.Transfer(asset, units)
	p.Transfer(cash, *consideration)
	return nil
}
