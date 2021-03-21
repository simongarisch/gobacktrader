package btutil

import (
	"math"
	"strings"
)

// CleanString cleans and returns some input string.
func CleanString(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

// Round2dp rounds a number to two decimal places.
func Round2dp(x float64) float64 {
	return math.Round(x*100) / 100
}

// Round4dp rounds a number to four decimal places.
func Round4dp(x float64) float64 {
	return math.Round(x*10000) / 10000
}

// AnyValidError takes a collection of errors and
// returns the first valid error if one exists, nil otherwise.
func AnyValidError(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}

// GetErrorString returns a string representation of the error.
// This will be blank if the error is nil.
func GetErrorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// PadRight pads a string on the right.
func PadRight(str, pad string, length uint) string {
	for {
		str += pad
		if len(str) > int(length) {
			return str[0:length]
		}
	}
}

// PadLeft pads a string on the left.
func PadLeft(str, pad string, length uint) string {
	for {
		str = pad + str
		if len(str) > int(length) {
			return str[(len(str) - int(length)):]
		}
	}
}
