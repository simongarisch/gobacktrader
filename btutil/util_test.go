package btutil

import "testing"

func TestCleanString(t *testing.T) {
	result := CleanString(" xYz ")
	expectedResult := "XYZ"
	if result != expectedResult {
		t.Fatalf("Bad formatting in CleanString: wanted '%s', got '%s'", expectedResult, result)
	}
}
