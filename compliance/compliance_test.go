package compliance

import (
	"errors"
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"testing"
)

type testRule struct {
	portfolio *asset.Portfolio
}

func (t *testRule) GetPortfolio() *asset.Portfolio {
	return t.portfolio
}

func (t *testRule) Passes() (bool, error) {
	return false, errors.New("this is a test error")
}

func TestCompliance(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	cash, err2 := asset.NewCash("AUD")
	stock, err3 := asset.NewStock("ZZB AU", "AUD")
	if err := btutil.AnyValidError(err1, err2, err3); err != nil {
		t.Errorf("Error in asset creation - %s", err)
	}

	// transfer 100 shares of stock and 100 AUD to the portfolio
	portfolio.ModifyPositions(&stock, 100)
	portfolio.ModifyPositions(&cash, 100)

	cashUnitLimit := NewUnitLimit(&portfolio, &cash, 100)
	stockUnitLimit := NewUnitLimit(&portfolio, &stock, 100)

	rules := NewRules(&portfolio)
	for _, rule := range []asset.IComplianceRule{&cashUnitLimit, &stockUnitLimit} {
		rules.AddRule(rule)
	}

	if rules.GetPortfolio() != &portfolio {
		t.Error("Unexpected portfolio")
	}

	// we are on the edge of our limits,
	// so compliance should pass
	pass, err := rules.Passes()
	if err != nil {
		t.Errorf("Error in Rules{}.Passes() - %s", err)
	}
	if !pass {
		t.Error("Expecting compliance rules to pass")
	}

	// add one share to our stock holding and rules
	// should now fail
	portfolio.ModifyPositions(&stock, 1)
	pass, err = rules.Passes()
	if err != nil {
		t.Errorf("Error in Rules{}.Passes() - %s", err)
	}
	if pass {
		t.Error("Expecting compliance rules to fail")
	}
}

func TestComplianceWrongPortfolio(t *testing.T) {
	portfolio1, err1 := asset.NewPortfolio("XXX", "AUD")
	portfolio2, err2 := asset.NewPortfolio("YYY", "AUD")
	cash, err3 := asset.NewCash("AUD")
	stock, err4 := asset.NewStock("ZZB AU", "AUD")
	if err := btutil.AnyValidError(err1, err2, err3, err4); err != nil {
		t.Errorf("Error in asset creation - %s", err)
	}

	// transfer 100 shares of stock and 100 AUD to the portfolio
	portfolio1.ModifyPositions(&stock, 100)
	portfolio1.ModifyPositions(&cash, 100)

	cashUnitLimit := NewUnitLimit(&portfolio1, &cash, 100)
	stockUnitLimit1 := NewUnitLimit(&portfolio1, &stock, 100)
	stockUnitLimit2 := NewUnitLimit(&portfolio2, &stock, 100)

	rules := NewRules(&portfolio1)
	for _, rule := range []asset.IComplianceRule{&cashUnitLimit, &stockUnitLimit1, &stockUnitLimit2} {
		err := rules.AddRule(rule)
		if rule == &stockUnitLimit2 {
			if btutil.GetErrorString(err) != "compliance rule relates to a different portfolio" {
				t.Error("Expecting an error when mixing portfolio rules")
			}
			continue
		}
		if err != nil {
			t.Errorf("Error in Rules{}.AddRule() - %s", err)
		}
	}
}

func TestCompliancePasses(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in asset creation - %s", err)
	}

	// transfer 100 shares of stock to the portfolio
	portfolio.ModifyPositions(&stock, 100)
	stockUnitLimit := NewUnitLimit(&portfolio, &stock, 100)

	rules := NewRules(&portfolio)
	if err := rules.AddRule(&stockUnitLimit); err != nil {
		t.Errorf("Error in Rules{}.AddRule() - %s", err)
	}

	pass, err := rules.Passes()
	if err != nil {
		t.Errorf("Error in Rules{}.Passes() - %s", err)
	}
	if !pass {
		t.Error("Expecting to be within unit limits")
	}

	// tip us over the unit limit
	portfolio.ModifyPositions(&stock, 1)
	pass, err = rules.Passes()
	if err != nil {
		t.Errorf("Error in Rules{}.Passes() - %s", err)
	}
	if pass {
		t.Error("Expecting to be outside unit limits")
	}

	// add a rule which will throw an error when running
	err = rules.AddRule(&testRule{portfolio: &portfolio})
	if err != nil {
		t.Errorf("Error in Rules{}.AddRule() - %s", err)
	}

	_, err = rules.Passes()
	if btutil.GetErrorString(err) != "this is a test error" {
		t.Error("Expecting a specific error when running Rules{}.Passes()")
	}
}
