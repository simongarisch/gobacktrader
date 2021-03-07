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
