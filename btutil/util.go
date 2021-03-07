package btutil

import "strings"

// CleanString cleans and returns some input string.
func CleanString(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}
