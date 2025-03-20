package utils

import (
	"strings"
	"time"
)

// ConcatStrings concatenates multiple strings efficiently using strings.Builder
func ConcatStrings(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

// GetCurrentDateFormatted returns the current date in format ddmmyy
func GetCurrentDateFormatted() string {
	return time.Now().Format("020106")
}
