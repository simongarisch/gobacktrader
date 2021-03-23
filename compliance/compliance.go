package compliance

import (
	"errors"
	"gobacktrader/asset"
)

// IComplianceRule defines the compliance rule interface
type IComplianceRule interface {
	GetPortfolio() *asset.Portfolio
	Passes() (bool, error)
}

// Rules represents a collection of compliance rules
// to be applied to some portfolio.
type Rules struct {
	rules     []IComplianceRule
	portfolio *asset.Portfolio
}

// NewRules returns a new instance of Rules.
func NewRules(portfolio *asset.Portfolio) Rules {
	return Rules{
		portfolio: portfolio,
	}
}

// GetPortfolio returns the portfolio for which
// these compliance rules are applied.
func (r *Rules) GetPortfolio() *asset.Portfolio {
	return r.portfolio
}

// AddRule adds a compliance rule to the compliance rules collection.
func (r *Rules) AddRule(rule IComplianceRule) error {
	if rule.GetPortfolio() != r.portfolio {
		return errors.New("compliance rule relates to a different portfolio")
	}
	r.rules = append(r.rules, rule)
	return nil
}

// Passes returns true if all compliance rules pass, false otherwise.
func (r *Rules) Passes() (bool, error) {
	allPasses := true
	for _, rule := range r.rules {
		rulePasses, err := rule.Passes()
		if err != nil {
			return allPasses, err
		}
		if !rulePasses {
			allPasses = false
		}
	}
	return allPasses, nil
}
