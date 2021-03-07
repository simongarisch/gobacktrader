package asset

import "testing"

func TestValidatePair(t *testing.T) {
	// start with a valid pair
	pair, err := ValidatePair(" audusd ")
	if err != nil {
		t.Errorf("Error in ValidatePair - %s", err)
	}
	if pair != "AUDUSD" {
		t.Errorf("Expecting 'AUDUSD', got '%s'", pair)
	}

	// then an invalid pair
	pair, err = ValidatePair("AUDUSDX")
	if err.Error() != "expecting a six character currency pair, got 'AUDUSDX'" {
		t.Error("Unexpected error string")
	}
}

func TestSplitPair(t *testing.T) {
	ccy1, ccy2, err := SplitPair("AUDUSD")
	if err != nil {
		t.Errorf("Error in SplitPair - %s", err)
	}
	if ccy1 != "AUD" {
		t.Errorf("Expecting 'AUD' as ccy1, got '%s'", ccy1)
	}
	if ccy2 != "USD" {
		t.Errorf("Expecting 'USD' as ccy2, got '%s'", ccy2)
	}

	_, _, err = SplitPair("AUDUSDX")
	if err.Error() != "expecting a six character currency pair, got 'AUDUSDX'" {
		t.Error("Unexpected error string")
	}
}
