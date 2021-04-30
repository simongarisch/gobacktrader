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
	if err == nil {
		t.Error("Expecting a valid error to be returned.")
	}
	if GetErrorString(err) != "this is a valid error" {
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

func TestPadRight(t *testing.T) {
	if PadRight("abc", "x", 0) != "" {
		t.Error("Expecting ''")
	}
	if PadRight("abc", "x", 2) != "ab" {
		t.Error("Expecting 'ab'")
	}
	if PadRight("abc", "x", 3) != "abc" {
		t.Error("Expecting 'abc'")
	}
	if PadRight("abc", "x", 5) != "abcxx" {
		t.Error("Expecting 'abcxx'")
	}
}

func TestPadLeft(t *testing.T) {
	if PadLeft("abc", "x", 0) != "" {
		t.Error("Expecting ''")
	}
	if PadLeft("abc", "x", 2) != "bc" {
		t.Error("Expecting 'bc'")
	}
	if PadLeft("abc", "x", 3) != "abc" {
		t.Error("Expecting 'abc'")
	}
	if PadLeft("abc", "x", 5) != "xxabc" {
		t.Error("Expecting 'xxabc'")
	}
}

func TestSgn(t *testing.T) {
	if Sgn(+10.0) != +1.0 {
		t.Error("Unexpected sign")
	}
	if Sgn(-10.0) != -1.0 {
		t.Error("Unexpected sign")
	}
	if Sgn(0.0) != 0.0 {
		t.Error("Unexpected sign")
	}
}

func TestDate(t *testing.T) {
	date := Date(2021, 4, 1)
	if date.Year() != 2021 {
		t.Error("Unexpected year")
	}
	if date.Month().String() != "April" {
		t.Error("Unexpected month")
	}
	if date.Day() != 1 {
		t.Error("Unexpected day")
	}
	if date.Hour() != 0 {
		t.Error("Unexpected hour")
	}
	if date.Minute() != 0 {
		t.Error("Unexpected minute")
	}
	if date.Second() != 0 {
		t.Error("Unexpected second")
	}
	if zone, _ := date.Zone(); zone != "UTC" {
		t.Error("Unexpected time zone")
	}
}

func TestReplaceStrings(t *testing.T) {
	s := ReplaceStrings("mystring", nil)
	if s != "mystring" {
		t.Errorf("Unexpected string - wanted 'mystring', got '%s'", s)
	}

	s = "{STOCK} and {API_KEY}"
	replacements := map[string]string{
		"{STOCK}":   "AAPL",
		"{API_KEY}": "demo",
	}

	s = ReplaceStrings(s, replacements)
	if s != "AAPL and demo" {
		t.Errorf("Unexpected string - wanted 'AAPL and demo', got '%s'", s)
	}
}
