package compliance

import (
	"errors"
	"gobacktrader/asset"
	"gobacktrader/btutil"
	"testing"
)

type testRule struct{}

func (t *testRule) Passes(portfolio *asset.Portfolio) (bool, error) {
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
	portfolio.Transfer(&stock, 100)
	portfolio.Transfer(&cash, 100)

	cashUnitLimit := NewUnitLimit(&cash, 100)
	stockUnitLimit := NewUnitLimit(&stock, 100)
	for _, rule := range []asset.IComplianceRule{&cashUnitLimit, &stockUnitLimit} {
		portfolio.AddComplianceRule(rule)
	}

	// we are on the edge of our limits,
	// so compliance should pass
	pass, err := portfolio.PassesCompliance()
	if err != nil {
		t.Errorf("Error in Portfolio{}.PassesCompliance() - %s", err)
	}
	if !pass {
		t.Error("Expecting compliance rules to pass")
	}

	// add one share to our stock holding and rules
	// should now fail
	portfolio.Transfer(&stock, 1)
	pass, err = portfolio.PassesCompliance()
	if err != nil {
		t.Errorf("Error in Portfolio{}.PassesCompliance() - %s", err)
	}
	if pass {
		t.Error("Expecting compliance rules to fail")
	}
}

func TestCompliancePasses(t *testing.T) {
	portfolio, err1 := asset.NewPortfolio("XXX", "AUD")
	stock, err2 := asset.NewStock("ZZB AU", "AUD")
	if err := btutil.AnyValidError(err1, err2); err != nil {
		t.Errorf("Error in asset creation - %s", err)
	}

	// transfer 100 shares of stock to the portfolio
	portfolio.Transfer(&stock, 100)
	stockUnitLimit := NewUnitLimit(&stock, 100)
	portfolio.AddComplianceRule(&stockUnitLimit)

	pass, err := portfolio.PassesCompliance()
	if err != nil {
		t.Errorf("Error in Portfolio{}.PassesCompliance() - %s", err)
	}
	if !pass {
		t.Error("Expecting to be within unit limits")
	}

	// tip us over the unit limit
	portfolio.Transfer(&stock, 1)
	pass, err = portfolio.PassesCompliance()
	if err != nil {
		t.Errorf("Error in Portfolio{}.PassesCompliance() - %s", err)
	}
	if pass {
		t.Error("Expecting to be outside unit limits")
	}

	// add a rule which will throw an error when running
	err = portfolio.AddComplianceRule(&testRule{})
	if err != nil {
		t.Errorf("Error in Portfolio{}.AddRule() - %s", err)
	}

	_, err = portfolio.PassesCompliance()
	if btutil.GetErrorString(err) != "this is a test error" {
		t.Error("Expecting a specific error when running Rules{}.Passes()")
	}
}
