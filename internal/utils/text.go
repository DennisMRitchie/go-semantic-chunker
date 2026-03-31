package utils

import "strings"

// WordCount returns approximate word count for a string.
func WordCount(s string) int {
	return len(strings.Fields(s))
}

// TruncateWords returns the first n words of s.
func TruncateWords(s string, n int) string {
	fields := strings.Fields(s)
	if len(fields) <= n {
		return s
	}
	return strings.Join(fields[:n], " ")
}
