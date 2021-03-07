package btutil

import "testing"

func TestCleanString(t *testing.T) {
	result := CleanString(" xYz ")
	expectedResult := "XYZ"
	if result != expectedResult {
		t.Fatalf("Bad formatting in CleanString: wanted '%s', got '%s'", expectedResult, result)
	}
}

func TestRound2dp(t *testing.T) {
	if Round2dp(1.222) != 1.22 {
		t.Error("Rounding down error.")
	}
	if Round2dp(1.216) != 1.22 {
		t.Error("Rounding up error.")
	}
}

func TestRound4dp(t *testing.T) {
	if Round4dp(1.22222) != 1.2222 {
		t.Error("Rounding down error.")
	}
	if Round4dp(1.21116) != 1.2112 {
		t.Error("Rounding up error.")
	}
}
