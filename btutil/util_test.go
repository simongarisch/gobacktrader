package btutil

import (
	"errors"
	"testing"
)

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

func TestAnyValidError(t *testing.T) {
	var e1, e2, e3 error
	if AnyValidError(e1, e2, e3) != nil {
		t.Error("Expecting no valid errors.")
	}

	e2 = errors.New("this is a valid error")
	err := AnyValidError(e1, e2, e3)
	if err.Error() != "this is a valid error" {
		t.Error("Unexpected error string.")
	}
}

func TestGetErrorString(t *testing.T) {
	if GetErrorString(nil) != "" {
		t.Error("Expecting a blank error string.")
	}

	s := "test one two three"
	err := errors.New(s)
	if GetErrorString(err) != s {
		t.Error("Unexpected error string.")
	}
}
